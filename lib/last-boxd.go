package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type Film struct {
	Title    string
	Year     string
	Director string
	Link     string
	Image    string
	Rating   string
	Rewatch  bool
	Like     bool
}

type TMDB struct {
	Poster string `json:"poster_path"`
}

func downloadFile(URL string) string {
	var imageData []byte
	if URL == "https://image.tmdb.org/t/p/w500" {
		var err error
		imageData, err = os.ReadFile("img/blank.jpg")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		response, err := http.Get(URL)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		imageData, err = io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return imageBase64
}

func GetLastBoxdTime(username string, dateLimit time.Time, qTitle, qDirector, qRating string) string {
	var (
		filmImages            []string = []string{}
		films                 []*Film  = []*Film{}
		image, year, director string
		directors             []string
		imageBase64           string
		done                  bool
		page                  int = 1
		col, row              int
		entryCount            int
		wg                    sync.WaitGroup
	)

	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
		colly.Async(false),
	)

	c.OnHTML(".table.film-table", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("tr.diary-entry-row", func(_ int, el *colly.HTMLElement) bool {
			if entryCount == 100 {
				done = true
				return false
			}

			date := strings.Join(strings.Split(el.ChildAttr(".td-day.diary-day.center > a", "href"), "/")[5:], "-")
			date = date[:len(date)-1]

			currentDate, err := time.Parse("2006-01-02", date)
			if err != nil {
				log.Fatal(err)
				done = true
				return false
			}
			if currentDate.Before(dateLimit) {
				done = true
				return false
			}

			title := el.ChildText(".-primary")
			href := el.ChildAttr(".-primary > a", "href")
			log.Println("Href:", href)
			splitHref := strings.Split(href, "/")
			link := "https://letterboxd.com/" + strings.Join(splitHref[2:4], "/")
			rating := el.ChildText("span.rating")
			rewatch := false
			like := false

			rewatchClasses := strings.Split(el.ChildAttr(".td-rewatch", "class"), " ")
			if len(rewatchClasses) == 2 {
				rewatch = true
			}

			likeChilds := el.DOM.Find(".td-like").Children().Length()
			if likeChilds == 2 {
				like = true
			}

			films = append(films, &Film{Title: title, Link: link, Rating: rating, Rewatch: rewatch, Like: like})
			entryCount += 1

			return true
		})
		wg.Done()
	})

	for !done {
		wg.Add(1)
		c.Visit("https://letterboxd.com/" + username + "/films/diary/page/" + fmt.Sprintf("%d", page))
		page += 1
	}
	wg.Wait()

	c.OnHTML(".text-link.text-footer", func(e *colly.HTMLElement) {
		siteLinks := e.ChildAttrs("a", "href")
		var tmdbURLSplit []string
		if len(siteLinks) > 2 {
			tmdbURLSplit = strings.Split(siteLinks[1], "/")
		} else {
			if siteLinks[0] == "" {
				image = "https://image.tmdb.org/t/p/w500"
				return
			}
			tmdbURLSplit = strings.Split(siteLinks[0], "/")
		}
		tmdbID := tmdbURLSplit[len(tmdbURLSplit)-2]
		tmdbType := tmdbURLSplit[len(tmdbURLSplit)-3]

		apiKey := fmt.Sprint(viper.Get("API_KEY"))
		resp, err := http.Get("https://api.themoviedb.org/3/" + tmdbType + "/" + tmdbID + "?api_key=" + apiKey)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var tmdb TMDB
		if err := json.NewDecoder(resp.Body).Decode(&tmdb); err != nil {
			log.Fatal(err)
		}
		image = "https://image.tmdb.org/t/p/w500" + tmdb.Poster
	})

	c.OnHTML("div.metablock", func(e *colly.HTMLElement) {
		year = e.ChildText("div.releaseyear > a")
		e.ForEach("span.prettify", func(_ int, elem *colly.HTMLElement) {
			dir := elem.Text
			directors = append(directors, dir)
		})
		if len(directors) > 1 {
			for i, dir := range directors {
				if i == len(directors)-1 {
					director += dir
				} else {
					director += dir + ", "
				}
			}
		} else if len(directors) == 1 {
			director = directors[0]
		}
	})

	for _, film := range films {
		c.Visit(film.Link)
		film.Image = image
		film.Year = year
		film.Director = director
		directors = nil
		director = ""
	}

	for _, film := range films {
		imageBase64 := downloadFile(film.Image)
		if imageBase64 == "" {
			continue
		}
		imageBase64 = DrawText(*film, imageBase64, qTitle, qDirector, qRating)
		filmImages = append(filmImages, imageBase64)
	}

	lenSqrt := math.Sqrt(float64(len(filmImages)))
	col = int(math.Round(lenSqrt))
	if col == int(math.Ceil(lenSqrt)) {
		row = col
	} else {
		row = col + 1
	}

	imageBase64 = MakeGrid(filmImages, col, row)

	return imageBase64
}

func GetLastBoxd(username string, col, row int, qTitle, qDirector, qRating string) string {
	var (
		filmImages            []string = []string{}
		films                 []*Film  = []*Film{}
		image, year, director string
		directors             []string
		entryCount            int
		page                  int = int(math.Ceil(float64(col*row) / 50.0))
		imageBase64           string
		wg                    sync.WaitGroup
	)

	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
		colly.Async(false),
	)

	c.OnHTML(".table.film-table", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("tr.diary-entry-row", func(_ int, el *colly.HTMLElement) bool {
			title := el.ChildText(".-primary")
			href := el.ChildAttr(".-primary > a", "href")
			log.Println("Href:", href)
			splitHref := strings.Split(href, "/")
			link := "https://letterboxd.com/" + strings.Join(splitHref[2:4], "/")
			rating := el.ChildText("span.rating")
			rewatch := false
			like := false

			rewatchClasses := strings.Split(el.ChildAttr(".td-rewatch", "class"), " ")
			if len(rewatchClasses) == 2 {
				rewatch = true
			}

			likeChilds := el.DOM.Find(".td-like").Children().Length()
			if likeChilds == 2 {
				like = true
			}

			films = append(films, &Film{Title: title, Link: link, Rating: rating, Rewatch: rewatch, Like: like})
			entryCount += 1

			return !(entryCount == col*row)
		})
		wg.Done()
	})

	for i := 1; i <= page; i++ {
		wg.Add(1)
		c.Visit("https://letterboxd.com/" + username + "/films/diary/page/" + fmt.Sprintf("%d", i))
	}
	wg.Wait()

	c.OnHTML(".text-link.text-footer", func(e *colly.HTMLElement) {
		siteLinks := e.ChildAttrs("a", "href")
		var tmdbURLSplit []string
		if len(siteLinks) > 2 {
			tmdbURLSplit = strings.Split(siteLinks[1], "/")
		} else {
			if siteLinks[0] == "" {
				image = "https://image.tmdb.org/t/p/w500"
				return
			}
			tmdbURLSplit = strings.Split(siteLinks[0], "/")
		}
		tmdbID := tmdbURLSplit[len(tmdbURLSplit)-2]
		tmdbType := tmdbURLSplit[len(tmdbURLSplit)-3]

		apiKey := fmt.Sprint(viper.Get("API_KEY"))
		resp, err := http.Get("https://api.themoviedb.org/3/" + tmdbType + "/" + tmdbID + "?api_key=" + apiKey)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var tmdb TMDB
		if err := json.NewDecoder(resp.Body).Decode(&tmdb); err != nil {
			log.Fatal(err)
		}
		image = "https://image.tmdb.org/t/p/w500" + tmdb.Poster
	})

	c.OnHTML(".productioninfo", func(e *colly.HTMLElement) {
		year = e.ChildText(".releasedate > a")
		e.ForEach("span.prettify", func(_ int, elem *colly.HTMLElement) {
			dir := elem.Text
			directors = append(directors, dir)
		})
		if len(directors) > 1 {
			for i, dir := range directors {
				if i == len(directors)-1 {
					director += dir
				} else {
					director += dir + ", "
				}
			}
		} else if len(directors) == 1 {
			director = directors[0]
		}
	})

	for _, film := range films {
		c.Visit(film.Link)
		film.Image = image
		film.Year = year
		film.Director = director
		directors = nil
		director = ""
	}

	for _, film := range films {
		imageBase64 := downloadFile(film.Image)
		if imageBase64 == "" {
			continue
		}
		imageBase64 = DrawText(*film, imageBase64, qTitle, qDirector, qRating)
		filmImages = append(filmImages, imageBase64)
	}

	imageBase64 = MakeGrid(filmImages, col, row)

	return imageBase64
}

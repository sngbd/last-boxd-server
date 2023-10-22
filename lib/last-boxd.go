package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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
	Rewatch   bool
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

func GetLastBoxd(username string, col, row int, qTitle, qDirector, qRating string) string {
	filmImages := []string{}
	films := []*Film{}

	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
		colly.Async(false),
	)

	c.OnHTML(".table.film-table", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("tr.diary-entry-row", func(i int, el *colly.HTMLElement) bool {
			title := el.ChildText("h3.headline-3.prettify")
			link := "https://letterboxd.com/" + strings.Join(strings.Split(el.ChildAttr("h3.headline-3.prettify > a", "href"), "/")[2:4], "/")
			rating := el.ChildText("span.rating")
			rewatch := false
			classes := strings.Split(el.ChildAttr(".td-rewatch", "class"), " ");
			if (len(classes) == 2) {
				rewatch = true
			}
			films = append(films, &Film{Title: title, Link: link, Rating: rating, Rewatch: rewatch})
			return !(i+1 == col*row)
		})
	})
	c.Visit("https://letterboxd.com/" + username + "/films/diary/")

	var image, year, director string
	var directors []string

	c.OnHTML(".text-link.text-footer", func(e *colly.HTMLElement) {
		siteLinks := e.ChildAttrs("a", "href")
		var tmdbURLSplit []string
		if len(siteLinks) > 2 {
			tmdbURLSplit = strings.Split(siteLinks[1], "/")
		} else {
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

	c.OnHTML("#featured-film-header", func(e *colly.HTMLElement) {
		year = e.ChildText("small.number")
		e.ForEach("span.prettify", func(_ int, elem *colly.HTMLElement) {
			dir := elem.Text
			directors = append(directors, dir)
		})
		if (len(directors) > 1) {
			for i, dir := range directors {
				if (i == len(directors) - 1) {
					director += dir
				} else {
					director += dir + ", "
				}
			}
		} else {
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

	return MakeGrid(filmImages, col, row)
}

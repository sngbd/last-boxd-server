package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
}

type TMDB struct {
	Poster string `json:"poster_path"`
}

func downloadFile(URL string) string {
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return imageBase64
}

func GetLastBoxd(username string, grid int, details string) string {
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
			films = append(films, &Film{Title: title, Link: link, Rating: rating})
			return !(i+1 == grid*grid)
		})
	})
	c.Visit("https://letterboxd.com/" + username + "/films/diary/")

	var image, year, director string

	c.OnHTML(".text-link.text-footer", func(e *colly.HTMLElement) {
		siteLinks := e.ChildAttrs("a", "href")
		var tmdbURLSplit []string
		if len(siteLinks) < 2 {
			tmdbURLSplit = strings.Split(siteLinks[1], "/")
		} else {
			tmdbURLSplit = strings.Split(siteLinks[0], "/")
		}
		tmdbMovieID := tmdbURLSplit[len(tmdbURLSplit)-2]

		apiKey := fmt.Sprint(viper.Get("API_KEY"))
		resp, err := http.Get("https://api.themoviedb.org/3/movie/" + tmdbMovieID + "?api_key=" + apiKey)
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
		director = e.ChildText("span.prettify")
	})

	for _, film := range films {
		c.Visit(film.Link)
		film.Image = image
		film.Year = year
		film.Director = director
	}

	for _, film := range films {
		imageBase64 := downloadFile(film.Image)
		if details != "off" {
			imageBase64 = DrawText(*film, imageBase64)
		}
		filmImages = append(filmImages, imageBase64)
	}

	return MakeGrid(filmImages, grid)
}

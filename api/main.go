package api

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"github.com/sngbd/last-boxd/lib"
)

func downloadFile(URL string) string {
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read the image data from the response body
	imageData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the image data as base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return imageBase64
}

func GetLastBoxd(username string) string {
	var grid int = 3
	filmImages := []string{}
	films := []*lib.Film{}
	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
	)
	c.OnHTML(".table.film-table", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("td.td-film-details", func(i int, el *colly.HTMLElement) bool {
			title := el.ChildText("h3.headline-3.prettify")
			link := "https://letterboxd.com/" + strings.Join(strings.Split(el.ChildAttr("a", "href"), "/")[2:4], "/")
			films = append(films, &lib.Film{Title: title, Link: link})
			return !(i+1 == grid*grid)
		})
	})
	c.Visit("https://letterboxd.com/" + username + "/films/diary/")

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var image string
	var year string
	var director string
	for _, film := range films {
		err := chromedp.Run(ctx,
			emulation.SetUserAgentOverride("WebScraper 1.0"),
			chromedp.Navigate(film.Link),
			chromedp.WaitVisible(`#poster-zoom`),
			chromedp.Evaluate(`(function() {return document.querySelector("img").getAttribute("src");})();`, &image),
			chromedp.Evaluate(`(function() {return document.querySelector("small.number").innerText;})();`, &year),
			chromedp.Evaluate(`(function() {return document.querySelector("span.prettify").innerText;})();`, &director),
		)
		if err != nil {
			log.Fatal(err)
		}
		film.Image = image
		film.Year = year
		film.Director = director
	}

	for _, film := range films {
		imageBase64 := lib.DrawText(*film, downloadFile(film.Image))
		filmImages = append(filmImages, imageBase64)
	}

	return lib.MakeGrid(filmImages)
}

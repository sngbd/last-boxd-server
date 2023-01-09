package lib

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

type Film struct {
	Title    string
	Year     string
	Director string
	Link     string
	Image    string
	Rating   string
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

func GetLastBoxd(username string, grid int, text string) string {
	filmImages := []string{}
	films := []*Film{}
	c := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
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
		imageBase64 := downloadFile(film.Image)
		if text != "off"{
			imageBase64 = DrawText(*film, imageBase64)
		}
		filmImages = append(filmImages, imageBase64)
	}

	return MakeGrid(filmImages, grid)
}

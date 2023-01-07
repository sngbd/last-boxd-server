package api

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sngbd/last-boxd/lib"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetLastBoxd(username string) {
	var grid int = 3
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

	for i, film := range films {
		fileName := strconv.Itoa(i) + ".jpg"
		err := downloadFile(film.Image, fileName)
		if err != nil {
			log.Fatal(err)
		}
		lib.DrawText(fileName, *film)
	}

	lib.MakeGrid()
}

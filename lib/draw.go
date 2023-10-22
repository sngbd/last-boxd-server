package lib

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var fontBytes []byte
var fnt *truetype.Font
var face font.Face
var dc *font.Drawer
var rgba *image.RGBA
var white color.RGBA
var img image.Image

func createOverlay(width, height int, colorRGBA color.RGBA) image.Image {
	overlay := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(overlay, overlay.Bounds(), &image.Uniform{colorRGBA}, image.Point{}, draw.Over)

	return overlay
}

func insertText(img *image.RGBA, text string, x, y int, face font.Face, c color.Color) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
}

func drawTextTriple(title, director, rating string) {
	titleWidth := dc.MeasureString(title).Ceil()
	directorWidth := dc.MeasureString(director).Ceil()
	ratingWidth := dc.MeasureString(rating).Ceil()

	maxTextWidth := titleWidth
	if directorWidth > maxTextWidth {
		maxTextWidth = directorWidth
	}
	if ratingWidth > maxTextWidth {
		maxTextWidth = ratingWidth
	}

	textHeight := (face.Metrics().Height.Ceil() + 5) * 3
	textWidth := maxTextWidth + 14

	overlayColor := color.RGBA{0, 0, 0, 100}

	overlay := createOverlay(textWidth, textHeight, overlayColor)

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	draw.Draw(rgba, overlay.Bounds().Add(image.Pt(0, 0)), overlay, image.Point{}, draw.Over)

	insertText(rgba, title, 7, 24, face, white)
	insertText(rgba, director, 7, 24+textHeight/3, face, white)
	insertText(rgba, rating, 7, 24+2*(textHeight/3), face, white)
}

func drawTextDouble(text1, text2 string) {
	text1Width := dc.MeasureString(text1).Ceil()
	text2Width := dc.MeasureString(text2).Ceil()

	maxTextWidth := text1Width
	if text2Width > maxTextWidth {
		maxTextWidth = text2Width
	}

	textHeight := (face.Metrics().Height.Ceil() + 5) * 2
	textWidth := maxTextWidth + 14

	overlayColor := color.RGBA{0, 0, 0, 100}

	overlay := createOverlay(textWidth, textHeight, overlayColor)

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	draw.Draw(rgba, overlay.Bounds().Add(image.Pt(0, 0)), overlay, image.Point{}, draw.Over)

	insertText(rgba, text1, 7, 24, face, white)
	insertText(rgba, text2, 7, 24+textHeight/2, face, white)
}

func drawTextSingle(text string) {
	textWidth := dc.MeasureString(text).Ceil() + 14
	textHeight := (face.Metrics().Height.Ceil() + 5)

	overlayColor := color.RGBA{0, 0, 0, 100}

	overlay := createOverlay(textWidth, textHeight, overlayColor)

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	draw.Draw(rgba, overlay.Bounds().Add(image.Pt(0, 0)), overlay, image.Point{}, draw.Over)

	insertText(rgba, text, 7, 24, face, white)
}

func DrawText(film Film, imageBase64, qTitle, qDirector, qRating string) string {
	imageByte, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err = image.Decode(bytes.NewReader(imageByte))
	if err != nil {
		log.Fatal(err)
	}

	fontBytes, err = os.ReadFile("font/FreeSerif.ttf")
	if err != nil {
		log.Fatal(err)
	}

	fnt, err = truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	face = truetype.NewFace(fnt, &truetype.Options{
		Size:    25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	dc = &font.Drawer{
		Face: face,
	}

	rgba = image.NewRGBA(img.Bounds())

	white = color.RGBA{255, 255, 255, 255}

	if (film.Rewatch) {
		film.Rating += " ↻"
	}

	if qTitle == "on" && qDirector == "on" && qRating == "on" {
		drawTextTriple(film.Title+" ("+film.Year+")", film.Director, film.Rating)
	} else if qTitle == "on" && qDirector == "on" {
		drawTextDouble(film.Title+" ("+film.Year+")", film.Director)
	} else if qTitle == "on" && qRating == "on" {
		drawTextDouble(film.Title+" ("+film.Year+")", film.Rating)
	} else if qDirector == "on" && qRating == "on" {
		drawTextDouble(film.Director, film.Rating)
	} else if qTitle == "on" {
		drawTextSingle(film.Title + " (" + film.Year + ")")
	} else if qDirector == "on" {
		drawTextSingle(film.Director)
	} else if qRating == "on" {
		drawTextSingle(film.Rating)
	} else {
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, nil)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

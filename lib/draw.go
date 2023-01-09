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

func drawOutlinedText(img *image.RGBA, text string, x, y int, face font.Face, c, outline color.Color) {
	insertText(img, text, x-1, y, face, outline)
	insertText(img, text, x+1, y, face, outline)
	insertText(img, text, x, y-1, face, outline)
	insertText(img, text, x, y+1, face, outline)
	insertText(img, text, x, y, face, c)
}

func DrawText(film Film, imageBase64 string) string {
	imageByte, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageByte))
	if err != nil {
		log.Fatal(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	fontBytes, err := os.ReadFile("font/FreeSerif.ttf")
	if err != nil {
		log.Fatal(err)
	}
	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(fnt, &truetype.Options{
		Size:    28,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	drawOutlinedText(rgba, film.Title + " (" + film.Year + ")", 7, 27, face, white, black)
	drawOutlinedText(rgba, film.Director, 7, 52, face, white, black)
	drawOutlinedText(rgba, film.Rating, 7, 77, face, white, black)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, nil)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

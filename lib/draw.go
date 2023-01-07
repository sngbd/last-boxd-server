package lib

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Film struct {
	Title    string
	Year     string
	Director string
	Link     string
	Image    string
}

func DrawText(film Film, imageBase64 string) string {
	// Decode the image
	imageByte, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageByte))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new RGBA image
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// Load the font
	fontBytes, err := ioutil.ReadFile("./Inconsolata.TTF")
	if err != nil {
		log.Fatal(err)
	}
	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Create a font face
	face := truetype.NewFace(fnt, &truetype.Options{
		Size:    28, // font size in points
		DPI:     72, // screen resolution in DPI
		Hinting: font.HintingFull,
	})

	// Draw the text
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(10 * 64), fixed.Int26_6(20 * 64)}
	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(film.Title + " (" + film.Year + ")")
	point = fixed.Point26_6{fixed.Int26_6(10 * 64), fixed.Int26_6(40 * 64)}
	d = &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(film.Director)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		log.Fatal(err)
	}

	imageBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())

	return imageBase64
}

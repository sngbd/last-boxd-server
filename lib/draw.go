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
    	"github.com/nfnt/resize"
    	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func createBlurredOverlay(width, height int, colorRGBA color.RGBA, opacity uint8) image.Image {
    overlay := image.NewRGBA(image.Rect(0, 0, width, height))

    for y := overlay.Bounds().Min.Y; y < overlay.Bounds().Max.Y; y++ {
        alpha := uint8(float64(y) / float64(height) * float64(opacity))
        for x := overlay.Bounds().Min.X; x < overlay.Bounds().Max.X; x++ {
            overlay.Set(x, y, color.RGBA{0, 0, 0, alpha})
        }
    }

    overlay2 := resize.Resize(uint(width), uint(height), overlay, resize.NearestNeighbor)
    overlay3 := imaging.Blur(overlay2, 5)

    draw.Draw(overlay3, overlay3.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, opacity}}, image.Point{}, draw.Over)

    return overlay3
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

	overlayColor := color.RGBA{23, 23, 23, 100}
	overlayOpacity := uint8(100)

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

	dc := &font.Drawer{
        	Face: face,
    	}

	titleWidth := dc.MeasureString(film.Title + " (" + film.Year + ")").Ceil()
	directorWidth := dc.MeasureString(film.Director).Ceil()
	ratingWidth := dc.MeasureString(film.Rating).Ceil()

	maxTextWidth := titleWidth
	if directorWidth > maxTextWidth {
		maxTextWidth = directorWidth
	}
	if ratingWidth > maxTextWidth {
		maxTextWidth = ratingWidth
	}

    	textHeight := (face.Metrics().Height.Ceil() + 5) * 3
	textWidth := maxTextWidth + 14

	overlay := createBlurredOverlay(textWidth, textHeight, overlayColor, overlayOpacity)

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	draw.Draw(rgba, overlay.Bounds().Add(image.Pt(0, 0)), overlay, image.Point{}, draw.Over)

	white := color.RGBA{255, 255, 255, 255}

	insertText(rgba, film.Title+" ("+film.Year+")", 7, 27, face, white)
	insertText(rgba, film.Director, 7, 27+textHeight/3, face, white)
	insertText(rgba, film.Rating, 7, 27+2*(textHeight/3), face, white)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, rgba, nil)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

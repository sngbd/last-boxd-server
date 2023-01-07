package lib

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
)

func MakeGrid(filmImages []string) string {
	img := image.NewRGBA(image.Rect(0, 0, 1500, 2250))

	var images []image.Image

	for _, imageBase64 := range filmImages {
		imageByte, err := base64.StdEncoding.DecodeString(imageBase64)
		if err != nil {
			log.Fatal(err)
		}

		img, _, err := image.Decode(bytes.NewReader(imageByte))
		if err != nil {
			log.Fatal(err)
		}

		images = append(images, img)
	}

	y0 := 0
	y1 := 750
	index := 0
	for i := 0; i < 3; i++ {
		x0 := 0
		x1 := 500
		for j := 0; j < 3; j++ {
			draw.Draw(img, image.Rect(x0, y0, x1, y1), images[index], image.Point{0, 0}, draw.Src)
			x0 += 500
			x1 += 500
			index += 1
		}
		y0 += 750
		y1 += 750
	}

	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, nil)
	if err != nil {
		log.Fatal(err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return imageBase64
}

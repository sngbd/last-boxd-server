package lib

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
)

func MakeGrid(filmImages []string, grid int) string {
	img := image.NewRGBA(image.Rect(0, 0, grid*500, grid*750))

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

merge:
	for i := 0; i < grid; i++ {
		x0 := 0
		x1 := 500
		for j := 0; j < grid; j++ {
			if len(images) == 0 {
				break merge
			}
			draw.Draw(img, image.Rect(x0, y0, x1, y1), images[0], image.Point{0, 0}, draw.Src)
			images = images[1:]
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

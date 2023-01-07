package lib

import (
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strconv"
)

func MakeGrid() {
	img := image.NewRGBA(image.Rect(0, 0, 1500, 2250))

	y0 := 0
	y1 := 750
	index := 0
	for i := 0; i < 3; i++ {
		x0 := 0
		x1 := 500
		for j := 0; j < 3; j++ {
			imageFile, err := os.Open(strconv.Itoa(index) + ".jpg")
			if err != nil {
				log.Fatal(err)
			}
			newImage, err := jpeg.Decode(imageFile)
			if err != nil {
				log.Fatal(err)
			}
			draw.Draw(img, image.Rect(x0, y0, x1, y1), newImage, image.Point{0, 0}, draw.Src)
			x0 += 500
			x1 += 500
			index += 1
		}
		y0 += 750
		y1 += 750
	}

	f, err := os.Create("grid.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}

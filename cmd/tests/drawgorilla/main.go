package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	imgdraw "github.com/arran4/gorillas/drawings/img"
)

func main() {
	scale := 2.0
	width, height := 90, 59
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	sky := color.RGBA{0, 0, 170, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, sky)
		}
	}

	orange := color.RGBA{255, 170, 85, 255}
	imgdraw.DrawBASGorilla(img, 15*scale+11, scale+12, scale, imgdraw.ArmsDown, orange)

	f, err := os.Create("gorilla.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

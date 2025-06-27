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
	width, height := 90, 90
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	sky := color.RGBA{0, 0, 170, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, sky)
		}
	}

	orange := color.RGBA{255, 170, 85, 255}
	x := (float64(width) + scale) / 2
	y := (float64(height) - 23*scale) / 2
	imgdraw.DrawBASGorilla(img, x, y, scale, imgdraw.ArmsDown, orange)

	f, err := os.Create("gorilla.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

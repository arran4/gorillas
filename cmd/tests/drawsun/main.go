package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	imgdraw "github.com/arran4/gorillas/drawings/img"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	imgdraw.DrawBASSun(img, 50, 50, 40, false, color.RGBA{255, 255, 0, 255})
	f, err := os.Create("sun.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

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
	img := image.NewRGBA(image.Rect(0, 0, int(30*scale), int(30*scale)))
	imgdraw.DrawBASGorilla(img, 15*scale, scale, scale, imgdraw.ArmsDown, color.RGBA{150, 75, 0, 255})
	f, err := os.Create("gorilla.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

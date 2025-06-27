package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"

	imgdraw "github.com/arran4/gorillas/drawings/img"
	"image/color/palette"
)

func makeFrame(scale float64, width, height int, arms int) *image.Paletted {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	sky := color.RGBA{0, 0, 170, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, sky)
		}
	}

	// Coordinates for the gorilla
	x := (float64(width) + scale) / 2
	y := (float64(height) - 23*scale) / 2

	// Draw building under gorilla feet
	bW := int(20 * scale)
	bH := height - int(y+23*scale)
	if bH < 1 {
		bH = 1
	}
	building := imgdraw.CreateBuildingSprite(float64(bW), float64(bH), color.RGBA{100, 100, 100, 255})
	bx := int(x - float64(bW)/2 + 0.5)
	by := height - bH
	draw.Draw(img, image.Rect(bx, by, bx+bW, by+bH), building, image.Point{}, draw.Over)

	orange := color.RGBA{255, 170, 85, 255}
	imgdraw.DrawBASGorilla(img, x, y, scale, arms, orange)

	pal := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(pal, img.Bounds(), img, image.Point{})
	return pal
}

func main() {
	scale := 2.0
	width, height := 90, 90

	frames := []int{imgdraw.ArmsDown, imgdraw.ArmsRightUp, imgdraw.ArmsDown, imgdraw.ArmsLeftUp}
	outGIF := &gif.GIF{}
	for _, arm := range frames {
		frame := makeFrame(scale, width, height, arm)
		outGIF.Image = append(outGIF.Image, frame)
		outGIF.Delay = append(outGIF.Delay, 100) // 100*10ms = 1s
	}

	// Also save the first frame as PNG for convenience
	fPNG, err := os.Create("gorilla.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(fPNG, outGIF.Image[0]); err != nil {
		panic(err)
	}
	fPNG.Close()

	fGIF, err := os.Create("gorilla.gif")
	if err != nil {
		panic(err)
	}
	defer fGIF.Close()
	if err := gif.EncodeAll(fGIF, outGIF); err != nil {
		panic(err)
	}
}

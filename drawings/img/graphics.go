package imgdraw

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"

	gorillas "github.com/arran4/gorillas"
	drawcommon "github.com/arran4/gorillas/drawings/common"
)

// DrawVectorLines joins a series of points with lines of the specified colour.
func DrawVectorLines(img draw.Image, pts []gorillas.VectorPoint, clr color.Color) {
	drawcommon.DrawVectorLines(img, pts, clr)
}

// DrawFilledRect fills a rectangle between two points.
func DrawFilledRect(img draw.Image, x1, y1, x2, y2 float64, clr color.Color) {
	drawcommon.DrawFilledRect(img, x1, y1, x2, y2, clr)
}

// DrawFilledCircle draws a filled circle at the provided centre and radius.
func DrawFilledCircle(img draw.Image, cx, cy, r float64, clr color.Color) {
	drawcommon.DrawFilledCircle(img, cx, cy, r, clr)
}

// DrawFilledEllipse draws a filled ellipse with horizontal radius rx and vertical radius ry.
func DrawFilledEllipse(img draw.Image, cx, cy, rx, ry float64, clr color.Color) {
	drawcommon.DrawFilledEllipse(img, cx, cy, rx, ry, clr)
}

// DrawArc renders an arc from startDeg to endDeg going clockwise from 0 degrees to the right.
func DrawArc(img draw.Image, cx, cy, r float64, startDeg, endDeg float64, clr color.Color) {
	drawcommon.DrawArc(img, cx, cy, r, startDeg, endDeg, clr)
}

// DrawThickArc draws a thicker arc using overlapping circles.
func DrawThickArc(img draw.Image, cx, cy, r, thickness, startDeg, endDeg float64, clr color.Color) {
	drawcommon.DrawThickArc(img, cx, cy, r, thickness, startDeg, endDeg, clr)
}

const (
	ArmsRightUp = drawcommon.ArmsRightUp
	ArmsLeftUp  = drawcommon.ArmsLeftUp
	ArmsDown    = drawcommon.ArmsDown
)

// DrawBASSun renders the classic QBASIC sun sprite at the given position and radius.
func DrawBASSun(img draw.Image, cx, cy, r float64, shocked bool, clr color.Color) {
	drawcommon.DrawBASSun(img, cx, cy, r, shocked, clr)
}

// DrawBASGorilla draws a simple approximation of the QBASIC gorilla sprite.
func DrawBASGorilla(img draw.Image, x, y, scale float64, arms int, clr color.Color) {
	drawcommon.DrawBASGorilla(img, x, y, scale, arms, clr)
}

// CreateBananaSprite converts an ASCII mask into an RGBA image using yellow pixels.
func CreateBananaSprite(mask []string) *image.RGBA {
	h := len(mask)
	w := len(mask[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	clr := color.RGBA{255, 255, 0, 255}
	for y, row := range mask {
		for x, c := range row {
			if c != '.' {
				img.Set(x, y, clr)
			}
		}
	}
	return img
}

// CreateBananaSprites builds small sprites for the four banana rotations.
func CreateBananaSprites() (left, right, up, down *image.RGBA) {
	left = CreateBananaSprite([]string{
		"..##.",
		".###.",
		"#####",
		".###.",
		"..##.",
	})
	right = CreateBananaSprite([]string{
		".##..",
		".###.",
		"#####",
		".###.",
		".##..",
	})
	up = CreateBananaSprite([]string{
		"..#..",
		".###.",
		"..#..",
		"..#..",
		"..#..",
	})
	down = CreateBananaSprite([]string{
		"..#..",
		"..#..",
		"..#..",
		".###.",
		"..#..",
	})
	return
}

// CreateGorillaSprite converts an ASCII mask into a coloured image.
func CreateGorillaSprite(mask []string, clr color.Color) *image.RGBA {
	h := len(mask)
	w := len(mask[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y, row := range mask {
		for x, c := range row {
			if c != '.' {
				img.Set(x, y, clr)
			}
		}
	}
	return img
}

// DefaultGorillaSprite returns a basic gorilla sprite rendered at the provided scale.
func DefaultGorillaSprite(scale float64) *image.RGBA {
	size := int(30 * scale)
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	clr := color.RGBA{150, 75, 0, 255}
	DrawBASGorilla(img, 15*scale, scale, scale, ArmsDown, clr)
	return img
}

// CreateBuildingSprite produces a simple building with random windows.
func CreateBuildingSprite(w, h float64, clr color.Color) *image.RGBA {
	iw := int(w)
	ih := int(h)
	img := image.NewRGBA(image.Rect(0, 0, iw, ih))
	for x := 0; x < iw; x++ {
		for y := 0; y < ih; y++ {
			img.Set(x, y, clr)
		}
	}
	winClr := color.RGBA{255, 255, 0, 255}
	for x := 3; x < iw-3; x += 6 {
		for y := ih - 3; y > 3; y -= 6 {
			if rand.Intn(3) != 0 {
				for dx := 0; dx < 3; dx++ {
					for dy := 0; dy < 3; dy++ {
						img.Set(x+dx, y+dy, winClr)
					}
				}
			}
		}
	}
	return img
}

// ClearRect clears a rectangle region on the image by setting pixels transparent.
func ClearRect(img draw.Image, x, y, w, h int) {
	drawcommon.ClearRect(img, x, y, w, h)
}

// ClearCircle clears a circle region on the image by setting pixels transparent.
func ClearCircle(img draw.Image, cx, cy int, r float64) {
	drawcommon.ClearCircle(img, cx, cy, r)
}

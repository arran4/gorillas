package imgdraw

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	gorillas "github.com/arran4/gorillas"
)

func drawLine(img image.Image, x1, y1, x2, y2 float64, clr color.Color) {
	drw, ok := img.(interface{ Set(int, int, color.Color) })
	if !ok {
		return
	}
	dx := x2 - x1
	dy := y2 - y1
	n := int(math.Max(math.Abs(dx), math.Abs(dy)))
	if n == 0 {
		drw.Set(int(math.Round(x1)), int(math.Round(y1)), clr)
		return
	}
	for i := 0; i <= n; i++ {
		x := int(math.Round(x1 + dx*float64(i)/float64(n)))
		y := int(math.Round(y1 + dy*float64(i)/float64(n)))
		drw.Set(x, y, clr)
	}
}

// DrawVectorLines joins a series of points with lines of the specified colour.
func DrawVectorLines(img image.Image, pts []gorillas.VectorPoint, clr color.Color) {
	if len(pts) == 0 {
		return
	}
	prev := pts[0]
	for _, p := range pts[1:] {
		drawLine(img, prev.X, prev.Y, p.X, p.Y, clr)
		prev = p
	}
}

// DrawFilledRect fills a rectangle between two points.
func DrawFilledRect(img image.Image, x1, y1, x2, y2 float64, clr color.Color) {
	drw, ok := img.(interface{ Set(int, int, color.Color) })
	if !ok {
		return
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	for x := int(math.Round(x1)); x <= int(math.Round(x2)); x++ {
		for y := int(math.Round(y1)); y <= int(math.Round(y2)); y++ {
			drw.Set(x, y, clr)
		}
	}
}

// DrawFilledCircle draws a filled circle at the provided centre and radius.
func DrawFilledCircle(img image.Image, cx, cy, r float64, clr color.Color) {
	drw, ok := img.(interface{ Set(int, int, color.Color) })
	if !ok {
		return
	}
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			if dx*dx+dy*dy <= r*r {
				x := int(math.Round(cx + dx))
				y := int(math.Round(cy + dy))
				drw.Set(x, y, clr)
			}
		}
	}
}

// DrawArc renders an arc from startDeg to endDeg going clockwise from 0 degrees to the right.
func DrawArc(img image.Image, cx, cy, r float64, startDeg, endDeg float64, clr color.Color) {
	if endDeg < startDeg {
		endDeg += 360
	}
	step := 2.0
	prevX := cx + r*math.Cos(startDeg*math.Pi/180)
	prevY := cy - r*math.Sin(startDeg*math.Pi/180)
	for a := startDeg + step; a <= endDeg; a += step {
		x := cx + r*math.Cos(a*math.Pi/180)
		y := cy - r*math.Sin(a*math.Pi/180)
		drawLine(img, prevX, prevY, x, y, clr)
		prevX, prevY = x, y
	}
}

const (
	ArmsRightUp = 1
	ArmsLeftUp  = 2
	ArmsDown    = 3
)

// DrawBASSun renders the classic QBASIC sun sprite at the given position and radius.
func DrawBASSun(img image.Image, cx, cy, r float64, shocked bool, clr color.Color) {
	DrawFilledCircle(img, cx, cy, r, clr)
	scale := r / 12
	lines := [][4]float64{
		{-20, 0, 20, 0},
		{0, -15, 0, 15},
		{-15, -10, 15, 10},
		{-15, 10, 15, -10},
		{-8, -13, 8, 13},
		{-8, 13, 8, -13},
		{-18, -5, 18, 5},
		{-18, 5, 18, -5},
	}
	for _, l := range lines {
		x1 := cx + l[0]*scale
		y1 := cy + l[1]*scale
		x2 := cx + l[2]*scale
		y2 := cy + l[3]*scale
		drawLine(img, x1, y1, x2, y2, clr)
	}
	eyeX := 3 * scale
	eyeY := -2 * scale
	DrawFilledCircle(img, cx-eyeX, cy+eyeY, 1*scale, color.Black)
	DrawFilledCircle(img, cx+eyeX, cy+eyeY, 1*scale, color.Black)
	if shocked {
		DrawFilledCircle(img, cx, cy+5*scale, 2.9*scale, color.Black)
	} else {
		DrawArc(img, cx, cy, 8*scale, 210, 330, color.Black)
	}
}

// DrawBASGorilla draws a simple approximation of the QBASIC gorilla sprite.
func DrawBASGorilla(img image.Image, x, y, scale float64, arms int, clr color.Color) {
	S := func(v float64) float64 { return v * scale }
	DrawFilledRect(img, x-S(4), y, x+S(2.9), y+S(6), clr)
	DrawFilledRect(img, x-S(5), y+S(2), x+S(4), y+S(4), clr)
	drawLine(img, x-S(3), y+S(2), x+S(2), y+S(2), color.Black)
	for i := -2.0; i <= -1.0; i++ {
		DrawFilledRect(img, x+S(i), y+S(4), x+S(i)+1, y+S(4)+1, color.Black)
		DrawFilledRect(img, x+S(i+3), y+S(4), x+S(i+3)+1, y+S(4)+1, color.Black)
	}
	drawLine(img, x-S(3), y+S(7), x+S(2), y+S(7), clr)
	DrawFilledRect(img, x-S(8), y+S(8), x+S(6.9), y+S(14), clr)
	DrawFilledRect(img, x-S(6), y+S(15), x+S(4.9), y+S(20), clr)
	for i := 0.0; i <= 4; i++ {
		DrawArc(img, x+S(i), y+S(25), S(10), 135, 202.5, clr)
		DrawArc(img, x-S(6)+S(i-0.1), y+S(25), S(10), 337.5, 45, clr)
	}
	DrawArc(img, x-S(4.9), y+S(10), S(4.9), 270, 360, color.Black)
	DrawArc(img, x+S(4.9), y+S(10), S(4.9), 180, 270, color.Black)
	for i := -5.0; i <= -1.0; i++ {
		switch arms {
		case ArmsRightUp:
			DrawArc(img, x+S(i-0.1), y+S(14), S(9), 135, 225, clr)
			DrawArc(img, x+S(4.9)+S(i), y+S(4), S(9), 315, 45, clr)
		case ArmsLeftUp:
			DrawArc(img, x+S(i-0.1), y+S(4), S(9), 135, 225, clr)
			DrawArc(img, x+S(4.9)+S(i), y+S(14), S(9), 315, 45, clr)
		default:
			DrawArc(img, x+S(i-0.1), y+S(14), S(9), 135, 225, clr)
			DrawArc(img, x+S(4.9)+S(i), y+S(14), S(9), 315, 45, clr)
		}
	}
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
func ClearRect(img image.Image, x, y, w, h int) {
	drw, ok := img.(interface{ Set(int, int, color.Color) })
	if !ok {
		return
	}
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < h; dy++ {
			drw.Set(x+dx, y+dy, color.RGBA{})
		}
	}
}

// ClearCircle clears a circle region on the image by setting pixels transparent.
func ClearCircle(img image.Image, cx, cy int, r float64) {
	drw, ok := img.(interface{ Set(int, int, color.Color) })
	if !ok {
		return
	}
	ir := int(r)
	r2 := r * r
	for dx := -ir; dx <= ir; dx++ {
		for dy := -ir; dy <= ir; dy++ {
			if float64(dx*dx+dy*dy) <= r2 {
				drw.Set(cx+dx, cy+dy, color.RGBA{})
			}
		}
	}
}

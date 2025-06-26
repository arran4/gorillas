package ebdraw

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/arran4/gorillas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// DrawVectorLines joins a series of points with lines of the specified colour.
func DrawVectorLines(img *ebiten.Image, pts []gorillas.VectorPoint, clr color.Color) {
	if len(pts) == 0 {
		return
	}
	prev := pts[0]
	for _, p := range pts[1:] {
		ebitenutil.DrawLine(img, prev.X, prev.Y, p.X, p.Y, clr)
		prev = p
	}
}

// DrawFilledRect fills a rectangle between two points.
func DrawFilledRect(img *ebiten.Image, x1, y1, x2, y2 float64, clr color.Color) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	ebitenutil.DrawRect(img, x1, y1, x2-x1, y2-y1, clr)
}

// DrawFilledCircle draws a filled circle at the provided centre and radius.
func DrawFilledCircle(img *ebiten.Image, cx, cy, r float64, clr color.Color) {
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			if dx*dx+dy*dy <= r*r {
				ebitenutil.DrawRect(img, cx+dx, cy+dy, 1, 1, clr)
			}
		}
	}
}

// DrawArc renders an arc from startDeg to endDeg going clockwise from 0 degrees to the right.
func DrawArc(img *ebiten.Image, cx, cy, r float64, startDeg, endDeg float64, clr color.Color) {
	step := 2.0
	prevX := cx + r*math.Cos(startDeg*math.Pi/180)
	prevY := cy - r*math.Sin(startDeg*math.Pi/180)
	for a := startDeg + step; a <= endDeg; a += step {
		x := cx + r*math.Cos(a*math.Pi/180)
		y := cy - r*math.Sin(a*math.Pi/180)
		ebitenutil.DrawLine(img, prevX, prevY, x, y, clr)
		prevX, prevY = x, y
	}
}

// Arms constants describe gorilla arm positions when drawing.
const (
	ArmsRightUp = 1
	ArmsLeftUp  = 2
	ArmsDown    = 3
)

// DrawBASSun renders the classic QBASIC sun sprite at the given position and radius.
func DrawBASSun(img *ebiten.Image, cx, cy, r float64, shocked bool, clr color.Color) {
	DrawFilledCircle(img, cx, cy, r, clr)
	rayLen := r / 2
	for i := 0; i < 8; i++ {
		ang := float64(i) * 45 * math.Pi / 180
		x1 := cx + r*math.Cos(ang)
		y1 := cy - r*math.Sin(ang)
		x2 := cx + (r+rayLen)*math.Cos(ang)
		y2 := cy - (r+rayLen)*math.Sin(ang)
		ebitenutil.DrawLine(img, x1, y1, x2, y2, clr)
	}
	scale := r / 12
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
func DrawBASGorilla(img *ebiten.Image, x, y, scale float64, arms int, clr color.Color) {
	S := func(v float64) float64 { return v * scale }
	// head
	DrawFilledRect(img, x-S(4), y, x+S(2.9), y+S(6), clr)
	DrawFilledRect(img, x-S(5), y+S(2), x+S(4), y+S(4), clr)
	ebitenutil.DrawLine(img, x-S(3), y+S(2), x+S(2), y+S(2), color.Black)
	for i := -2.0; i <= -1.0; i++ {
		ebitenutil.DrawRect(img, x+S(i), y+S(4), 1, 1, color.Black)
		ebitenutil.DrawRect(img, x+S(i+3), y+S(4), 1, 1, color.Black)
	}
	// neck
	ebitenutil.DrawLine(img, x-S(3), y+S(7), x+S(2), y+S(7), clr)
	// body
	DrawFilledRect(img, x-S(8), y+S(8), x+S(6.9), y+S(14), clr)
	DrawFilledRect(img, x-S(6), y+S(15), x+S(4.9), y+S(20), clr)
	// legs
	for i := 0.0; i <= 4; i++ {
		DrawArc(img, x+S(i), y+S(25), S(10), 135, 202.5, clr)
		DrawArc(img, x-S(6)+S(i-0.1), y+S(25), S(10), 337.5, 45, clr)
	}
	// chest outline
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

// CreateBananaSprite converts an ASCII mask into an ebiten image using yellow pixels.
func CreateBananaSprite(mask []string) *ebiten.Image {
	h := len(mask)
	w := len(mask[0])
	img := ebiten.NewImage(w, h)
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
func CreateBananaSprites() (left, right, up, down *ebiten.Image) {
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

// CreateGorillaSprite converts an ASCII mask into a coloured ebiten image.
func CreateGorillaSprite(mask []string, clr color.Color) *ebiten.Image {
	h := len(mask)
	w := len(mask[0])
	img := ebiten.NewImage(w, h)
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
func DefaultGorillaSprite(scale float64) *ebiten.Image {
	size := int(30 * scale)
	img := ebiten.NewImage(size, size)
	clr := color.RGBA{150, 75, 0, 255}
	DrawBASGorilla(img, 15*scale, scale, scale, ArmsDown, clr)
	return img
}

// CreateBuildingSprite produces a simple building with random windows.
func CreateBuildingSprite(w, h float64, clr color.Color) *ebiten.Image {
	iw := int(w)
	ih := int(h)
	img := ebiten.NewImage(iw, ih)
	img.Fill(clr)
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
func ClearRect(img *ebiten.Image, x, y, w, h int) {
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < h; dy++ {
			ix := x + dx
			iy := y + dy
			if ix >= 0 && ix < img.Bounds().Dx() && iy >= 0 && iy < img.Bounds().Dy() {
				img.Set(ix, iy, color.RGBA{})
			}
		}
	}
}

// ClearCircle clears a circle region on the image by setting pixels transparent.
func ClearCircle(img *ebiten.Image, cx, cy int, r float64) {
	ir := int(r)
	r2 := r * r
	for dx := -ir; dx <= ir; dx++ {
		for dy := -ir; dy <= ir; dy++ {
			if float64(dx*dx+dy*dy) <= r2 {
				ix := cx + dx
				iy := cy + dy
				if ix >= 0 && ix < img.Bounds().Dx() && iy >= 0 && iy < img.Bounds().Dy() {
					img.Set(ix, iy, color.RGBA{})
				}
			}
		}
	}
}

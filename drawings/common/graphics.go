package drawcommon

import (
	"image/color"
	"image/draw"
	"math"

	gorillas "github.com/arran4/gorillas"
)

// drawLine renders a straight line between two points using simple pixel plots.
func drawLine(img draw.Image, x1, y1, x2, y2 float64, clr color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	n := int(math.Max(math.Abs(dx), math.Abs(dy)))
	if n == 0 {
		img.Set(int(math.Round(x1)), int(math.Round(y1)), clr)
		return
	}
	for i := 0; i <= n; i++ {
		x := int(math.Round(x1 + dx*float64(i)/float64(n)))
		y := int(math.Round(y1 + dy*float64(i)/float64(n)))
		img.Set(x, y, clr)
	}
}

// DrawVectorLines joins a series of points with lines of the specified colour.
func DrawVectorLines(img draw.Image, pts []gorillas.VectorPoint, clr color.Color) {
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
func DrawFilledRect(img draw.Image, x1, y1, x2, y2 float64, clr color.Color) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	for x := int(math.Round(x1)); x <= int(math.Round(x2)); x++ {
		for y := int(math.Round(y1)); y <= int(math.Round(y2)); y++ {
			img.Set(x, y, clr)
		}
	}
}

// DrawFilledCircle draws a filled circle at the provided centre and radius.
func DrawFilledCircle(img draw.Image, cx, cy, r float64, clr color.Color) {
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			if dx*dx+dy*dy <= r*r {
				x := int(math.Round(cx + dx))
				y := int(math.Round(cy + dy))
				img.Set(x, y, clr)
			}
		}
	}
}

// DrawFilledEllipse draws a filled ellipse with horizontal radius rx and vertical radius ry.
func DrawFilledEllipse(img draw.Image, cx, cy, rx, ry float64, clr color.Color) {
	maxX := int(math.Ceil(rx))
	maxY := int(math.Ceil(ry))
	for dx := -maxX; dx <= maxX; dx++ {
		for dy := -maxY; dy <= maxY; dy++ {
			x := float64(dx)
			y := float64(dy)
			if (x*x)/(rx*rx)+(y*y)/(ry*ry) <= 1 {
				img.Set(int(math.Round(cx+x)), int(math.Round(cy+y)), clr)
			}
		}
	}
}

// DrawArc renders an arc from startDeg to endDeg going clockwise from 0 degrees to the right.
func DrawArc(img draw.Image, cx, cy, r float64, startDeg, endDeg float64, clr color.Color) {
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

// DrawThickArc draws a thicker arc using overlapping circles.
func DrawThickArc(img draw.Image, cx, cy, r, thickness, startDeg, endDeg float64, clr color.Color) {
	if endDeg < startDeg {
		endDeg += 360
	}
	step := 1.0
	radius := thickness / 2
	for a := startDeg; a <= endDeg; a += step {
		x := cx + r*math.Cos(a*math.Pi/180)
		y := cy - r*math.Sin(a*math.Pi/180)
		DrawFilledCircle(img, x, y, radius, clr)
	}
}

const (
	ArmsRightUp = 1
	ArmsLeftUp  = 2
	ArmsDown    = 3
)

// DrawBASSun renders the classic QBASIC sun sprite at the given position and radius.
func DrawBASSun(img draw.Image, cx, cy, r float64, shocked bool, clr color.Color) {
	DrawFilledCircle(img, cx, cy, r, clr)
	rayLen := r / 2
	for i := 0; i < 8; i++ {
		ang := float64(i) * 45 * math.Pi / 180
		x1 := cx + r*math.Cos(ang)
		y1 := cy - r*math.Sin(ang)
		x2 := cx + (r+rayLen)*math.Cos(ang)
		y2 := cy - (r+rayLen)*math.Sin(ang)
		drawLine(img, x1, y1, x2, y2, clr)
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
func DrawBASGorilla(img draw.Image, x, y, scale float64, arms int, clr color.Color) {
	S := func(v float64) float64 { return v * scale }
	DrawFilledCircle(img, x, y+S(3.5), S(4.5), clr)
	drawLine(img, x-S(3), y+S(2), x+S(2), y+S(2), color.Black)
	for i := -2.0; i <= -1.0; i++ {
		DrawFilledRect(img, x+S(i), y+S(4), x+S(i)+1, y+S(4)+1, color.Black)
		DrawFilledRect(img, x+S(i+3), y+S(4), x+S(i+3)+1, y+S(4)+1, color.Black)
	}
	drawLine(img, x-S(3), y+S(7), x+S(2), y+S(7), clr)
	// body uses stacked rectangles like the BASIC original
	DrawFilledRect(img, x-S(8), y+S(8), x+S(6.9), y+S(14), clr)
	DrawFilledRect(img, x-S(6), y+S(14), x+S(4.9), y+S(20), clr)
	// round the torso edges
	DrawFilledCircle(img, x-S(6), y+S(14), S(2), clr)
	DrawFilledCircle(img, x+S(4.9), y+S(14), S(2), clr)
       thick := S(4)
       // legs rendered similar to the arms using thick arcs
       DrawThickArc(img, x-S(3), y+S(13), S(8), thick, 135, 225, clr)
       DrawThickArc(img, x+S(2), y+S(13), S(8), thick, -45, 45, clr)
       // extend the feet downward so they match the ASCII art legs
       DrawFilledRect(img, x-S(3)-thick/2, y+S(20), x-S(3)+thick/2, y+S(25), clr)
       DrawFilledRect(img, x+S(2)-thick/2, y+S(20), x+S(2)+thick/2, y+S(25), clr)
       DrawArc(img, x-S(4.9), y+S(10), S(4.9), 270, 360, color.Black)
       DrawArc(img, x+S(4.9), y+S(10), S(4.9), 180, 270, color.Black)
	switch arms {
	case ArmsRightUp:
		DrawThickArc(img, x-S(3), y+S(15), S(9), thick, 135, 225, clr)
		DrawThickArc(img, x+S(2), y+S(5), S(9), thick, 315, 45, clr)
	case ArmsLeftUp:
		DrawThickArc(img, x-S(3), y+S(5), S(9), thick, 135, 225, clr)
		DrawThickArc(img, x+S(2), y+S(15), S(9), thick, 315, 45, clr)
	default:
		DrawThickArc(img, x-S(3), y+S(15), S(9), thick, 135, 225, clr)
		DrawThickArc(img, x+S(2), y+S(15), S(9), thick, 315, 45, clr)
	}
}

// ClearRect clears a rectangle region on the image by setting pixels transparent.
func ClearRect(img draw.Image, x, y, w, h int) {
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < h; dy++ {
			img.Set(x+dx, y+dy, color.RGBA{})
		}
	}
}

// ClearCircle clears a circle region on the image by setting pixels transparent.
func ClearCircle(img draw.Image, cx, cy int, r float64) {
	ir := int(r)
	r2 := r * r
	for dx := -ir; dx <= ir; dx++ {
		for dy := -ir; dy <= ir; dy++ {
			if float64(dx*dx+dy*dy) <= r2 {
				img.Set(cx+dx, cy+dy, color.RGBA{})
			}
		}
	}
}

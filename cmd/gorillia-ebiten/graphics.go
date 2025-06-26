//go:build !test

package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func drawFilledRect(img *ebiten.Image, x1, y1, x2, y2 float64, clr color.Color) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	ebitenutil.DrawRect(img, x1, y1, x2-x1, y2-y1, clr)
}

func drawFilledCircle(img *ebiten.Image, cx, cy, r float64, clr color.Color) {
	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			if dx*dx+dy*dy <= r*r {
				ebitenutil.DrawRect(img, cx+dx, cy+dy, 1, 1, clr)
			}
		}
	}
}

func drawArc(img *ebiten.Image, cx, cy, r float64, startDeg, endDeg float64, clr color.Color) {
	step := 2.0
	prevX := cx + r*math.Cos(startDeg*math.Pi/180)
	prevY := cy + r*math.Sin(startDeg*math.Pi/180)
	for a := startDeg + step; a <= endDeg; a += step {
		x := cx + r*math.Cos(a*math.Pi/180)
		y := cy + r*math.Sin(a*math.Pi/180)
		ebitenutil.DrawLine(img, prevX, prevY, x, y, clr)
		prevX, prevY = x, y
	}
}

// drawBASSun renders the classic BAS sun at the given position and radius.
// If shocked is true the sun uses an "O" mouth like in the original game.
func drawBASSun(img *ebiten.Image, cx, cy, r float64, shocked bool, clr color.Color) {
	drawFilledCircle(img, cx, cy, r, clr)
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
		ebitenutil.DrawLine(img, x1, y1, x2, y2, clr)
	}
	eyeX := 3 * scale
	eyeY := -2 * scale
	drawFilledCircle(img, cx-eyeX, cy+eyeY, 1*scale, color.Black)
	drawFilledCircle(img, cx+eyeX, cy+eyeY, 1*scale, color.Black)
	if shocked {
		drawFilledCircle(img, cx, cy+5*scale, 2.9*scale, color.Black)
	} else {
		drawArc(img, cx, cy, 8*scale, 210, 330, color.Black)
	}
}

// drawBASGorilla draws a simple approximation of the BAS gorilla sprite.
func drawBASGorilla(img *ebiten.Image, x, y, scale float64, arms int, clr color.Color) {
	S := func(v float64) float64 { return v * scale }
	// head
	drawFilledRect(img, x-S(4), y, x+S(2.9), y+S(6), clr)
	drawFilledRect(img, x-S(5), y+S(2), x+S(4), y+S(4), clr)
	ebitenutil.DrawLine(img, x-S(3), y+S(2), x+S(2), y+S(2), color.Black)
	for i := -2.0; i <= -1.0; i++ {
		ebitenutil.DrawRect(img, x+S(i), y+S(4), 1, 1, color.Black)
		ebitenutil.DrawRect(img, x+S(i+3), y+S(4), 1, 1, color.Black)
	}
	// neck
	ebitenutil.DrawLine(img, x-S(3), y+S(7), x+S(2), y+S(7), clr)
	// body
	drawFilledRect(img, x-S(8), y+S(8), x+S(6.9), y+S(14), clr)
	drawFilledRect(img, x-S(6), y+S(15), x+S(4.9), y+S(20), clr)
	// legs
	for i := 0.0; i <= 4; i++ {
		drawArc(img, x+S(i), y+S(25), S(10), 135, 202.5, clr)
		drawArc(img, x-S(6)+S(i-0.1), y+S(25), S(10), 337.5, 45, clr)
	}
	// chest outline
	drawArc(img, x-S(4.9), y+S(10), S(4.9), 270, 360, color.Black)
	drawArc(img, x+S(4.9), y+S(10), S(4.9), 180, 270, color.Black)
	for i := -5.0; i <= -1.0; i++ {
		switch arms {
		case armsRightUp:
			drawArc(img, x+S(i-0.1), y+S(14), S(9), 135, 225, clr)
			drawArc(img, x+S(4.9)+S(i), y+S(4), S(9), 315, 45, clr)
		case armsLeftUp:
			drawArc(img, x+S(i-0.1), y+S(4), S(9), 135, 225, clr)
			drawArc(img, x+S(4.9)+S(i), y+S(14), S(9), 315, 45, clr)
		default:
			drawArc(img, x+S(i-0.1), y+S(14), S(9), 135, 225, clr)
			drawArc(img, x+S(4.9)+S(i), y+S(14), S(9), 315, 45, clr)
		}
	}
}

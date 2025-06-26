package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// playState implements the main gameplay loop.
type playState struct{}

func (playState) Update(g *Game) error {
	if !g.Banana.Active && !g.Explosion.Active {
		if g.enteringAng || g.enteringPow {
			for _, r := range ebiten.AppendInputChars(nil) {
				if r == '*' {
					if g.enteringAng && len(g.angleInput) == 0 {
						g.angleInput = "*"
					} else if g.enteringPow && len(g.powerInput) == 0 {
						g.powerInput = "*"
					}
					continue
				}
				if r >= '0' && r <= '9' {
					if g.enteringAng && len(g.angleInput) < 3 {
						g.angleInput += string(r)
					} else if g.enteringPow && len(g.powerInput) < 3 {
						g.powerInput += string(r)
					}
				}
			}
			for _, k := range inpututil.AppendJustPressedKeys(nil) {
				switch k {
				case ebiten.KeyEnter:
					if g.enteringAng {
						if strings.HasPrefix(g.angleInput, "*") {
							g.Angle = g.LastAngle[g.Current]
						} else if v, err := strconv.Atoi(g.angleInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 360 {
								v = 360
							}
							g.Angle = float64(v)
						}
						g.enteringAng = false
						g.angleInput = ""
						g.enteringPow = true
					} else {
						if strings.HasPrefix(g.powerInput, "*") {
							g.Power = g.LastPower[g.Current]
						} else if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Power = float64(v)
						}
						g.enteringPow = false
						g.powerInput = ""
						g.Throw()
					}
				case ebiten.KeyEscape:
					if g.enteringAng || g.enteringPow {
						g.enteringAng = false
						g.enteringPow = false
						g.angleInput = ""
						g.powerInput = ""
					} else {
						g.State = newScoreState(g.StatsString())
					}
				case ebiten.KeyBackspace:
					if g.enteringAng && len(g.angleInput) > 0 {
						g.angleInput = g.angleInput[:len(g.angleInput)-1]
					} else if g.enteringPow && len(g.powerInput) > 0 {
						g.powerInput = g.powerInput[:len(g.powerInput)-1]
					}
				default:
					if k >= ebiten.Key0 && k <= ebiten.Key9 {
						r := '0' + rune(k-ebiten.Key0)
						if g.enteringAng && len(g.angleInput) < 3 {
							g.angleInput += string(r)
						} else if g.enteringPow && len(g.powerInput) < 3 {
							g.powerInput += string(r)
						}
					}
				}
			}
			return nil
		}
		for _, r := range ebiten.AppendInputChars(nil) {
			if r == '*' {
				g.enteringAng = true
				g.angleInput = "*"
				return nil
			}
			if r >= '0' && r <= '9' {
				g.enteringAng = true
				g.angleInput = string(r)
				return nil
			}
		}
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			if k >= ebiten.Key0 && k <= ebiten.Key9 {
				g.enteringAng = true
				g.angleInput = string('0' + rune(k-ebiten.Key0))
				return nil
			}
		}
		if g.Current == 0 {
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				g.Angle += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				g.Angle -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				g.Power += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				g.Power -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				g.Throw()
			}
		} else {
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				g.Angle += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				g.Angle -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				g.Power += 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				g.Power -= 1
			}
			if ebiten.IsKeyPressed(ebiten.KeyF) {
				g.Throw()
			}
		}
	} else {
		g.Step()
		if g.Banana.Active {
			if g.Banana.X >= g.sunX-sunRadius && g.Banana.X <= g.sunX+sunRadius &&
				g.Banana.Y >= g.sunY-sunRadius && g.Banana.Y <= g.sunY+sunRadius {
				g.sunHitTicks = 10
			}
		}
	}
	if g.sunHitTicks > 0 {
		g.sunHitTicks--
	}
	return nil
}

func (playState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	for i, b := range g.buildings {
		ebitenutil.DrawRect(screen, b.x, float64(g.Height)-b.h, b.w-1, b.h, b.color)
		for _, w := range b.windows {
			ebitenutil.DrawRect(screen, w.x, w.y, w.w, w.h, color.RGBA{255, 255, 0, 255})
		}
		_ = i
	}
	for i := range g.Gorillas {
		g.drawGorilla(screen, i)
	}
	if g.Banana.Active {
		dir := 0
		if math.Abs(g.Banana.VX) > math.Abs(g.Banana.VY) {
			if g.Banana.VX < 0 {
				dir = 0
			} else {
				dir = 1
			}
		} else {
			if g.Banana.VY < 0 {
				dir = 2
			} else {
				dir = 3
			}
		}
		var img *ebiten.Image
		switch dir {
		case 0:
			img = g.bananaLeft
		case 1:
			img = g.bananaRight
		case 2:
			img = g.bananaUp
		case 3:
			img = g.bananaDown
		}
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			w, h := img.Size()
			op.GeoM.Translate(g.Banana.X-float64(w)/2, g.Banana.Y-float64(h)/2)
			screen.DrawImage(img, op)
		}
	}
	if g.Explosion.Active {
		clr := color.RGBA{255, 255, 0, 255}
		if len(g.Explosion.Colors) > g.Explosion.Frame {
			clr = color.RGBAModel.Convert(g.Explosion.Colors[g.Explosion.Frame]).(color.RGBA)
		}
		drawFilledCircle(screen, g.Explosion.X, g.Explosion.Y, g.Explosion.Radii[g.Explosion.Frame], clr)
	}
	g.drawSun(screen)
	g.drawWindArrow(screen)
	angleStr := fmt.Sprintf("%3.0f", g.Angle)
	if g.enteringAng {
		if g.angleInput == "" {
			angleStr = "_"
		} else {
			angleStr = g.angleInput
		}
	}
	powerStr := fmt.Sprintf("%3.0f", g.Power)
	if g.enteringPow {
		if g.powerInput == "" {
			powerStr = "_"
		} else {
			powerStr = g.powerInput
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("A:%3s P:%3s W:%+2.0f P%d %d-%d", angleStr, powerStr, g.Wind, g.Current+1, g.Wins[0], g.Wins[1]))
}

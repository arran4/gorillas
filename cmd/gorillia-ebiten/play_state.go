//go:build !test

package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	gorillas "github.com/arran4/gorillas"
	ebdraw "github.com/arran4/gorillas/drawings/ebiten"
)

// playState implements the main gameplay loop.
type playState struct{}

func (playState) Update(g *Game) error {
	if g.abortPrompt {
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			switch k {
			case ebiten.KeyY:
				g.State = newAbortState()
			case ebiten.KeyN:
				g.abortPrompt = false
				if g.resumeAng {
					g.enteringAng = true
				}
				if g.resumePow {
					g.enteringPow = true
				}
				g.angleInput = ""
				g.powerInput = ""
				g.resumeAng = false
				g.resumePow = false
			}
		}
		return nil
	}
	if !g.Banana.Active && !g.Explosion.Active {
		if g.AI && g.Current == 1 {
			g.Game.AutoShot()
			return nil
		}
		if g.enteringAng || g.enteringPow {
			now := time.Now()
			for _, r := range ebiten.AppendInputChars(nil) {
				if r == '*' {
					if g.enteringAng && len(g.angleInput) == 0 {
						g.angleInput = "*"
					} else if g.enteringPow && len(g.powerInput) == 0 {
						g.powerInput = "*"
					}
					g.lastDigit = now
					continue
				}
				if r == ',' {
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
					} else if g.enteringPow {
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
					continue
				}
				if r >= '0' && r <= '9' {
					if now.Sub(g.lastDigit) > digitBufferTimeout {
						if g.enteringAng {
							g.angleInput = string(r)
						} else {
							g.powerInput = string(r)
						}
					} else {
						if g.enteringAng && len(g.angleInput) < 3 {
							g.angleInput += string(r)
						} else if g.enteringPow && len(g.powerInput) < 3 {
							g.powerInput += string(r)
						}
					}
					g.lastDigit = now
				}
			}
			for _, k := range inpututil.AppendJustPressedKeys(nil) {
				switch k {
				case ebiten.KeyEnter, ebiten.KeyComma:
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
						g.abortPrompt = true
						g.resumeAng = g.enteringAng
						g.resumePow = g.enteringPow
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
				}
			}
			return nil
		}
		for _, r := range ebiten.AppendInputChars(nil) {
			if r == '*' {
				g.enteringAng = true
				g.angleInput = "*"
				g.lastDigit = time.Now()
				return nil
			}
			if r >= '0' && r <= '9' {
				g.enteringAng = true
				g.angleInput = string(r)
				g.lastDigit = time.Now()
				return nil
			}
		}
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			if k >= ebiten.Key0 && k <= ebiten.Key9 {
				g.enteringAng = true
				g.angleInput = string('0' + rune(k-ebiten.Key0))
				g.lastDigit = time.Now()
				return nil
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.Angle += 0.5
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.Angle -= 0.5
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.Power += 0.5
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.Power -= 0.5
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Throw()
		}

		g.gamepads = ebiten.AppendGamepadIDs(g.gamepads[:0])
		for _, id := range g.gamepads {
			if ebiten.IsStandardGamepadLayoutAvailable(id) {
				lx := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
				ly := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
				if lx < -0.2 {
					g.Angle += 0.5
				}
				if lx > 0.2 {
					g.Angle -= 0.5
				}
				if ly < -0.2 {
					g.Power += 0.5
				}
				if ly > 0.2 {
					g.Power -= 0.5
				}
				if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) {
					g.Throw()
				}
			} else {
				if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton0) {
					g.Throw()
				}
			}
		}
	} else {
		g.Step()
		if g.Banana.Active && g.sunIntegrity > 0 {
			r := float64(g.sunIntegrity) * sunRadius / sunMaxIntegrity
			if g.Banana.X >= g.sunX-r && g.Banana.X <= g.sunX+r &&
				g.Banana.Y >= g.sunY-r && g.Banana.Y <= g.sunY+r {
				g.sunHitTicks = 10
				if g.sunIntegrity > 0 {
					g.sunIntegrity--
				}
			}
		}
	}
	if g.sunHitTicks > 0 {
		g.sunHitTicks--
	}
	return nil
}

func (playState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 255, 255})
	bw := float64(g.Width) / float64(g.Game.BuildingCount)
	for i := 0; i < g.Game.BuildingCount; i++ {
		h := g.Buildings[i].H
		intH := int(h)
		img := g.buildingImg[i]
		img.Fill(color.RGBA{})
		img.DrawImage(g.buildingBase[i], nil)
		for _, d := range g.Buildings[i].Damage {
			rx := int(d.X - float64(i)*bw)
			ry := int(d.Y - float64(g.Height-intH))
			ebdraw.ClearCircle(img, rx, ry, d.R)
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i)*bw, float64(g.Height-intH))
		screen.DrawImage(img, op)
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
			op.GeoM.Scale(bananaScale, bananaScale)
			op.GeoM.Translate(g.Banana.X-float64(w)*bananaScale/2, g.Banana.Y-float64(h)*bananaScale/2)
			screen.DrawImage(img, op)
		}
	}
	if g.Explosion.Active {
		clr := color.RGBA{255, 255, 0, 255}
		if len(g.Explosion.Colors) > g.Explosion.Frame {
			clr = color.RGBAModel.Convert(g.Explosion.Colors[g.Explosion.Frame]).(color.RGBA)
		}
		frame := g.Explosion.Frame
		if g.Settings.UseVectorExplosions && frame > 0 && frame-1 < len(g.Explosion.Vectors) {
			drawVectorLines(screen, g.Explosion.Vectors[frame-1], clr)
		} else {
			ebdraw.DrawFilledCircle(screen, g.Explosion.X, g.Explosion.Y, g.Explosion.Radii[frame], clr)
		}
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
	info := fmt.Sprintf("Player %d (%s) - Angle:%s° Power:%s Wind:%+2.0f Score:%d-%d",
		g.Current+1, g.Players[g.Current], angleStr, powerStr, g.Wind, g.Wins[0], g.Wins[1])
	x := 0
	if g.Current == 1 {
		x = g.Width - len(info)*charW
		if x < 0 {
			x = 0
		}
	}
	ebitenutil.DebugPrintAt(screen, info, x, 0)
	if g.abortPrompt {
		msg := "Abort game? [Y/N]"
		x := (g.Width - len(msg)*charW) / 2
		y := g.Height/2 - charH/2
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	} else if g.LastEvent != gorillas.EventNone {
		msg := g.LastEventMsg
		x := (g.Width - len(msg)*charW) / 2
		y := g.Height / 3
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	}
}

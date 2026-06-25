//go:build !test

package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"os"
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
	// Debug save state key
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		b, err := json.MarshalIndent(g.Game, "", "  ")
		if err == nil {
			_ = os.WriteFile("dump_state.json", b, 0644)
			fmt.Println("Saved dump_state.json")
		} else {
			fmt.Printf("Error saving state: %v\n", err)
		}
	}

	if g.abortPrompt {
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			switch k {
			case ebiten.KeyY:
				g.State = newAbortState()
			case ebiten.KeyN:
				g.abortPrompt = false
				g.angleInput = ""
				g.powerInput = ""
			}
		}
		return nil
	}
	if !g.Game.Banana.Active && !g.Game.Explosion.Active {
		if g.AI && g.Game.Current == 1 {
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
							g.Game.Angle = g.Game.LastAngle[g.Game.Current]
						} else if v, err := strconv.Atoi(g.angleInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 360 {
								v = 360
							}
							g.Game.Angle = float64(v)
						}
						g.enteringAng = false
						g.angleInput = ""
						g.enteringPow = true
					} else if g.enteringPow {
						if strings.HasPrefix(g.powerInput, "*") {
							g.Game.Power = g.Game.LastPower[g.Game.Current]
						} else if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Game.Power = float64(v)
						}
						g.enteringPow = false
						g.powerInput = ""
						g.Game.Throw()
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
				case ebiten.KeyEnter:
					if g.enteringAng {
						if strings.HasPrefix(g.angleInput, "*") {
							g.Game.Angle = g.Game.LastAngle[g.Game.Current]
						} else if v, err := strconv.Atoi(g.angleInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 360 {
								v = 360
							}
							g.Game.Angle = float64(v)
						}
						g.enteringAng = false
						g.angleInput = ""
						g.enteringPow = true
					} else {
						if strings.HasPrefix(g.powerInput, "*") {
							g.Game.Power = g.Game.LastPower[g.Game.Current]
						} else if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Game.Power = float64(v)
						}
						g.enteringPow = false
						g.powerInput = ""
						g.Game.Throw()
					}
				case ebiten.KeyEscape:
					if g.enteringAng || g.enteringPow {
						g.abortPrompt = true
						g.enteringAng = false
						g.enteringPow = false
						g.angleInput = ""
						g.powerInput = ""
					} else {
						g.State = newScoreState(g.Game.StatsString())
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
				if g.selAngle {
					g.enteringAng = true
					g.angleInput = "*"
				} else {
					g.enteringPow = true
					g.powerInput = "*"
				}
				g.lastDigit = time.Now()
				return nil
			}
			if r >= '0' && r <= '9' {
				if g.selAngle {
					g.enteringAng = true
					g.angleInput = string(r)
				} else {
					g.enteringPow = true
					g.powerInput = string(r)
				}
				g.lastDigit = time.Now()
				return nil
			}
		}
		for _, k := range inpututil.AppendJustPressedKeys(nil) {
			if k >= ebiten.Key0 && k <= ebiten.Key9 {
				if g.selAngle {
					g.enteringAng = true
					g.angleInput = string('0' + rune(k-ebiten.Key0))
				} else {
					g.enteringPow = true
					g.powerInput = string('0' + rune(k-ebiten.Key0))
				}
				g.lastDigit = time.Now()
				return nil
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.selAngle = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.selAngle = false
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			if g.selAngle {
				g.Game.Angle += 0.5
			} else {
				g.Game.Power += 0.5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			if g.selAngle {
				g.Game.Angle -= 0.5
			} else {
				g.Game.Power -= 0.5
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Game.Throw()
		}

		g.gamepads = ebiten.AppendGamepadIDs(g.gamepads[:0])
		for _, id := range g.gamepads {
			if ebiten.IsStandardGamepadLayoutAvailable(id) {
				lx := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
				ly := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
				if lx < -0.2 {
					g.selAngle = true
				}
				if lx > 0.2 {
					g.selAngle = false
				}
				if ly < -0.2 {
					if g.selAngle {
						g.Game.Angle += 0.5
					} else {
						g.Game.Power += 0.5
					}
				}
				if ly > 0.2 {
					if g.selAngle {
						g.Game.Angle -= 0.5
					} else {
						g.Game.Power -= 0.5
					}
				}
				if inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightBottom) {
					g.Game.Throw()
				}
			} else {
				if inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton0) {
					g.Game.Throw()
				}
			}
		}
	} else {
		g.Game.Step()
		if g.Game.Banana.Active && g.sunIntegrity > 0 {
			r := float64(g.sunIntegrity) * sunRadius / sunMaxIntegrity
			if g.Game.Banana.X >= g.sunX-r && g.Game.Banana.X <= g.sunX+r &&
				g.Game.Banana.Y >= g.sunY-r && g.Game.Banana.Y <= g.sunY+r {
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
	bw := float64(g.Game.Width) / float64(g.Game.BuildingCount)
	for i := 0; i < g.Game.BuildingCount; i++ {
		h := g.Game.Buildings[i].H
		intH := int(h)
		img := g.buildingImg[i]
		img.Fill(color.RGBA{})
		img.DrawImage(g.buildingBase[i], nil)
		for _, d := range g.Game.Buildings[i].Damage {
			rx := int(d.X - float64(i)*bw)
			ry := int(d.Y - float64(g.Game.Height-intH))
			ebdraw.ClearCircle(img, rx, ry, d.R)
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i)*bw, float64(g.Game.Height-intH))
		screen.DrawImage(img, op)
	}
	for i := range g.Game.Gorillas {
		g.drawGorilla(screen, i)
	}
	if g.Game.Banana.Active {
		dir := 0
		if math.Abs(g.Game.Banana.VX) > math.Abs(g.Game.Banana.VY) {
			if g.Game.Banana.VX < 0 {
				dir = 0
			} else {
				dir = 1
			}
		} else {
			if g.Game.Banana.VY < 0 {
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
			op.GeoM.Translate(g.Game.Banana.X-float64(w)*bananaScale/2, g.Game.Banana.Y-float64(h)*bananaScale/2)
			screen.DrawImage(img, op)
		}
	}
	if g.Game.Explosion.Active {
		clr := color.RGBA{255, 255, 0, 255}
		if len(g.Game.Explosion.Colors) > g.Game.Explosion.Frame {
			clr = color.RGBAModel.Convert(g.Game.Explosion.Colors[g.Game.Explosion.Frame]).(color.RGBA)
		}
		frame := g.Game.Explosion.Frame
		if g.Game.Settings.UseVectorExplosions && frame > 0 && frame-1 < len(g.Game.Explosion.Vectors) {
			drawVectorLines(screen, g.Game.Explosion.Vectors[frame-1], clr)
		} else {
			ebdraw.DrawFilledCircle(screen, g.Game.Explosion.X, g.Game.Explosion.Y, g.Game.Explosion.Radii[frame], clr)
		}
	}
	g.drawSun(screen)
	g.drawWindArrow(screen)
	angleStr := fmt.Sprintf("%3.0f", g.Game.Angle)
	if g.enteringAng {
		if g.angleInput == "" {
			angleStr = "_"
		} else {
			angleStr = g.angleInput
		}
	}
	powerStr := fmt.Sprintf("%3.0f", g.Game.Power)
	if g.enteringPow {
		if g.powerInput == "" {
			powerStr = "_"
		} else {
			powerStr = g.powerInput
		}
	}
	if g.selAngle {
		angleStr = "[" + angleStr + "]"
	} else {
		powerStr = "[" + powerStr + "]"
	}
	info := fmt.Sprintf("Player %d (%s) - Angle:%s° Power:%s Wind:%+2.0f Score:%d-%d",
		g.Game.Current+1, g.Game.Players[g.Game.Current], angleStr, powerStr, g.Game.Wind, g.Game.Wins[0], g.Game.Wins[1])
	x := 0
	if g.Game.Current == 1 {
		x = g.Game.Width - len(info)*charW
		if x < 0 {
			x = 0
		}
	}
	ebitenutil.DebugPrintAt(screen, info, x, 0)
	if g.abortPrompt {
		msg := "Abort game? [Y/N]"
		x := (g.Game.Width - len(msg)*charW) / 2
		y := g.Game.Height/2 - charH/2
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	} else if g.Game.LastEvent != gorillas.EventNone {
		msg := g.Game.LastEventMsg
		x := (g.Game.Width - len(msg)*charW) / 2
		y := g.Game.Height / 3
		ebitenutil.DebugPrintAt(screen, msg, x, y)
	}
}

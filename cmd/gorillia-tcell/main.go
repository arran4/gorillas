//go:build !test

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
	"unicode"

	"github.com/arran4/gorillas"
	"github.com/gdamore/tcell/v2"
)

type damageRect struct{ x, y, w, h int }

type building struct {
	h       int
	windows []int
	damage  []damageRect
}

type Game struct {
	*gorillas.Game
	buildings    []building
	screen       tcell.Screen
	sunX, sunY   int
	sunHitTicks  int
	sunIntegrity int
	angleInput   string
	powerInput   string
	enteringAng  bool
	enteringPow  bool
	abortPrompt  bool
	selAngle     bool
	gorillaArt   [][]string
	js           *joystick
	lastDigit    time.Time
}

const (
	buildingWidth      = 8
	sunMaxIntegrity    = 4
	digitBufferTimeout = 3 * time.Second
)

func drawLine(s tcell.Screen, x0, y0, x1, y1 int, r rune) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		s.SetContent(x0, y0, r, nil, tcell.StyleDefault)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			if x0 == x1 {
				break
			}
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			if y0 == y1 {
				break
			}
			err += dx
			y0 += sy
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func newGame(settings gorillas.Settings, buildings int, wind float64) *Game {
	g := &Game{Game: gorillas.NewGame(80, 24, buildings)}
	if !math.IsNaN(wind) {
		g.Game.Wind = wind
	}
	g.Game.Settings = settings
	if art, err := gorillas.LoadGorillaArt("assets/gorilla.txt"); err == nil {
		g.gorillaArt = art
	} else {
		g.gorillaArt = [][]string{{" O ", "/|\\", "/ \\"}}
	}
	g.LoadScores()
	rand.Seed(time.Now().UnixNano())
	for _, b := range g.Buildings {
		var wins []int
		top := g.Height - int(b.H) + 2
		for y := g.Height - 2; y > top; y -= 2 {
			if rand.Intn(3) != 0 {
				wins = append(wins, y)
			}
		}
		g.buildings = append(g.buildings, building{h: int(b.H), windows: wins})
	}
	// position the sun in the horizontal centre
	g.sunX = g.Width/2 - 1
	g.sunY = 1
	if js, err := openJoystick(); err == nil {
		g.js = js
	}
	g.sunIntegrity = sunMaxIntegrity
	g.Game.ResetHook = func() { g.sunIntegrity = sunMaxIntegrity }
	return g
}

var (
	sunHappy = []string{`\|/`, `-o-`, `/|\`}
	sunShock = []string{`\|/`, `-O-`, `/|\`}
)

func (g *Game) drawSun() {
	if g.sunIntegrity <= 0 {
		return
	}
	art := sunHappy
	if g.sunHitTicks > 0 {
		art = sunShock
		g.sunHitTicks--
	}
	switch g.sunIntegrity {
	case 1:
		r := rune(art[1][1])
		g.screen.SetContent(g.sunX+1, g.sunY+1, r, nil, tcell.StyleDefault)
	case 2:
		line := art[1]
		for dx, r := range line {
			if r != ' ' {
				g.screen.SetContent(g.sunX+dx, g.sunY+1, r, nil, tcell.StyleDefault)
			}
		}
	case 3:
		for dy, line := range art[:2] {
			for dx, r := range line {
				if r != ' ' {
					g.screen.SetContent(g.sunX+dx, g.sunY+dy, r, nil, tcell.StyleDefault)
				}
			}
		}
	default:
		for dy, line := range art {
			for dx, r := range line {
				if r != ' ' {
					g.screen.SetContent(g.sunX+dx, g.sunY+dy, r, nil, tcell.StyleDefault)
				}
			}
		}
	}
}

func (g *Game) draw() {
	g.screen.Clear()
	for i := range g.buildings {
		g.buildings[i].h = int(g.Buildings[i].H)
		g.buildings[i].damage = g.buildings[i].damage[:0]
		for _, d := range g.Buildings[i].Damage {
			r := int(math.Round(d.R))
			g.buildings[i].damage = append(g.buildings[i].damage, damageRect{
				x: int(math.Round(d.X)) - r,
				y: int(math.Round(d.Y)) - r,
				w: 2 * r,
				h: 2 * r,
			})
		}
	}
	for i, b := range g.buildings {
		x := i*buildingWidth + 4
		for y := g.Height - 1; y >= g.Height-b.h; y-- {
			g.screen.SetContent(x, y, '#', nil, tcell.StyleDefault)
		}
		for _, wy := range b.windows {
			g.screen.SetContent(x, wy, 'o', nil, tcell.StyleDefault)
		}
		for _, d := range b.damage {
			for dx := 0; dx < d.w; dx++ {
				for dy := 0; dy < d.h; dy++ {
					g.screen.SetContent(d.x+dx, d.y+dy, ' ', nil, tcell.StyleDefault)
				}
			}
		}
	}
	g.drawGorilla(0)
	g.drawGorilla(1)
	if g.Banana.Active {
		ch := 'o'
		if math.Abs(g.Banana.VX) > math.Abs(g.Banana.VY) {
			if g.Banana.VX < 0 {
				ch = '<'
			} else {
				ch = '>'
			}
		} else {
			if g.Banana.VY < 0 {
				ch = '^'
			} else {
				ch = 'v'
			}
		}
		g.screen.SetContent(int(g.Banana.X), int(g.Banana.Y), ch, nil, tcell.StyleDefault)
	}
	if g.Explosion.Active {
		char := '*'
		if !g.Settings.UseOldExplosions {
			chars := []rune{'#', '@', 'O', 'o', '.'}
			if g.Explosion.Frame < len(chars) {
				char = chars[g.Explosion.Frame]
			} else {
				char = chars[len(chars)-1]
			}
		}
		frame := g.Explosion.Frame
		if g.Settings.UseVectorExplosions && frame > 0 && frame-1 < len(g.Explosion.Vectors) {
			pts := g.Explosion.Vectors[frame-1]
			for i := 1; i < len(pts); i++ {
				drawLine(g.screen, int(pts[i-1].X), int(pts[i-1].Y), int(pts[i].X), int(pts[i].Y), char)
			}
		} else {
			r := int(g.Explosion.Radii[frame])
			ex := int(g.Explosion.X)
			ey := int(g.Explosion.Y)
			for dx := -r; dx <= r; dx++ {
				for dy := -r; dy <= r; dy++ {
					if dx*dx+dy*dy <= r*r {
						x := ex + dx
						y := ey + dy
						if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
							g.screen.SetContent(x, y, char, nil, tcell.StyleDefault)
						}
					}
				}
			}
		}
	}
	g.drawSun()
	g.drawWindArrow()
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
	if g.selAngle {
		angleStr = "[" + angleStr + "]"
	} else {
		powerStr = "[" + powerStr + "]"
	}
	info := fmt.Sprintf("Player %d (%s) - Angle:%s° Power:%s Wind:%+2.0f Score:%d-%d",
		g.Current+1, g.Players[g.Current], angleStr, powerStr, g.Wind, g.Wins[0], g.Wins[1])
	x := 0
	if g.Current == 1 {
		x = g.Width - len(info)
		if x < 0 {
			x = 0
		}
	}
	drawString(g.screen, x, 0, info)
	if g.abortPrompt {
		msg := "Abort game? [Y/N]"
		drawString(g.screen, (g.Width-len(msg))/2, 1, msg)
	} else if g.LastEvent != gorillas.EventNone {
		msg := g.LastEventMsg
		drawString(g.screen, (g.Width-len(msg))/2, g.Height/3, msg)
	}
	g.screen.Show()
}

func (g *Game) drawWindArrow() {
	if g.Wind == 0 {
		return
	}
	length := int(math.Round(g.Wind * 3 * float64(g.Width) / 320))
	// Draw near the top instead of bottom for better visibility
	y := 1
	x := g.Width / 2
	dir := 1
	if length < 0 {
		dir = -1
	}
	for i := dir; i != length; i += dir {
		pos := x + i
		if pos >= 0 && pos < g.Width {
			g.screen.SetContent(pos, y, '-', nil, tcell.StyleDefault)
		}
	}
	headX := x + length
	if headX >= 0 && headX < g.Width {
		head := '>'
		if length < 0 {
			head = '<'
		}
		g.screen.SetContent(headX, y, head, nil, tcell.StyleDefault)
	}
}

func (g *Game) drawGorilla(idx int) {
	if len(g.gorillaArt) == 0 {
		return
	}
	frame := g.gorillaArt[0]
	width := gorillas.FrameWidth(frame)
	x := int(g.Gorillas[idx].X) - width/2
	y := int(g.Gorillas[idx].Y) - len(frame)
	style := tcell.StyleDefault
	for dy, line := range frame {
		for dx, r := range line {
			if r != ' ' {
				g.screen.SetContent(x+dx, y+dy, r, nil, style)
			}
		}
	}
}

func (g *Game) startVictoryDance(idx int) {
	g.Dance = gorillas.NewDance(idx, []float64{-3, 0, -3, 0}, g.Gorillas[idx].Y)
}

func (g *Game) throw() {
	g.Throw()
}

func (g *Game) run(s tcell.Screen, ai bool) error {
	g.screen = s

	ticker := time.NewTicker(50 * time.Millisecond)
	prevExplosion := g.Explosion.Active
	for {
		g.draw()
		<-ticker.C
		g.Step()
		if !prevExplosion && g.Explosion.Active {
			g.startVictoryDance(g.Current)
		}
		prevExplosion = g.Explosion.Active
		if g.Banana.Active && g.sunIntegrity > 0 {
			if int(g.Banana.X) >= g.sunX && int(g.Banana.X) < g.sunX+3 && int(g.Banana.Y) >= g.sunY && int(g.Banana.Y) < g.sunY+3 {
				g.sunHitTicks = 10
				if g.sunIntegrity > 0 {
					g.sunIntegrity--
				}
			}
		}
		if g.Banana.Active || g.Explosion.Active || g.Dance.Active {
			continue
		}

		if g.js != nil {
			g.js.poll()
			if g.js.axis[0] < -10000 {
				g.Angle += 0.5
			}
			if g.js.axis[0] > 10000 {
				g.Angle -= 0.5
			}
			if g.js.axis[1] < -10000 {
				g.Power += 0.5
			}
			if g.js.axis[1] > 10000 {
				g.Power -= 0.5
			}
			if g.js.btn[0] {
				g.throw()
			}
		}

		if ai && g.Current == 1 {
			g.AutoShot()
			continue
		}

		ev := s.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			if g.abortPrompt {
				r := unicode.ToUpper(key.Rune())
				if r == 'Y' {
					g.Aborted = true
					return nil
				}
				if r == 'N' {
					g.abortPrompt = false
				}
				continue
			}
			now := time.Now()
			switch key.Key() {
			case tcell.KeyEscape:
				g.abortPrompt = true
			case tcell.KeyLeft, tcell.KeyRight:
				g.selAngle = !g.selAngle
			case tcell.KeyUp:
				if g.selAngle {
					g.Angle += 0.5
				} else {
					g.Power += 0.5
				}
			case tcell.KeyDown:
				if g.selAngle {
					g.Angle -= 0.5
				} else {
					g.Power -= 0.5
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if g.selAngle {
					if len(g.angleInput) > 0 {
						g.angleInput = g.angleInput[:len(g.angleInput)-1]
						if v, err := strconv.Atoi(g.angleInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 360 {
								v = 360
							}
							g.Angle = float64(v)
						}
					}
				} else {
					if len(g.powerInput) > 0 {
						g.powerInput = g.powerInput[:len(g.powerInput)-1]
						if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Power = float64(v)
						}
					}
				}
			case tcell.KeyEnter:
				g.throw()
			default:
				r := key.Rune()
				if r == '*' {
					if g.selAngle {
						g.Angle = g.LastAngle[g.Current]
						g.angleInput = ""
					} else {
						g.Power = g.LastPower[g.Current]
						g.powerInput = ""
					}
					g.lastDigit = now
				} else if r >= '0' && r <= '9' {
					if g.selAngle {
						if now.Sub(g.lastDigit) > digitBufferTimeout {
							g.angleInput = string(r)
						} else if len(g.angleInput) < 3 {
							g.angleInput += string(r)
						}
						if v, err := strconv.Atoi(g.angleInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 360 {
								v = 360
							}
							g.Angle = float64(v)
						}
					} else {
						if now.Sub(g.lastDigit) > digitBufferTimeout {
							g.powerInput = string(r)
						} else if len(g.powerInput) < 3 {
							g.powerInput += string(r)
						}
						if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Power = float64(v)
						}
					}
					g.lastDigit = now
				}
			}
			continue
		}
	}
}

// setupScreen presents an interactive form allowing the player names,
// round count and gravity to be edited. It returns the updated values
// once the user starts the game by pressing Enter on "Start" or
// pressing Escape.
func setupScreen(s tcell.Screen, league *gorillas.League, p1, p2 string, rounds int, gravity float64) (string, string, int, float64, bool) {
	fields := []string{p1, p2, strconv.Itoa(rounds), fmt.Sprintf("%.0f", gravity)}
	players := league.Names()
	cur := 0
	editing := false
	editingPlayer := -1
	newPlayer := false
	oldName := ""
	assignField := 0

	updateAssignField := func() {
		if cur < 2 {
			assignField = cur
		} else if cur < len(fields) {
			assignField = -1
		}
	}
	updateAssignField()
	labels := []string{"Player 1:", "Player 2:", "Rounds:", "Gravity:"}
	selectedPlayer := -1
	for {
		s.Clear()
		_, h := s.Size()
		baseY := h/2 - 2
		opts := []string{"New Player", "Rename Player", "Delete Player", "Start"}
		total := len(fields) + len(players) + len(opts)
		newIdx := len(fields) + len(players)
		renameIdx := newIdx + 1
		deleteIdx := renameIdx + 1
		startIdx := deleteIdx + 1
		drawString(s, 2, baseY-2, "Game Setup")
		for i, lbl := range labels {
			style := tcell.StyleDefault
			if i == cur {
				style = style.Reverse(true)
			}
			line := fmt.Sprintf("%s [%s]", lbl, fields[i])
			for x, r := range line {
				s.SetContent(2+x, baseY+i, r, nil, style)
			}
		}
		py := baseY + len(labels) + 1
		drawString(s, 2, py, "Players:")
		for i, name := range players {
			style := tcell.StyleDefault
			if cur == len(fields)+i {
				style = style.Reverse(true)
			}
			drawString(s, 4, py+1+i, fmt.Sprintf("[%s]", name))
		}
		optY := py + 1 + len(players)
		newIdx = len(fields) + len(players)
		for i, opt := range opts {
			style := tcell.StyleDefault
			if cur == newIdx+i {
				style = style.Reverse(true)
			}
			drawString(s, 2, optY+i, opt)
		}

		if editing {
			var cx, cy int
			if editingPlayer >= 0 {
				cx = 4 + len(players[editingPlayer])
				cy = py + 1 + editingPlayer
			} else {
				lbl := labels[cur]
				cx = 2 + len(lbl) + 1 + len(fields[cur])
				cy = baseY + cur
			}
			s.ShowCursor(cx, cy)
		} else {
			s.HideCursor()
		}
		s.Show()

		ev := s.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			if editing {
				switch key.Key() {
				case tcell.KeyEnter:
					if editingPlayer >= 0 {
						name := players[editingPlayer]
						if newPlayer {
							league.AddPlayer(name)
						} else {
							league.RenamePlayer(oldName, name)
							if fields[0] == oldName {
								fields[0] = name
							}
							if fields[1] == oldName {
								fields[1] = name
							}
						}
						league.Save()
						selectedPlayer = editingPlayer
						editingPlayer = -1
						newPlayer = false
					}
					editing = false
					cur = (cur + 1) % total
					updateAssignField()
				case tcell.KeyEsc:
					if editingPlayer >= 0 {
						if newPlayer {
							players = players[:len(players)-1]
							selectedPlayer = -1
						} else {
							players[editingPlayer] = oldName
							selectedPlayer = editingPlayer
						}
						editingPlayer = -1
						newPlayer = false
					}
					editing = false
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if editingPlayer >= 0 {
						if len(players[editingPlayer]) > 0 {
							players[editingPlayer] = players[editingPlayer][:len(players[editingPlayer])-1]
						}
					} else if len(fields[cur]) > 0 {
						fields[cur] = fields[cur][:len(fields[cur])-1]
					}
				default:
					if key.Rune() != 0 {
						if editingPlayer >= 0 {
							players[editingPlayer] += string(key.Rune())
						} else {
							if cur >= 2 {
								if key.Rune() >= '0' && key.Rune() <= '9' {
									fields[cur] += string(key.Rune())
								}
							} else {
								fields[cur] += string(key.Rune())
							}
						}
					}
				}
				continue
			}

			// Automatically enter editing mode when typing or
			// pressing backspace on a selected field or player.
			if key.Key() == tcell.KeyBackspace || key.Key() == tcell.KeyBackspace2 || key.Rune() != 0 {
				if cur < len(fields) {
					editing = true
					editingPlayer = -1
					if key.Key() == tcell.KeyBackspace || key.Key() == tcell.KeyBackspace2 {
						if len(fields[cur]) > 0 {
							fields[cur] = fields[cur][:len(fields[cur])-1]
						}
					} else {
						r := key.Rune()
						if cur >= 2 {
							if r >= '0' && r <= '9' {
								fields[cur] += string(r)
							}
						} else {
							fields[cur] += string(r)
						}
					}
					continue
				} else if cur >= len(fields) && cur < len(fields)+len(players) {
					editing = true
					editingPlayer = cur - len(fields)
					oldName = players[editingPlayer]
					newPlayer = false
					selectedPlayer = editingPlayer
					if key.Key() == tcell.KeyBackspace || key.Key() == tcell.KeyBackspace2 {
						if len(players[editingPlayer]) > 0 {
							players[editingPlayer] = players[editingPlayer][:len(players[editingPlayer])-1]
						}
					} else {
						players[editingPlayer] += string(key.Rune())
					}
					continue
				}
			}

			switch key.Key() {
			case tcell.KeyEsc:
				r, _ := strconv.Atoi(fields[2])
				g, _ := strconv.ParseFloat(fields[3], 64)
				return fields[0], fields[1], r, g, true
			case tcell.KeyCtrlC:
				r, _ := strconv.Atoi(fields[2])
				g, _ := strconv.ParseFloat(fields[3], 64)
				return fields[0], fields[1], r, g, false
			case tcell.KeyUp:
				if cur > 0 {
					cur--
				} else {
					cur = total - 1
				}
				updateAssignField()
				if cur >= len(fields) && cur < len(fields)+len(players) {
					selectedPlayer = cur - len(fields)
				}
			case tcell.KeyDown, tcell.KeyTab:
				cur = (cur + 1) % total
				updateAssignField()
				if cur >= len(fields) && cur < len(fields)+len(players) {
					selectedPlayer = cur - len(fields)
				}
			case tcell.KeyEnter:
				if cur == startIdx {
					r, _ := strconv.Atoi(fields[2])
					g, _ := strconv.ParseFloat(fields[3], 64)
					return fields[0], fields[1], r, g, true
				} else if cur == newIdx {
					players = append(players, "")
					cur = len(fields) + len(players) - 1
					editing = true
					editingPlayer = len(players) - 1
					newPlayer = true
					selectedPlayer = editingPlayer
				} else if cur == renameIdx {
					if selectedPlayer >= 0 {
						editing = true
						editingPlayer = selectedPlayer
						oldName = players[selectedPlayer]
						newPlayer = false
						cur = len(fields) + selectedPlayer
					}
				} else if cur == deleteIdx {
					if selectedPlayer >= 0 {
						name := players[selectedPlayer]
						league.DeletePlayer(name)
						league.Save()
						players = append(players[:selectedPlayer], players[selectedPlayer+1:]...)
						if fields[0] == name {
							fields[0] = ""
						}
						if fields[1] == name {
							fields[1] = ""
						}
						if selectedPlayer >= len(players) {
							selectedPlayer = len(players) - 1
						}
						if cur >= total-1 {
							cur--
						}
					}
				} else if cur >= len(fields) && cur < len(fields)+len(players) && assignField >= 0 {
					name := players[cur-len(fields)]
					other := 1 - assignField
					if fields[other] == name {
						fields[other], fields[assignField] = fields[assignField], name
					} else {
						fields[assignField] = name
					}
					cur = assignField
				} else if cur < len(fields) {
					editing = true
				} else {
					editing = true
					editingPlayer = cur - len(fields)
					oldName = players[editingPlayer]
					selectedPlayer = editingPlayer
				}
			case tcell.KeyRune:
				switch key.Rune() {
				case 'n':
					players = append(players, "")
					editing = true
					editingPlayer = len(players) - 1
					newPlayer = true
					selectedPlayer = editingPlayer
					cur = len(fields) + editingPlayer
				case 'd':
					if selectedPlayer >= 0 {
						name := players[selectedPlayer]
						league.DeletePlayer(name)
						league.Save()
						players = append(players[:selectedPlayer], players[selectedPlayer+1:]...)
						if fields[0] == name {
							fields[0] = ""
						}
						if fields[1] == name {
							fields[1] = ""
						}
						if selectedPlayer >= len(players) {
							selectedPlayer = len(players) - 1
						}
						if cur >= startIdx {
							cur--
						}
					}
				case 'r':
					if selectedPlayer >= 0 {
						editing = true
						editingPlayer = selectedPlayer
						oldName = players[editingPlayer]
						newPlayer = false
						cur = len(fields) + editingPlayer
					}
				}
			}
		}
	}
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		// When TERM is unset or tcell cannot figure out the terminal
		// capabilities, NewScreen() returns ErrTermNotFound with the
		// underlying error from infocmp. Provide a hint for the user
		// in this case as the error message isn't very descriptive.
		if errors.Is(err, tcell.ErrTermNotFound) {
			log.Printf("Unable to detect terminal. Ensure $TERM is set and 'infocmp' is installed")
		} else {
			log.Printf("Error: %s", err.Error())
		}
		panic(err)
	}
	if err = s.Init(); err != nil {
		panic(fmt.Errorf("screen init: %w", err))
	}
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	defer s.Fini()

	settings := gorillas.LoadSettings()
	wind := flag.Float64("wind", math.NaN(), "initial wind")
	gravity := flag.Float64("gravity", settings.DefaultGravity, "gravity")
	rounds := flag.Int("rounds", settings.DefaultRoundQty, "round count")
	buildings := flag.Int("buildings", gorillas.DefaultBuildingCount, "building count")
	p1 := flag.String("player1", "Player 1", "name of player 1")
	p2 := flag.String("player2", "Player 2", "name of player 2")
	flag.BoolVar(&settings.UseSound, "sound", settings.UseSound, "enable sound")
	flag.BoolVar(&settings.WinnerFirst, "winnerfirst", settings.WinnerFirst, "winner starts next round")
	ai := flag.Bool("ai", false, "enable computer opponent")
	flag.Parse()
	settings.DefaultGravity = *gravity
	settings.DefaultRoundQty = *rounds

	if settings.ShowIntro {
		showIntroMovie(s, settings.UseSound, settings.UseSlidingText)
	}

	if !introScreen(s, settings.UseSound, settings.UseSlidingText) {
		return
	}

	league := gorillas.LoadLeague("gorillas.lge")
	var ok bool
	*p1, *p2, *rounds, *gravity, ok = setupScreen(s, league, *p1, *p2, *rounds, *gravity)
	if !ok {
		return
	}
	settings.DefaultGravity = *gravity
	settings.DefaultRoundQty = *rounds

	g := newGame(settings, *buildings, *wind)
	g.Players = [2]string{*p1, *p2}
	g.League = league
	winsBackup := g.TotalWins
	var playersBackup map[string]*gorillas.PlayerStats
	if g.League != nil {
		playersBackup = make(map[string]*gorillas.PlayerStats, len(g.League.Players))
		for n, ps := range g.League.Players {
			cp := *ps
			playersBackup[n] = &cp
		}
	}
	if err := g.run(s, *ai); err != nil {
		panic(fmt.Errorf("run game: %w", err))
	}
	if g.Aborted {
		g.TotalWins = winsBackup
		if g.League != nil {
			g.League.Players = playersBackup
			g.League.Save()
		}
		showGameAborted(s)
		return
	}
	g.SaveScores()
	showStats(s, g.StatsString())
	if g.League != nil {
		showLeague(s, g.League)
	}
	fmt.Println(g.StatsString())
	showExtro(s)
}

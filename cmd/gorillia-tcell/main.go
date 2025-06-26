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

type building struct {
	h       int
	windows []int
}

type Game struct {
	*gorillas.Game
	buildings   []building
	screen      tcell.Screen
	sunX, sunY  int
	sunHitTicks int
	angleInput  string
	powerInput  string
	enteringAng bool
	enteringPow bool
	abortPrompt bool
	resumeAng   bool
	resumePow   bool
	gorillaArt  [][]string
}

const buildingWidth = 8

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
	g.sunX = g.Width - 4
	g.sunY = 1
	return g
}

var (
	sunHappy = []string{`\|/`, `-o-`, `/|\`}
	sunShock = []string{`\|/`, `-O-`, `/|\`}
)

func (g *Game) drawSun() {
	art := sunHappy
	if g.sunHitTicks > 0 {
		art = sunShock
		g.sunHitTicks--
	}
	for dy, line := range art {
		for dx, r := range line {
			if r != ' ' {
				g.screen.SetContent(g.sunX+dx, g.sunY+dy, r, nil, tcell.StyleDefault)
			}
		}
	}
}

func (g *Game) draw() {
	g.screen.Clear()
	for i, b := range g.buildings {
		x := i*buildingWidth + 4
		for y := g.Height - 1; y >= g.Height-b.h; y-- {
			g.screen.SetContent(x, y, '#', nil, tcell.StyleDefault)
		}
		for _, wy := range b.windows {
			g.screen.SetContent(x, wy, 'o', nil, tcell.StyleDefault)
		}
	}
	g.drawGorilla(0)
	g.drawGorilla(1)
	// draw a simple sun
	g.screen.SetContent(g.Width-2, 1, 'O', nil, tcell.StyleDefault)
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
		r := int(g.Explosion.Radii[g.Explosion.Frame])
		ex := int(g.Explosion.X)
		ey := int(g.Explosion.Y)
		char := '*'
		if !g.Settings.UseOldExplosions {
			chars := []rune{'#', '@', 'O', 'o', '.'}
			if g.Explosion.Frame < len(chars) {
				char = chars[g.Explosion.Frame]
			} else {
				char = chars[len(chars)-1]
			}
		}
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
	s := fmt.Sprintf("A:%3s P:%3s W:%+2.0f P%d %d-%d", angleStr, powerStr, g.Wind, g.Current+1, g.Wins[0], g.Wins[1])
	for i, r := range s {
		g.screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
	}
	if g.abortPrompt {
		msg := "Abort game? [Y/N]"
		drawString(g.screen, (g.Width-len(msg))/2, 1, msg)
	}
	g.screen.Show()
}

func (g *Game) drawWindArrow() {
	if g.Wind == 0 {
		return
	}
	length := int(math.Round(g.Wind * 3 * float64(g.Width) / 320))
	y := g.Height - 1
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
	x := int(g.Gorillas[idx].X) - len(frame[0])/2
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

func (g *Game) throw() {
	g.Throw()
}

func (g *Game) run(s tcell.Screen, ai bool) error {
	g.screen = s

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		g.draw()
		if g.Banana.Active || g.Explosion.Active {
			<-ticker.C
			g.Step()
			if int(g.Banana.X) >= g.sunX && int(g.Banana.X) < g.sunX+3 && int(g.Banana.Y) >= g.sunY && int(g.Banana.Y) < g.sunY+3 {
				g.sunHitTicks = 10
			}
			continue
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
					return nil
				}
				if r == 'N' {
					g.abortPrompt = false
					if g.resumeAng {
						g.enteringAng = true
						g.angleInput = ""
					}
					if g.resumePow {
						g.enteringPow = true
						g.powerInput = ""
					}
					g.resumeAng = false
					g.resumePow = false
				}
				continue
			}
			if g.enteringAng || g.enteringPow {
				switch key.Key() {
				case tcell.KeyEnter:
					if g.enteringAng {
						if v, err := strconv.Atoi(g.angleInput); err == nil {
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
						if v, err := strconv.Atoi(g.powerInput); err == nil {
							if v < 0 {
								v = 0
							} else if v > 200 {
								v = 200
							}
							g.Power = float64(v)
						}
						g.enteringPow = false
						g.powerInput = ""
						g.throw()
					}
				case tcell.KeyEsc:
					g.abortPrompt = true
					g.resumeAng = g.enteringAng
					g.resumePow = g.enteringPow
					g.enteringAng = false
					g.enteringPow = false
					g.angleInput = ""
					g.powerInput = ""
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if g.enteringAng {
						if len(g.angleInput) > 0 {
							g.angleInput = g.angleInput[:len(g.angleInput)-1]
						}
					} else if g.enteringPow {
						if len(g.powerInput) > 0 {
							g.powerInput = g.powerInput[:len(g.powerInput)-1]
						}
					}
				default:
					if key.Rune() >= '0' && key.Rune() <= '9' {
						if g.enteringAng {
							if len(g.angleInput) < 3 {
								g.angleInput += string(key.Rune())
							}
						} else if g.enteringPow {
							if len(g.powerInput) < 3 {
								g.powerInput += string(key.Rune())
							}
						}
					}
				}
				continue
			}
			if key.Rune() >= '0' && key.Rune() <= '9' {
				g.enteringAng = true
				g.angleInput = string(key.Rune())
				continue
			}
			switch key.Key() {
			case tcell.KeyEscape:
				return nil
			case tcell.KeyLeft:
				g.Angle += 1
			case tcell.KeyRight:
				g.Angle -= 1
			case tcell.KeyUp:
				g.Power += 1
			case tcell.KeyDown:
				g.Power -= 1
			case tcell.KeyEnter:
				g.throw()
			}
		}
	}
}

// setupScreen presents an interactive form allowing the player names,
// round count and gravity to be edited. It returns the updated values
// once the user presses Escape to start the game.
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
	for {
		s.Clear()
		_, h := s.Size()
		baseY := h/2 - 2
		drawString(s, 2, baseY-2, "Game Setup (Esc to start)")
		labels := []string{"Player 1:", "Player 2:", "Rounds:", "Gravity:"}
		for i, lbl := range labels {
			style := tcell.StyleDefault
			if i == cur {
				style = style.Reverse(true)
			}
			line := fmt.Sprintf("%s %s", lbl, fields[i])
			for x, r := range line {
				s.SetContent(2+x, baseY+i, r, nil, style)
			}
		}
		py := baseY + len(labels) + 1
		drawString(s, 2, py, "Players (n=new r=rename d=del):")
		for i, name := range players {
			style := tcell.StyleDefault
			if cur == len(fields)+i {
				style = style.Reverse(true)
			}
			drawString(s, 4, py+1+i, name)
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
						editingPlayer = -1
						newPlayer = false
					}
					editing = false
				case tcell.KeyEsc:
					if editingPlayer >= 0 {
						if newPlayer {
							players = players[:len(players)-1]
						} else {
							players[editingPlayer] = oldName
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
					cur = len(fields) + len(players) - 1
				}
				updateAssignField()
			case tcell.KeyDown, tcell.KeyTab:
				cur = (cur + 1) % (len(fields) + len(players))
				updateAssignField()
			case tcell.KeyEnter:
				if cur >= len(fields) && assignField >= 0 {
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
				}
			case tcell.KeyRune:
				switch key.Rune() {
				case 'n':
					players = append(players, "")
					cur = len(fields) + len(players) - 1
					editing = true
					editingPlayer = cur - len(fields)
					newPlayer = true
				case 'd':
					if cur >= len(fields) {
						idx := cur - len(fields)
						name := players[idx]
						league.DeletePlayer(name)
						league.Save()
						players = append(players[:idx], players[idx+1:]...)
						if fields[0] == name {
							fields[0] = ""
						}
						if fields[1] == name {
							fields[1] = ""
						}
						if cur >= len(fields)+len(players) {
							cur--
						}
					}
				case 'r':
					if cur >= len(fields) {
						editing = true
						editingPlayer = cur - len(fields)
						oldName = players[editingPlayer]
						newPlayer = false
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
	if err := g.run(s, *ai); err != nil {
		panic(fmt.Errorf("run game: %w", err))
	}
	g.SaveScores()
	showStats(s, g.StatsString())
	if g.League != nil {
		showLeague(s, g.League)
	}
	fmt.Println(g.StatsString())
	showExtro(s)
}

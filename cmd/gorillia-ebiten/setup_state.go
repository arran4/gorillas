//go:build !test

package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// setupState allows editing players, rounds and gravity before starting.
type setupState struct {
	game          *Game
	fields        []string
	players       []string
	cur           int
	editing       bool
	editingPlayer int
	newPlayer     bool
	oldName       string
	assignField   int
}

func newSetupState(g *Game) *setupState {
	s := &setupState{
		game:          g,
		fields:        []string{g.Players[0], g.Players[1], strconv.Itoa(g.Settings.DefaultRoundQty), fmt.Sprintf("%.0f", g.Settings.DefaultGravity)},
		players:       g.League.Names(),
		editingPlayer: -1,
	}
	s.updateAssignField()
	return s
}

func (s *setupState) updateAssignField() {
	if s.cur < 2 {
		s.assignField = s.cur
	} else if s.cur < len(s.fields) {
		s.assignField = -1
	}
}

func keyToRune(k ebiten.Key) rune {
	if k >= ebiten.Key0 && k <= ebiten.Key9 {
		return '0' + rune(k-ebiten.Key0)
	}
	if k >= ebiten.KeyA && k <= ebiten.KeyZ {
		return 'A' + rune(k-ebiten.KeyA)
	}
	switch k {
	case ebiten.KeySpace:
		return ' '
	case ebiten.KeyMinus:
		return '-'
	case ebiten.KeyPeriod:
		return '.'
	}
	return 0
}

func (s *setupState) Update(g *Game) error {
	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		if s.editing {
			switch k {
			case ebiten.KeyEnter:
				if s.editingPlayer >= 0 {
					name := s.players[s.editingPlayer]
					if s.newPlayer {
						s.game.League.AddPlayer(name)
					} else {
						s.game.League.RenamePlayer(s.oldName, name)
						if s.fields[0] == s.oldName {
							s.fields[0] = name
						}
						if s.fields[1] == s.oldName {
							s.fields[1] = name
						}
					}
					s.game.League.Save()
					s.editingPlayer = -1
					s.newPlayer = false
				}
				s.editing = false
			case ebiten.KeyEscape:
				if s.editingPlayer >= 0 {
					if s.newPlayer {
						s.players = s.players[:len(s.players)-1]
					} else {
						s.players[s.editingPlayer] = s.oldName
					}
					s.editingPlayer = -1
					s.newPlayer = false
				}
				s.editing = false
			case ebiten.KeyBackspace:
				if s.editingPlayer >= 0 {
					if len(s.players[s.editingPlayer]) > 0 {
						s.players[s.editingPlayer] = s.players[s.editingPlayer][:len(s.players[s.editingPlayer])-1]
					}
				} else if len(s.fields[s.cur]) > 0 {
					s.fields[s.cur] = s.fields[s.cur][:len(s.fields[s.cur])-1]
				}
			default:
				r := keyToRune(k)
				if r != 0 {
					if s.editingPlayer >= 0 {
						s.players[s.editingPlayer] += string(r)
					} else {
						if s.cur >= 2 {
							if (r >= '0' && r <= '9') || (s.cur == 3 && r == '.') {
								s.fields[s.cur] += string(r)
							}
						} else {
							s.fields[s.cur] += string(r)
						}
					}
				}
			}
			return nil
		}

		// Automatically start editing when typing or pressing backspace
		// on the selected field or player.
		if k != ebiten.KeyN && k != ebiten.KeyD && k != ebiten.KeyR {
			if k == ebiten.KeyBackspace || keyToRune(k) != 0 {
				if s.cur < len(s.fields) {
					s.editing = true
					s.editingPlayer = -1
					if k == ebiten.KeyBackspace {
						if len(s.fields[s.cur]) > 0 {
							s.fields[s.cur] = s.fields[s.cur][:len(s.fields[s.cur])-1]
						}
					} else {
						r := keyToRune(k)
						if s.cur >= 2 {
							if (r >= '0' && r <= '9') || (s.cur == 3 && r == '.') {
								s.fields[s.cur] += string(r)
							}
						} else {
							s.fields[s.cur] += string(r)
						}
					}
					continue
				} else if s.cur >= len(s.fields) && s.cur < len(s.fields)+len(s.players) {
					s.editing = true
					s.editingPlayer = s.cur - len(s.fields)
					s.oldName = s.players[s.editingPlayer]
					s.newPlayer = false
					if k == ebiten.KeyBackspace {
						if len(s.players[s.editingPlayer]) > 0 {
							s.players[s.editingPlayer] = s.players[s.editingPlayer][:len(s.players[s.editingPlayer])-1]
						}
					} else {
						s.players[s.editingPlayer] += string(keyToRune(k))
					}
					continue
				}
			}
		}

		switch k {
		case ebiten.KeyEscape:
			r, _ := strconv.Atoi(s.fields[2])
			gval, _ := strconv.ParseFloat(s.fields[3], 64)
			s.game.Players = [2]string{s.fields[0], s.fields[1]}
			s.game.Settings.DefaultRoundQty = r
			s.game.Settings.DefaultGravity = gval
			s.game.Gravity = gval
			s.game.State = playState{}
			return nil
		case ebiten.KeyQ:
			s.game.State = newMenuState(s.game.Settings.UseSound, s.game.Settings.UseSlidingText)
			return nil
		case ebiten.KeyUp:
			if s.cur > 0 {
				s.cur--
			} else {
				s.cur = len(s.fields) + len(s.players) - 1
			}
			s.updateAssignField()
		case ebiten.KeyDown, ebiten.KeyTab:
			s.cur = (s.cur + 1) % (len(s.fields) + len(s.players))
			s.updateAssignField()
		case ebiten.KeyEnter:
			if s.cur >= len(s.fields) && s.assignField >= 0 {
				name := s.players[s.cur-len(s.fields)]
				other := 1 - s.assignField
				if s.fields[other] == name {
					s.fields[other], s.fields[s.assignField] = s.fields[s.assignField], name
				} else {
					s.fields[s.assignField] = name
				}
				s.cur = s.assignField
			} else if s.cur < len(s.fields) {
				s.editing = true
				s.editingPlayer = -1
			} else {
				s.editing = true
				s.editingPlayer = s.cur - len(s.fields)
				s.oldName = s.players[s.editingPlayer]
			}
		case ebiten.KeyN:
			s.players = append(s.players, "")
			s.cur = len(s.fields) + len(s.players) - 1
			s.editing = true
			s.editingPlayer = len(s.players) - 1
			s.newPlayer = true
		case ebiten.KeyD:
			if s.cur >= len(s.fields) {
				idx := s.cur - len(s.fields)
				name := s.players[idx]
				s.game.League.DeletePlayer(name)
				s.game.League.Save()
				s.players = append(s.players[:idx], s.players[idx+1:]...)
				if s.fields[0] == name {
					s.fields[0] = ""
				}
				if s.fields[1] == name {
					s.fields[1] = ""
				}
				if s.cur >= len(s.fields)+len(s.players) {
					s.cur--
				}
			}
		case ebiten.KeyR:
			if s.cur >= len(s.fields) {
				s.editing = true
				s.editingPlayer = s.cur - len(s.fields)
				s.oldName = s.players[s.editingPlayer]
				s.newPlayer = false
			}
		}
	}
	return nil
}

func (s *setupState) Draw(g *Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	baseY := g.Height/2 - 2*charH
	ebitenutil.DebugPrintAt(screen, "Game Setup (Esc to start)", 2*charW, baseY-2*charH)
	labels := []string{"Player 1:", "Player 2:", "Rounds:", "Gravity:"}
	for i, lbl := range labels {
		line := fmt.Sprintf("%s %s", lbl, s.fields[i])
		prefix := "  "
		if i == s.cur {
			prefix = "> "
		}
		ebitenutil.DebugPrintAt(screen, prefix+line, 2*charW, baseY+i*charH)
	}
	py := baseY + len(labels)*charH + charH
	ebitenutil.DebugPrintAt(screen, "Players (n=new r=rename d=del):", 2*charW, py)
	for i, name := range s.players {
		prefix := "  "
		if len(s.fields)+i == s.cur {
			prefix = "> "
		}
		ebitenutil.DebugPrintAt(screen, prefix+name, 4*charW, py+(i+1)*charH)
	}
}

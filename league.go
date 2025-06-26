package gorillas

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// PlayerStats holds accumulated statistics for a player.
type PlayerStats struct {
	Rounds   int     `json:"rounds"`
	Wins     int     `json:"wins"`
	Accuracy float64 `json:"accuracy"`
}

// League manages a set of PlayerStats loaded from disk.
type League struct {
	Players map[string]*PlayerStats `json:"players"`
	file    string
}

// LoadLeague reads statistics from the given file.
func LoadLeague(path string) *League {
	l := &League{Players: map[string]*PlayerStats{}, file: path}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &l.Players)
	}
	return l
}

// Save writes the league statistics back to disk.
func (l *League) Save() {
	if l == nil {
		return
	}
	if l.file == "" {
		return
	}
	if b, err := json.Marshal(l.Players); err == nil {
		_ = os.WriteFile(l.file, b, 0644)
	}
}

// RecordRound updates the league for a round between p1 and p2.
// winner indicates which player won (0 or 1). shots is how many throws
// the winner took to achieve the hit.
func (l *League) RecordRound(p1, p2 string, winner, shots int) {
	if l == nil {
		return
	}
	ps1 := l.getPlayer(p1)
	ps2 := l.getPlayer(p2)
	ps1.Rounds++
	ps2.Rounds++
	if winner == 0 {
		ps := ps1
		ps.Wins++
		if shots > 0 {
			if ps.Accuracy > 0 {
				ps.Accuracy = (ps.Accuracy + float64(shots)) / 2
			} else {
				ps.Accuracy = float64(shots)
			}
		}
	} else if winner == 1 {
		ps := ps2
		ps.Wins++
		if shots > 0 {
			if ps.Accuracy > 0 {
				ps.Accuracy = (ps.Accuracy + float64(shots)) / 2
			} else {
				ps.Accuracy = float64(shots)
			}
		}
	}
}

func (l *League) getPlayer(name string) *PlayerStats {
	if ps, ok := l.Players[name]; ok {
		return ps
	}
	ps := &PlayerStats{}
	l.Players[name] = ps
	return ps
}

// Standings returns the league table sorted by win ratio then accuracy.
func (l *League) Standings() []struct {
	Name string
	PlayerStats
} {
	var list []struct {
		Name string
		PlayerStats
	}
	for name, ps := range l.Players {
		list = append(list, struct {
			Name string
			PlayerStats
		}{name, *ps})
	}
	sort.Slice(list, func(i, j int) bool {
		ratio1 := 0.0
		if list[i].Rounds > 0 {
			ratio1 = float64(list[i].Wins) / float64(list[i].Rounds)
		}
		ratio2 := 0.0
		if list[j].Rounds > 0 {
			ratio2 = float64(list[j].Wins) / float64(list[j].Rounds)
		}
		if ratio1 == ratio2 {
			return list[i].Accuracy < list[j].Accuracy
		}
		return ratio1 > ratio2
	})
	return list
}

// String returns a printable league table.
func (l *League) String() string {
	if l == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("Player           Rounds Wins Accuracy\n")
	for i, s := range l.Standings() {
		_ = i
		b.WriteString(fmt.Sprintf("%-15s %6d %4d %8.1f\n", s.Name, s.Rounds, s.Wins, s.Accuracy))
	}
	return b.String()
}

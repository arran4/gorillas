//go:build !test

package gorillas

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

var (
	audioOnce   sync.Once
	audioCtx    *audio.Context
	beepSample  []byte
	introOnce   sync.Once
	introSample []byte
)

func initAudio() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "audio initialization failed: %v\n", r)
			audioCtx = nil
		}
	}()
	audioCtx = audio.NewContext(sampleRate)
	n := sampleRate / 10
	beepSample = make([]byte, n*4)
	for i := 0; i < n; i++ {
		v := math.Sin(2 * math.Pi * 440 * float64(i) / sampleRate)
		s := int16(v * 0.3 * 32767)
		beepSample[i*4] = byte(s)
		beepSample[i*4+1] = byte(s >> 8)
		beepSample[i*4+2] = byte(s)
		beepSample[i*4+3] = byte(s >> 8)
	}
}

type qbNote struct {
	freq float64
	dur  time.Duration
}

func noteDuration(tempo, l int) time.Duration {
	if l <= 0 {
		return 0
	}
	sec := (60.0 / float64(tempo)) * (4.0 / float64(l))
	return time.Duration(sec * float64(time.Second))
}

func noteFreq(octave, pitch int) float64 {
	// A4 index is 4*12 + 9 = 57
	n := octave*12 + pitch
	diff := n - 57
	return 440 * math.Pow(2, float64(diff)/12)
}

func parsePlayString(seq string) []qbNote {
	tempo := 120
	octave := 4
	length := 4
	var notes []qbNote
	i := 0
	toInt := func(s string) (int, int) {
		n := 0
		j := 0
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			n = n*10 + int(s[j]-'0')
			j++
		}
		return n, j
	}
	pitchMap := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	seq = strings.ToLower(seq)
	for i < len(seq) {
		switch seq[i] {
		case 't':
			v, n := toInt(seq[i+1:])
			if v > 0 {
				tempo = v
			}
			i += 1 + n
		case 'o':
			v, n := toInt(seq[i+1:])
			octave = v
			i += 1 + n
		case 'l':
			v, n := toInt(seq[i+1:])
			if v > 0 {
				length = v
			}
			i += 1 + n
		case 'n':
			v, n := toInt(seq[i+1:])
			d := noteDuration(tempo, length)
			if v == 0 {
				notes = append(notes, qbNote{dur: d})
			}
			i += 1 + n
		case 'p':
			v, n := toInt(seq[i+1:])
			notes = append(notes, qbNote{dur: noteDuration(tempo, v)})
			i += 1 + n
		case '>', '<':
			if seq[i] == '>' {
				octave++
			} else {
				octave--
			}
			i++
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g':
			note := seq[i]
			i++
			adj := 0
			if i < len(seq) {
				switch seq[i] {
				case '#', '+':
					adj = 1
					i++
				case '-':
					adj = -1
					i++
				}
			}
			v, n := toInt(seq[i:])
			if n > 0 {
				i += n
			}
			l := length
			if v > 0 {
				l = v
			}
			pitch := pitchMap[note] + adj
			notes = append(notes, qbNote{freq: noteFreq(octave, pitch), dur: noteDuration(tempo, l)})
		default:
			i++
		}
	}
	return notes
}

func synthesize(notes []qbNote) []byte {
	var out []byte
	for _, n := range notes {
		count := int(float64(sampleRate) * n.dur.Seconds())
		if count <= 0 {
			continue
		}
		if n.freq == 0 {
			out = append(out, make([]byte, count*4)...)
			continue
		}
		for i := 0; i < count; i++ {
			v := math.Sin(2 * math.Pi * n.freq * float64(i) / sampleRate)
			s := int16(v * 0.3 * 32767)
			out = append(out, byte(s), byte(s>>8), byte(s), byte(s>>8))
		}
	}
	return out
}

func initIntro() {
	seqs := []string{
		"t120o1l16b9n0baan0bn0bn0baaan0b9n0baan0b",
		"o2l16e-9n0e-d-d-n0e-n0e-n0e-d-d-d-n0e-9n0e-d-d-n0e-",
		"o2l16g-9n0g-een0g-n0g-n0g-eeen0g-9n0g-een0g-",
		"o2l16b9n0baan0g-n0g-n0g-eeen0o1b9n0baan0b",
	}
	var notes []qbNote
	for _, s := range seqs {
		notes = append(notes, parsePlayString(s)...)
	}
	snippet := parsePlayString("T160O0L32EFGEFDC")
	for i := 0; i < 4; i++ {
		notes = append(notes, snippet...)
		notes = append(notes, qbNote{dur: 100 * time.Millisecond})
	}
	introSample = synthesize(notes)
}

func PlayBeep() {
	audioOnce.Do(initAudio)
	if audioCtx != nil {
		p, err := audioCtx.NewPlayer(bytes.NewReader(beepSample))
		if err != nil {
			panic(fmt.Errorf("new player: %w", err))
		}
		p.Play()
	} else {
		fmt.Print("\a")
	}
}

func PlayIntroMusic() {
	introOnce.Do(initIntro)
	audioOnce.Do(initAudio)
	if audioCtx != nil {
		p, err := audioCtx.NewPlayer(bytes.NewReader(introSample))
		if err != nil {
			panic(fmt.Errorf("new player: %w", err))
		}
		p.Play()
	} else {
		for i := 0; i < 3; i++ {
			fmt.Print("\a")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

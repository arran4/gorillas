package gorillas

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

var (
	audioOnce  sync.Once
	audioCtx   *audio.Context
	beepSample []byte
)

func initAudio() {
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

func PlayBeep() {
	audioOnce.Do(initAudio)
	if audioCtx != nil {
		p := audioCtx.NewPlayer(bytes.NewReader(beepSample))
		p.Play()
	} else {
		fmt.Print("\a")
	}
}

func PlayIntroMusic() {
	for i := 0; i < 3; i++ {
		PlayBeep()
		time.Sleep(100 * time.Millisecond)
	}
}

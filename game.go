package gorillas

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
)

type DamageCircle struct {
	X, Y, R float64
}

type Building struct {
	X, W, H float64
	Damage  []DamageCircle
}

type Gorilla struct {
	X, Y float64
}

type Banana struct {
	X, Y   float64
	VX, VY float64
	Active bool
}

type Settings struct {
	UseSound            bool
	UseOldExplosions    bool
	UseVectorExplosions bool
	NewExplosionRadius  float64
	UseSlidingText      bool
	DefaultGravity      float64
	DefaultRoundQty     int
	ShowIntro           bool
	ForceCGA            bool
	WinnerFirst         bool
	VariableWind        bool
	WindFluctuations    bool
}

type Explosion struct {
	X, Y    float64
	Radii   []float64
	Colors  []color.Color
	Vectors [][]VectorPoint
	Frame   int
	Active  bool
}

// ShotEvent indicates special outcomes for a banana throw.
type ShotEvent int

const (
	EventNone ShotEvent = iota
	EventWeak
	EventBackwards
	EventSelf
)

// EventMessage returns the display text for a given ShotEvent.
func EventMessage(e ShotEvent) string {
	switch e {
	case EventWeak:
		msgs := []string{
			"Your little muscles not strong enough?",
			"Now that was feeble.",
			"You can do better than that!",
		}
		return msgs[rand.Intn(len(msgs))]
	case EventBackwards:
		return "Don't throw it that way!"
	case EventSelf:
		return "Now that was pretty dumb."
	}
	return ""
}

// Dance holds temporary state for the winner's victory animation.
type Dance struct {
	idx    int
	frames []float64
	frame  int
	baseY  float64
	Active bool
}

// NewDance creates a Dance with the provided parameters and initializes
// internal state so the dance can be used by external packages without
// accessing unexported fields.
func NewDance(idx int, frames []float64, baseY float64) Dance {
	d := Dance{
		idx:    idx,
		frames: frames,
		baseY:  baseY,
		Active: true,
	}
	d.frame = 0
	return d
}

// ShotRecord stores the angle and power for a single throw.
type ShotRecord struct {
	Angle float64 `json:"angle"`
	Power float64 `json:"power"`
}

func DefaultSettings() Settings {
	return Settings{
		UseSound:            true,
		UseVectorExplosions: false,
		NewExplosionRadius:  40,
		UseSlidingText:      false,
		DefaultGravity:      17,
		DefaultRoundQty:     4,
		ShowIntro:           true,
		ForceCGA:            false,
		WinnerFirst:         false,
		VariableWind:        false,
		WindFluctuations:    false,
	}
}

// LoadScores reads the persistent win totals from disk.
func (g *Game) LoadScores() {
	file := g.ScoreFile
	if file == "" {
		file = defaultScoreFile
	}
	b, err := os.ReadFile(file)
	if err == nil {
		if err := json.Unmarshal(b, &g.TotalWins); err != nil {
			fmt.Fprintf(os.Stderr, "load scores: %v\n", err)
		}
	}
}

// SaveScores writes the accumulated win totals to disk.
func (g *Game) SaveScores() {
	file := g.ScoreFile
	if file == "" {
		file = defaultScoreFile
	}
	b, err := json.Marshal(g.TotalWins)
	if err == nil {
		if err := os.WriteFile(file, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "save scores: %v\n", err)
		}
	}
}

// LoadShots reads the shot history from disk.
func (g *Game) LoadShots() {
	file := g.ShotsFile
	if file == "" {
		file = defaultShotsFile
	}
	b, err := os.ReadFile(file)
	if err == nil {
		if err := json.Unmarshal(b, &g.ShotHistory); err != nil {
			fmt.Fprintf(os.Stderr, "load shots: %v\n", err)
		}
	}
}

// SaveShots writes the shot history to disk.
func (g *Game) SaveShots() {
	file := g.ShotsFile
	if file == "" {
		file = defaultShotsFile
	}
	b, err := json.Marshal(g.ShotHistory)
	if err == nil {
		if err := os.WriteFile(file, b, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "save shots: %v\n", err)
		}
	}
}

// StatsString returns a printable summary of wins this session and overall.
func (g *Game) StatsString() string {
	session := fmt.Sprintf("Session - P1:%d P2:%d", g.Wins[0], g.Wins[1])
	total := fmt.Sprintf("Overall - P1:%d P2:%d", g.TotalWins[0], g.TotalWins[1])
	if g.League != nil {
		return session + "\n" + total + "\n\n" + g.League.String()
	}
	return session + "\n" + total
}

type Game struct {
	Width, Height int
	Buildings     []Building
	Gorillas      [2]Gorilla
	Banana        Banana
	Explosion     Explosion
	Dance         Dance
	Settings      Settings
	Angle         float64
	Power         float64
	Angles        [2]float64
	Powers        [2]float64
	Current       int
	Wins          [2]int
	TotalWins     [2]int
	Shots         [2]int
	LastAngle     [2]float64
	LastPower     [2]float64
	Players       [2]string
	League        *League
	ScoreFile     string
	ShotsFile     string
	ShotHistory   []ShotRecord
	Wind          float64
	BuildingCount int
	Gravity       float64
	HitMap        *HitMap

	// LastEvent records the outcome of the most recent shot.
	LastEvent ShotEvent
	// LastEventTicks counts down the display duration of LastEvent.
	LastEventTicks int
	// LastEventMsg stores the random message associated with LastEvent.
	LastEventMsg string

	lastStartX float64
	lastOtherX float64
	lastVX     float64
	ResetHook  func()

	// roundOver indicates whether the current explosion ends the round.
	roundOver bool

	// Aborted indicates whether the game was aborted mid-play.
	Aborted bool
}

const DefaultBuildingCount = 10
const defaultScoreFile = "gorillas_scores.json"
const defaultShotsFile = "gorillas_shots.json"
const defaultLeagueFile = "gorillas.lge"
const groundBounceFactor = 0.4
const groundBounceThreshold = 5.0
const eventDisplayTicks = 40

func NewGame(width, height, buildingCount int) *Game {
	if buildingCount <= 0 {
		buildingCount = DefaultBuildingCount
	}
	g := &Game{Width: width, Height: height, Angle: 45, Power: 50, ScoreFile: defaultScoreFile, ShotsFile: defaultShotsFile, BuildingCount: buildingCount, Aborted: false}
	g.roundOver = true
	g.Angles = [2]float64{45, 45}
	g.Powers = [2]float64{50, 50}
	g.League = LoadLeague(defaultLeagueFile)
	g.Players = [2]string{"Player 1", "Player 2"}
	g.Settings = DefaultSettings()
	g.Gravity = g.Settings.DefaultGravity
	g.Wind = basicWind()
	bw := float64(width) / float64(g.BuildingCount)

	// create a sloping skyline similar to the original BASIC version
	slope := rand.Intn(6) + 1
	newHt := float64(height) * 0.2
	if slope == 2 || slope == 6 {
		newHt = float64(height) * 0.6
	}
	htInc := float64(height) / 35

	for i := 0; i < g.BuildingCount; i++ {
		x := float64(i) * bw
		switch slope {
		case 1:
			newHt += htInc
		case 2:
			newHt -= htInc
		case 3, 5:
			if x > float64(width)/2 {
				newHt -= 2 * htInc
			} else {
				newHt += 2 * htInc
			}
		case 4:
			if x > float64(width)/2 {
				newHt += 2 * htInc
			} else {
				newHt -= 2 * htInc
			}
		}

		h := newHt + rand.Float64()*float64(height)/6 - float64(height)/12
		if h < float64(height)*0.15 {
			h = float64(height) * 0.15
		}
		if h > float64(height)*0.6 {
			h = float64(height) * 0.6
		}

		g.Buildings = append(g.Buildings, Building{
			X: x,
			W: bw,
			H: h,
		})
	}
	g.Gorillas[0] = Gorilla{g.Buildings[1].X + bw/2, float64(height) - g.Buildings[1].H}
	g.Gorillas[1] = Gorilla{g.Buildings[g.BuildingCount-2].X + bw/2, float64(height) - g.Buildings[g.BuildingCount-2].H}

	g.HitMap = NewHitMap(width, height)
	for _, b := range g.Buildings {
		x1 := int(b.X)
		x2 := int(b.X + b.W)
		y1 := height - int(b.H)
		g.HitMap.AddBuilding(x1, y1, x2, height)
	}
	for i, gr := range g.Gorillas {
		g.HitMap.DrawGorilla(int(gr.X), int(gr.Y), i, 4)
	}

	return g
}

func (g *Game) Reset() {
	wins := g.Wins
	totals := g.TotalWins
	file := g.ScoreFile
	shotsFile := g.ShotsFile
	players := g.Players
	league := g.League
	settings := g.Settings
	gravity := g.Gravity
	hook := g.ResetHook
	*g = *NewGame(g.Width, g.Height, g.BuildingCount)
	g.Wins = wins
	g.TotalWins = totals
	g.ScoreFile = file
	g.ShotsFile = shotsFile
	g.Players = players
	g.League = league
	g.Settings = settings
	g.Gravity = gravity
	g.ResetHook = hook
	if g.ResetHook != nil {
		g.ResetHook()
	}
}

func (g *Game) setCurrent(idx int) {
	g.Current = idx
	g.Angle = g.Angles[idx]
	g.Power = g.Powers[idx]
}

func fnRan(x int) int {
	return rand.Intn(x) + 1
}

func basicWind() float64 {
	w := float64(fnRan(10) - 5)
	if fnRan(3) == 1 {
		if w > 0 {
			w += float64(fnRan(10))
		} else {
			w -= float64(fnRan(10))
		}
	}
	return w
}

func (g *Game) recordExplosionDamage(x, y, r float64) {
	for i := range g.Buildings {
		b := &g.Buildings[i]
		bx1 := b.X
		bx2 := b.X + b.W
		by1 := float64(g.Height) - b.H
		by2 := float64(g.Height)
		if x+r <= bx1 || x-r >= bx2 || y+r <= by1 || y-r >= by2 {
			continue
		}
		b.Damage = append(b.Damage, DamageCircle{x, y, r})
	}
	if g.HitMap != nil {
		g.HitMap.ClearBuildingArea(int(math.Round(x)), int(math.Round(y)), int(math.Ceil(r)))
	}
}

func (g *Game) pointInDamage(idx int, x, y float64) bool {
	b := &g.Buildings[idx]
	for _, d := range b.Damage {
		dx := x - d.X
		dy := y - d.Y
		if dx*dx+dy*dy <= d.R*d.R {
			return true
		}
	}
	return false
}

func (g *Game) explosionBase() float64 {
	base := g.Settings.NewExplosionRadius
	if base <= 0 {
		base = 16
	}
	if g.Settings.ForceCGA {
		base /= 2
	}
	return base
}

func (g *Game) handleGorillaKill(idx int) {
	shooter := g.Current
	winner := shooter
	event := EventNone
	if idx == shooter {
		winner = (shooter + 1) % 2
		event = EventSelf
		g.LastEvent = event
		g.LastEventTicks = eventDisplayTicks
		g.LastEventMsg = EventMessage(event)
	}
	g.Wins[winner]++
	g.TotalWins[winner]++
	if g.League != nil {
		g.League.RecordRound(g.Players[0], g.Players[1], winner, g.Shots[shooter])
		g.League.Save()
	}
	g.Shots = [2]int{}
	g.SaveScores()
	if g.HitMap != nil {
		gr := g.Gorillas[idx]
		g.HitMap.ClearGorilla(int(gr.X), int(gr.Y), idx, 4)
	}
	g.startVictoryDance(winner)
	g.setCurrent(winner)
	if g.Settings.UseSound && event != EventNone {
		PlayBeep()
	}
	g.roundOver = true
}

func (g *Game) killGorillaIfInRadius(x, y, r float64) bool {
	if g.HitMap != nil {
		idx := g.HitMap.GorillaHitInCircle(int(math.Round(x)), int(math.Round(y)), int(math.Ceil(r)))
		if idx >= 0 {
			g.handleGorillaKill(idx)
			return true
		}
		return false
	}
	for i, gr := range g.Gorillas {
		dx := gr.X - x
		dy := gr.Y - y
		if dx*dx+dy*dy <= r*r {
			g.handleGorillaKill(i)
			return true
		}
	}
	return false
}

// gorillaHitBetween checks if the line from (x1,y1) to (x2,y2) intersects any
// gorilla. It returns the index of the gorilla hit or -1.
func (g *Game) gorillaHitBetween(x1, y1, x2, y2 float64) int {
	steps := int(math.Ceil(math.Max(math.Abs(x2-x1), math.Abs(y2-y1))))
	if steps < 1 {
		steps = 1
	}
	vx := x2 - x1
	vy := y2 - y1
	forward := (g.Current == 0 && vx > 0) || (g.Current == 1 && vx < 0)
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		if g.HitMap != nil {
			idx := g.HitMap.GorillaHitAt(int(math.Round(x)), int(math.Round(y)))
			if idx >= 0 {
				if idx == g.Current && forward && vy <= 0 && i <= 5 {
					continue
				}
				return idx
			}
		}
		for j, gr := range g.Gorillas {
			if math.Abs(gr.X-x) < 5 && math.Abs(gr.Y-y) < 10 {
				if j == g.Current && forward && vy <= 0 && i <= 5 {
					continue
				}
				return j
			}
		}
	}
	return -1
}

func (g *Game) startExplosion(x, y float64) {
	base := g.explosionBase()
	if g.Settings.UseSound {
		PlayExplosionMelody()
	}
	g.Explosion = Explosion{X: x, Y: y}
	if g.Settings.UseOldExplosions {
		for i := 1; i <= int(base); i++ {
			g.Explosion.Radii = append(g.Explosion.Radii, float64(i))
		}
		for i := int(base * 1.5); i >= 1; i-- {
			g.Explosion.Radii = append(g.Explosion.Radii, float64(i))
		}
	} else {
		g.Explosion.Radii = []float64{base * 1.175, base, base * 0.9, base * 0.6, base * 0.45, 0}
		g.Explosion.Colors = []color.Color{
			color.RGBA{128, 128, 128, 255},
			color.RGBA{255, 0, 0, 255},
			color.RGBA{255, 165, 0, 255},
			color.RGBA{255, 255, 0, 255},
			color.RGBA{255, 255, 255, 255},
			color.Black,
		}
		if g.Settings.UseVectorExplosions {
			frames := []float64{base, base * 0.9, base * 0.6, base * 0.45}
			for _, r := range frames {
				w := 2 * r
				h := 2 * r * 0.825
				offX := x - r
				offY := y - r*0.825
				g.Explosion.Vectors = append(g.Explosion.Vectors, scaleVector(vectorData, w, h, offX, offY))
			}
		}
	}
	g.Explosion.Active = true
	maxR := base
	for _, r := range g.Explosion.Radii {
		if r > maxR {
			maxR = r
		}
	}
	if g.killGorillaIfInRadius(x, y, maxR) {
		// handleGorillaKill sets roundOver and other state
	} else {
		g.roundOver = false
	}
	g.recordExplosionDamage(x, y, base)
}

func (g *Game) startGorillaExplosion(idx int) {
	base := g.Settings.NewExplosionRadius
	if base <= 0 {
		base = 16
	}
	if g.Settings.ForceCGA {
		base /= 2
	}
	if g.Settings.UseSound {
		PlayExplosionMelody()
	}
	g.Explosion = Explosion{X: g.Gorillas[idx].X, Y: g.Gorillas[idx].Y}
	if g.Settings.UseOldExplosions {
		for i := 1; i <= int(base); i++ {
			g.Explosion.Radii = append(g.Explosion.Radii, float64(i))
		}
		for i := int(base * 1.5); i >= 1; i-- {
			g.Explosion.Radii = append(g.Explosion.Radii, float64(i))
		}
	} else {
		g.Explosion.Radii = []float64{base * 1.175, base, base * 0.9, base * 0.6, base * 0.45, 0}
		g.Explosion.Colors = []color.Color{
			color.RGBA{128, 128, 128, 255},
			color.RGBA{255, 0, 0, 255},
			color.RGBA{255, 165, 0, 255},
			color.RGBA{255, 255, 0, 255},
			color.RGBA{255, 255, 255, 255},
			color.Black,
		}
		if g.Settings.UseVectorExplosions {
			frames := []float64{base, base * 0.9, base * 0.6, base * 0.45}
			for _, r := range frames {
				w := 2 * r
				h := 2 * r * 0.825
				offX := g.Explosion.X - r
				offY := g.Explosion.Y - r*0.825
				g.Explosion.Vectors = append(g.Explosion.Vectors, scaleVector(vectorData, w, h, offX, offY))
			}
		}
	}
	g.Explosion.Active = true
	g.roundOver = true
	g.recordExplosionDamage(g.Explosion.X, g.Explosion.Y, base)
}

func (g *Game) startVictoryDance(idx int) {
	g.Dance = NewDance(idx, []float64{-3, 0, -3, 0}, g.Gorillas[idx].Y)
}

func (g *Game) stepVictoryDance() {
	if !g.Dance.Active {
		return
	}
	if g.Dance.frame >= len(g.Dance.frames) {
		g.Gorillas[g.Dance.idx].Y = g.Dance.baseY
		g.Dance.Active = false
		return
	}
	offset := g.Dance.frames[g.Dance.frame]
	g.Gorillas[g.Dance.idx].Y = g.Dance.baseY + offset
	g.Dance.frame++
	if g.Settings.UseSound {
		PlayDanceMelody()
	}
}

func (g *Game) Throw() {
	if g.Settings.UseSound {
		PlayBeep()
	}
	if g.Settings.WindFluctuations {
		g.Wind += float64(rand.Intn(5) - 2)
		if g.Wind > 10 {
			g.Wind = 10
		} else if g.Wind < -10 {
			g.Wind = -10
		}
	}
	g.LastAngle[g.Current] = g.Angle
	g.LastPower[g.Current] = g.Power
	g.Angles[g.Current] = g.Angle
	g.Powers[g.Current] = g.Power
	g.Shots[g.Current]++
	start := g.Gorillas[g.Current]
	g.lastStartX = start.X
	g.lastOtherX = g.Gorillas[(g.Current+1)%2].X
	radians := g.Angle * math.Pi / 180
	speed := g.Power / 2
	g.Banana.X = start.X
	g.Banana.Y = start.Y
	if g.Current == 1 {
		g.Banana.VX = -math.Cos(radians) * speed
	} else {
		g.Banana.VX = math.Cos(radians) * speed
	}
	g.Banana.VY = -math.Sin(radians) * speed
	g.lastVX = g.Banana.VX
	g.LastEvent = EventNone
	g.LastEventTicks = 0
	g.LastEventMsg = ""
	g.Banana.Active = true
}

func (g *Game) Step() ShotEvent {
	g.stepVictoryDance()
	if g.LastEventTicks > 0 {
		g.LastEventTicks--
		if g.LastEventTicks == 0 {
			g.LastEvent = EventNone
			g.LastEventMsg = ""
		}
	}
	if g.Explosion.Active {
		if g.Explosion.Frame < len(g.Explosion.Radii)-1 {
			g.Explosion.Frame++
		} else {
			g.Explosion.Active = false
			if g.roundOver {
				cur := g.Current
				g.Reset()
				if g.Settings.VariableWind {
					g.Wind = basicWind()
				}
				if g.Settings.WinnerFirst {
					g.setCurrent(cur)
				} else {
					g.setCurrent((cur + 1) % 2)
				}
			}
		}
		return EventNone
	}

	if !g.Banana.Active {
		return EventNone
	}
	oldX := g.Banana.X
	oldY := g.Banana.Y
	g.Banana.X += g.Banana.VX
	g.Banana.Y += g.Banana.VY
	// apply gravity scaled to the configured constant
	// default behaviour uses DefaultGravity which equals Gravity initially
	g.Banana.VY += g.Gravity / 34
	g.Banana.VX += g.Wind / 20
	if g.Banana.Y > float64(g.Height) {
		if g.Banana.VY > groundBounceThreshold {
			g.Banana.Y = float64(g.Height)
			g.Banana.VY = -g.Banana.VY * groundBounceFactor
		} else {
			g.Banana.Active = false
			g.evaluateMiss()
			g.setCurrent((g.Current + 1) % 2)
			return g.LastEvent
		}
	}
	if hit := g.gorillaHitBetween(oldX, oldY, g.Banana.X, g.Banana.Y); hit >= 0 {
		g.Banana.Active = false
		g.handleGorillaKill(hit)
		g.startGorillaExplosion(hit)
		return g.LastEvent
	}
	for i, gr := range g.Gorillas {
		hit := false
		if g.HitMap != nil {
			if g.HitMap.GorillaHitAt(int(math.Round(g.Banana.X)), int(math.Round(g.Banana.Y))) == i {
				hit = true
			}
		}
		if !hit && math.Abs(gr.X-g.Banana.X) < 5 && math.Abs(gr.Y-g.Banana.Y) < 10 {
			hit = true
		}
		if hit {
			g.Banana.Active = false
			shooter := g.Current
			winner := shooter
			event := EventNone
			if i == shooter {
				winner = (shooter + 1) % 2
				event = EventSelf
				g.LastEvent = event
				g.LastEventTicks = eventDisplayTicks
				g.LastEventMsg = EventMessage(event)
			}
			g.Wins[winner]++
			g.TotalWins[winner]++
			if g.League != nil {
				g.League.RecordRound(g.Players[0], g.Players[1], winner, g.Shots[shooter])
				g.League.Save()
			}
			g.Shots = [2]int{}
			g.SaveScores()
			g.startGorillaExplosion(i)
			g.startVictoryDance(winner)
			g.setCurrent(winner)
			if g.Settings.UseSound && event != EventNone {
				PlayBeep()
			}
			return event
		}
	}
	bw := float64(g.Width) / float64(g.BuildingCount)
	idx := int(g.Banana.X / bw)
	if idx >= 0 && idx < g.BuildingCount && g.Banana.Y < float64(g.Height) &&
		g.Banana.Y > float64(g.Height)-g.Buildings[idx].H {
		if !g.pointInDamage(idx, g.Banana.X, g.Banana.Y) {
			g.Banana.Active = false
			g.startExplosion(g.Banana.X, g.Banana.Y)
			if g.roundOver {
				return g.LastEvent
			}
			g.evaluateMiss()
			g.setCurrent((g.Current + 1) % 2)
			return g.LastEvent
		}
	}
	if g.Banana.Y > float64(g.Height) || g.Banana.X < 0 || g.Banana.X >= float64(g.Width) {
		g.Banana.Active = false
		g.evaluateMiss()
		g.setCurrent((g.Current + 1) % 2)
	}
	return g.LastEvent
}
func (g *Game) testShot(angle, power float64) bool {
	sim := *g
	sim.Angle = angle
	sim.Power = power
	sim.Throw()
	for i := 0; i < 500 && (sim.Banana.Active || sim.Explosion.Active); i++ {
		sim.Step()
	}
	return sim.Wins[g.Current] > g.Wins[g.Current]
}

// FindShot searches for an angle and power likely to hit the opponent.
func (g *Game) FindShot() (angle, power float64) {
	for a := 15.0; a <= 75; a += 1 {
		for p := 20.0; p <= 100; p += 2 {
			if g.testShot(a, p) {
				return a, p
			}
		}
	}
	return 45, 50
}

// AutoShot selects a shot using FindShot and throws the banana.
func (g *Game) AutoShot() {
	g.Angle, g.Power = g.FindShot()
	g.Throw()
}

// evaluateMiss analyses a non-scoring shot and sets LastEvent if it was weak or backwards.
func (g *Game) evaluateMiss() {
	dxToOther := g.lastOtherX - g.lastStartX
	dxShot := g.Banana.X - g.lastStartX
	if g.lastVX*dxToOther < 0 {
		g.LastEvent = EventBackwards
		g.LastEventTicks = eventDisplayTicks
		g.LastEventMsg = EventMessage(EventBackwards)
		if g.Settings.UseSound {
			PlayBeep()
		}
		return
	}
	if math.Abs(dxShot) < math.Abs(dxToOther)/3 {
		g.LastEvent = EventWeak
		g.LastEventTicks = eventDisplayTicks
		g.LastEventMsg = EventMessage(EventWeak)
		if g.Settings.UseSound {
			PlayBeep()
		}
	}
}

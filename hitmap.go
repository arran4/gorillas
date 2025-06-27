package gorillas

import "image"

const (
	hitMapEmpty    = 0
	hitMapBuilding = 1
	hitMapGorilla0 = 2
	hitMapGorilla1 = 3
)

// HitMap is a simple bitmap identifying hittable objects.
type HitMap struct {
	width, height int
	data          []byte
}

// NewHitMap allocates a HitMap for the given dimensions.
func NewHitMap(w, h int) *HitMap {
	return &HitMap{width: w, height: h, data: make([]byte, w*h)}
}

func (m *HitMap) index(x, y int) int { return y*m.width + x }

// At returns the value stored at the given coordinates.
func (m *HitMap) At(x, y int) byte {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return hitMapEmpty
	}
	return m.data[m.index(x, y)]
}

// Set assigns val to the pixel at x,y if in range.
func (m *HitMap) Set(x, y int, val byte) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return
	}
	m.data[m.index(x, y)] = val
}

// DrawRect fills a rectangular area with val.
func (m *HitMap) DrawRect(x1, y1, x2, y2 int, val byte) {
	if x1 < 0 {
		x1 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	if x2 > m.width {
		x2 = m.width
	}
	if y2 > m.height {
		y2 = m.height
	}
	for y := y1; y < y2; y++ {
		idx := y * m.width
		for x := x1; x < x2; x++ {
			m.data[idx+x] = val
		}
	}
}

// DrawCircle fills a circle with val.
func (m *HitMap) DrawCircle(cx, cy, r int, val byte) {
	r2 := r * r
	for y := cy - r; y <= cy+r; y++ {
		if y < 0 || y >= m.height {
			continue
		}
		dy := y - cy
		for x := cx - r; x <= cx+r; x++ {
			if x < 0 || x >= m.width {
				continue
			}
			dx := x - cx
			if dx*dx+dy*dy <= r2 {
				m.Set(x, y, val)
			}
		}
	}
}

// ClearCircle clears pixels in a circular area.
func (m *HitMap) ClearCircle(cx, cy, r int) {
	r2 := r * r
	for y := cy - r; y <= cy+r; y++ {
		if y < 0 || y >= m.height {
			continue
		}
		dy := y - cy
		for x := cx - r; x <= cx+r; x++ {
			if x < 0 || x >= m.width {
				continue
			}
			dx := x - cx
			if dx*dx+dy*dy <= r2 {
				m.Set(x, y, hitMapEmpty)
			}
		}
	}
}

// AnyValueInCircle reports whether val occurs within the circle.
func (m *HitMap) AnyValueInCircle(cx, cy, r int, val byte) bool {
	r2 := r * r
	for y := cy - r; y <= cy+r; y++ {
		if y < 0 || y >= m.height {
			continue
		}
		dy := y - cy
		for x := cx - r; x <= cx+r; x++ {
			if x < 0 || x >= m.width {
				continue
			}
			dx := x - cx
			if dx*dx+dy*dy <= r2 && m.At(x, y) == val {
				return true
			}
		}
	}
	return false
}

// DrawGorilla marks the gorilla location with the appropriate value.
func (m *HitMap) DrawGorilla(x, y int, idx int, r int) {
	val := hitMapGorilla0
	if idx == 1 {
		val = hitMapGorilla1
	}
	m.DrawCircle(x, y, r, byte(val))
}

// DrawGorillaImage marks non-transparent pixels of img using the same anchor
// position as DrawGorilla (bottom centre). This allows hit detection to match
// the rendered gorilla sprite.
func (m *HitMap) DrawGorillaImage(x, y int, idx int, img image.Image) {
	val := byte(hitMapGorilla0)
	if idx == 1 {
		val = hitMapGorilla1
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	baseX := x - w/2
	baseY := y - h
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			_, _, _, a := img.At(b.Min.X+xx, b.Min.Y+yy).RGBA()
			if a != 0 {
				m.Set(baseX+xx, baseY+yy, val)
			}
		}
	}
}

// GorillaValue returns the gorilla index stored at the coordinate or -1.
func (m *HitMap) GorillaValue(x, y int) int {
	v := m.At(x, y)
	switch v {
	case hitMapGorilla0:
		return 0
	case hitMapGorilla1:
		return 1
	default:
		return -1
	}
}

// ClearGorilla removes gorilla pixels for idx.
func (m *HitMap) ClearGorilla(x, y int, idx int, r int) {
	val := byte(hitMapGorilla0)
	if idx == 1 {
		val = hitMapGorilla1
	}
	r2 := r * r
	for yy := y - r; yy <= y+r; yy++ {
		if yy < 0 || yy >= m.height {
			continue
		}
		dy := yy - y
		for xx := x - r; xx <= x+r; xx++ {
			if xx < 0 || xx >= m.width {
				continue
			}
			dx := xx - x
			if dx*dx+dy*dy <= r2 && m.At(xx, yy) == val {
				m.Set(xx, yy, hitMapEmpty)
			}
		}
	}
}

// GorillaHitInCircle returns the index of a gorilla found within the circle, or -1.
func (m *HitMap) GorillaHitInCircle(cx, cy, r int) int {
	if m.AnyValueInCircle(cx, cy, r, hitMapGorilla0) {
		return 0
	}
	if m.AnyValueInCircle(cx, cy, r, hitMapGorilla1) {
		return 1
	}
	return -1
}

// GorillaHitAt checks if a gorilla occupies the exact coordinate.
func (m *HitMap) GorillaHitAt(x, y int) int {
	v := m.At(x, y)
	switch v {
	case hitMapGorilla0:
		return 0
	case hitMapGorilla1:
		return 1
	default:
		return -1
	}
}

// AddBuilding draws the building rectangle on the hit map.
func (m *HitMap) AddBuilding(x1, y1, x2, y2 int) {
	m.DrawRect(x1, y1, x2, y2, hitMapBuilding)
}

// ClearBuildingArea clears a circular area from the hit map.
func (m *HitMap) ClearBuildingArea(cx, cy, r int) {
	m.ClearCircle(cx, cy, r)
}

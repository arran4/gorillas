package gorillas

// VectorPoint represents a relative coordinate in a vector shape.
type VectorPoint struct {
	X, Y float64
}

// vectorData holds the points used by the original BASIC vector explosion.
var vectorData = []VectorPoint{
	{0.582, 0.988}, {0.608, 0.850}, {0.663, 0.788}, {0.738, 0.800},
	{0.863, 0.838}, {0.813, 0.713}, {0.819, 0.650}, {0.875, 0.588},
	{1.000, 0.563}, {0.850, 0.450}, {0.825, 0.400}, {0.830, 0.340},
	{0.925, 0.238}, {0.775, 0.243}, {0.694, 0.225}, {0.650, 0.188}, {0.630, 0.105},
	{0.625, 0.025}, {0.535, 0.150}, {0.475, 0.175}, {0.425, 0.150},
	{0.325, 0.044}, {0.325, 0.150}, {0.315, 0.208}, {0.288, 0.250}, {0.225, 0.275},
	{0.053, 0.288}, {0.150, 0.392}, {0.175, 0.463}, {0.144, 0.525},
	{0.025, 0.638}, {0.163, 0.650}, {0.225, 0.693}, {0.250, 0.775},
	{0.225, 0.905}, {0.360, 0.825}, {0.450, 0.823}, {0.525, 0.863},
	{0.582, 0.988},
}

// scaleVector returns the given data scaled to width and height and translated
// by offset.
func scaleVector(data []VectorPoint, width, height, offX, offY float64) []VectorPoint {
	pts := make([]VectorPoint, len(data))
	for i, p := range data {
		pts[i] = VectorPoint{
			X: offX + p.X*width,
			Y: offY + p.Y*height,
		}
	}
	return pts
}

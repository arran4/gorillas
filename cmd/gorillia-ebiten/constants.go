//go:build !test

package main

// charW and charH define the width and height of ASCII characters used
// in the intro movie and menu screens.
const (
	charW = 6
	charH = 16
	// gorillaScale controls how large the gorilla sprite appears in game
	// The previous value of 4 produced gorillas far larger than the
	// buildings drawn at 800x600, so reduce this to better match the
	// original QBASIC proportions.
	gorillaScale = 2
	// bananaScale is larger than gorillaScale so bananas remain visible
	// as the gorilla sprite was reduced in size.
	bananaScale = gorillaScale * 2
	sunScale    = gorillaScale
)

// gorillaFrames represents the ASCII gorilla animation frames shared by
// the intro movie and menu state.
var gorillaFrames = [][]string{
	{
		" O ",
		"/|\\",
		"/ \\",
	},
	{
		" O ",
		"/| ",
		"/ \\",
	},
	{
		" O ",
		" |\\",
		"/ \\",
	},
}

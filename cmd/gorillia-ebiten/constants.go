//go:build !test

package main

// charW and charH define the width and height of ASCII characters used
// in the intro movie and menu screens.
const (
	charW = 6
	charH = 16
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

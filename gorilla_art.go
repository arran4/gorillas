package gorillas

import (
	"bufio"
	"os"
	"strings"
)

// LoadGorillaArt reads ASCII art frames from a file. Frames are separated by
// a line containing only "===". Whitespace is preserved.
func LoadGorillaArt(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var frames [][]string
	var cur []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "===" {
			cur = trimBlankLines(cur)
			if len(cur) > 0 {
				frames = append(frames, cur)
			}
			cur = []string{}
			continue
		}
		cur = append(cur, strings.TrimRight(line, " \t"))
	}
	cur = trimBlankLines(cur)
	if len(cur) > 0 {
		frames = append(frames, cur)
	}
	return frames, scanner.Err()
}

func trimBlankLines(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// FrameWidth returns the maximum line width of a frame of ASCII art.
func FrameWidth(frame []string) int {
	w := 0
	for _, line := range frame {
		if len(line) > w {
			w = len(line)
		}
	}
	return w
}

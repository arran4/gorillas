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
			if len(cur) > 0 {
				frames = append(frames, cur)
				cur = []string{}
			}
			continue
		}
		cur = append(cur, line)
	}
	if len(cur) > 0 {
		frames = append(frames, cur)
	}
	return frames, scanner.Err()
}

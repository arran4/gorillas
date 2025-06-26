package gorillas

import (
	"bufio"
	"os"
	"strings"
)

// LoadInstructions parses the SlidyText section from gorillas.bas and
// returns the credit and instruction lines.
func LoadInstructions() ([]string, error) {
	f, err := os.Open("gorillas.bas")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []string
	found := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !found {
			if strings.HasPrefix(line, "SlidyText:") {
				found = true
			}
			continue
		}
		if strings.HasPrefix(line, "PartingMessage:") {
			break
		}
		if strings.HasPrefix(line, "DATA") {
			fields := strings.SplitN(strings.TrimSpace(line[4:]), ",", 2)
			if len(fields) > 0 {
				token := strings.TrimSpace(fields[0])
				if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
					lines = append(lines, strings.Trim(token, "\""))
				}
			}
		}
	}
	return lines, scanner.Err()
}

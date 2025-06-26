package gorillas

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// LoadSettings reads configuration from environment variables.
// The flag package can override these values using its BoolVar API.
// Recognised variables:
//
//		GORILLAS_SOUND - set to 'true' or 'false' to enable or disable sound.
//		GORILLAS_OLD_EXPLOSIONS - 'true' to use the old explosion style.
//	     GORILLAS_WINNER_FIRST - 'true' if round winner starts the next round.
//		GORILLAS_EXPLOSION_RADIUS - floating point radius for new explosions.
func loadSettingsFile(path string, s *Settings) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		switch strings.ToUpper(key) {
		case "USESOUND":
			if b, err := strconv.ParseBool(val); err == nil {
				s.UseSound = b
			} else if strings.EqualFold(val, "YES") {
				s.UseSound = true
			} else if strings.EqualFold(val, "NO") {
				s.UseSound = false
			}
		case "USEOLDEXPLOSIONS":
			if b, err := strconv.ParseBool(val); err == nil {
				s.UseOldExplosions = b
			} else if strings.EqualFold(val, "YES") {
				s.UseOldExplosions = true
			} else if strings.EqualFold(val, "NO") {
				s.UseOldExplosions = false
			}
		case "NEWEXPLOSIONRADIUS":
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				s.NewExplosionRadius = f
			}
		case "DEFAULTGRAVITY":
			if f, err := strconv.ParseFloat(val, 64); err == nil && f > 0 {
				s.DefaultGravity = f
			}
		case "DEFAULTROUNDQTY":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				s.DefaultRoundQty = n
			}
		case "WINNERFIRST":
			if b, err := strconv.ParseBool(val); err == nil {
				s.WinnerFirst = b
			} else if strings.EqualFold(val, "YES") {
				s.WinnerFirst = true
			} else if strings.EqualFold(val, "NO") {
				s.WinnerFirst = false
			}
		}
	}
}

func LoadSettings() Settings {
	s := DefaultSettings()
	loadSettingsFile("gorillas.ini", &s)
	if v, ok := os.LookupEnv("GORILLAS_SOUND"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.UseSound = b
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_OLD_EXPLOSIONS"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.UseOldExplosions = b
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_EXPLOSION_RADIUS"); ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			s.NewExplosionRadius = f
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_GRAVITY"); ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			s.DefaultGravity = f
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_ROUNDS"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			s.DefaultRoundQty = n
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_WINNER_FIRST"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.WinnerFirst = b
		}
	}
	return s
}

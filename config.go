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
//	GORILLAS_SOUND - set to 'true' or 'false' to enable or disable sound.
//	GORILLAS_OLD_EXPLOSIONS - 'true' to use the old explosion style.
//	GORILLAS_EXPLOSION_RADIUS - floating point radius for new explosions.
//	GORILLAS_GRAVITY - gravitational constant used in game physics.
//	GORILLAS_ROUNDS - default number of rounds to play.
//	GORILLAS_SLIDING_TEXT - 'true' to enable sliding text effects.
//	GORILLAS_SHOW_INTRO - 'true' to display the intro sequence.
//	GORILLAS_FORCE_CGA - 'true' to force CGA mode graphics.
//	     GORILLAS_WINNER_FIRST - 'true' if round winner starts the next round.
//	GORILLAS_VARIABLE_WIND - 'true' to mimic BASIC wind changes each round.
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
		case "USESLIDINGTEXT":
			if b, err := strconv.ParseBool(val); err == nil {
				s.UseSlidingText = b
			} else if strings.EqualFold(val, "YES") {
				s.UseSlidingText = true
			} else if strings.EqualFold(val, "NO") {
				s.UseSlidingText = false
			}
		case "DEFAULTGRAVITY":
			if f, err := strconv.ParseFloat(val, 64); err == nil && f > 0 {
				s.DefaultGravity = f
			}
		case "DEFAULTROUNDQTY":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				s.DefaultRoundQty = n
			}
		case "SHOWINTRO":
			if b, err := strconv.ParseBool(val); err == nil {
				s.ShowIntro = b
			} else if strings.EqualFold(val, "YES") {
				s.ShowIntro = true
			} else if strings.EqualFold(val, "NO") {
				s.ShowIntro = false
			}
		case "FORCECGA":
			if b, err := strconv.ParseBool(val); err == nil {
				s.ForceCGA = b
			} else if strings.EqualFold(val, "YES") {
				s.ForceCGA = true
			} else if strings.EqualFold(val, "NO") {
				s.ForceCGA = false
			}
		case "WINNERFIRST":
			if b, err := strconv.ParseBool(val); err == nil {
				s.WinnerFirst = b
			} else if strings.EqualFold(val, "YES") {
				s.WinnerFirst = true
			} else if strings.EqualFold(val, "NO") {
				s.WinnerFirst = false
			}
		case "VARIABLEWIND":
			if b, err := strconv.ParseBool(val); err == nil {
				s.VariableWind = b
			} else if strings.EqualFold(val, "YES") {
				s.VariableWind = true
			} else if strings.EqualFold(val, "NO") {
				s.VariableWind = false
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
	if v, ok := os.LookupEnv("GORILLAS_SLIDING_TEXT"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.UseSlidingText = b
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_SHOW_INTRO"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.ShowIntro = b
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_FORCE_CGA"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.ForceCGA = b
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_WINNER_FIRST"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.WinnerFirst = b
		}
	}
	if v, ok := os.LookupEnv("GORILLAS_VARIABLE_WIND"); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			s.VariableWind = b
		}
	}
	return s
}

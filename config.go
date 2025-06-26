package gorillas

import (
	"os"
	"strconv"
)

// LoadSettings reads configuration from environment variables.
// The flag package can override these values using its BoolVar API.
// Recognised variables:
//
//	GORILLAS_SOUND - set to 'true' or 'false' to enable or disable sound.
//	GORILLAS_OLD_EXPLOSIONS - 'true' to use the old explosion style.
//	GORILLAS_EXPLOSION_RADIUS - floating point radius for new explosions.
func LoadSettings() Settings {
	s := DefaultSettings()
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
	return s
}

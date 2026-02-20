package gorillas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettingsRejectsInvalidNumericEnvValues(t *testing.T) {
	t.Setenv("GORILLAS_EXPLOSION_RADIUS", "-1")
	t.Setenv("GORILLAS_GRAVITY", "0")
	t.Setenv("GORILLAS_ROUNDS", "-5")

	defaults := DefaultSettings()
	s := LoadSettings()

	if s.NewExplosionRadius != defaults.NewExplosionRadius {
		t.Fatalf("expected default explosion radius %f, got %f", defaults.NewExplosionRadius, s.NewExplosionRadius)
	}
	if s.DefaultGravity != defaults.DefaultGravity {
		t.Fatalf("expected default gravity %f, got %f", defaults.DefaultGravity, s.DefaultGravity)
	}
	if s.DefaultRoundQty != defaults.DefaultRoundQty {
		t.Fatalf("expected default rounds %d, got %d", defaults.DefaultRoundQty, s.DefaultRoundQty)
	}
}

func TestLoadSettingsTrimsNumericEnvValues(t *testing.T) {
	t.Setenv("GORILLAS_EXPLOSION_RADIUS", " 25.5 ")
	t.Setenv("GORILLAS_GRAVITY", " 15 ")
	t.Setenv("GORILLAS_ROUNDS", " 9 ")

	s := LoadSettings()

	if s.NewExplosionRadius != 25.5 {
		t.Fatalf("expected explosion radius 25.5, got %f", s.NewExplosionRadius)
	}
	if s.DefaultGravity != 15 {
		t.Fatalf("expected gravity 15, got %f", s.DefaultGravity)
	}
	if s.DefaultRoundQty != 9 {
		t.Fatalf("expected rounds 9, got %d", s.DefaultRoundQty)
	}
}

func TestLoadSettingsFileRejectsInvalidNumericValues(t *testing.T) {
	dir := t.TempDir()
	ini := filepath.Join(dir, "gorillas.ini")
	data := []byte("" +
		"NewExplosionRadius=-20.5\n" +
		"DefaultGravity=0\n" +
		"DefaultRoundQty=-7\n")
	if err := os.WriteFile(ini, data, 0644); err != nil {
		t.Fatal(err)
	}
	defaults := DefaultSettings()
	s := defaults
	loadSettingsFile(ini, &s)
	if s.NewExplosionRadius != defaults.NewExplosionRadius {
		t.Fatalf("expected default explosion radius %f, got %f", defaults.NewExplosionRadius, s.NewExplosionRadius)
	}
	if s.DefaultGravity != defaults.DefaultGravity {
		t.Fatalf("expected default gravity %f, got %f", defaults.DefaultGravity, s.DefaultGravity)
	}
	if s.DefaultRoundQty != defaults.DefaultRoundQty {
		t.Fatalf("expected default rounds %d, got %d", defaults.DefaultRoundQty, s.DefaultRoundQty)
	}
}

func TestLoadSettingsFile(t *testing.T) {
	dir := t.TempDir()
	ini := filepath.Join(dir, "gorillas.ini")
	data := []byte("" +
		"UseSound=no\n" +
		"UseOldExplosions=yes\n" +
		"NewExplosionRadius=20.5\n" +
		"DefaultGravity=30\n" +
		"DefaultRoundQty=7\n" +
		"UseSlidingText=yes\n" +
		"ShowIntro=no\n" +
		"ForceCGA=yes\n" +
		"WinnerFirst=yes\n" +
		"VariableWind=yes\n" +
		"WindFluctuations=yes\n" +
		"UseVectorExplosions=yes\n")
	if err := os.WriteFile(ini, data, 0644); err != nil {
		t.Fatal(err)
	}
	s := DefaultSettings()
	loadSettingsFile(ini, &s)
	if s.UseSound != false {
		t.Errorf("expected UseSound=false got %v", s.UseSound)
	}
	if !s.UseOldExplosions {
		t.Errorf("expected UseOldExplosions=true")
	}
	if s.NewExplosionRadius != 20.5 {
		t.Errorf("unexpected radius %f", s.NewExplosionRadius)
	}
	if s.DefaultGravity != 30 {
		t.Errorf("unexpected gravity %f", s.DefaultGravity)
	}
	if s.DefaultRoundQty != 7 {
		t.Errorf("unexpected round qty %d", s.DefaultRoundQty)
	}
	if !s.UseSlidingText {
		t.Errorf("expected UseSlidingText=true")
	}
	if s.ShowIntro {
		t.Errorf("expected ShowIntro=false")
	}
	if !s.ForceCGA {
		t.Errorf("expected ForceCGA=true")
	}
	if !s.WinnerFirst {
		t.Errorf("expected WinnerFirst=true")
	}
	if !s.VariableWind {
		t.Errorf("expected VariableWind=true")
	}
	if !s.WindFluctuations {
		t.Errorf("expected WindFluctuations=true")
	}
	if !s.UseVectorExplosions {
		t.Errorf("expected UseVectorExplosions=true")
	}
}

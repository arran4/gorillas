package gorillas

import (
	"os"
	"path/filepath"
	"testing"
)

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
		"WindFluctuations=yes\n")
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
}

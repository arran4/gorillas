//go:build !linux

package main

type joystick struct {
    axis [2]int16
    btn  [16]bool
}

func openJoystick() (*joystick, error) { return nil, nil }
func (j *joystick) poll()              {}

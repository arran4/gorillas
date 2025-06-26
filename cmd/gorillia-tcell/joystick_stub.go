//go:build !linux

package main

type joystick struct{}

func openJoystick() (*joystick, error) { return nil, nil }
func (j *joystick) poll()              {}

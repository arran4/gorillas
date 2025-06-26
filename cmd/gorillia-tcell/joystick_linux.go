//go:build linux

package main

import (
	"encoding/binary"
	"errors"
	"golang.org/x/sys/unix"
	"os"
	"strconv"
)

type joystickEvent struct {
	Time   uint32
	Value  int16
	Type   uint8
	Number uint8
}

const (
	jsEventButton = 0x01
	jsEventAxis   = 0x02
	jsEventInit   = 0x80
)

type joystick struct {
	f    *os.File
	axis [2]int16
	btn  [16]bool
}

func openJoystick() (*joystick, error) {
	for i := 0; i < 4; i++ {
		p := "/dev/input/js" + strconv.Itoa(i)
		if f, err := os.Open(p); err == nil {
			unix.SetNonblock(int(f.Fd()), true)
			return &joystick{f: f}, nil
		}
	}
	return nil, os.ErrNotExist
}

func (j *joystick) poll() {
	if j == nil || j.f == nil {
		return
	}
	for {
		var e joystickEvent
		err := binary.Read(j.f, binary.LittleEndian, &e)
		if err != nil {
			if errors.Is(err, unix.EAGAIN) || errors.Is(err, os.ErrClosed) || err == os.ErrClosed {
				break
			}
			return
		}
		et := e.Type &^ jsEventInit
		switch et {
		case jsEventAxis:
			if int(e.Number) < len(j.axis) {
				j.axis[e.Number] = e.Value
			}
		case jsEventButton:
			if int(e.Number) < len(j.btn) {
				j.btn[e.Number] = e.Value != 0
			}
		}
	}
}

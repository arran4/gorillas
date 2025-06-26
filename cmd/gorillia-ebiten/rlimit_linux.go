//go:build linux && !test

package main

import "golang.org/x/sys/unix"

func increaseRLimit() error {
	var r unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &r); err != nil {
		return err
	}
	if r.Cur < r.Max {
		r.Cur = r.Max
		if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &r); err != nil {
			return err
		}
	}
	return nil
}

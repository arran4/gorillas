//go:build linux && !test

package main

import "golang.org/x/sys/unix"

func increaseRLimit() {
	var r unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &r); err != nil {
		return
	}
	if r.Cur < r.Max {
		r.Cur = r.Max
		_ = unix.Setrlimit(unix.RLIMIT_NOFILE, &r)
	}
}

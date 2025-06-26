//go:build !linux && !test

package main

func increaseRLimit() error { return nil }

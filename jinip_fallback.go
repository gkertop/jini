//go:build !darwin && !linux && !freebsd && !windows

package jini

func getLineBreak() string {
	return "\n"
}

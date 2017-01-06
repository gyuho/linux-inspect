package psn

import "os/user"

// SSEntry is a socket entry.
// Simplied from 'NetTCP'.
type SSEntry struct {
	Protocol string

	Program string
	State   string
	PID     int64

	LocalIP   string
	LocalPort int64

	RemoteIP   string
	RemotePort int64

	User user.User
}

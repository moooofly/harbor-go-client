package term

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)

// Termios is the Unix API for terminal I/O.
type Termios unix.Termios

// +build darwin freebsd openbsd netbsd

package term

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TIOCGETA
	setTermios = unix.TIOCSETA
)

// Termios is the Unix API for terminal I/O.
type Termios unix.Termios

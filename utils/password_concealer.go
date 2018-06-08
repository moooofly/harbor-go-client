package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)

var (
	// ErrInvalidState is returned if the state of the terminal is invalid.
	ErrInvalidState = errors.New("Invalid terminal state")
)

// Termios is the Unix API for terminal I/O.
type Termios unix.Termios

// State represents the state of the terminal.
type State struct {
	termios Termios
}

func tcget(fd uintptr, p *Termios) syscall.Errno {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(getTermios), uintptr(unsafe.Pointer(p)))
	return err
}

func tcset(fd uintptr, p *Termios) syscall.Errno {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, setTermios, uintptr(unsafe.Pointer(p)))
	return err
}

// RestoreTerminal restores the terminal connected to the given file descriptor
// to a previous state.
func RestoreTerminal(fd uintptr, state *State) error {
	if state == nil {
		return ErrInvalidState
	}
	if err := tcset(fd, &state.termios); err != 0 {
		return err
	}
	return nil
}

func handleInterrupt(fd uintptr, state *State) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		for range sigchan {
			// quit cleanly and the new terminal item is on a new line
			fmt.Println()
			signal.Stop(sigchan)
			close(sigchan)
			RestoreTerminal(fd, state)
			os.Exit(1)
		}
	}()
}

func readInput(in io.Reader, out io.Writer) string {
	reader := bufio.NewReader(in)
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Fprintln(out, err.Error())
		os.Exit(1)
	}
	return string(line)
}

// ReadPasswordFromTerm gets user password from stdin without showing on screen
func ReadPasswordFromTerm() (string, error) {

	var oldState State
	if err := tcget(os.Stdin.Fd(), &oldState.termios); err != 0 {
		return "", err
	}

	fmt.Fprintf(os.Stdout, "Password: ")

	newState := oldState.termios
	newState.Lflag &^= unix.ECHO

	if err := tcset(os.Stdin.Fd(), &newState); err != 0 {
		return "", err
	}
	handleInterrupt(os.Stdin.Fd(), &oldState)

	passwd := readInput(os.Stdin, os.Stdout)
	fmt.Fprint(os.Stdout, "\n")

	if err := tcset(os.Stdin.Fd(), &oldState.termios); err != 0 {
		return "", err
	}

	return passwd, nil
}

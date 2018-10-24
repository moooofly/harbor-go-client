package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/moooofly/harbor-go-client/utils/term"
)

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

	oldState, err := term.SaveState(os.Stdin.Fd())
	if err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stdout, "Password: ")

	err = term.DisableEcho(os.Stdin.Fd(), oldState)
	if err != nil {
		return "", err
	}

	passwd := readInput(os.Stdin, os.Stdout)
	fmt.Fprint(os.Stdout, "\n")

	err = term.RestoreTerminal(os.Stdin.Fd(), oldState)
	if err != nil {
		return "", err
	}

	return passwd, nil
}

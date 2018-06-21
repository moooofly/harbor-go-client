package utils

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func init() {
	Parser.AddCommand("version",
		"Show version info.",
		"Format: <client-version> (<golang-runtime-version> on <GOOS>/<GOARCH>; <Compiler>)",
		&verinfo)
}

type verInfo struct {
}

var verinfo verInfo

func (x *verInfo) Execute(args []string) error {
	PrintVersion()
	return nil
}

func version() string {

	f, err := os.Open(versionfile)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	v, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}

	return fmt.Sprintf("%s (%s on %s/%s; %s)",
		strings.Trim(v, "\n"), runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
}

// PrintVersion print version info.
func PrintVersion() {
	fmt.Println(version())
}

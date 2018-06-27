package utils

import (
	"fmt"
)

func init() {
	Parser.AddCommand("version",
		"Show version info.",
		"Show version infos as \"| Type | Value |\"",
		&verinfo)
}

type verInfo struct {
}

var verinfo verInfo

func (x *verInfo) Execute(args []string) error {
	PrintVersion()
	return nil
}

// PrintVersion print version info.
func PrintVersion() {
	PrintLogo()
	fmt.Println("+----------------------+------------------------------------------+")
	fmt.Printf("| % -20s | % -40s |\n", "Client Version", ClientVersion)
	fmt.Printf("| % -20s | % -40s |\n", "Go Version", GoVersion)
	fmt.Printf("| % -20s | % -40s |\n", "UTC Build Time", UTCBuildTime)
	fmt.Printf("| % -20s | % -40s |\n", "Git Branch", GitBranch)
	fmt.Printf("| % -20s | % -40s |\n", "Git Tag", GitTag)
	fmt.Printf("| % -20s | % -40s |\n", "Git Hash", GitHash)
	fmt.Println("+----------------------+------------------------------------------+")
}

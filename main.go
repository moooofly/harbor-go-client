package main

import (
	"os"

	"git.llsapp.com/fei.sun/harbor-go-client/utils"

	"github.com/jessevdk/go-flags"

	_ "git.llsapp.com/fei.sun/harbor-go-client/api"
)

func main() {
	if _, err := utils.Parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

package main

import (
	"os"

	"github.com/arch-err/yqp/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		// error is discarded as cobra already reported it
		os.Exit(1)
	}
}

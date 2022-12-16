package main

import (
	"github.com/qinsheng99/issue/cmd"
	"log"
	"os"
)

func main() {
	issueCmd := cmd.Cmd()
	issueCmd.SetErr(os.Stderr)
	issueCmd.SetOut(os.Stdout)
	issueCmd.SetIn(os.Stdin)

	err := issueCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

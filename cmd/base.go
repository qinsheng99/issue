package cmd

import (
	"fmt"
	"github.com/qinsheng99/issue/util"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var issue = &cobra.Command{
	Use:  "issue",
	Long: "issue command can create openeuler issue",
	Run: func(cmd *cobra.Command, args []string) {
		checkErr(cmd.Help())
	},
}

func checkErr(err error) {
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}

type CommandGroup struct {
	Message  string
	Commands []*cobra.Command
}

type CommandGroups []CommandGroup

func (g CommandGroups) Add(c *cobra.Command) {
	for _, group := range g {
		c.AddCommand(group.Commands...)
	}
}

func (g CommandGroups) Has(c *cobra.Command) bool {
	for _, group := range g {
		for _, command := range group.Commands {
			if command == c {
				return true
			}
		}
	}
	return false
}

type base struct {
	Streams
	util.ReqImpl
}

type baseResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Streams struct {
	// In default os.Stdin
	In io.Reader

	// Out default os.Stdout
	Out io.Writer

	// ErrOut default os.Stderr
	ErrOut io.Writer
}

func Cmd() *cobra.Command {
	s := base{
		Streams: Streams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		},
		ReqImpl: util.NewRequest(nil),
	}

	group := CommandGroups{
		{
			Commands: []*cobra.Command{
				newCmdGet(s),
			},
		},
		{
			Commands: []*cobra.Command{
				newCmdCreate(s),
			},
		},
	}

	group.Add(issue)

	return issue
}

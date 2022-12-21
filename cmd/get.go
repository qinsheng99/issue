package cmd

import (
	"github.com/qinsheng99/issue/util"
	"github.com/spf13/cobra"
)

type getOption struct {
	Streams
	h util.ReqImpl
}

func newGetOption(s base) *getOption {
	return &getOption{Streams: s.Streams, h: s.ReqImpl}
}

const getExample = `
		# List all repo in output
		issue get repo [options]
		# List all issue_type in output
		issue get issue-type [options]
	`

func newCmdGet(s base) *cobra.Command {
	o := newGetOption(s)

	cmd := &cobra.Command{
		Use:     "get [repo|issue_type]",
		Short:   "get resource",
		Example: getExample,
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(cmd.Help())
		},
	}

	cmd.AddCommand(newRepoCmd(o.Streams, o.h))
	cmd.AddCommand(newIssueTypeCmd(o.Streams, o.h))

	return cmd
}

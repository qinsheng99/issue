package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"sort"

	"github.com/qinsheng99/issue/util"
	"github.com/spf13/cobra"
)

const basefile = "%s.txt"

type issueOption struct {
	Streams
	h util.ReqImpl

	name     string
	file     bool
	filename string
}

func newIssueOption(s Streams, h util.ReqImpl) *issueOption {
	return &issueOption{Streams: s, h: h}
}

func newIssueTypeCmd(s Streams, h util.ReqImpl) *cobra.Command {
	o := newIssueOption(s, h)

	cmd := &cobra.Command{
		Use:     "issue_type",
		Aliases: []string{"it"},
		Short:   "get openeuler community issue type",
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.name, "name", "n", o.name, "issue type name")
	cmd.Flags().StringVar(&o.filename, "filename", o.filename, "output file name default[issue.txt]")
	cmd.Flags().BoolVarP(&o.file, "file", "f", o.file, "output the content to a file")

	return cmd
}

func (i *issueOption) Run() error {
	if len(i.name) > 0 {
		return i.uniqueOne()
	}

	u := "https://quickissue.openeuler.org/api-issues/issues/types"
	var res = struct {
		baseResp
		Data []struct {
			UniqueId int64  `json:"id"`
			Name     string `json:"name"`
		}
	}{}

	_, err := i.h.CustomRequest(u, "GET", nil, nil, nil, &res)
	if err != nil {
		return err
	}

	if res.Code != 200 {
		return fmt.Errorf(res.Msg)
	}

	err = i.printContextHeaders(i.Out)
	if err != nil {
		return err
	}
	var data = res.Data
	sort.Slice(data, func(i, j int) bool {
		return data[i].UniqueId < data[j].UniqueId
	})
	for _, v := range data {
		_, err = fmt.Fprintf(i.Out, "%-15d\t%s\n", v.UniqueId, v.Name)
	}

	return err
}

func (i *issueOption) uniqueOne() error {
	u := "https://quickissue.openeuler.org/api-issues/issues/types"
	var v = url.Values{}
	v.Add("name", i.name)
	var res = struct {
		baseResp
		Data []struct {
			UniqueId int64  `json:"id"`
			Name     string `json:"name"`
			Template string `json:"template"`
		}
	}{}

	_, err := i.h.CustomRequest(u, "GET", nil, nil, v, &res)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		return fmt.Errorf(res.Msg)
	}

	if len(res.Data) == 0 {
		return fmt.Errorf("not found issue type : %s", i.name)
	}

	if i.file {
		return i.writeFile(res.Data[0].Template)
	}

	_, err = fmt.Fprintln(i.Out, res.Data[0].Template)
	return err
}

func (i *issueOption) writeFile(content string) error {
	var file = fmt.Sprintf(basefile, "issue")

	if len(i.filename) > 0 {
		file = fmt.Sprintf(basefile, i.filename)
	}

	return ioutil.WriteFile(file, []byte(content), fs.ModePerm)
}

func (i *issueOption) printContextHeaders(out io.Writer) error {
	columnNames := []any{"UNIQUEID", "NAME"}
	_, err := fmt.Fprintf(out, "%-15s\t%s\n", columnNames...)
	return err
}

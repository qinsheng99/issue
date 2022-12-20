package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"

	"github.com/qinsheng99/issue/util"
	"github.com/spf13/cobra"
)

const basefile = "%s.txt"

type issueOption struct {
	Streams
	h util.ReqImpl

	uniqueId int
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

	cmd.Flags().IntVarP(&o.uniqueId, "unique", "u", o.uniqueId, "issue type id")
	cmd.Flags().StringVar(&o.filename, "filename", o.filename, "output file name default[issue.txt]")
	cmd.Flags().BoolVarP(&o.file, "file", "f", o.file, "output the content to a file")

	return cmd
}

func (i *issueOption) Run() error {
	if i.uniqueId != 0 {
		return i.uniqueOne()
	}

	u := "http://localhost:8000/v1/issue-type/list"
	var res = struct {
		baseResp
		Result []struct {
			UniqueId int64  `json:"uniqueId"`
			Name     string `json:"name"`
		}
	}{}

	_, err := i.h.CustomRequest(u, "GET", nil, nil, nil, &res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}

	err = i.printContextHeaders(i.Out)
	if err != nil {
		return err
	}
	var data = res.Result
	sort.Slice(data, func(i, j int) bool {
		return data[i].UniqueId < data[j].UniqueId
	})
	for _, v := range data {
		_, err = fmt.Fprintf(i.Out, "%d\t%s\n", v.UniqueId, v.Name)
	}

	return err
}

func (i *issueOption) uniqueOne() error {
	u := "http://localhost:8000/v1/issue-type/one"
	var v = url.Values{}
	v.Add("unique", strconv.Itoa(i.uniqueId))
	var res = struct {
		baseResp
		Result struct {
			UniqueId int64  `json:"uniqueId"`
			Name     string `json:"name"`
			Template string `json:"template"`
		}
	}{}

	_, err := i.h.CustomRequest(u, "GET", nil, nil, v, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}

	if i.file {
		return i.writeFile(res.Result.Template)
	}

	_, err = fmt.Fprintln(i.Out, res.Result.Template)
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
	_, err := fmt.Fprintf(out, "%s\t%s\n", columnNames...)
	return err
}

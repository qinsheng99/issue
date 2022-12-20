package cmd

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/qinsheng99/issue/util"
	"github.com/spf13/cobra"
)

type repoOption struct {
	Streams
	h util.ReqImpl

	page int
	size int

	name string

	accurate bool
}

func newRepoOption(s Streams, h util.ReqImpl) *repoOption {
	return &repoOption{Streams: s, h: h}
}

func newRepoCmd(s Streams, h util.ReqImpl) *cobra.Command {
	o := newRepoOption(s, h)

	cmd := &cobra.Command{
		Use:   "repo [options]",
		Short: "obtain information about the repository that openeuler can use to create an issue",
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(o.Validate())
			checkErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.name, "name", "n", "", "repo name")
	cmd.Flags().IntVarP(&o.page, "page", "p", 1, "get the number of pages for the warehouse")
	cmd.Flags().IntVarP(&o.size, "size", "s", 20, "get the number of sizes for the warehouse")

	cmd.Flags().BoolVarP(&o.accurate, "accurate", "a", o.accurate, "whether accurate search")

	return cmd
}

func (o *repoOption) Run() error {
	if o.accurate {
		return o.findAccurate()
	}

	u := "http://localhost:8000/v1/repo/repo-names"
	var v = url.Values{}
	v.Add("page", strconv.Itoa(o.page))
	v.Add("size", strconv.Itoa(o.size))
	v.Add("name", o.name)

	var res = struct {
		baseResp
		Result []struct {
			Name string `json:"fullRepoName"`
			Id   int64  `json:"repoId"`
		} `json:"result"`
	}{}

	_, err := o.h.CustomRequest(u, "GET", nil, nil, v, &res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}

	if err = o.printContextHeaders(o.Out); err != nil {
		return err
	}
	for _, s := range res.Result {
		_, err = fmt.Fprintf(o.Out, "%d\t%s\n", s.Id, s.Name)
	}
	return err
}

func (o *repoOption) findAccurate() error {
	u := "http://localhost:8000/v1/repo/repo"
	var v = url.Values{}
	v.Add("name", o.name)
	var res = struct {
		baseResp
		Result struct {
			Name string `json:"fullRepoName"`
			Id   int64  `json:"repoId"`
		} `json:"result"`
	}{}
	_, err := o.h.CustomRequest(u, "GET", nil, nil, v, &res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}
	if err = o.printContextHeaders(o.Out); err != nil {
		return err
	}
	_, err = fmt.Fprintf(o.Out, "%d\t%s\n", res.Result.Id, res.Result.Name)
	return err
}

func (o *repoOption) Validate() error {
	if o.accurate && len(o.name) <= 0 {
		return fmt.Errorf("name must be specified if exact lookup")
	}

	return nil
}

func (o *repoOption) printContextHeaders(out io.Writer) error {
	columnNames := []any{"REPOID", "REPONAME"}
	_, err := fmt.Fprintf(out, "%-10s\t%s\n", columnNames...)
	return err
}

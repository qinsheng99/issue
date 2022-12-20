package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/qinsheng99/issue/util"
	"github.com/spf13/cobra"
	"os"
)

type createOption struct {
	Streams
	h util.ReqImpl

	filepath string
	title    string
	repoid   int64

	email, code string
}

func newCreateOption(s base) *createOption {
	return &createOption{Streams: s.Streams, h: s.ReqImpl}
}

func newCmdCreate(s base) *cobra.Command {
	o := newCreateOption(s)

	cmd := &cobra.Command{
		Use:   "create [options]",
		Short: "create issue for openeuler",
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(o.Validate())
			checkErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.filepath, "file", "f", o.filepath, "issue body file path")
	cmd.Flags().StringVarP(&o.title, "title", "t", o.filepath, "issue title")
	cmd.Flags().Int64VarP(&o.repoid, "repoid", "i", o.repoid, "create an issue in that repository")

	return cmd
}

func (c *createOption) Validate() error {
	if len(c.filepath) <= 0 {
		return fmt.Errorf("please enter file path")
	}

	if len(c.title) <= 0 {
		return fmt.Errorf("please enter the issue title")
	}

	if c.repoid <= 0 {
		return fmt.Errorf("please specify the repo id,you can use `issue get repo`")
	}

	return nil
}

func (c *createOption) Run() error {
	var email string
	fmt.Println("请输入邮箱:")
	_, err := fmt.Fscan(c.In, &email)
	if err != nil {
		return err
	}

	u := "http://localhost:8000/v1/send/%s"
	var res = struct {
		baseResp
	}{}

	_, err = c.h.CustomRequest(fmt.Sprintf(u, email), "GET", nil, nil, nil, &res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}
	c.email = email

	fmt.Println("验证码已经发送至邮箱,请输入验证码:")
	var code string
	_, err = fmt.Fscan(c.In, &code)
	if err != nil {
		return err
	}
	c.code = code
	return c.createIssue()
}

func (c *createOption) createIssue() error {
	bys, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}
	var req = struct {
		Id    int64  `json:"id"`
		Title string `json:"title"`
		Body  string `json:"body"`
		Email string `json:"email"`
		Code  string `json:"code"`
	}{Id: c.repoid, Title: c.title, Body: string(bys), Email: c.email, Code: c.code}

	url := "http://localhost:8000/v1/create-issue"

	bys, err = json.Marshal(req)
	if err != nil {
		return err
	}

	var res = struct {
		baseResp
	}{}

	_, err = c.h.CustomRequest(url, "POST", bys, nil, nil, &res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}

	_, err = fmt.Fprintln(c.Out, "success")
	return err
}

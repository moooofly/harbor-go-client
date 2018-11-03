package api

import (
	"fmt"
	"os"
	"strconv"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("logs",
		"Get recent logs of the projects which the user is a member of.",
		"This endpoint let user see the recent operation logs of the projects which he is member of.",
		&logs)
}

type recentLogs struct {
	Username       string `short:"u" long:"username" description:"Username of the operator."`
	Repository     string `short:"r" long:"repository" description:"The name of repository."`
	Tag            string `short:"t" long:"tag" description:"The name of tag."`
	Operation      string `short:"o" long:"operation" description:"The operation. ([create|delete|push|pull])"`
	BeginTimestamp string `short:"b" long:"begin_timestamp" description:"The begin timestamp. (format: yyyymmdd)"`
	EndTimestamp   string `short:"e" long:"end_timestamp" description:"The end timestamp. (format: yyyymmdd)"`
	Page           int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize       int    `short:"s" long:"page_size" description:"The size of per page, default is 10, maximum is 100." default:"10"`
}

var logs recentLogs

func (x *recentLogs) Execute(args []string) error {
	GetOPLogs(utils.URLGen("/api/logs"))
	return nil
}

// GetOPLogs ...
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/logs?username=admin&repository=prj2%2Fphoton&tag=v3&operation=push&begin_timestamp=20171102&page=1&page_size=10'
func GetOPLogs(baseURL string) {
	if logs.Operation != "" &&
		logs.Operation != "create" &&
		logs.Operation != "delete" &&
		logs.Operation != "push" &&
		logs.Operation != "pull" {
		fmt.Println("error: operation must be one of [create|delete|push|pull]")
		os.Exit(1)
	}

	targetURL := baseURL + "?username=" + logs.Username +
		"&repository=" + logs.Repository +
		"&tag=" + logs.Tag +
		"&operation=" + logs.Operation +
		"&begin_timestamp=" + logs.BeginTimestamp +
		"&end_timestamp=" + logs.EndTimestamp +
		"&page=" + strconv.Itoa(logs.Page) +
		"&page_size=" + strconv.Itoa(logs.PageSize)

	fmt.Println("==> GET", targetURL)

	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

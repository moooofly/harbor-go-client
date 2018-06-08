package api

import (
	"fmt"

	"git.llsapp.com/fei.sun/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("search",
		"Search for projects and repositories.",
		"The Search endpoint returns information about the projects and repositories offered at public status or related to the current logged in user. The response includes the project and repository list in a proper display order.",
		&searching)
}

type search struct {
	Q string `short:"q" long:"query" description:"(REQUIRED) Search parameter for project and repository name." required:"yes"`
}

var searching search

func (x *search) Execute(args []string) error {
	SearchPrjAndRepo(utils.URLGen("/api/search"))
	return nil
}

// SearchPrjAndRepo returns information about the projects and repositories offered at public status or related to the current logged in user. The response includes the project and repository list in a proper display order.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/search?q=hello-world'
func SearchPrjAndRepo(baseURL string) {
	targetURL := baseURL + "?q=" + searching.Q
	fmt.Println("==> GET", targetURL)

	// NOTE:
	// 实验表明该 API 在没有 cookie 的情况下也可以使用
	// 文档中 "offered at public status or related to the current logged in user" 覆盖到了这层含义
	// 此处将 cookie 的使用设置成必须，可以酌情调整
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

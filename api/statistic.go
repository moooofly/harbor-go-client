package api

import (
	"fmt"

	"git.llsapp.com/fei.sun/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("statistics",
		"Get projects number and repositories number relevant to the user.",
		"This endpoint is aimed to statistic all of the projects number and repositories number relevant to the logined user, also the public projects number and repositories number. If the user is admin, he can also get total projects number and total repositories number.",
		&stats)
}

type statistics struct {
}

var stats statistics

func (x *statistics) Execute(args []string) error {
	GetStats(utils.URLGen("/api/statistics"))
	return nil
}

// GetStats is aimed to statistic all of the projects number and repositories number relevant to the logined user, also the public projects number and repositories number. If the user is admin, he can also get total projects number and total repositories number.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/statistics'
func GetStats(baseURL string) {
	targetURL := baseURL
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

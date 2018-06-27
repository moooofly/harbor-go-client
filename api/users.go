package api

import (
	"fmt"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	// NOTE:
	// 由于 users_list 命令是是用于列出当前 login 用户相关信息
	// 故将其改名为 whoami
	utils.Parser.AddCommand("whoami",
		"Show info about current login user only.", "Maybe 'whoami' is a better name.", &userslist)
}

type usersList struct {
	Apikey string `short:"k" long:"api_key" description:"have no idea about this option usage." default:"top"`
}

var userslist usersList

func (x *usersList) Execute(args []string) error {
	GetUsers(utils.URLGen("/api/users"))
	return nil
}

// GetUsers gets the current user information.
//
// GET /users/current
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/users/current?api_key=top'
func GetUsers(baseURL string) {
	targetURL := baseURL + "/current" + "?api_key=" + userslist.Apikey
	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		// NOTE:
		// 若后续需要根据用户权限做文章，则需要将用户信息进行维护
		// 可以定制一个新的回调函数
		End(utils.PrintStatus)
}

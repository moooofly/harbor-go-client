package api

import (
	"fmt"
	"net/url"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("login",
		"Log in to Harbor.", "Log in to Harbor with username and password.", &li)
	utils.Parser.AddCommand("logout",
		"Log out from Harbor.", "Log out current user from Harbor.", &lo)
}

type login struct {
	Username string `short:"u" long:"username" description:"(REQUIRED) Current login username." required:"yes"`
	Password string `short:"p" long:"password" description:"Current login password." default:""`
	// FIXME:
	// 需要设计一种可以覆盖 config.yaml 配置文件中 dstip 的方式
	//Address  string `short:"a" long:"address" description:"The specified ip address of the harbor service." default:""`
}

var li login

func (x *login) Execute(args []string) error {
	LoginHarbor(utils.URLGen("/login"))
	return nil
}

type logout struct {
}

var lo logout

func (x *logout) Execute(args []string) error {
	LogoutHarbor(utils.URLGen("/log_out"))
	return nil
}

// LoginHarbor log in to Harbor.
//
// params:
// 	username - Current login username.
//  password - Current login password.
//
// e.g. curl -X POST --header 'Content-Type: application/x-www-form-urlencoded;param=value' 'https://11.11.11.12/login' -i -k -d "principal=admin&password=Harbor12345"
func LoginHarbor(baseURL string) {

	if li.Password == "" {
		// 支持密码隐藏功能
		passwd, err := utils.ReadPasswordFromTerm()
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		if passwd == "" {
			fmt.Println("error: Password Required.")
			return
		}

		li.Password = passwd
	} else {
		fmt.Println("WARNING! Using --password via the CLI is insecure.")
	}

	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	//fmt.Printf("==> username: %s   password: %s   escape: %s\n", li.Username, li.Password, url.QueryEscape(li.Password))

	utils.Request.Post(targetURL).
		Set("Content-Type", "application/x-www-form-urlencoded;param=value").
		// NOTE:
		// After some experiments, conclude that the value of Cookie has two forms:
		// 1. Cookie:rem-username=admin; harbor-lang=zh-cn; beegosessionID=720210**3d76d6
		// 2. Cookie:rem-username=admin; harbor-lang=zh-cn;
		//
		// The fist one reuses the value of beegosessionID in Set-Cookie from response headers.
		// The second one is equivalent to a fresh login.
		//
		// Taking the second form just for long-live coding.
		Set("Cookie", "harbor-lang=zh-cn").
		Send("principal=" + li.Username + "&password=" + url.QueryEscape(li.Password)).
		End(utils.LoginProc)
}

// LogoutHarbor log out from Harbor.
//
// params:
//
// e.g. curl -X GET 'https://11.11.11.12/log_out' -i -k
func LogoutHarbor(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> GET", targetURL)

	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.LogoutProc)
}

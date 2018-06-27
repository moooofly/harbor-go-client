package api

import (
	"fmt"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("sysinfo_general",
		"Get general system info.",
		"This API is for retrieving general system info, this can be called by anonymous request.",
		&sysGeneral)
	utils.Parser.AddCommand("sysinfo_volumes",
		"Get system volume info (total/free size).",
		"This endpoint is for retrieving system volume info that only provides for admin user.",
		&sysVolumes)
	utils.Parser.AddCommand("sysinfo_rootcert",
		"Get default root certificate under OVA deployment.",
		"This endpoint is for downloading a default root certificate that only provides for admin user under OVA deployment.",
		&sysRootCert)
}

type sysInfoGeneral struct {
}

var sysGeneral sysInfoGeneral

func (x *sysInfoGeneral) Execute(args []string) error {
	GetSysGeneral(utils.URLGen("/api/systeminfo"))
	return nil
}

type sysInfoVolumes struct {
}

var sysVolumes sysInfoVolumes

func (x *sysInfoVolumes) Execute(args []string) error {
	GetSysVolumes(utils.URLGen("/api/systeminfo/volumes"))
	return nil
}

type sysInfoRootCert struct {
}

var sysRootCert sysInfoRootCert

func (x *sysInfoRootCert) Execute(args []string) error {
	GetSysRootCert(utils.URLGen("/api/systeminfo/getcert"))
	return nil
}

// GetSysGeneral is for retrieving general system info, this can be called by anonymous request.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/systeminfo'
func GetSysGeneral(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> GET", targetURL)

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn").
		End(utils.PrintStatus)
}

// GetSysVolumes is for retrieving system volume info that only provides for admin user.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/systeminfo/volumes'
func GetSysVolumes(baseURL string) {
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

// GetSysRootCert is for downloading a default root certificate that only provides for admin user under OVA deployment.
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/systeminfo/getcert'
func GetSysRootCert(baseURL string) {
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

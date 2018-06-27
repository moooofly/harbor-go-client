package api

import (
	"encoding/json"
	"fmt"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("configurations_get",
		"Get system configurations.",
		"This endpoint is for retrieving system configurations that only provides for admin user.",
		&scGet)
	utils.Parser.AddCommand("configurations_create",
		"Modify system configurations. (set configuration in conf/config.yaml)",
		"This endpoint is for modifying system configurations that only provides for admin user.",
		&scCreate)
	utils.Parser.AddCommand("configurations_reset",
		"Reset system configurations.",
		"Reset system configurations from environment variables. Can only be accessed by admin user.",
		&scReset)
}

type sysConfigGet struct {
}

var scGet sysConfigGet

func (x *sysConfigGet) Execute(args []string) error {
	GetSysConfig(utils.URLGen("/api/configurations"))
	return nil
}

type sysConfigCreate struct {
}

var scCreate sysConfigCreate

func (x *sysConfigCreate) Execute(args []string) error {
	PutSysConfigCreate(utils.URLGen("/api/configurations"))
	return nil
}

type sysConfigReset struct {
}

var scReset sysConfigReset

func (x *sysConfigReset) Execute(args []string) error {
	PostSysConfigReset(utils.URLGen("/api/configurations/reset"))
	return nil
}

// GetSysConfig is for retrieving system configurations that only provides for admin user.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/configurations'
func GetSysConfig(baseURL string) {
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

// PutSysConfigCreate is for modifying system configurations that only provides for admin user.
//
/* e.g.
  curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{
  "auth_mode": "db_auth",
  "email_from": "admin <sample_admin@mydomain.com>",
  "email_host": "smtp.mydomain.com",
  "email_port": 1200,
  "email_identity": "",
  "email_username": "sample_admin@mydomain.com",
  "email_ssl": false,
  "email_insecure": true,
  "ldap_url": "ldaps://ldap.mydomain.com",
  "ldap_base_dn": "ou=people,dc=mydomain,dc=com",
  "ldap_filter": "",
  "ldap_scope": 3,
  "ldap_uid": "uid",
  "ldap_search_dn": "",
  "ldap_timeout": 5,
  "project_creation_restriction": "everyone",
  "self_registration": true,
  "token_expiration": 30,
  "verify_remote_cert": true,
  "scan_all_policy": {
    "type": "daily",
    "parameter": {
      "daily_time": 0
    }
  }
}' 'https://localhost/api/configurations'
*/
func PutSysConfigCreate(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> PUT", targetURL)

	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	sc, err := utils.SysConfigLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	msc, err := json.Marshal(sc)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(msc)).
		End(utils.PrintStatus)
}

// PostSysConfigReset resets system configurations from environment variables. Can only be accessed by admin user.
//
// e.g. curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' 'https://localhost/api/configurations/reset'
func PostSysConfigReset(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

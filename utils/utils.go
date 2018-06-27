package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/parnurzeal/gorequest"
	yaml "gopkg.in/yaml.v2"
)

// These variables are populated via the Go linker.
var (
	UTCBuildTime  = "unknown"
	ClientVersion = "unknown"
	GoVersion     = "unknown"
	GitBranch     = "unknown"
	GitTag        = "unknown"
	GitHash       = "unknown"
)

var errMalCookies = errors.New("get malformed cookies")
var errCookiesNotAvailable = errors.New("target cookies are not available")

// Parser is a command registry
var Parser = flags.NewParser(nil, flags.Default)

// Request is a new SuperAgent object with a setting of not verifying
// server's certificate chain and host name.
var Request = gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true})

var configfile = "conf/config.yaml"
var secretfile = "conf/.cookie.yaml"

// Beegocookie is for beegosessionID storage
type Beegocookie struct {
	BeegosessionID string `yaml:"beegosessionID"`
}

type generalConfig struct {
	Scheme string `yaml:"scheme"`
	Dstip  string `yaml:"dstip"`
}

// SysConfig defines system configurations
type SysConfig struct {
	AuthMode                   string `yaml:"auth_mode" json:"auth_mode"`
	EmailFrom                  string `yaml:"email_from" json:"email_from"`
	EmailHost                  string `yaml:"email_host" json:"email_host"`
	EmailPort                  int    `yaml:"email_port" json:"email_port"`
	EmailIdentity              string `yaml:"email_identity" json:"email_identity"`
	EmailUsername              string `yaml:"email_username" json:"email_username"`
	EmailSsl                   bool   `yaml:"email_ssl" json:"email_ssl"`
	EmailInsecure              bool   `yaml:"email_insecure" json:"email_insecure"`
	LdapURL                    string `yaml:"ldap_url" json:"ldap_url"`
	LdapBaseDN                 string `yaml:"ldap_base_dn" json:"ldap_base_dn"`
	LdapFilter                 string `yaml:"ldap_filter" json:"ldap_filter"`
	LdapScope                  int    `yaml:"ldap_scope" json:"ldap_scope"`
	LdapUID                    string `yaml:"ldap_uid" jsonb:"ldap_uid"`
	LdapSearchDN               string `yaml:"ldap_search_dn" json:"ldap_search_dn"`
	LdapTimeout                int    `yaml:"ldap_timeout" json:"ldap_timeout"`
	ProjectCreationRestriction string `yaml:"project_creation_restriction" json:"project_creation_restriction"`
	SelfRegistration           bool   `yaml:"self_registration" json:"self_registration"`
	TokenExpiration            int    `yaml:"token_expiration" json:"token_expiration"`
	VerifyRemoteCert           bool   `yaml:"verify_remote_cert" json:"verify_remote_cert"`
	ScanAllPolicy              struct {
		Type      string `yaml:"type" json:"type"`
		Parameter struct {
			DailyTime int `yaml:"daily_time" json:"daily_time"`
		} `yaml:"parameter" json:"parameter"`
	} `yaml:"scan_all_policy" json:"scan_all_policy"`
}

// cookieFilter filters specific cookie string.
func cookieFilter(cookies []*http.Cookie, filter string) (string, error) {

	for _, cookie := range cookies {
		parts := strings.Split(strings.TrimSpace(cookie.String()), ";")

		if len(parts) == 1 && parts[0] == "" {
			return "", errMalCookies
		}

		for _, part := range parts {
			part = strings.TrimSpace(part)
			j := strings.Index(part, "=")
			if j < 0 {
				if part == filter {
					return "", errMalCookies
				}
				fmt.Println("name=", part)
				continue
			}
			name, value := part[:j], part[j+1:]
			if 0 == strings.Compare(name, filter) {
				return value, nil
			}
		}
	}
	return "", errCookiesNotAvailable
}

// cookieSave saves beegosessionID into .cookie.yaml
//
// This function is called only in stage of login, and will reset the content of
// .cookie.yaml no matter whether it exists or not.
func cookieSave(beegosessionID string) error {

	var cookie Beegocookie
	cookie.BeegosessionID = beegosessionID

	c, err := yaml.Marshal(&cookie)
	if err != nil {
		return err
	}
	//fmt.Printf("--- c dump:\n%s\n\n", string(c))

	if err = ioutil.WriteFile(secretfile, []byte(c), 0644); err != nil {
		return err
	}

	return nil
}

// CookieLoad loads beegosessionID from .cookie.yaml.
func CookieLoad() (*Beegocookie, error) {
	var cookie Beegocookie

	dataBytes, err := ioutil.ReadFile(secretfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(dataBytes), &cookie)
	if err != nil {
		return nil, err
	}

	return &cookie, nil
}

// SysConfigLoad loads system configuration from conf/config.yaml.
func SysConfigLoad() (*SysConfig, error) {
	var config SysConfig

	dataBytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(dataBytes), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// generalConfigLoad load general configuration from conf/config.yaml
func generalConfigLoad() (*generalConfig, error) {
	var config generalConfig

	dataBytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(dataBytes), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// URLGen generates target URL.
func URLGen(uri string) string {
	config, err := generalConfigLoad()
	if err != nil {
		fmt.Println("URLGen:", err)
		os.Exit(1)
	}
	url := config.Scheme + "://" + config.Dstip + uri

	return url
}

// LoginProc is the callback function for login.
func LoginProc(resp gorequest.Response, body string, errs []error) {
	for _, e := range errs {
		if e != nil {
			fmt.Println(e)
			return
		}
	}

	cookies := (*http.Response)(resp).Cookies()
	fmt.Println("<== Cookies:", cookies)

	sid, err := cookieFilter(cookies, "beegosessionID")
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: 根据状态码进行 .cookie.yaml 文件处理，以及用户友好提示
	fmt.Println("<== Rsp Status:", resp.Status)
	fmt.Println("<== Rsp Body:", body)

	err = cookieSave(sid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

// LogoutProc is the callback function for logout.
func LogoutProc(resp gorequest.Response, body string, errs []error) {
	for _, e := range errs {
		if e != nil {
			fmt.Println(e)
			return
		}
	}

	fmt.Println("<== Rsp Cookies:", (*http.Response)(resp).Cookies())
	fmt.Println("<== Rsp Status:", resp.Status)
	fmt.Println("<== Rsp Body:", body)

	os.Remove(secretfile)
}

// PrintStatus is a regular callback function.
func PrintStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println("<== ")
	for _, e := range errs {
		if e != nil {
			fmt.Println(e)
			return
		}
	}

	fmt.Println("<== Rsp Status:", resp.Status)
	fmt.Printf("<== Rsp Body: %s\n", body)
}

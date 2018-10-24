package api

import (
	"encoding/json"
	"fmt"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("syncregistry",
		"Sync repositories from registry to DB.",
		"This endpoint is for syncing all repositories of registry with database.",
		&syncregistry)
	utils.Parser.AddCommand("email_ping",
		"Test connection and authentication with email server.",
		"Test connection and authentication with email server.",
		&emailping)
}

type syncRegistry struct {
}

var syncregistry syncRegistry

func (x *syncRegistry) Execute(args []string) error {
	PostSyncRegistry(utils.URLGen("/api/internal/syncregistry"))
	return nil
}

// PostSyncRegistry is for syncing all repositories of registry with database.
//
// params:
//
// format:
//   POST /internal/syncregistry
//
// e.g. curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' 'https://localhost/api/internal/syncregistry'
func PostSyncRegistry(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

type emailPing struct {
	EmailHost     string `short:"h" long:"email_host" description:"The host of email server." default:"smtp.mydomain.com" json:"email_host"`
	EmailPort     int    `short:"t" long:"email_port" description:"The port of email server." default:"25" json:"email_port"`
	EmailUsername string `short:"u" long:"email_username" description:"The username of email server." default:"sample_admin@mydomain.com" json:"email_username"`
	EmailPassword string `short:"p" long:"email_password" description:"The password of email server." default:"" json:"email_password"`
	EmailSsl      bool   `short:"s" long:"email_ssl" description:"Use ssl/tls or not." json:"email_ssl"`
	EmailIdentity string `short:"i" long:"email_identity" description:"The identity of email server." default:"" json:"email_identity"`
}

var emailping emailPing

func (x *emailPing) Execute(args []string) error {
	PostEmailPing(utils.URLGen("/api/email/ping"))
	return nil
}

// PostEmailPing tests connection and authentication with email server.
//
// params:
//  email_host     - The host of email server.
//  email_port     - The port of email server.
//  email_username - The username of email server.
//  email_password - The password of email server.
//  email_ssl      - Use ssl/tls or not.
//  email_identity - The dentity of email server.
//
// format:
//   POST /email/ping
//
// e.g.
/*
curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "email_host": "string", \
   "email_port": 0, \
   "email_username": "string", \
   "email_password": "string", \
   "email_ssl": true, \
   "email_identity": "string" \
 }' 'https://localhost/api/email/ping'
)*/
func PostEmailPing(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&emailping)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> email ping:", string(t))

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

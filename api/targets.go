package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("targets_list",
		"List targets filtered by name.",
		"This endpoint let user list targets filtered by name, if name is nil, list returns all targets.",
		&tl)
	utils.Parser.AddCommand("targets_create",
		"Create a new replication target.",
		"This endpoint is for user to create a new replication target.",
		&tc)
	utils.Parser.AddCommand("targets_ping",
		"Ping validates target.",
		"This endpoint is for ping validates whether the target is reachable and whether the credential is valid.",
		&tping)
	utils.Parser.AddCommand("targets_ping_by_tid",
		"Ping target.",
		"This endpoint is for ping target.",
		&tpingByID)
	utils.Parser.AddCommand("targets_delete_by_tid",
		"Delete specific replication's target.",
		"This endpoint is for to delete specific replication's target.",
		&tdByID)
	utils.Parser.AddCommand("targets_get_by_tid",
		"Get replication's target.",
		"This endpoint is for get specific replication's target.",
		&tgByID)
	utils.Parser.AddCommand("targets_update_by_tid",
		"Update replication's target.",
		"This endpoint is for update specific replication's target.",
		&tuByID)
	utils.Parser.AddCommand("targets_policies_by_tid",
		"List the target relevant policies.",
		"This endpoint list policies filter with specific replication's target ID.",
		&tpoliciesByID)
}

type targetsList struct {
	Name string `short:"n" long:"name" description:"The replication's target name (for filter)." default:""`
}

var tl targetsList

func (x *targetsList) Execute(args []string) error {
	GetTargetsList(utils.URLGen("/api/targets"))
	return nil
}

type targetsCreate struct {
	EndpointURL  string `short:"e" long:"endpoint" description:"(REQUIRED) The target address URL string. (Should be globally unique)" required:"yes" json:"endpoint"`
	EndpointName string `short:"n" long:"name" description:"(REQUIRED) The target name. (Should be globally unique)" required:"yes" json:"name"`
	Username     string `short:"u" long:"username" description:"(REQUIRED) The target server username." required:"yes" json:"username"`
	Password     string `short:"p" long:"password" description:"(REQUIRED) The target server password." required:"yes" json:"password"`
	Insecure     bool   `short:"x" long:"insecure" description:"(REQUIRED) Whether or not the certificate will be verified when Harbor tries to access the server." required:"yes" json:"insecure"`
}

var tc targetsCreate

func (x *targetsCreate) Execute(args []string) error {
	PostTargetsCreate(utils.URLGen("/api/targets"))
	return nil
}

type targetsPing struct {
	EndpointURL string `short:"e" long:"endpoint" description:"(REQUIRED) The target address URL string." required:"yes" json:"endpoint"`
	Username    string `short:"u" long:"username" description:"(REQUIRED) The target server username." required:"yes" json:"username"`
	Password    string `short:"p" long:"password" description:"(REQUIRED) The target server password." required:"yes" json:"password"`
	Insecure    bool   `short:"x" long:"insecure" description:"(REQUIRED) Whether or not the certificate will be verified when Harbor tries to access the server." required:"yes" json:"insecure"`
}

var tping targetsPing

func (x *targetsPing) Execute(args []string) error {
	PostTargetsPing(utils.URLGen("/api/targets/ping"))
	return nil
}

type targetsPingByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) The replication's target ID." required:"yes"`
}

var tpingByID targetsPingByID

func (x *targetsPingByID) Execute(args []string) error {
	PostTargetsPingByID(utils.URLGen("/api/targets"))
	return nil
}

type targetsDeleteByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) The replication's target ID." required:"yes"`
}

var tdByID targetsDeleteByID

func (x *targetsDeleteByID) Execute(args []string) error {
	DeleteTargetsByID(utils.URLGen("/api/targets"))
	return nil
}

type targetsGetByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) The replication's target ID." required:"yes"`
}

var tgByID targetsGetByID

func (x *targetsGetByID) Execute(args []string) error {
	GetTargetsByID(utils.URLGen("/api/targets"))
	return nil
}

type targetsUpdateByID struct {
	ID           int    `short:"i" long:"id" description:"(REQUIRED) The replication's target ID." required:"yes" json:"-"`
	EndpointURL  string `short:"e" long:"endpoint" description:"(REQUIRED) The target address URL string." required:"yes" json:"endpoint"`
	EndpointName string `short:"n" long:"name" description:"(REQUIRED) The target name." required:"yes" json:"name"`
	Username     string `short:"u" long:"username" description:"(REQUIRED) The target server username." required:"yes" json:"username"`
	Password     string `short:"p" long:"password" description:"(REQUIRED) The target server password." required:"yes" json:"password"`
	Insecure     bool   `short:"x" long:"insecure" description:"(REQUIRED) Whether or not the certificate will be verified when Harbor tries to access the server." required:"yes" json:"insecure"`
}

var tuByID targetsUpdateByID

func (x *targetsUpdateByID) Execute(args []string) error {
	UpdateTargetsByID(utils.URLGen("/api/targets"))
	return nil
}

type targetsPoliciesByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) The replication's target ID." required:"yes"`
}

var tpoliciesByID targetsPoliciesByID

func (x *targetsPoliciesByID) Execute(args []string) error {
	GetPoliciesByID(utils.URLGen("/api/targets"))
	return nil
}

// GetTargetsList let user list filters targets by name, if name is nil, list returns all targets.
//
// params:
//  name - The replication's target name (for filter).
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/targets?name=remote'
func GetTargetsList(baseURL string) {
	targetURL := baseURL + "?name=" + tl.Name
	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

// PostTargetsCreate is for user to create a new replication target.
//
// params:
//  endpoint - The target address URL string
//  name - The target name.
//  username - The target server username.
//  password - The target server password.
//  insecure - Whether or not the certificate will be verified when Harbor tries to access the server.
//
// e.g.
/*
curl -X POST --header 'Content-Type: text/plain' --header 'Accept: text/plain' -d '{
  "endpoint": "https://11.11.11.14",
  "name": "rule4",
  "username": "admin",
  "password": "Harbor12345",
  "insecure": true
}' 'https://11.11.11.12/api/targets'
*/
func PostTargetsCreate(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&tc)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

// PostTargetsPing is for ping validates whether the target is reachable and whether the credential is valid.
//
// params:
//  endpoint - The target address URL string
//  username - The target server username.
//  password - The target server password.
//  insecure - Whether or not the certificate will be verified when Harbor tries to access the server.
//
// e.g.
/*
curl -X POST --header 'Content-Type: text/plain' --header 'Accept: text/plain' -d '{
  "endpoint": "https://11.11.11.11",
  "username": "admin",
  "password": "Harbor12345",
  "insecure": true
}' 'https://11.11.11.12/api/targets/ping'
*/
func PostTargetsPing(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	p, err := json.Marshal(&tping)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(p)).
		End(utils.PrintStatus)
}

// PostTargetsPingByID is for ping target.
//
// params:
//  id - (REQUIRED) The replication's target ID.
//
// e.g. curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' 'https://11.11.11.12/api/targets/1/ping'
func PostTargetsPingByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(tpingByID.ID) + "/ping"
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

// DeleteTargetsByID is for to delete specific replication's target.
//
// params:
//  id - (REQUIRED) The replication's target ID.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://11.11.11.12/api/targets/2'
func DeleteTargetsByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(tdByID.ID)
	fmt.Println("==> DELETE", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Delete(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

// GetTargetsByID is for get specific replication's target.
//
// params:
//  id - (REQUIRED) The replication's target ID.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/targets/1'
func GetTargetsByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(tgByID.ID)
	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

// UpdateTargetsByID is for update specific replication's target.
//
// params:
//  id - (REQUIRED) The replication's target ID.
//  endpoint - The target address URL string
//  name - The target name.
//  username - The target server username.
//  password - The target server password.
//  insecure - Whether or not the certificate will be verified when Harbor tries to access the server.
//
// e.g.
/*
curl -X PUT --header 'Content-Type: text/plain' --header 'Accept: text/plain' -d '{
  "endpoint": "https://11.11.11.14",
  "name": "change from rule4 to RULE4",
  "username": "admin",
  "password": "12345",
  "insecure": true
}' 'https://11.11.11.12/api/targets/4'
*/
func UpdateTargetsByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(tuByID.ID)
	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&tuByID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	//fmt.Println("===>", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

// GetPoliciesByID lists policies filter with specific replication's target ID.
//
// params:
//  id - (REQUIRED) The replication's target ID.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/targets/1/policies/'
func GetPoliciesByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(tpoliciesByID.ID) + "/policies/"
	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

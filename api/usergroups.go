package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("usergroups_list",
		"Get all user groups information",
		"Get all user groups information",
		&ugList)
	utils.Parser.AddCommand("usergroup_create",
		"Create user group",
		"Create user group information",
		&ugCreate)
	utils.Parser.AddCommand("usergroup_del",
		"Delete user group",
		"Delete user group",
		&ugDel)
	utils.Parser.AddCommand("usergroup_get",
		"Get user group information",
		"Get user group information",
		&ugGet)
	utils.Parser.AddCommand("usergroup_update",
		"Update group information",
		"Update group information",
		&ugUpdate)
}

type usergroupsList struct {
}

var ugList usergroupsList

func (x *usergroupsList) Execute(args []string) error {
	GetUsergroupsList(utils.URLGen("/api/usergroups"))
	return nil
}

// GetUsergroupsList get all user groups information
//
// params:
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/usergroups'
func GetUsergroupsList(baseURL string) {
	targetURL := baseURL
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

type usergroupCreate struct {
	ID          int    `short:"i" long:"id" description:"The ID of the user group" default:"0" json:"id"`
	GroupName   string `short:"n" long:"group_name" description:"The name of the user group" default:"tmp-group" json:"group_name"`
	GroupType   int    `short:"t" long:"group_type" description:"The group type, 1 for LDAP group." default:"1" json:"group_type"`
	LDAPGroupDN string `short:"l" long:"ldap_group_dn" description:"The DN of the LDAP group if group type is 1 (LDAP group)." default:"" json:"ldap_group_dn"`
}

var ugCreate usergroupCreate

func (x *usergroupCreate) Execute(args []string) error {
	PostUsergroupCreate(utils.URLGen("/api/usergroups"))
	return nil
}

// PostUsergroupCreate create user group information
//
// params:
//  id            - The ID of the user group
//  group_name    - The name of the user group
//  group_type    - The group type, 1 for LDAP group
//  ldap_group_dn - The DN of the LDAP group if group type is 1 (LDAP group)
//
// e.g.
/*
curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "id": 100, \
   "group_name": "tmp-group", \
   "group_type": 1, \
   "ldap_group_dn": "" \
 }' 'https://localhost/api/usergroups'
*/
func PostUsergroupCreate(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&ugCreate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("===> usergroup create:", string(t))

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type usergroupDel struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) The ID of the user group" required:"yes"`
}

var ugDel usergroupDel

func (x *usergroupDel) Execute(args []string) error {
	DeleteUsergroup(utils.URLGen("/api/usergroups"))
	return nil
}

// DeleteUsergroup delete user group
//
// params:
//  id - (REQUIRED) The ID of the user group
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/usergroups/1'
func DeleteUsergroup(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(ugDel.ID)
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

type usergroupGet struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) The ID of the user group" required:"yes"`
}

var ugGet usergroupGet

func (x *usergroupGet) Execute(args []string) error {
	GetUsergroup(utils.URLGen("/api/usergroups"))
	return nil
}

// GetUsergroup get user group information
//
// params:
//  id - (REQUIRED) The ID of the user group
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/usergroups/1'
func GetUsergroup(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(ugGet.ID)
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

type usergroupUpdate struct {
	ID          int    `short:"i" long:"id" description:"The ID of the user group" default:"0" json:"id"`
	GroupName   string `short:"n" long:"group_name" description:"The name of the user group" default:"tmp-group" json:"group_name"`
	GroupType   int    `short:"t" long:"group_type" description:"The group type, 1 for LDAP group." default:"1" json:"group_type"`
	LDAPGroupDN string `short:"l" long:"ldap_group_dn" description:"The DN of the LDAP group if group type is 1 (LDAP group)." default:"" json:"ldap_group_dn"`
}

var ugUpdate usergroupUpdate

func (x *usergroupUpdate) Execute(args []string) error {
	PutUsergroup(utils.URLGen("/api/usergroups"))
	return nil
}

// PutUsergroup update user group information
//
// params:
//  id            - The ID of the user group
//  group_name    - The name of the user group
//  group_type    - The group type, 1 for LDAP group
//  ldap_group_dn - The DN of the LDAP group if group type is 1 (LDAP group)
//
// e.g.
/*
curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "id": 1, \
   "group_name": "tmp-group", \
   "group_type": 1, \
   "ldap_group_dn": "" \
 }' 'https://localhost/api/usergroups/1'
*/
func PutUsergroup(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(ugUpdate.ID)
	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&ugUpdate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("===> usergroup update:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

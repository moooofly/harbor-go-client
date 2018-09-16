package api

import (
	"fmt"
	"strconv"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("policy_update_by_id",
		"Modify name, description, target and enablement of a policy.",
		"This endpoint let user update policy's name, description, target and enablement.",
		&poUpdateByID)
	utils.Parser.AddCommand("policy_get_by_id",
		"Get a policy.",
		"This endpoint let user search a policy by specific ID.",
		&poGetByID)
	utils.Parser.AddCommand("policy_create",
		"Create a policy.",
		"This endpoint let user creates a policy, and if it is enabled, the replication will be triggered right now.",
		&poCreate)
	utils.Parser.AddCommand("policies_list",
		"Filter policies by name and project_id.",
		"This endpoint let user filter policies by name and project_id, if name and project_id are nil, list returns all policies.",
		&poList)
}

type policyUpdateByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) policy ID" required:"yes"`
	// add more
}

var poUpdateByID policyUpdateByID

func (x *policyUpdateByID) Execute(args []string) error {
	PutPolicyUpdateByID(utils.URLGen("/api/policies"))
	return nil
}

// PutPolicyUpdateByID let user update policy name, description, target and enablement.
//
//  params:
//
//  format:
//
// e.g.
func PutPolicyUpdateByID(baseURL string) {}

type policyGetByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) policy ID" required:"yes"`
}

var poGetByID policyGetByID

func (x *policyGetByID) Execute(args []string) error {
	GetPolicyByID(utils.URLGen("/api/policies"))
	return nil
}

// GetPolicyByID let user search replication policy by specific ID.
//
// params:
//   id - (REQUIRED) policy ID
//
// format:
//   GET /policies/replication/{id}
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/policies/replication/1'
func GetPolicyByID(baseURL string) {
	targetURL := baseURL + "/replication/" + strconv.Itoa(poGetByID.ID)
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

type policyCreate struct {
}

var poCreate policyCreate

func (x *policyCreate) Execute(args []string) error {
	PostPolicyCreate(utils.URLGen("/api/policies"))
	return nil
}

// PostPolicyCreate let user creates a policy, and if it is enabled, the replication will be triggered right now.
//
//  params:
//
//  format:
//
// e.g.
func PostPolicyCreate(baseURL string) {}

type policiesList struct {
	Name      string `short:"n" long:"name" description:"The replication's policy name." default:""`
	ProjectID int    `short:"j" long:"project_id" description:"The ID of project." default:"0"`
	Page      int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize  int    `short:"s" long:"page_size" description:"The size of per page, default is 10, maximum is 100." default:"10"`
}

var poList policiesList

func (x *policiesList) Execute(args []string) error {
	GetPoliciesList(utils.URLGen("/api/policies"))
	return nil
}

// GetPoliciesList let user list filters policies by name and project_id, if name and project_id are nil, list returns all policies.
//
// params:
//   name       - The replication's policy name.
//   project_id - The ID of project.
//   page       - The page nubmer, default is 1.
//   page_size  - The size of per page, default is 10, maximum is 100.
//
// format:
//   GET /policies/replication
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/policies/replication?name=repl_policy_name&project_id=86&page=1&page_size=10'
func GetPoliciesList(baseURL string) {
	targetURL := baseURL + "/replication?name=" + poList.Name +
		"&project_id=" + strconv.Itoa(poList.ProjectID) +
		"&page=" + strconv.Itoa(poList.Page) +
		"&page_size=" + strconv.Itoa(poList.PageSize)
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

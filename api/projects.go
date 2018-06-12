package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("prj_create",
		"Create a new project.",
		"This endpoint is for user to create a new project.",
		&prjCreate)
	utils.Parser.AddCommand("prj_get",
		"Return specific project detail information.",
		"This endpoint returns specific project information by project ID.",
		&prjGet)
	utils.Parser.AddCommand("prj_del",
		"Delete a project by project_id.",
		"This endpoint is aimed to delete a project by project_id.",
		&prjDel)
	utils.Parser.AddCommand("prjs_list",
		"List projects.",
		"This endpoint returns all projects created by Harbor, and can be filtered by project name.",
		&prjsList)
}

type projectCreate struct {
	ProjectName                                string `short:"n" long:"project_name" description:"(REQUIRED) The name of the project." required:"yes" json:"project_name"`
	Public                                     int    `short:"k" long:"public" description:"(REQUIRED) The public status of the project, public(1) or private(0)." required:"yes" json:"public"`
	EnablelontentTrust                         bool   `short:"t" long:"enable_content_trust" description:"Whether content trust is enabled or not. If it is enabled, user cann't pull unsigned images from this project." json:"enable_content_trust"`
	PreventVulnerableImagesFromRunning         bool   `short:"r" long:"prevent_vulnerable_images_from_running" description:"Whether prevent the vulnerable images from running." json:"prevent_vulnerable_images_from_running"`
	PreventVulnerableImagesFromRunningSeverity string `short:"s" long:"prevent_vulnerable_images_from_running_severity" description:"If the vulnerability is high than severity defined here, the images cann't be pulled." default:"" json:"prevent_vulnerable_images_from_running_severity"`
	AutomaticallyScanImagesOnPush              bool   `short:"a" long:"automatically_scan_images_on_push" description:"Whether scan images automatically when pushing." json:"automatically_scan_images_on_push"`
}

var prjCreate projectCreate

func (x *projectCreate) Execute(args []string) error {
	PostPrjCreate(utils.URLGen("/api/projects"))
	return nil
}

type projectGet struct {
	ProjectID int `short:"j" long:"project_id" description:"(REQUIRED) Project ID of project which will be get." required:"yes"`
}

var prjGet projectGet

func (x *projectGet) Execute(args []string) error {
	GetPrjByPrjID(utils.URLGen("/api/projects"))
	return nil
}

type projectDel struct {
	ProjectID int `short:"j" long:"project_id" description:"(REQUIRED) Project ID of project which will be deleted." required:"yes"`
}

var prjDel projectDel

func (x *projectDel) Execute(args []string) error {
	DelPrjByPrjID(utils.URLGen("/api/projects"))
	return nil
}

type projectsList struct {
	Name string `short:"n" long:"name" description:"The name of the project (for filtering)." default:""`
	// NOTE:
	// 这里将 public 的类型从 bool 变更为 string ，因为bool 类型只有 true 和 false 二值语义，而实际使用中需要第三种语义
	// 1. 若为 true 则仅返回 public 项目；
	// 2. 若为 false 则仅返回 private 项目；
	// 3. 若不指定 public 参数，则应该同时返回 public 和 private 项目；
	Public string `short:"k" long:"public" description:"The project is public or private. (default: \"\")" default:""`
	// FIXME:
	// harbor 中基于 owner 过滤的功能似乎存在问题；
	Owner    string `short:"o" long:"owner" description:"The name of project owner." default:""`
	Page     int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize int    `short:"s" long:"page_size" description:"The size of per page, default is 10, maximum is 100." default:"10"`
}

var prjsList projectsList

func (x *projectsList) Execute(args []string) error {
	GetPrjsList(utils.URLGen("/api/projects"))
	return nil
}

// PostPrjCreate is for user to create a new project.
//
// params:
//  project_name - The name of the project.
//  public - The public status of the project, public(1) or private(0).
//  enable_content_trust - Whether content trust is enabled or not. If it is enabled, user cann't pull unsigned images from this project.
//  prevent_vulnerable_images_from_running - Whether prevent the vulnerable images from running.
//  prevent_vulnerable_images_from_running_severity - If the vulnerability is high than severity defined here, the images cann't be pulled.
//  automatically_scan_images_on_push - Whether scan images automatically when pushing.
//
// e.g.
/*
curl -X POST --header 'Content-Type: text/plain' --header 'Accept: text/plain' -d '{
  "project_name": "t1",
  "public": 0,
  "enable_content_trust": false,
  "prevent_vulnerable_images_from_running": false,
  "prevent_vulnerable_images_from_running_severity": "",
  "automatically_scan_images_on_push": false
}' 'https://11.11.11.12/api/projects'
*/
func PostPrjCreate(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	p, err := json.Marshal(&prjCreate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(p)).
		End(utils.PrintStatus)
}

// GetPrjByPrjID returns specific project information by project ID.
//
// params:
//  project_id - (REQUIRED) Project ID of project which will be get.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/projects/100'
func GetPrjByPrjID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(prjGet.ProjectID)
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

// DelPrjByPrjID is aimed to delete project by project ID
//
// params:
//  project_id - (REQUIRED) Project ID of project which will be deleted.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://11.11.11.12/api/projects/100'
func DelPrjByPrjID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(prjDel.ProjectID)
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

// GetPrjsList returns all projects created by Harbor, and can be filtered by project name.
//
// params:
//  name - The name of project.
//  public - The project is public or private. default is "", return both public and private prjs.
//  owner - The name of project owner.
//  page - The page nubmer, default is 1.
//  page_size - The size of per page, default is 10, maximum is 100.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/projects?name=prj&public=true&owner=moooofly&page=1&page_size=10'
func GetPrjsList(baseURL string) {
	targetURL := baseURL + "?name=" + prjsList.Name +
		"&public=" + prjsList.Public +
		"&owner=" + prjsList.Owner +
		"&page=" + strconv.Itoa(prjsList.Page) +
		"&page_size=" + strconv.Itoa(prjsList.PageSize)
	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		// TODO:
		// 可以通过解析 Rsp Heaer 中的 X-Total-Count 直接得到返回的 projects 数量
		End(utils.PrintStatus)
}

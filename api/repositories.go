package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("repo_image_label_del",
		"Delete label from the image under specific repository.",
		"This endpoint deletes the label from the image specified by the repo_name and tag.",
		&repoImageLabelDel)
	utils.Parser.AddCommand("repo_image_label_add",
		"Add a label to the image under specific repository.",
		"This endpoint adds a label to the image under specific repository.",
		&repoImageLabelAdd)
	utils.Parser.AddCommand("repo_image_labels_get",
		"Get labels of an image under specific repository.",
		"This endpoint gets labels of an image under specific repository specified by the repo_name and tag.",
		&repoImageLabelsGet)
	utils.Parser.AddCommand("repo_label_del",
		"Delete a label from the repository.",
		"This endpoint deletes the label from the repository specified by the repo_name.",
		&repoLabelDel)
	utils.Parser.AddCommand("repo_label_add",
		"Add a label to the repository.",
		"This endpoint adds an already existing label (global or project specific) to the repository.",
		&repoLabelAdd)
	utils.Parser.AddCommand("repo_labels_get",
		"Get labels of a repository.",
		"This endpoint gets labels of a repository specified by the repo_name. NOTE: This API gets '401 Unauthorized' all the time, even when logging in as admin user.",
		&repoLabelsGet)
	utils.Parser.AddCommand("repo_desp_update",
		"Update description of the repository.",
		"This endpoint is used to update description of the repository.",
		&repoUpdate)
	utils.Parser.AddCommand("repo_del",
		"Delete a repository by repo_name.",
		"This endpoint let user delete a repository by repo_name.",
		&repoDel)
	utils.Parser.AddCommand("repos_list",
		"Get repositories accompany with relevant project and repo name.",
		"This endpoint let user search repositories accompanying with relevant project ID and repo name.",
		&reposList)
	utils.Parser.AddCommand("repos_top",
		"Get public repositories which are accessed most.",
		"This endpoint aims to let users see the most popular public repositories",
		&reposTop)
}

type repositoryImageLabelDel struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository." required:"yes"`
	Tag      string `short:"t" long:"tag" description:"(REQUIRED) The tag of the image." required:"yes"`
	LabelID  int    `short:"i" long:"label_id" description:"(REQUIRED) The ID of label." required:"yes"`
}

var repoImageLabelDel repositoryImageLabelDel

func (x *repositoryImageLabelDel) Execute(args []string) error {
	DeleteRepoImageLabel(utils.URLGen("/api/repositories"))
	return nil
}

// DeleteRepoImageLabel deletes the label from the image specified by the repo_name and tag.
//
// params:
//   repo_name - (REQUIRED) The name of repository.
//   tag       - (REQUIRED) The tag of the image.
//   id        - (REQUIRED) The ID of label.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/repositories/temp_3%2Fhello-world/tags/v1/labels/2'
func DeleteRepoImageLabel(baseURL string) {
	targetURL := baseURL + "/" + repoImageLabelDel.RepoName +
		"/tags/" + repoImageLabelDel.Tag +
		"/labels/" + strconv.Itoa(repoImageLabelDel.LabelID)
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

type repositoryImageLabelAdd struct {
	RepoName     string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository that you want to add a label." required:"yes"`
	Tag          string `short:"t" long:"tag" description:"(REQUIRED) The tag of the image." required:"yes"`
	ID           int    `short:"i" long:"id" description:"(REQUIRED) The ID of the already existing label." required:"yes" json:"id"`
	Name         string `long:"name" description:"The name of this label." default:"" json:"name"`
	Description  string `long:"description" description:"The description of this label." default:"" json:"description"`
	Color        string `long:"color" description:"The color code of this label. (e.g. Format: #A9B6BE)" default:"" json:"color"`
	Scope        string `long:"scope" description:"The scope of this label. ('p' indicats project scope, 'g' indicates global scope)" default:"" json:"scope"`
	ProjectID    int    `long:"project_id" description:"Which project (id) this label belongs to when created. ('0' indicates global label, others indicate specific project)" default:"" json:"project_id"`
	CreationTime string `long:"creation_time" description:"The creation time of this label. default time.Now()" default:"" json:"creation_time"`
	UpdateTime   string `long:"update_time" description:"The update time of this label. default time.Now()" default:"" json:"update_time"`
	Deleted      bool   `long:"deleted" description:"not sure" json:"deleted"`
}

var repoImageLabelAdd repositoryImageLabelAdd

func (x *repositoryImageLabelAdd) Execute(args []string) error {
	PostRepoImageLabelAdd(utils.URLGen("/api/repositories"))
	return nil
}

// PostRepoImageLabelAdd adds a label to the image under specific repository.
//
// params:
//   repo_name     - (REQUIRED) The name of repository that you want to add a label.
//   tag           - (REQUIRED) The tag of the image.
//   id            - (REQUIRED) The ID of the already existing label.
//   name          - The name of this label.
//   description   - The description of this label.
//   color         - The color code of this label. (e.g. Format: #A9B6BE)
//   scope         - The scope of this label. ('p' indicats project scope, 'g' indicates global scope)
//   project_id    - Which project (id) this label belongs to when created. ('0' indicates global label, others indicate specific project)
//   creation_time - The creation time of this label. default time.Now()
//   update_time   - The update time of this label. default time.Now()
//   deleted       - not sure
//
// format:
//   POST /repositories/{repo_name}/tags/{tag}/labels
/*
curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "id": 3, \
   "name": "test-name", \
   "description": "test-description", \
   "color": "test-color", \
   "scope": "test-scope", \
   "project_id": 10, \
   "deleted": true \
 }' 'https://localhost/api/repositories/temp_3%2Fhello-world/tags/v1/labels'
*/
func PostRepoImageLabelAdd(baseURL string) {
	if repoImageLabelAdd.CreationTime == "" || repoImageLabelAdd.UpdateTime == "" {
		now := time.Now().Format("2006-01-02T15:04:05Z")
		repoImageLabelAdd.CreationTime = now
		repoImageLabelAdd.UpdateTime = now
	}

	targetURL := baseURL + "/" + repoImageLabelAdd.RepoName +
		"/tags/" + repoImageLabelAdd.Tag + "/labels"
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&repoImageLabelAdd)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> label add:", string(t))

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type repositoryImageLabelsGet struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository." required:"yes"`
	Tag      string `short:"t" long:"tag" description:"(REQUIRED) The tag of the image." required:"yes"`
}

var repoImageLabelsGet repositoryImageLabelsGet

func (x *repositoryImageLabelsGet) Execute(args []string) error {
	GetRepoImageLabel(utils.URLGen("/api/repositories"))
	return nil
}

// GetRepoImageLabel gets labels of an image specified by the repo_name and tag.
//
// params:
//   repo_name - (REQUIRED) The name of repository.
//   tag       - (REQUIRED) The tag of the image.
//
// format:
//   GET /repositories/{repo_name}/tags/{tag}/labels
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/repositories/temp_3%2Fhello-world/tags/v1/labels'
func GetRepoImageLabel(baseURL string) {
	targetURL := baseURL + "/" + repoImageLabelsGet.RepoName +
		"/tags/" + repoImageLabelsGet.Tag + "/labels"
	fmt.Println("==> GET", targetURL)

	utils.Request.Get(targetURL).End(utils.PrintStatus)
}

type repositoryLabelDel struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository that you want to delete a label from." required:"yes"`
	ID       int    `short:"i" long:"id" description:"(REQUIRED) The ID of label." required:"yes"`
}

var repoLabelDel repositoryLabelDel

func (x *repositoryLabelDel) Execute(args []string) error {
	DeleteRepoLabel(utils.URLGen("/api/repositories"))
	return nil
}

// DelRepoByRepoName deletes the label from the repository specified by the repo_name.
//
// params:
//   repo_name - (REQUIRED) The name of repository that you want to delete a label from.
//   id        - (REQUIRED) The ID of label.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/repositories/temp_3%2Fhello-world/labels/2'
func DeleteRepoLabel(baseURL string) {
	targetURL := baseURL + "/" + repoLabelDel.RepoName +
		"/labels/" + strconv.Itoa(repoLabelDel.ID)
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

type repositoryLabelAdd struct {
	RepoName     string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository that you want to add a label." required:"yes"`
	ID           int    `short:"i" long:"id" description:"(REQUIRED) The ID of the already existing label." required:"yes" json:"id"`
	Name         string `long:"name" description:"The name of this label." default:"" json:"name"`
	Description  string `long:"description" description:"The description of this label." default:"" json:"description"`
	Color        string `long:"color" description:"The color code of this label. (e.g. Format: #A9B6BE)" default:"" json:"color"`
	Scope        string `long:"scope" description:"The scope of this label. ('p' indicats project scope, 'g' indicates global scope)" default:"" json:"scope"`
	ProjectID    int    `long:"project_id" description:"Which project (id) this label belongs to when created. ('0' indicates global label, others indicate specific project)" default:"" json:"project_id"`
	CreationTime string `long:"creation_time" description:"The creation time of this label. default time.Now()" default:"" json:"creation_time"`
	UpdateTime   string `long:"update_time" description:"The update time of this label. default time.Now()" default:"" json:"update_time"`
	Deleted      bool   `long:"deleted" description:"not sure" json:"deleted"`
}

var repoLabelAdd repositoryLabelAdd

func (x *repositoryLabelAdd) Execute(args []string) error {
	PostRepoLabelAdd(utils.URLGen("/api/repositories"))
	return nil
}

// PostRepoLabelAdd add a label to the repository.
//
// params:
//   repo_name     - (REQUIRED) The name of repository that you want to add a label.
//   id            - (REQUIRED) The ID of the already existing label.
//   name          - The name of this label.
//   description   - The description of this label.
//   color         - The color code of this label. (e.g. Format: #A9B6BE)
//   scope         - The scope of this label. ('p' indicats project scope, 'g' indicates global scope)
//   project_id    - Which project (id) this label belongs to when created. ('0' indicates global label, others indicate specific project)
//   creation_time - The creation time of this label. default time.Now()
//   update_time   - The update time of this label. default time.Now()
//   deleted       - not sure
//
// format:
//   POST /repositories/{repo_name}/labels
/*
curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "id": 2, \
   "name": "test-name", \
   "description": "test-description", \
   "color": "test-color", \
   "scope": "test-scope", \
   "project_id": 10, \
   "deleted": true \
 }' 'https://localhost/api/repositories/temp_5%2Fhello-world/labels'
*/
func PostRepoLabelAdd(baseURL string) {
	if repoLabelAdd.CreationTime == "" || repoLabelAdd.UpdateTime == "" {
		now := time.Now().Format("2006-01-02T15:04:05Z")
		repoLabelAdd.CreationTime = now
		repoLabelAdd.UpdateTime = now
	}

	targetURL := baseURL + "/" + repoLabelAdd.RepoName + "/labels"
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&repoLabelAdd)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> label add:", string(t))

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type repositoryLabelsGet struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository." required:"yes"`
}

var repoLabelsGet repositoryLabelsGet

func (x *repositoryLabelsGet) Execute(args []string) error {
	GetRepoLabels(utils.URLGen("/api/repositories"))
	return nil
}

// GetRepoLabels get labels of a repository specified by the repo_name.
//
// params:
//   repo_name - (REQUIRED) The name of repository.
//
// format:
//   GET /repositories/{repo_name}/labels
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/repositories/temp_5%2Fhello-world/labels'
func GetRepoLabels(baseURL string) {
	targetURL := baseURL + "/" + repoLabelsGet.RepoName + "/labels"
	fmt.Println("==> GET", targetURL)

	utils.Request.Get(targetURL).End(utils.PrintStatus)
}

type repoDescriptionUpdate struct {
	RepoName    string `short:"n" long:"repo_name" description:"(REQUIRED) Repo name for filtering results." required:"yes" json:"-"`
	Description string `short:"d" long:"description" description:"(REQUIRED) The description of the repository." required:"yes" json:"description"`
}

var repoUpdate repoDescriptionUpdate

func (x *repoDescriptionUpdate) Execute(args []string) error {
	PutRepoDescriptionUpdate(utils.URLGen("/api/repositories"))
	return nil
}

// PutRepoDescriptionUpdate is used to update description of the repository.
//
// params:
//   repo_name - (REQUIRED) The name of repository which will be deleted.
//   description - (REQUIRED) The description of the repository.
//
// format:
//   PUT /repositories/{repo_name}
//
// e.g.
/*
curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "description": "change" \
 }' 'https://localhost/api/repositories/temp_5%2Fhello-world'
*/

func PutRepoDescriptionUpdate(baseURL string) {
	targetURL := baseURL + "/" + repoUpdate.RepoName
	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&repoUpdate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> description:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type repositoriesList struct {
	ProjectID int    `short:"j" long:"project_id" description:"(REQUIRED) Relevant project ID." required:"yes"`
	RepoName  string `short:"n" long:"repo_name" description:"Repo name for filtering results." default:""`
	LabelID   int    `short:"l" long:"label_id" description:"The ID of label used to filter the result." default:"0"`
	Page      int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize  int    `short:"s" long:"page_size" description:"The size of per page, default is 10, maximum is 100." default:"10"`
}

var reposList repositoriesList

func (x *repositoriesList) Execute(args []string) error {
	GetReposByPrjID(utils.URLGen("/api/repositories"))
	return nil
}

type repositoriesTop struct {
	Count int `short:"c" long:"count" description:"The number of the requested public repositories, default is 10 if not provided." default:"10"`
}

var reposTop repositoriesTop

func (x *repositoriesTop) Execute(args []string) error {
	GetTopRepos(utils.URLGen("/api/repositories/top"))
	return nil
}

type repositoryDel struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository which will be deleted." required:"yes"`
}

var repoDel repositoryDel

func (x *repositoryDel) Execute(args []string) error {
	DelRepoByRepoName(utils.URLGen("/api/repositories"))
	return nil
}

// GetReposByPrjID let user search repositories accompanying with relevant project ID and repo name.
//
// params:
//   project_id - (REQUIRED) Relevant project ID.
//   q          - Repo name for filtering results.
//   label_id   - The ID of label used to filter the result.
//   page       - The page nubmer, default is 1.
//   pageSize   - The size of per page, default is 10, maximum is 100.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/repositories?project_id=1&q=prj&label_id=100&page=1&page_size=10'
func GetReposByPrjID(baseURL string) {
	targetURL := baseURL + "?project_id=" + strconv.Itoa(reposList.ProjectID) +
		"&q=" + reposList.RepoName +
		"&label_id=" + strconv.Itoa(reposList.LabelID) +
		"&page=" + strconv.Itoa(reposList.Page) +
		"&page_size=" + strconv.Itoa(reposList.PageSize)
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

// GetTopRepos aims to let users see the most popular public repositories
//
// params:
//   count - The number of the requested public repositories, default is 10 if not provided.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/repositories/top?count=3'
func GetTopRepos(baseURL string) {
	targetURL := baseURL + "?count=" + strconv.Itoa(reposTop.Count)
	fmt.Println("==> GET", targetURL)

	utils.Request.Get(targetURL).End(utils.PrintStatus)
}

// DelRepoByRepoName let user delete a repository with name.
//
// params:
//   repo_name - (REQUIRED) The name of repository which will be deleted.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/repositories/prj1%2Fhello-world'
func DelRepoByRepoName(baseURL string) {
	targetURL := baseURL + "/" + repoDel.RepoName
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

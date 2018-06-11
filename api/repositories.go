package api

import (
	"fmt"
	"strconv"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("repos_list",
		"Get repositories accompany with relevant project and repo name.",
		"This endpoint let user search repositories accompanying with relevant project ID and repo name.",
		&reposList)
	utils.Parser.AddCommand("repos_top",
		"Get public repositories which are accessed most.",
		"This endpoint aims to let users see the most popular public repositories",
		&reposTop)
	utils.Parser.AddCommand("repo_del",
		"Delete a repository by repo_name.",
		"This endpoint let user delete a repository by repo_name.",
		&repodel)
}

type repositoriesList struct {
	ProjectID int    `short:"j" long:"project_id" description:"(REQUIRED) Relevant project ID." required:"yes"`
	RepoName  string `short:"n" long:"repo_name" description:"Repo name for filtering results." default:""`
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

var repodel repositoryDel

func (x *repositoryDel) Execute(args []string) error {
	DelRepoByRepoName(utils.URLGen("/api/repositories"))
	return nil
}

// GetReposByPrjID let user search repositories accompanying with relevant project ID and repo name.
//
// params:
//  project_id - (REQUIRED) Relevant project ID.
//  q          - Repo name for filtering results.
//  page       - The page nubmer, default is 1.
//  pageSize   - The size of per page, default is 10, maximum is 100.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/repositories?project_id=1&q=prj&page=1&page_size=10'
func GetReposByPrjID(baseURL string) {
	targetURL := baseURL + "?project_id=" + strconv.Itoa(reposList.ProjectID) +
		"&q=" + reposList.RepoName +
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
//  count - The number of the requested public repositories, default is 10 if not provided.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/repositories/top?count=3'
func GetTopRepos(baseURL string) {
	targetURL := baseURL + "?count=" + strconv.Itoa(reposTop.Count)
	fmt.Println("==> GET", targetURL)

	utils.Request.Get(targetURL).End(utils.PrintStatus)
}

// DelRepoByRepoName let user delete a repository with name.
//
// params:
//  repo_name - (REQUIRED) The name of repository which will be deleted.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://11.11.11.12/api/repositories/prj1%2Fhello-world'
func DelRepoByRepoName(baseURL string) {
	targetURL := baseURL + "/" + repodel.RepoName
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

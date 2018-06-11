package api

import (
	"fmt"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("tag_get",
		"Get the tag of the repository.",
		"This endpoint aims to retrieve the tag of the repository. If deployed with Notary, the signature property of response represents whether the image is singed or not. If the property is null, the image is unsigned.",
		&tagget)
	utils.Parser.AddCommand("tag_del",
		"Delete a tag in a repository.",
		"This endpoint let user delete tags with repo name and tag.",
		&tagdel)
	utils.Parser.AddCommand("tags_list",
		"Get tags of a relevant repository.",
		"This endpoint aims to retrieve tags from a relevant repository. If deployed with Notary, the signature property of response represents whether the image is singed or not. If the property is null, the image is unsigned.",
		&tagslist)
}

type tagGet struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) Relevant repository name." required:"yes"`
	Tag      string `short:"t" long:"tag" description:"(REQUIRED) Tag of the repository." required:"yes"`
}

var tagget tagGet

func (x *tagGet) Execute(args []string) error {
	GetTaginfoOfRepo(utils.URLGen("/api/repositories"))
	return nil
}

type tagDel struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) The name of repository which will be deleted." required:"yes"`
	Tag      string `short:"t" long:"tag" description:"(REQUIRED) Tag of a repository." required:"yes"`
}

var tagdel tagDel

func (x *tagDel) Execute(args []string) error {
	DelTaginfoOfRepo(utils.URLGen("/api/repositories"))
	return nil
}

type tagsList struct {
	RepoName string `short:"n" long:"repo_name" description:"(REQUIRED) Relevant repository name." required:"yes"`
}

var tagslist tagsList

func (x *tagsList) Execute(args []string) error {
	GetTagsByRepoName(utils.URLGen("/api/repositories"))
	return nil
}

// GetTaginfoOfRepo aims to retrieve the tag of the repository. If deployed with Notary, the signature property of
// response represents whether the image is singed or not. If the property is null, the image is unsigned.
//
// params:
//  repo_name - (REQUIRED) Relevant repository name.
//  tag       - (REQUIRED) Tag of the repository.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/repositories/prj2%2Fphoton/tags/v2'
func GetTaginfoOfRepo(baseURL string) {
	targetURL := baseURL + "/" + tagget.RepoName + "/tags/" + tagget.Tag
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

// DelTaginfoOfRepo let user delete tags with repo name and tag.
//
// params:
//  repo_name - (REQUIRED) The name of repository which will be deleted.
//  tag       - (REQUIRED) Tag of a repository.
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://11.11.11.12/api/repositories/prj2%2Fphoton/tags/v2'
func DelTaginfoOfRepo(baseURL string) {
	targetURL := baseURL + "/" + tagdel.RepoName + "/tags/" + tagdel.Tag
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

// GetTagsByRepoName aims to retrieve tags from a relevant repository. If deployed with Notary, the signature property of response represents whether the image is singed or not. If the property is null, the image is unsigned.
//
// params:
//  repo_name - (REQUIRED) Relevant repository name.
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://11.11.11.12/api/repositories/prj2%2Fphoton/tags'
func GetTagsByRepoName(baseURL string) {
	targetURL := baseURL + "/" + tagslist.RepoName + "/tags"
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

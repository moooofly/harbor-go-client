package utils

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type statistics struct {
	PrivateProjectCount int `json:"private_project_count"`
	PrivateRepoCount    int `json:"private_repo_count"`
	PublicProjectCount  int `json:"public_project_count"`
	PublicRepoCount     int `json:"public_repo_count"`
	TotalProjectCount   int `json:"total_project_count"`
	TotalRepoCount      int `json:"total_repo_count"`
}

var stats statistics

// ---

type repoTop struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ProjectID    int    `json:"project_id"`
	Description  string `json:"description"`
	PullCount    int    `json:"pull_count"`
	StarCount    int    `json:"star_count"`
	TagsCount    int    `json:"tags_count"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
}

var repos []*repoTop

// ---

type repoSearch struct {
	ProjectID      int    `json:"project_id"`
	ProjectName    string `json:"project_name"`
	ProjectPublic  int    `json:"project_public"`
	PullCount      int    `json:"pull_count"`
	RepositoryName string `json:"repository_name"`
	TagsCount      int    `json:"tags_count"`
}

type searchRsp struct {
	Repository []*repoSearch `json:"repository"`
	Project    []interface{} `json:"project"`
}

var scRsp searchRsp

// ---

type tagInfo struct {
	Digest        string `json:"digest"`
	Name          string `json:"name"`
	Architecture  string `json:"architecture"`
	DockerVersion string `json:"docker_version"`
	Author        string `json:"author"`
	Created       string `json:"created"`
	Signature     string `json:"signature"`
}

type tagListRsp []*tagInfo

var tlRsp tagListRsp

func init() {
	Parser.AddCommand("rp_repos",
		"Delete repos by retention policy.",
		"Run retention policy analysis on Repositories, do soft deletion as you command, prompt user performing a GC.",
		&reposRP)
	Parser.AddCommand("rp_tags",
		"Delete tags of repo by retention policy.",
		"Run retention policy analysis on tags, and do deletion as you command.",
		&tagsRP)
}

type reposRetentionPolicy struct {
}

var reposRP reposRetentionPolicy

func (x *reposRetentionPolicy) Execute(args []string) error {
	if err := repoAnalyse(); err != nil {
		os.Exit(1)
	}
	if err := repoErase(); err != nil {
		os.Exit(1)
	}
	rpGCHint()
	return nil
}

type tagsRetentionPolicy struct {
	Day      int    `short:"d" long:"day" description:"(REQUIRED) The tags of a repository created less than N days should not be deleted." required:"yes"`
	Max      int    `short:"m" long:"max" description:"(REQUIRED) The maximum quantity of tags created more than N days of a repository should keep untouched." required:"yes"`
	RepoName string `short:"n" long:"repo_name" description:"Repo name for specific target. If not set, rp_tags will do jobs on all repos." default:""`
}

var tagsRP tagsRetentionPolicy

func (x *tagsRetentionPolicy) Execute(args []string) error {
	if err := tagAnalyseAndErase(); err != nil {
		os.Exit(1)
	}
	return nil
}

func tagAnalyseAndErase() error {
	fmt.Println("=========================")
	fmt.Println("==  开始 tags RP 分析  ==")
	fmt.Println("=========================")
	fmt.Println()

	c, err := CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	// 基于 search 接口获取全部 projects 和 repositories 信息
	// 设置 "q=" 可以获取全部信息
	// 设置 "q=xxx" 可以过滤指定信息，但是目前发现该功能有 bug ，故暂时无法基于该接口针对指定 repo 进行处理
	searchURL := URLGen("/api/search") + "?q=" + tagsRP.RepoName
	fmt.Println("--------------------")
	fmt.Println("==> GET", searchURL)

	_, _, errs := Request.Get(searchURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		EndStruct(&scRsp)
	for _, e := range errs {
		if e != nil {
			fmt.Println("error:", e)
			return e
		}
	}

	if tagsRP.RepoName == "" {
		fmt.Printf("==> on all Repos, max-days-untouched: %d   max-keep-num-after-Ndays: %d\n",
			tagsRP.Day, tagsRP.Max)
	} else {
		fmt.Printf("==> only on Repo [%s], max-days-untouched: %d   max-keep-num-after-Ndays: %d\n",
			tagsRP.RepoName, tagsRP.Day, tagsRP.Max)
	}
	fmt.Println("--------------------")

	// 遍历全部 repositories 信息
	for _, r := range scRsp.Repository {
		fmt.Println("\n")
		fmt.Println("------------------------------------------------------")
		fmt.Printf("| repo_name: %s | tags_count: %d |\n", r.RepositoryName, r.TagsCount)
		fmt.Println("------------------------------------------------------")

		// 获取每个 repo 下的 tags 信息
		tagsListURL := URLGen("/api/repositories") + "/" + r.RepositoryName + "/tags"
		//fmt.Println("==> GET", tagsListURL)

		_, _, errs := Request.Get(tagsListURL).
			Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
			EndStruct(&tlRsp)
		for _, e := range errs {
			if e != nil {
				fmt.Println("error:", e)
				return e
			}
		}

		// 用于针对每个 repo 下的 tags 进行排序
		tagmh = tagminheap{}
		heap.Init(&tagmh)
		for _, t := range tlRsp {
			//fmt.Printf("==> name: %s    created: %s\n", t.Name, t.Created)

			tagC := rfc3339Transform(t.Created)
			dayPast := time.Now().Sub(tagC).Hours() / 24

			// a. 针对每个 repo ，创建于最近 N 天之内的所有 tag 不做处理
			if tagsRP.Day < int(dayPast) {
				it := &tagItem{
					tagName:   t.Name,
					timestamp: tagC.Unix(),
				}
				// 超过 N 天的 tag 保存到 minheap
				fmt.Printf("[PUSH] %s <==> create: %s    dayPast: %f\n", it.tagName, t.Created, dayPast)
				heap.Push(&tagmh, it)
			} else {
				fmt.Printf("[noPUSH] %s <==> create: %s    dayPast: %f\n", t.Name, t.Created, dayPast)
			}

		}
		// b. 针对创建于 N 天之外的 tag ，每个 repo 最多保留 Max 个
		gtNdays := tagmh.Len()
		fmt.Println("---")
		fmt.Printf("--> # of tags less than %d days: %d , # of tags more than %d days: %d\n",
			tagsRP.Day, r.TagsCount-gtNdays, tagsRP.Day, gtNdays)
		if gtNdays <= tagsRP.Max {
			fmt.Printf("--> max-keep-num-after-Ndays (%d) more than actual num (%d), so DO NOTHING.\n", tagsRP.Max, gtNdays)
		} else {
			fmt.Printf("--> max-keep-num-after-Ndays (%d) less than actual num (%d), so START DELETING.\n", tagsRP.Max, gtNdays)
			fmt.Println("---")
			for gtNdays > tagsRP.Max {
				it := heap.Pop(&tagmh).(*tagItem)
				fmt.Printf("[POP] %s <==> %d\n", it.tagName, it.timestamp)

				// 如果 len(minheap) > max 则删除 len(minheap) - max 个 tag
				targetURL := URLGen("/api/repositories") + "/" + r.RepositoryName + "/tags/" + it.tagName
				fmt.Println("==> DELETE", targetURL)

				Request.Delete(targetURL).
					Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
					End(PrintStatus)

				gtNdays--
			}
		}
	}

	fmt.Printf("\n=== 完成 tags RP 分析 ===\n\n")

	return nil
}

type retentionPolicy struct {
	UpdateTime struct {
		Base    float32 `yaml:"base" json:"base"`
		Factors []struct {
			Weight float32 `yaml:"weight" json:"weight"`
			Range  struct {
				Low  int `yaml:"low" json:"low"`
				High int `yaml:"high" json:"high"`
			} `yaml:"range" json:"range"`
		} `yaml:"factors" json:"factors"`
	} `yaml:"update_time" json:"update_time"`
	PullCount struct {
		Base    float32 `yaml:"base" json:"base"`
		Factors []struct {
			Weight float32 `yaml:"weight" json:"weight"`
			Range  struct {
				Low  int `yaml:"low" json:"low"`
				High int `yaml:"high" json:"high"`
			} `yaml:"range" json:"range"`
		} `yaml:"factors" json:"factors"`
	} `yaml:"pull_count" json:"pull_count"`
	TagsCount struct {
		Base    float32 `yaml:"base" json:"base"`
		Factors []struct {
			Weight float32 `yaml:"weight" json:"weight"`
			Range  struct {
				Low  int `yaml:"low" json:"low"`
				High int `yaml:"high" json:"high"`
			} `yaml:"range" json:"range"`
		} `yaml:"factors" json:"factors"`
	} `yaml:"tags_count" json:"tags_count"`
}

var rpfile = "./rp.yaml"

// rpLoad loads retention policy settings from rp.yaml
func rpLoad() (*retentionPolicy, error) {
	var rp retentionPolicy

	dataBytes, err := ioutil.ReadFile(rpfile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(dataBytes), &rp)
	if err != nil {
		return nil, err
	}

	return &rp, nil
}

// format exhibits current retention policy settings
func format(rp *retentionPolicy) {
	rp, err := rpLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	rps, err := json.MarshalIndent(rp, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("===>", string(rps))
}

// rfc3339Transform parses timestamp string as RFC3339 layout
func rfc3339Transform(in string) time.Time {
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	return t
}

// grade calculates the score of each repo according to retention policy
func grade(r *repoTop, rp *retentionPolicy) float32 {
	var uf, pf, tf float32

	day := time.Now().Sub(rfc3339Transform(r.UpdateTime)).Hours() / 24
	for _, f := range rp.UpdateTime.Factors {
		if f.Range.Low <= int(day) && int(day) < f.Range.High {
			uf = f.Weight
			break
		}
	}
	if uf == float32(0) {
		fmt.Printf("Out of range: day = %.2f, uf is 0.0\n", day)
		uf = float32(0)
	}

	for _, f := range rp.PullCount.Factors {
		if f.Range.Low <= r.PullCount && r.PullCount < f.Range.High {
			pf = f.Weight
			break
		}
	}
	if pf == float32(0) {
		fmt.Printf("Out of range: pull_count = %d, pf is 1.0\n", r.PullCount)
		pf = float32(1)
	}

	for _, f := range rp.TagsCount.Factors {
		if f.Range.Low <= r.TagsCount && r.TagsCount < f.Range.High {
			tf = f.Weight
			break
		}
	}
	if tf == float32(0) {
		fmt.Printf("Out of range: tags_count = %d, tf is 1.0\n", r.TagsCount)
		tf = float32(1)
	}

	score := rp.UpdateTime.Base*uf + rp.PullCount.Base*pf + rp.TagsCount.Base*tf
	fmt.Printf("[factors] ==> score = UpdateTimeBase*uf + PullCountBase*pf + TagsCountBase*tf = %.2f * %.2f + %.2f * %.2f + %.2f * %.2f = %.2f   repo_id: %d\n",
		rp.UpdateTime.Base, uf, rp.PullCount.Base, pf, rp.TagsCount.Base, tf, score, r.ID)

	return score
}

// repoAnalyse calculates scores and output topN element by minheap sort
func repoAnalyse() error {

	rp, err := rpLoad()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	// 格式化输出当前 RP 设置
	//format(rp)

	statsURL := URLGen("/api/statistics")

	c, err := CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	resp, _, statsErrs := Request.Get(statsURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		EndStruct(&stats)
	if resp.StatusCode != 200 {
		fmt.Printf("error: Expected StatusCode=200, actual StatusCode=%v\n", resp.StatusCode)
	}
	for _, e := range statsErrs {
		if e != nil {
			fmt.Println("error:", e)
			return e
		}
	}

	fmt.Println("=====> current Public Repo Count:", stats.PublicRepoCount)

	topURL := URLGen("/api/repositories/top") + "?count=" + strconv.Itoa(stats.PublicRepoCount)
	fmt.Println("==> GET", topURL)
	_, _, topErrs := Request.Get(topURL).EndStruct(&repos)
	for _, e := range topErrs {
		if e != nil {
			fmt.Println("error:", e)
			return e
		}
	}

	heap.Init(&minh)
	heap.Init(&mhBk)
	for _, r := range repos {
		sc := grade(r, rp)

		// NOTE:
		// 以下代码仅做调试使用
		// ===========
		/*
			rs, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
				return err
			}
			fmt.Printf("score (%f) =>\n%s\n", sc, rs)
		*/
		// ===========

		it := &repoItem{
			data:  r,
			score: sc,
		}
		heap.Push(&minh, it)
		heap.Push(&mhBk, it)
	}

	fmt.Printf("\n========== 根据分数排名（由低到高）建议删除 public repo 信息如下 ========\n\n")

	for mhBk.Len() > 0 {
		it := heap.Pop(&mhBk).(*repoItem)
		fmt.Printf("%.2f <==> %+v\n", it.score, *it.data)
	}

	return nil
}

// repoErase implements soft deletion
func repoErase() error {

	var num int
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n请输入希望删除的 repo 数量: ")
	for scanner.Scan() {
		in, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Print("输入数字不合法，请重新输入: ")
			continue
		}
		fmt.Println("您输入的数字为:", in)

		fmt.Print("确认吗 [y/n]: ")

		if !scanner.Scan() {
			break
		}
		confirm := scanner.Text()

		if strings.EqualFold(confirm, "y") {
			num = in
			break
		} else {
			fmt.Print("请重新输入希望删除的 repo 数量: ")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
		return err
	}

	if num <= 0 || num > 50 {
		fmt.Println("[Warning] The valid number range is (0, 50].")
		fmt.Println("[Warning] Sorry, you're not allowed proceeding... Abort.")
		return fmt.Errorf("error: the number is out of range")
	}

	fmt.Printf("\n=== 开始 soft deletion ===\n\n")

	for num > 0 {
		if minh.Len() > 0 {
			it := heap.Pop(&minh).(*repoItem)

			// NOTE:
			// 进行删除动作前，必须成功登陆，这里没有进行判定，而是直接发出 delete 动作

			// 对应 repo_del 的调用
			targetURL := URLGen("/api/repositories") + "/" + it.data.Name
			fmt.Println("==> DELETE", targetURL)

			c, err := CookieLoad()
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}

			Request.Delete(targetURL).
				Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
				End(PrintStatus)
		}
		num--
	}

	fmt.Printf("\n=== 完成 soft deletion ===\n\n")
	return nil
}

// rpGCHint gives a hint about hard deletion.
func rpGCHint() {

	fmt.Println("-----------------------------")
	fmt.Println("您已成功完成 soft deletion ，若想真正释放磁盘空间，还需要:")
	fmt.Println("1. 切换到 harbor 的安装主目录（例如 /opt/apps/harbor/）")
	fmt.Println("2. 运行如下命令以 preview 哪些 files/images 会被删除：")
	fmt.Println("    a. docker-compose stop")
	fmt.Println("    b. docker run -it --name gc --rm --volumes-from registry vmware/registry:2.6.2-photon garbage-collect --dry-run /etc/registry/config.yml")
	fmt.Println("3. 运行如下命令以真正触发 GC 动作：")
	fmt.Println("    a. docker run -it --name gc --rm --volumes-from registry vmware/registry:2.6.2-photon garbage-collect  /etc/registry/config.yml")
	fmt.Println("    b. docker-compose start")
	fmt.Println("")
	fmt.Println("WARNING:\nMake sure that no one is pushing images or Harbor is not running at all before you perform a GC. If someone were pushing an image while GC is running, there is a risk that the image's layers will be mistakenly deleted which results in a corrupted image. So before running GC, a preferred approach is to stop Harbor first.")
	fmt.Println("-----------------------------")
}

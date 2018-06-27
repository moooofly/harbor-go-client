package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("jobs_repl_list_by_filters",
		"List jobs filtered by specific policy and repository.",
		"This endpoint let user list jobs filtered by specific policy and repository. (if start_time and end_time are both null, list jobs of last 10 days)",
		&rplistbyfilter)
	utils.Parser.AddCommand("jobs_repl_stop_by_policy",
		"Update status of jobs. Only \"stop\" is supported for now.",
		"The endpoint is used to stop the replication jobs of a policy.",
		&replstopbypolicy)
	utils.Parser.AddCommand("jobs_repl_job_del_by_jid",
		"Delete replication job with specific ID.",
		"This endpoint is aimed to remove job with specific ID from jobservice.",
		&repljobdelbyid)
	utils.Parser.AddCommand("jobs_repl_log_get_by_jid",
		"Get replication job logs by specific job ID.",
		"This endpoint let user search job replication logs filtered by specific job ID.",
		&repllogbyid)
	utils.Parser.AddCommand("jobs_scan_log_get_by_jid",
		"Get scan job logs by specific job ID.",
		"This endpoint let user get scan job logs filtered by specific ID.",
		&scanlogbyid)
}

type replListByFilters struct {
	PolicyID   int    `short:"i" long:"policy_id" description:"(REQUIRED) The ID of the policy that triggered this job. (by targets_list and targets_policies_by_tid)" required:"yes"`
	Num        int    `short:"n" long:"num" description:"The length of return list." default:"50"`
	StartTime  string `short:"s" long:"start_time" description:"The start time of jobs. (format: yyyymmdd)" default:""`
	EndTime    string `short:"e" long:"end_time" description:"The end time of jobs. (format: yyyymmdd)" default:""`
	Repository string `short:"r" long:"repository" description:"The repository name to be filtered."`
	Status     string `short:"t" long:"status" description:"The status to be filtered. ([running|error|pending|retrying|stopped|finished|canceled])" default:""`
	Page       int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize   int    `short:"z" long:"page_size" description:"The size of per page, default is 10, maximum is 100." default:"10"`
}

var rplistbyfilter replListByFilters

func (x *replListByFilters) Execute(args []string) error {
	GetReplListByFilters(utils.URLGen("/api/jobs/replication"))
	return nil
}

type replStopByPolicy struct {
	PolicyID int    `short:"i" long:"policy_id" description:"(REQUIRED) The ID of replication policy." required:"yes" json:"policy_id"`
	Status   string `short:"s" long:"status" description:"(REQUIRED) The status of jobs to be changed into. The only valid value is \"stop\" for now." required:"yes" json:"status"`
}

var replstopbypolicy replStopByPolicy

func (x *replStopByPolicy) Execute(args []string) error {
	PutReplStopByPolicy(utils.URLGen("/api/jobs/replication"))
	return nil
}

type replJobDelByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) Replication job ID to delete." required:"yes" default:""`
}

var repljobdelbyid replJobDelByID

func (x *replJobDelByID) Execute(args []string) error {
	DelReplJobByID(utils.URLGen("/api/jobs/replication"))
	return nil
}

type replLogByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) Relevant job ID." required:"yes" default:""`
}

var repllogbyid replLogByID

func (x *replLogByID) Execute(args []string) error {
	GetReplLogByID(utils.URLGen("/api/jobs/replication"))
	return nil
}

type scanLogByID struct {
	ID int `short:"i" long:"id" description:"(REQUIRED) Relevant job ID." required:"yes" default:""`
}

var scanlogbyid scanLogByID

func (x *scanLogByID) Execute(args []string) error {
	GetScanLogByID(utils.URLGen("/api/jobs/scan"))
	return nil
}

// GetReplListByFilters list filtered jobs according to the policy and repository
//
// params:
//  policy_id  - (REQUIRED) The ID of the policy that triggered this job.
//  num        - The return list length number.
//  end_time   - The end time of jobs done. (Timestamp)
//  start_time - The start time of jobs. (Timestamp)
//  repository - The jobs list filtered by repository name.
//  status     - The jobs list filtered by status.
//               status must be one of [running|error|pending|retrying|stopped|finished|canceled].
//               If not set, means 'all' by default.
//  page       - The page nubmer, default is 1.
//  page_size  - The size of per page, default is 10, maximum is 100.
//
// operation format:
//  GET /jobs/replication/{id}/log
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/jobs/replication?page=1&page_size=15&status=finished&start_time=1529884800&end_time=1530057600&policy_id=6'
//
func GetReplListByFilters(baseURL string) {
	if rplistbyfilter.StartTime == "" || rplistbyfilter.EndTime == "" {
		// if start_time and end_time are both null, list jobs of last 10 days
		now := time.Now()
		rplistbyfilter.StartTime = now.AddDate(0, 0, -10).Format("20060102")
		rplistbyfilter.EndTime = now.Format("20060102")
	}

	//fmt.Println("StartTime:", rplistbyfilter.StartTime)
	//fmt.Println("EndTime:", rplistbyfilter.EndTime)

	st, err := time.Parse("20060102", rplistbyfilter.StartTime)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	et, err := time.Parse("20060102", rplistbyfilter.EndTime)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if rplistbyfilter.Status != "" &&
		rplistbyfilter.Status != "running" &&
		rplistbyfilter.Status != "error" &&
		rplistbyfilter.Status != "pending" &&
		rplistbyfilter.Status != "retrying" &&
		rplistbyfilter.Status != "stopped" &&
		rplistbyfilter.Status != "finished" &&
		rplistbyfilter.Status != "canceled" {
		fmt.Println("error: status must be one of [running|error|pending|retrying|stopped|finished|canceled].")
		os.Exit(1)
	}

	targetURL := baseURL + "?policy_id=" + strconv.Itoa(rplistbyfilter.PolicyID) +
		"&page=" + strconv.Itoa(rplistbyfilter.Page) +
		"&page_size=" + strconv.Itoa(rplistbyfilter.PageSize) +
		"&status=" + rplistbyfilter.Status +
		"&start_time=" + strconv.FormatInt(st.Unix(), 10) +
		"&end_time=" + strconv.FormatInt(et.Unix(), 10) +
		"&repository=" + rplistbyfilter.Repository +
		"&num=" + strconv.Itoa(rplistbyfilter.Num)

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

// PutReplStopByPolicy is used to stop the replication jobs of a policy.
//
// params:
//  policy_id - (REQUIRED) The ID of replication policy.
//  status    - (REQUIRED) The status of jobs. The only valid value is "stop" for now.
//
// operation format:
//  PUT /jobs/replication
//
// e.g.
/*
curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
   "policy_id": 1, \
   "status": "stop" \
}' 'https://localhost/api/jobs/replication'
*/
func PutReplStopByPolicy(baseURL string) {
	targetURL := baseURL

	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&replstopbypolicy)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> policyinfo:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

// DelReplJobByID is aimed to remove job with specific ID from jobservice.
//
// params:
//  id - (REQUIRED) Replication job ID to delete.
//
// operation format:
//  DELETE /jobs/replication/{id}
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/jobs/replication/1'
//
func DelReplJobByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(repljobdelbyid.ID)

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

// GetReplLogByID let user search job logs filtered by specific ID.
//
// params:
//  id - (REQUIRED) Relevant job ID.
//
// operation format:
//  GET /jobs/replication/{id}/log
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/jobs/replication/1/log'
//
func GetReplLogByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(repllogbyid.ID) + "/log"

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

// GetScanLogByID let user get scan job logs filtered by specific ID.
//
// params:
//  id - (REQUIRED) Relevant job ID.
//
// operation format:
//  GET /jobs/scan/{id}/log
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/jobs/scan/1/log'
//
func GetScanLogByID(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(scanlogbyid.ID) + "/log"

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

package api

import (
	"encoding/json"
	"fmt"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("replication_trigger_by_id",
		"Trigger the replication according to the specified policy.",
		"This endpoint is used to trigger a replication.",
		&replTriByID)
}

type replicationTriByID struct {
	PolicyID int `short:"i" long:"policy_id" description:"(REQUIRED) The ID of replication policy" required:"yes" json:"policy_id"`
}

var replTriByID replicationTriByID

func (x *replicationTriByID) Execute(args []string) error {
	PostReplTriByID(utils.URLGen("/api/replications"))
	return nil
}

// PostReplTriByID is used to trigger a replication.
//
// params:
//   policy_id - (REQUIRED) The ID of replication policy
//
// format:
//   POST /replications
//
// e.g.
/*
  curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
     "policy_id": 1 \
   }' 'https://localhost/api/replications'
*/
func PostReplTriByID(baseURL string) {
	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&replTriByID)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

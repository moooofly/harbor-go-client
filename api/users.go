package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/moooofly/harbor-go-client/utils"
)

func init() {
	utils.Parser.AddCommand("user_update_role",
		"Update a registered user to change to be an administrator of Harbor.",
		"This endpoint let a registered user change to be an administrator of Harbor.",
		&usrUpdateRole)
	utils.Parser.AddCommand("user_update_password",
		"Change the password on a user that already exists.",
		"This endpoint is for user to update password. Users with the admin role can change any user's password. Guest users can change only their own password.",
		&usrUpdatePassword)
	utils.Parser.AddCommand("user_update",
		"Update a registered user to change his profile.",
		"This endpoint let a registered user change his profile.",
		&usrUpdate)
	utils.Parser.AddCommand("user_get",
		"Get a user's profile.",
		"Get user's profile with user id.",
		&usrGet)
	utils.Parser.AddCommand("user_delete",
		"Mark a registered user as be removed.",
		"This endpoint let administrator of Harbor mark a registered user as be removed. It actually won't be deleted from DB.",
		&usrDelete)
	utils.Parser.AddCommand("user_create",
		"Creates a new user account.",
		"This endpoint is to create a user if the user does not already exist.",
		&usrCreate)
	utils.Parser.AddCommand("users_search",
		"Get registered users of Harbor.",
		"This endpoint is for user to search registered users, support for filtering results with username. Notice, by now this operation is only for administrator.",
		&usrSearch)
	// NOTE:
	// 由于 user_current 命令是是用于列出当前 login 用户相关信息
	// 故将其改名为 whoami
	utils.Parser.AddCommand("whoami",
		"Show info about current login user only.",
		"Maybe 'whoami' is a better name.",
		&usrCurrent)
}

type userUpdateRole struct {
	UserID       int `short:"i" long:"user_id" description:"(REQUIRED) Registered user ID." required:"yes" json:"-"`
	HasAdminRole int `short:"r" long:"has_admin_role" description:"(REQUIRED) Toggle a user to admin or not." required:"yes" json:"has_admin_role"`
}

var usrUpdateRole userUpdateRole

func (x *userUpdateRole) Execute(args []string) error {
	PutUserUpdateRole(utils.URLGen("/api/users"))
	return nil
}

// PutUserUpdateRole let a registered user change to be an administrator of Harbor.
//
// params:
//  id - (REQUIRED) Registered user ID.
//  has_admin_role - (REQUIRED) Toggle a user to admin or not.
//
// format:
//  PUT /users/{user_id}/password
//
// e.g.
// curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
//    "has_admin_role": 1 \
//  }' 'https://localhost/api/users/1/sysadmin'
//
func PutUserUpdateRole(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(usrUpdateRole.UserID) + "/sysadmin"

	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&usrUpdateRole)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> user_update_role:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type userUpdatePassword struct {
	UserID      int    `short:"i" long:"user_id" description:"(REQUIRED) Registered user ID." required:"yes" json:"-"`
	OldPassword string `short:"o" long:"old_password" description:"(REQUIRED) Old password." required:"yes" json:"old_password"`
	NewPassword string `short:"n" long:"new_password" description:"(REQUIRED) New password." required:"yes" json:"new_password"`
}

var usrUpdatePassword userUpdatePassword

func (x *userUpdatePassword) Execute(args []string) error {
	PutUserUpdatePassword(utils.URLGen("/api/users"))
	return nil
}

// PutUserUpdatePassword is for user to update password. Users with the admin role can change any user's password. Guest users can change only their own password.
//
// params:
//  id - (REQUIRED) Registered user ID.
//  old_password - (REQUIRED) Old password.
//  new_password - (REQUIRED) New password.
//
// format:
//  PUT /users/{user_id}/password
//
// e.g.
// curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
//    "old_password": "old password", \
//    "new_password": "new password" \
//  }' 'https://localhost/api/users/1/password'
//
func PutUserUpdatePassword(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(usrUpdatePassword.UserID) + "/password"

	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&usrUpdatePassword)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> user_update_password:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type userUpdate struct {
	UserID int `short:"i" long:"user_id" description:"(REQUIRED) Registered user ID." required:"yes" json:"-"`
	// Only email, realname and comment can be modified.
	Email    string `short:"e" long:"email" description:"(REQUIRED) User email." required:"yes" json:"email"`
	RealName string `short:"r" long:"realname" description:"(REQUIRED) User's realname." required:"yes" json:"realname"`
	Comment  string `short:"m" long:"comment" description:"(REQUIRED) Custom comment." required:"yes" json:"comment"`
}

var usrUpdate userUpdate

func (x *userUpdate) Execute(args []string) error {
	PutUserUpdate(utils.URLGen("/api/users"))
	return nil
}

// PutUserUpdate let a registered user change his profile.
//
// params:
//  id - (REQUIRED) Registered user ID.
//  email - (REQUIRED) User email.
//  realname - (REQUIRED) User's realname.
//  comment - (REQUIRED) Custom comment.
//
// format:
//  PUT /users/{user_id}
//
// e.g.
// curl -X PUT --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
//    "email": "si.li%40163.com", \
//    "realname": "si.li", \
//    "comment": "I'm Li Si" \
//  }' 'https://localhost/api/users/1'
//
func PutUserUpdate(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(usrUpdate.UserID)

	fmt.Println("==> PUT", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&usrUpdate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> user_update:", string(t))

	utils.Request.Put(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type userGet struct {
	UserID int `short:"i" long:"user_id" description:"(REQUIRED) Registered user ID." required:"yes"`
}

var usrGet userGet

func (x *userGet) Execute(args []string) error {
	GetUserProfile(utils.URLGen("/api/users"))
	return nil
}

// GetUserProfile gets user's profile with user id.
//
// params:
//  id - (REQUIRED) Registered user ID.
//
// format:
//  GET /users/{user_id}
//
// e.g. curl -X GET --header 'Accept: text/plain' 'https://localhost/api/users/1'
//
func GetUserProfile(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(usrGet.UserID)

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

type userDelete struct {
	UserID int `short:"i" long:"user_id" description:"(REQUIRED) User ID for marking as to be removed." required:"yes"`
}

var usrDelete userDelete

func (x *userDelete) Execute(args []string) error {
	DeleteUser(utils.URLGen("/api/users"))
	return nil
}

// DeleteUser let administrator of Harbor mark a registered user as be removed.It actually won't be deleted from DB.
//
// params:
//  id - (REQUIRED) User ID for marking as to be removed.
//
// format:
//  DELETE /users/{user_id}
//
// e.g. curl -X DELETE --header 'Accept: text/plain' 'https://localhost/api/users/1'
//
func DeleteUser(baseURL string) {
	targetURL := baseURL + "/" + strconv.Itoa(usrDelete.UserID)

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

type userCreate struct {
	UserID       int    `long:"user_id" description:"(REQUIRED) Registered user ID. Must be unique." required:"yes" json:"user_id"`
	Username     string `long:"username" description:"(REQUIRED) User name." required:"yes" json:"username"`
	Password     string `long:"password" description:"(REQUIRED) User password. (not support consealing here)" required:"yes" json:"password"`
	Email        string `long:"email" description:"(REQUIRED) User's email." required:"yes" json:"email"`
	HasAdminRole int    `long:"has_admin_role" description:"(REQUIRED) Mark a user whether is admin or not." required:"yes" json:"has_admin_role"`
	// realname can not be "", at least one character needed.
	RealName     string `long:"realname" description:"User's realname." default:" " json:"realname"`
	Comment      string `long:"comment" description:"Custom comment." default:"" json:"comment"`
	Deleted      int    `long:"deleted" description:"Deleted (no idea about this)." default:"0" json:"deleted"`
	RoleName     string `long:"role_name" description:"User's role name." default:"" json:"role_name"`
	RoleID       int    `long:"role_id" description:"User's role id." default:"0" json:"role_id"`
	ResetUUID    string `long:"reset_uuid" description:"Reset UUID (no idea about this)." default:"" json:"reset_uuid"`
	Salt         string `long:"salt" description:"Salt for password encryption." default:"" json:"salt"`
	CreationTime string `short:"c" long:"creation_time" description:"User's creation time. Default time.Now()." default:"" json:"creation_time"`
	UpdateTime   string `short:"u" long:"update_time" description:"User's update time. Default time.Now()." default:"" json:"update_time"`
}

var usrCreate userCreate

func (x *userCreate) Execute(args []string) error {
	PostUserCreate(utils.URLGen("/api/users"))
	return nil
}

// PostUserCreate Creates a new user account.
//
// params:
//  user_id - (REQUIRED) Registered user ID. Must be unique.
//  username - (REQUIRED) User name.
//  password - (REQUIRED) User password. (not support consealing here)
//  email - (REQUIRED) User's email.
//  has_admin_role - (REQUIRED) Mark a user whether is admin or not.
//  realname - User's realname.
//  comment - Custom comment.
//  deleted - Deleted (no idea about this).
//  role_name - User's role name.
//  role_id - User's role id.
//  reset_uuid - Reset UUID (no idea about this).
//  salt - Salt for password encryption.
//  creation_time - User's creation time. Default time.Now().
//  update_time - User's update time. Default time.Now().
//
// format:
//  POST /users
//
// e.g.
// curl -X POST --header 'Content-Type: application/json' --header 'Accept: text/plain' -d '{ \
//    "user_id": 100, \
//    "username": "san.zhang", \
//    "email": "san.zhang%40163.com", \
//    "password": "Harbor12345", \
//    "realname": "san.zhang", \
//    "comment": "I%27m Zhang San", \
//    "deleted": 0, \
//    "role_name": "", \
//    "role_id": 0, \
//    "has_admin_role": 0, \
//    "reset_uuid": "", \
//    "Salt": "", \
//    "creation_time": "2018-07-23T05:59:26Z", \
//    "update_time": "2018-07-23T05:59:26Z" \
//  }' 'https://localhost/api/users'
//
func PostUserCreate(baseURL string) {
	if usrCreate.CreationTime == "" || usrCreate.UpdateTime == "" {
		now := time.Now().Format("2006-01-02T15:04:05Z")
		usrCreate.CreationTime = now
		usrCreate.UpdateTime = now
	}

	targetURL := baseURL
	fmt.Println("==> POST", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	t, err := json.Marshal(&usrCreate)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("==> user_create:", string(t))

	utils.Request.Post(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		Send(string(t)).
		End(utils.PrintStatus)
}

type usersSearch struct {
	Username string `short:"u" long:"username" description:"Username for filtering results." default:""`
	Email    string `short:"e" long:"email" description:"Email for filtering results." default:""`
	Page     int    `short:"p" long:"page" description:"The page nubmer, default is 1." default:"1"`
	PageSize int    `short:"s" long:"page_size" description:"The size of per page, default is 10." default:"10"`
}

var usrSearch usersSearch

func (x *usersSearch) Execute(args []string) error {
	GetUsersSearch(utils.URLGen("/api/users"))
	return nil
}

// GetUsersSearch Get registered users of Harbor.
//
// params:
//  username - Username for filtering results.
//  email - Email for filtering results.
//  page - The page nubmer, default is 1.
//  page_size - The size of per page, default is 10.
//
// format:
//  GET /users
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/users?username=san.zhang&email=san.zhang@163.com&page=1&page_size=10'
//
func GetUsersSearch(baseURL string) {
	targetURL := baseURL + "?username=" + usrSearch.Username +
		"&email=" + usrSearch.Email +
		"&page=" + strconv.Itoa(usrSearch.Page) +
		"&page_size=" + strconv.Itoa(usrSearch.PageSize)

	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		End(utils.PrintStatus)
}

type userCurrent struct {
}

var usrCurrent userCurrent

func (x *userCurrent) Execute(args []string) error {
	GetUserCurrent(utils.URLGen("/api/users"))
	return nil
}

// GetUserCurrent gets the current user information.
//
// format:
//  GET /users/current
//
// e.g. curl -X GET --header 'Accept: application/json' 'https://localhost/api/users/current?api_key=top'
//
func GetUserCurrent(baseURL string) {
	targetURL := baseURL + "/current"
	fmt.Println("==> GET", targetURL)

	// Read beegosessionID from .cookie.yaml
	c, err := utils.CookieLoad()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	utils.Request.Get(targetURL).
		Set("Cookie", "harbor-lang=zh-cn; beegosessionID="+c.BeegosessionID).
		// NOTE:
		// 若后续需要根据用户权限做文章，则需要将用户信息进行维护
		// 可以定制一个新的回调函数
		End(utils.PrintStatus)
}

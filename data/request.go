package data

// MINPASSWORDLENTH to validate password
const MINPASSWORDLENTH = 3

/*
 * API Model
 */

// RequestLogin for api request to log in
type RequestLogin struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Superadmin bool   `json:"superadmin,oemitempty"`
	Active     bool   `json:"active,oemitempty"`
}

// ChangePasswordRequest for api request of a new password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentpassword"`
	NewPassword     string `json:"newpassword"`
}

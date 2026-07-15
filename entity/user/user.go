package user

// UserCreateReq request body for POST /demo/usr
type UserCreateReq struct {
	Realname string `json:"realname"` // optional
	Username string `json:"username"` // optional
	Phone    string `json:"phone"`    // required, unique among non-deleted records
	Age      *int   `json:"age"`      // optional, nil means not provided
}

// UserPathReq path param for /demo/usr/:id
type UserPathReq struct {
	ID int `param:"id" json:"-"`
}

// UserUpdateReq request body for PUT /demo/usr/:id
type UserUpdateReq struct {
	Realname *string `json:"realname"` // optional, nil means no update
	Username *string `json:"username"` // optional, nil means no update
	Age      *int    `json:"age"`      // optional, nil means no update
}

package sqliteusr

// TestUserCreateReq is the request body for POST /sqlite/testuser.
type TestUserCreateReq struct {
	Name  string `json:"name"`  // required
	Phone string `json:"phone"` // required
	Age   *int   `json:"age"`   // optional, nil means 0
}

// TestUserUpdateReq is the request body for PUT /sqlite/testuser/:id.
type TestUserUpdateReq struct {
	Name  *string `json:"name"`  // optional
	Phone *string `json:"phone"` // optional
	Age   *int    `json:"age"`   // optional
}

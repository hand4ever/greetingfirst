package demo

// Filter GET /search?tag=go&tag=web&tag=api
type Filter struct {
	Tags []string `query:"tag" json:"tag"`
}

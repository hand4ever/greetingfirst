package demo

// Filter GET /search?tag=go&tag=web&tag=api
type Filter struct {
	Tags []string `query:"tag" json:"tag"`
}

// Echo path param binding for /demo/err/debug/:str
type Echo struct {
	Str string `param:"str" json:"str"`
}

// Sha256Request query parameter binding for /demo/sha256
type Sha256Request struct {
	Text string `query:"text"`
}

// Sha256Response payload for SHA256 computation result
type Sha256Response struct {
	Input string `json:"input"`
	Hash  string `json:"hash"`
}

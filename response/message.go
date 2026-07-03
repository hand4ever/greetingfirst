package response

type ErrMsg struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
	Cost    string `json:"cost"`
}

type Code int

const (
	ErrCodeOk      Code = 0
	ErrCodeCustom  Code = 100001
	ErrCodeNetwork Code = 100002 // 网络
)

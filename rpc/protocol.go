package rpc

type Request struct {
	ID uint64 `json:"id"`
	Method string `json:"method"`
	Params map[string]any `json:"params"`
}
type Cancel struct{
	ID uint64  `json:"id"`
}
type Response struct {
	ID uint64 `json:"id"`
	Result any `json:"result"`
	Error  string `json:"error,omitempty"`
}
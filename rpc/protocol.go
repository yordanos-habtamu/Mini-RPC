package rpc

type Request struct {
	Method string `json:"method"`
	Params map[string]any `json:"params"`
}

type Response struct {
	Result any `json:"result"`
	Error  string `json:"error,omitempty"`
}
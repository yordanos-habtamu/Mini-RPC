package rpc

type Request struct {
	Method string `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type Response struct {
	Result interface{} `json:"result"`
	Error  string `json:"error,omitempty"`
}
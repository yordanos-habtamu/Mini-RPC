package rpc

type RPCType int
const (
	Unary RPCType = iota
	Streaming
	Bidirectional
)
type MessageType string

const (
	MsgCall MessageType ="call"
	MsgStream MessageType ="stream"
)



type Request struct {
	ID uint64 `json:"id"`
	Type MessageType `json:"type"`
	Method string `json:"method,omitempty"`
	Params map[string]any `json:"params,omitempty"`
	Data any `json:"data,omitempty"`
}
type Cancel struct{
	ID uint64  `json:"id"`
}
type Response struct {
	ID uint64 `json:"id"`
	Result any `json:"result"`
	Error  string `json:"error,omitempty"`
}

type Stream struct {
	In <- chan any
	Out chan<- any
}
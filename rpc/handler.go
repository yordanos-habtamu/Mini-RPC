package rpc

type HandlerFunc func(params map[string]any) (any,error)
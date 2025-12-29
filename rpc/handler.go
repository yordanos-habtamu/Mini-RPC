package rpc

import (
	"context"
	"sync"
)


type UnairyHandler func(ctx context.Context,params map[string]any)(any,error)

type StreamHandler func(ctx context.Context, params map[string]any,stream chan<- any)(error)

type Method struct {
	Kind RPCType 
	uniaryHandler UnairyHandler
	streamHanlder StreamHandler
}

type Server struct {
	mu sync.RWMutex
	methods map[string]Method
}
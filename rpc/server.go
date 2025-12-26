package rpc

import (
	"context"
	"fmt"
)


func NewServer() *Server{
	return &Server{
		handlers: make(map[string]Handler),
	}
}

func (s *Server) Register(name string,handler Handler){
	if _,exists := s.handlers[name];exists{
		panic(fmt.Sprint("method %v already registered", name))
	}
	s.handlers[name] = handler
}

func (s *Server) Call(ctx context.Context,method string,params map[string]any) (any,error){
	handler,ok:= s.handlers[method]
	if !ok{
		return nil,fmt.Errorf("unknown method %v",method)
	}
	return handler(ctx,params)
}
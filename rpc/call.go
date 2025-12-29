package rpc

import (
	"context"
	"fmt"
)

func(s *Server) Call(
	ctx context.Context,
	methodName string,
	params map[string]any,
)(Method,any,error){
	s.mu.RLock()
	m,ok := s.methods[methodName]
	s.mu.RUnlock()
	if !ok{
		return Method{},nil,fmt.Errorf("Unknown method")
	}
	if m.Kind == Unary{
		result, err := m.uniaryHandler(ctx,params)
		return m, result, err
	}
	 stream := make(chan any)
	 go func (){
		defer close(stream)
		_= m.streamHanlder(ctx, params,stream)
	 }()
	 return m,stream,nil

}
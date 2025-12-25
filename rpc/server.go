package rpc

import "fmt"
type Server struct {
	methods map[string]HandlerFunc
}

func NewServer() *Server{
	return &Server{
		methods: make(map[string]HandlerFunc),
	}
}

func (s *Server) Register(name string,handler HandlerFunc){
	if _,exists := s.methods[name];exists{
		panic(fmt.Sprint("method %v already registered", name))
	}
	s.methods[name] = handler
}

func (s *Server) Call(method string,params map[string]any) (any,error){
	handler,ok:= s.methods[method]
	if !ok{
		return nil,fmt.Errorf("unknown method %v",method)
	}
	return handler(params)
}
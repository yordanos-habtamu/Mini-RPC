 package rpc

 func (s *Server) RegisterUniary(name string, handler UnairyHandler){
 s.mu.Lock()
 defer s.mu.Unlock()
 s.methods[name] = Method{
	Kind : Unary,
	uniaryHandler:handler,
 }
 }
 func(s *Server) RegisterStreaming(name string,handler StreamHandler){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.methods[name] = Method{
		Kind:Streaming,
		streamHanlder: handler,
	}
 }

 
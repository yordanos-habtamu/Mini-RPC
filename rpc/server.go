package rpc

import (
	"sync"
)


func NewServer() *Server{
	return &Server{
	   mu: sync.RWMutex{},
	   methods: make(map[string]Method),
	}
}


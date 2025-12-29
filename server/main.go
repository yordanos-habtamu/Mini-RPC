package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/yordanos-habtamu/mini-rpc/rpc"
)

func main() {
	server := rpc.NewServer()
 server.RegisterUniary("Add", func(ctx context.Context, params map[string]any) (any, error) {
	a := int(params["a"].(float64))
	b := int(params["b"].(float64))
	return a + b, nil
})
server.RegisterStreaming("CountToN",
	func(ctx context.Context, params map[string]any, stream chan<- any) error {
		n := int(params["n"].(float64))

		for i := 1; i <= n; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case stream <- i:
				time.Sleep(500 * time.Millisecond)
			}
		}
		return nil
	})


	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is listening on port  9000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, server)
	}
}
func handleConnection(c net.Conn, server *rpc.Server) {
	defer c.Close()
	decoder := json.NewDecoder(c)
	encoder := json.NewEncoder(c)
	const maxInflight = 5
	sema := make(chan struct{}, maxInflight)
	var writeMu sync.Mutex
	var inflightMu sync.Mutex
	inflight := make(map[uint64]context.CancelFunc)

	defer func() {
		inflightMu.Lock()
		for _, cancel := range inflight {
			cancel()
		}
		inflightMu.Unlock()
	}()

	for {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return
		}

		var req rpc.Request
		if json.Unmarshal(raw, &req) == nil && req.Method != "" {

			select {
			case sema <- struct{}{}:
			default:
				writeMu.Lock()
				_ = encoder.Encode(
					rpc.Response{
						ID:    req.ID,
						Error: "server busy, try again later",
					})
				writeMu.Unlock()
				continue
			}
			ctx, cancel := context.WithCancel(context.Background())
			inflightMu.Lock()
			inflight[req.ID] = cancel
			inflightMu.Unlock()
			go func(req rpc.Request, ctx context.Context) {
				defer func() {
					inflightMu.Lock()
					delete(inflight, req.ID)
					inflightMu.Unlock()
					<-sema
				}()
				method,result, err := server.Call(ctx, req.Method, req.Params)
				select {
				case <-ctx.Done():
					return
				default:
				}
				switch method.Kind {
				case rpc.Unary:
					res := rpc.Response{ID:req.ID}
					if err != nil{
						res.Error = err.Error()
					}else{
						res.Result = result
					}
					writeMu.Lock()
					_= encoder.Encode(res)
					writeMu.Unlock()
				case rpc.Streaming:
					stream := result.(chan any)
					for  v := range stream{
						select {
						case <- ctx.Done():
							return
						default:
							writeMu.Lock()
							_=encoder.Encode(
								rpc.Response{
									ID:req.ID,
									Result:v,
								})
								writeMu.Unlock()
						}
					}
				}
			}(req, ctx)
			continue
		}
		var cancel rpc.Cancel
		if json.Unmarshal(raw, &cancel) == nil {
			inflightMu.Lock()
			if cfn := inflight[cancel.ID]; cfn != nil {
				cfn()
			}
			inflightMu.Unlock()
		}
	}
}

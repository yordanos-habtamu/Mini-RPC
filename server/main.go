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

func main(){
	server := rpc.NewServer()
	server.Register("Add", func(params map[string]any)(any,error){
		a :=int(params["a"].(float64))
		b:= int(params["b"].(float64))
		return a +b,nil
	})
	server.Register("Multiply", func(params map[string]any)(any,error){
		time.Sleep(3 * time.Second )
		a :=int(params["a"].(float64))
		b:= int(params["b"].(float64))
		return a * b,nil
	})
	server.Register("Substract", func(params map[string]any)(any,error){
		a := int(params["a"].(float64))
		b := int(params["b"].(float64))
		return a-b, nil
	})
	
    ln,err := net.Listen("tcp",":9000")
	if err != nil{
		panic(err)
	}
	
	fmt.Println("Server is listening on port  9000")
	for {
		conn, err := ln.Accept()
		if err !=nil{
			panic(err)
		}
		go handleConnection(conn,server)
	}
}
func handleConnection(c net.Conn, server *rpc.Server){
		defer c.Close()
		decoder := json.NewDecoder(c)
		encoder := json.NewEncoder(c)
		var writeMu sync.Mutex 
		inflight := make(map[uint64]context.CancelFunc)
		for {
			var raw json.RawMessage
			if err := decoder.Decode(&raw); err !=nil{
				return 
			}
		
			var req rpc.Request
			if  json.Unmarshal(raw,&req)==nil && req.Method !=""{
				ctx,cancel := context.WithCancel(context.Background())
				writeMu.Lock()
				inflight[req.ID] = cancel
				writeMu.Unlock()
				go func(req rpc.Request, ctx context.Context){
					defer func(){
						writeMu.Lock()
						delete(inflight,req.ID)
						writeMu.Unlock()
					}()
					result, err := server.Call(req.Method,req.Params)
					select{
					case <- ctx.Done():
						return
					default:
					}
				res := rpc.Response{ID:req.ID}
				if err != nil{
					res.Error = err.Error()
				}else{
					res.Result = result
				}
				writeMu.Lock()
				encoder.Encode(res)
				writeMu.Unlock()
			}(req,ctx)
			continue
			}
			var cancel rpc.Cancel
			if json.Unmarshal(raw,&cancel) == nil{
				writeMu.Lock()
				if cfn := inflight[cancel.ID]; cfn != nil{
					cfn()
				}
				writeMu.Unlock()
			}
		}
}

		
package main

import (
	"encoding/json"	
	"time"
	"sync"
	"fmt"
	"net"
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
		time.Sleep(30 * time.Second )
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
		for{
			var req rpc.Request
			if err := decoder.Decode(&req); err!=nil{
				return
			}

			go func(req rpc.Request){
				result, err := server.Call(req.Method,req.Params)
				res := rpc.Response{
					ID: req.ID,
				}
				if err != nil{
					res.Error = err.Error()
				}else{
					res.Result = result
				}

				writeMu.Lock()
				encoder.Encode(res)
				writeMu.Unlock()		
			}(req)
			
		}
	}

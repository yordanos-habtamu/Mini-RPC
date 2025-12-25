package main

import (
	"encoding/json"
	"net"
	"fmt"
	"github.com/yordanos-habtamu/mini-rpc/rpc"
	"sync/atomic"
)

var nextID uint64 =1

func main(){
	conn,err := net.Dial("tcp",":9000")
	if err != nil{
		panic(err)
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	call := func(method string, params map[string]any){
		id := atomic.AddUint64(&nextID,1)
		req := rpc.Request{
			ID:id,
			Method: method,
			Params:params,
		}
		encoder.Encode(req)
		var res rpc.Response
		decoder.Decode(&res)
		if res.Error != ""{
			fmt.Printf("RPC %d error: %s\n", res.ID,res.Error)
		}else{
			fmt.Printf("RPC %d result: %d ",res.ID,res.Result)
		}
	}
	call("Add",map[string]any{"a":2,"b":4})
	call("Substract",map[string]any{"a":2,"b":4})
	call("Multiply",map[string]any{"a":2,"b":4})
	
	}



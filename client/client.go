package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"github.com/yordanos-habtamu/mini-rpc/rpc"
)

var nextID uint64 =0

func main(){
	conn,err := net.Dial("tcp",":9000")
	if err != nil{
		panic(err)
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	pending := make(map[uint64]chan rpc.Response)
	var mu sync.Mutex
	go func(){
		for{
			var res rpc.Response
			if err := decoder.Decode(&res); err != nil{
				return
			}
			mu.Lock()
			ch:=pending[res.ID]
			delete(pending,res.ID)
			mu.Unlock()
			ch<-res
		}
	}()
	call := func(method string, params map[string]any){
		id := atomic.AddUint64(&nextID,1)
		resCh := make(chan rpc.Response,1)
		mu.Lock()
		pending[id] = resCh
		mu.Unlock()
		req := rpc.Request{
			ID:id,
			Method: method,
			Params:params,
		}
		encoder.Encode(req)
		res := <- resCh
		if res.Error != ""{
			fmt.Printf("RPC %v error : %v \n",id,res.Error)
		}else{
			fmt.Printf("RPC %v result: %v \n", id , res.Result)
		}
	}
	for i:=0; i < 10 ; i++ {
	call("Add",map[string]any{"a":2,"b":4+i})
	call("Substract",map[string]any{"a":2+i,"b":4})
	call("Multiply",map[string]any{"a":2-i,"b":4})
	}
	select {}
	}



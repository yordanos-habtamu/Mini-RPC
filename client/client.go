package main

import (
	"encoding/json"
	"time"
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
	call := func(method string, params map[string]any, timeout time.Duration){
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
		select {
		case res := <- resCh:
			if res.Error != ""{
				fmt.Printf("RPC %v error : %v \n",id,res.Error)
			}else{
				fmt.Printf("RPC %v result: %v \n", id , res.Result)
			}
	case <- time.After(timeout):
		fmt.Printf("RPC %v timed out \n", id)
		cancel := rpc.Cancel{ID:id}
		encoder.Encode(cancel)
		mu.Lock()
		delete(pending,id)
		mu.Unlock()
		}
	}
	call("Add",map[string]any{"a":2,"b":4},2 * time.Second)
	call("Substract",map[string]any{"a":2,"b":4},1  *time.Second)
	call("Multiply",map[string]any{"a":2,"b":4},2 *time.Second)

	select {}
	}



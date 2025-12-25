package main

import (
	"encoding/json"
	"net"
	"fmt"
	"github.com/yordanos-habtamu/mini-rpc/rpc"
)


func main(){
	conn,err := net.Dial("tcp",":9000")
	if err != nil{
		panic(err)
	}
	defer conn.Close()
	req := rpc.Request{
		Method:"Add",
		Params:map[string]any{
			"a":5,
			"b":4,
		},	
		}
    if err := json.NewEncoder(conn).Encode(&req); err != nil{
		panic(err)
	}
	var res rpc.Response 
	if err := json.NewDecoder(conn).Decode(&res); err != nil {
		panic(err)
	}
	 if res.Error != "" {
		fmt.Println("Error from server:", res.Error)
		return 
	 }
	fmt.Println("Result from server:", res.Result)
	
	}



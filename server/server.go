package server

import (
	"encoding/json"
	"fmt"
	"net"
	"github.com/yordanos-habtamu/mini-rpc/rpc"
)

func Add(params map[string]any) int {
	a := int(params["a"].(float64))
	b := int(params["b"].(float64))
	return a + b
}


func main(){
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
		go func(c net.Conn){
			defer c.Close()
            var req rpc.Request
			if err := json.NewDecoder(c).Decode(&req); err != nil{
				return
			}
			var res rpc.Response
			switch req.Method{
			case "Add":
				result := Add(req.Params)
				res.Result = result
			default:
				res.Error = "unknown method"
			}
			json.NewEncoder(c).Encode(&res)
		}(conn)
	}

}
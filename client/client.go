package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"github.com/yordanos-habtamu/mini-rpc/rpc"
)

var nextID uint64 = 0

func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	pending := make(map[uint64]chan rpc.Response)
	var mu sync.Mutex

	go func() {
		for {
			var res rpc.Response
			if err := decoder.Decode(&res); err != nil {
				return
			}
			mu.Lock()
			ch := pending[res.ID]
			mu.Unlock()
			if ch != nil {
				select {
				case ch <- res:
				default:
					// Channel full or closed, skip
				}
			}
		}
	}()
	call := func(method string, params map[string]any, timeout time.Duration) {
		id := atomic.AddUint64(&nextID, 1)
		resCh := make(chan rpc.Response, 100) 
		mu.Lock()
		pending[id] = resCh
		mu.Unlock()
		req := rpc.Request{
			ID:     id,
			Method: method,
			Params: params,
		}
		encoder.Encode(req)
		
		// Loop to receive multiple streaming responses
		timer := time.NewTimer(timeout)
		for {
			select {
			case res := <-resCh:
				if res.Error != "" {
					fmt.Printf("RPC %v error : %v \n", id, res.Error)
					mu.Lock()
					delete(pending, id)
					mu.Unlock()
					return
				}
				fmt.Printf("RPC %v result: %v \n", id, res.Result)
				timer.Reset(timeout) // Reset timeout on each response
			case <-timer.C:
				fmt.Printf("RPC %v timed out or stream ended\n", id)
				cancel := rpc.Cancel{ID: id}
				encoder.Encode(cancel)
				mu.Lock()
				delete(pending, id)
				mu.Unlock()
				return
			}
		}
	}
	// call("Add",map[string]any{"a":2,"b":4},2 * time.Second)
	// call("Substract",map[string]any{"a":2,"b":4},1  *time.Second)
	// call("Multiply",map[string]any{"a":2,"b":4},2 *time.Second)
	call("countToN", map[string]any{"n": 5}, 10*time.Second)

	select {}
}

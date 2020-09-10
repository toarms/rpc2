package main

import (
		"github.com/toarms/rpc2/rpc"
	   )

func main() {
	var e rpc.Echo
	rpc.HandlesFunc("0x01", e)
	rpc.ListenAndServe(":7890")
}

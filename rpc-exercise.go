package main

import (
	//"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Result struct {
	Result int
}

type Arith struct {
	foo int
}

func (this *Arith) Add(arg *Args, ret *Result) error {
	ret.Result = arg.A + arg.B
	return nil
}

func InitServer() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen:", err)
		return
	}
	go http.Serve(listen, nil)
}

func ClientConnect() (*rpc.Client, error) {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("Dial", err)
		return nil, err
	}
	return client, err
}

func RpcClientAdd(A, B int) (int, error) {
	client, err := ClientConnect()
	if err != nil {
		return 0, err
	}
	args := Args{A, B}
	ret := Result{}
	err = client.Call("Arith.Add", &args, &ret)
	if err != nil {
		return 0, err
	}

	client.Close()
	return ret.Result, nil
}

func main() {
	InitServer()

	ret, err := RpcClientAdd(1, 2)
	if err != nil {
		fmt.Println("error", err)
	} else {
		fmt.Println("result", ret)
	}
}

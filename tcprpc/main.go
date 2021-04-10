package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

const (
	DefaultPort = 6060
)

type NilArgs struct{}

type ServerOp struct{}

func (o *ServerOp) Hostname(args *NilArgs, reply *string) error {
	var err error
	var h string
	h, err = os.Hostname()
	*reply = h
	if err != nil {
		return err
	}
	return nil
}

func (o *ServerOp) Ping(args *NilArgs, reply *int) error {
	*reply = 1
	return nil
}

func main() {
	var args = os.Args
	if len(args) < 2 {
		fmt.Printf("main: invalid number of commandline arguments \n")
		os.Exit(1)
	}

	var cmd = args[1]
	if cmd == "start" && len(args) >= 2 {
		start(args[2:])
		return
	} else if cmd == "send" && len(args) >= 2 {
		send(args[2:])
		return
	} else {
		fmt.Printf("[main] invalid commandline arguments \n")
		os.Exit(1)
	}
}

func start(args []string) {
	var err error

	var port int
	flag.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flag.Parse()

	fmt.Printf("[start] port: %d \n", port)

	// creates instance of the RPC servers
	var rpcHandler = rpc.NewServer()
	// registers the operations for the RPC
	rpcHandler.Register(new(ServerOp))

	// creates instance of the TCP listener
	var l net.Listener
	l, err = net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}

	// listener to start accepting RPC calls
	fmt.Printf("[start] service started \n")
	rpcHandler.Accept(l)
}

func send(args []string) {
	var err error

	var port int
	flag.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flag.Parse()

	// client connects to RPC server
	var c *rpc.Client
	c, err = rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}

	var start = time.Now()

	// Calls methos from the registered operations in the RPC server
	var r int
	err = c.Call("ServerOp.Ping", new(NilArgs), &r)
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("[send] ping: %d \n", r)
	fmt.Printf("[send] elapsed: %v \n", time.Since(start))
}

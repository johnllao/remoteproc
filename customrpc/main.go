package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/johnllao/remoteproc/pkg/hmac"
	"github.com/johnllao/remoteproc/pkg/security"
)

const (
	DefaultPort = 6060
)

type NilArgs struct{}

type ServerOp struct{}

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
	if cmd == "server" && len(args) >= 2 {
		server(args[2:])
		return
	} else if cmd == "client" && len(args) >= 2 {
		client(args[2:])
		return
	} else if cmd == "token" && len(args) >= 2 {
		gentoken(args[2:])
		return
	} else {
		fmt.Printf("[main] invalid commandline arguments \n")
		os.Exit(1)
	}
}

func server(args []string) {
	var err error

	var port int
	var key string
	var flagset = flag.NewFlagSet("server", flag.ContinueOnError)
	flagset.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flagset.StringVar(&key, "key", "secret", "secret hmac key")
	flagset.Parse(args)

	fmt.Printf("[server] port: %d \n", port)

	// creates instance of the TCP listener
	var l net.Listener
	l, err = net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("[server] %s \n", err.Error())
		os.Exit(1)
	}

	// listener to start accepting RPC calls
	fmt.Printf("[server] service started \n")
	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			fmt.Printf("[server] %s \n", err.Error())
			continue
		}

		go func(conn net.Conn) {
			// security handshake with client connection
			var connErr = security.ServerHandshake(conn, conn, key)
			if connErr != nil {
				fmt.Printf("[server] %s \n", connErr.Error())
				_ = conn.Close()
				return
			}

			// creates instance of the RPC servers
			var rpcHandler = rpc.NewServer()
			// registers the operations for the RPC
			connErr = rpcHandler.Register(new(ServerOp))
			if connErr != nil {
				fmt.Printf("[server] %s \n", connErr.Error())
				return
			}

			rpcHandler.ServeConn(conn)
		}(conn)
	}
}

func client(args []string) {
	var err error

	var port int
	var token string
	var flagset = flag.NewFlagSet("client", flag.ContinueOnError)
	flagset.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flagset.StringVar(&token, "token", "", "authentication token")
	flagset.Parse(args)

	// client connects to RPC server
	var conn net.Conn
	conn, err = net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("[client] %s \n", err.Error())
		os.Exit(1)
	}

	err = security.ClientHandshake(conn, conn, token)
	if err != nil {
		fmt.Printf("[client] %s \n", err.Error())
		os.Exit(1)
	}

	var c = rpc.NewClient(conn)

	// Calls methos from the registered operations in the RPC server
	for i := 0; i < 5; i++ {
		var start = time.Now()

		var reply int
		err = c.Call("ServerOp.Ping", new(NilArgs), &reply)
		if err != nil {
			fmt.Printf("[client] %s \n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("[client] ping: %d, elapsed: %v \n", reply, time.Since(start))
	}
}

func gentoken(args []string) {
	var err error

	var key, name string
	var exp int
	var flagset = flag.NewFlagSet("gentoken", flag.ContinueOnError)
	flagset.StringVar(&key, "key", "secret", "secret hmac key")
	flagset.StringVar(&name, "name", "", "name of the app")
	flagset.IntVar(&exp, "exp", 24, "expiry in hours")
	flagset.Parse(args)

	var token string
	token, err = hmac.Token(key, &hmac.Claim{
		Name:   name,
		Expiry: time.Now().Add(time.Duration(exp) * time.Hour).UnixNano(),
	})
	if err != nil {
		fmt.Printf("[gentoken] %s \n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("TOKEN:%s\n", token)
}

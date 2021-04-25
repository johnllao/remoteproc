package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/johnllao/remoteproc/customrpc/hmac"
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
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("[server] %s \n", err.Error())
			continue
		}

		// start the handshake
		var tokenRequest string
		var tokenResponse = "OK\n"

		var r = bufio.NewReader(conn)
		tokenRequest, err = r.ReadString('\n')
		if err != nil {
			fmt.Printf("[server] %s \n", err.Error())
			continue
		}
		if strings.HasPrefix(tokenRequest, "TOKEN:") {
			var token = tokenRequest[6 : len(tokenRequest)-1]

			var validToken bool
			if validToken, err = hmac.VerifyToken(key, token); !validToken || err != nil {
				tokenResponse = "INVALID_TOKEN\n"
			}

		} else {
			tokenResponse = "INVALID_TOKEN\n"
			continue
		}

		var w = bufio.NewWriter(conn)
		_, err = w.WriteString(tokenResponse)
		if err != nil {
			fmt.Printf("[server] %s \n", err.Error())
			continue
		}
		err = w.Flush()
		if err != nil {
			fmt.Printf("[server] %s \n", err.Error())
			continue
		}
		// end of the handshake

		// creates instance of the RPC servers
		var rpcHandler = rpc.NewServer()
		// registers the operations for the RPC
		rpcHandler.Register(new(ServerOp))

		go rpcHandler.ServeConn(conn)
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

	// start the handshake
	var tokenRequest = "TOKEN:" + token + "\n"
	var w = bufio.NewWriter(conn)
	_, err = w.WriteString(tokenRequest)
	if err != nil {
		fmt.Printf("[client] %s \n", err.Error())
		os.Exit(1)
	}
	err = w.Flush()
	if err != nil {
		fmt.Printf("[client] %s \n", err.Error())
		os.Exit(1)
	}
	var tokenReply string
	var r = bufio.NewReader(conn)
	tokenReply, err = r.ReadString('\n')
	if err != nil {
		fmt.Printf("[client] %s \n", err.Error())
		os.Exit(1)
	}
	if tokenReply[:2] != "OK" {
		fmt.Printf("[client] not authorized \n")
		os.Exit(1)
	}
	// end of the handshake

	var c = rpc.NewClient(conn)

	// Calls methos from the registered operations in the RPC server
	for i := 0; i < 100; i++ {
		var start = time.Now()

		var reply int
		err = c.Call("ServerOp.Ping", new(NilArgs), &reply)
		if err != nil {
			fmt.Printf("[client] %s \n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("[client] ping: %d \n", reply)
		fmt.Printf("[client] elapsed: %v \n", time.Since(start))
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

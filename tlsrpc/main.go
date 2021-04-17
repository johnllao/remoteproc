package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/johnllao/remoteproc/tlsrpc/ops"
)

const (
	DefaultPort = 6060
)

func main() {
	var args = os.Args
	if len(args) < 2 {
		fmt.Printf("[main] invalid number of commandline arguments \n")
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
	var flagset = flag.NewFlagSet("start", flag.ContinueOnError)
	flagset.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	err = flagset.Parse(args)
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("[start] port: %d \n", port)

	var certificate, privateKey, caBytes []byte
	certificate, err = os.ReadFile("./cert1.pem")
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}
	privateKey, err = os.ReadFile("./cert1.key")
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}
	caBytes, err = os.ReadFile("./ca.pem")
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}

	// configure TLS settings here
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certificate, privateKey)
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}
	var ca = x509.NewCertPool()
	ca.AppendCertsFromPEM(caBytes)

	// configure TLS settings here
	var cfg = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	var l net.Listener
	l, err = tls.Listen("tcp", ":"+strconv.Itoa(port), cfg)
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}

	// creates instance of the RPC servers
	var rpcHandler = rpc.NewServer()
	// registers the operations for the RPC
	rpcHandler.Register(new(ops.ServerOp))

	// listener to start accepting RPC calls
	fmt.Printf("[start] service started \n")

	// alternative to rpcHandler.Accept(l)
	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			fmt.Printf("[start] %s \n", err.Error())
			continue
		}
		if c, ok := conn.(*tls.Conn); ok {
			var remoteAddr = c.RemoteAddr().String()
			fmt.Printf("[start] received from %s \n", remoteAddr)
		}
		go rpcHandler.ServeConn(conn)
	}
}

func send(args []string) {
	var err error

	var port int
	var certName string
	var flagset = flag.NewFlagSet("send", flag.ContinueOnError)
	flagset.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flagset.StringVar(&certName, "cert", "cert1", "name of the certificate file")
	err = flagset.Parse(args)
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("[send] loading certificate %s \n", certName)
	var certificate, privateKey, caBytes []byte
	certificate, err = os.ReadFile("./" + certName + ".pem")
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	privateKey, err = os.ReadFile("./" + certName + ".key")
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	caBytes, err = os.ReadFile("./ca.pem")
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}

	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certificate, privateKey)
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	var ca = x509.NewCertPool()
	ca.AppendCertsFromPEM(caBytes)

	// configure TLS settings here
	var cfg = &tls.Config{
		ServerName:   "localhost",
		RootCAs:      ca,
		Certificates: []tls.Certificate{cert},
	}

	// create a TLS connection with TLS config
	var conn *tls.Conn
	conn, err = tls.Dial("tcp", ":"+strconv.Itoa(port), cfg)
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	// create RPC client with the TLS config
	var c = rpc.NewClient(conn)

	var start = time.Now()

	// Calls methos from the registered operations in the RPC server
	var r int
	err = c.Call("ServerOp.Ping", new(ops.NilArgs), &r)
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("[send] ping: %d \n", r)
	fmt.Printf("[send] elapsed: %v \n", time.Since(start))
}

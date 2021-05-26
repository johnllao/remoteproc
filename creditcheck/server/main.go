package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/johnllao/remoteproc/creditcheck/ops"
	"github.com/johnllao/remoteproc/pkg/security"
)

var (
	dbpath  string
	keypath string
	port    int

	key string
)

func main() {
	var err error

	flag.StringVar(&dbpath, "dbpath", "./cchk.db", "path of the database file")
	flag.StringVar(&keypath, "keypath", "./cchk.key", "path of the security key")
	flag.IntVar(&port, "port", 6060, "service port number")
	flag.Parse()

	if dbpath == "" {
		log.Fatalf("ERR: main() missing dbpath from the argument")
	}
	if keypath == "" {
		log.Fatalf("ERR: main() missing key path from the argument")
	}

	log.Printf("reading the key file. path: %s", keypath)
	var keyb []byte
	keyb, err = os.ReadFile(keypath)
	if err != nil {
		log.Fatalf("ERR: main() %s", err.Error())
	}
	key = string(bytes.TrimSpace(keyb))

	log.Printf("loading database file. path: %s", dbpath)
	var db *bolt.DB
	db, err = bolt.Open(dbpath, 0600, nil)
	if err != nil {
		log.Fatalf("ERR: main() %s", err.Error())
	}
	defer db.Close()

	// handle any process interrupts
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Close()
		os.Exit(1)
	}()

	var custOp = ops.NewCustomerOps(db)

	var l net.Listener
	l, err = net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("ERR: main() %s", err.Error())
	}
	defer l.Close()

	log.Printf("service started. port: %d \n", port)
	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			log.Printf("WARN: main() %s", err.Error())
			continue
		}

		go func(conn net.Conn) {
			// security handshake with client connection
			var connErr = security.ServerHandshake(conn, conn, key)
			if connErr != nil {
				fmt.Printf("WARN: main() %s \n", connErr.Error())
				_ = conn.Close()
				return
			}

			// creates instance of the RPC servers
			var rpcHandler = rpc.NewServer()
			// registers the operations for the RPC
			connErr = rpcHandler.Register(custOp)
			if connErr != nil {
				fmt.Printf("WARN: main() %s \n", connErr.Error())
				return
			}
			rpcHandler.ServeConn(conn)
		}(conn)
	}
}

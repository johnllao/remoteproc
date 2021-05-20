package main

import (
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
	dbpath string
	key    string
	port   int
)

func main() {
	var err error
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

	var l net.Listener
	l, err = net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("ERR: main() %s", err.Error())
	}
	defer l.Close()

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
			connErr = rpcHandler.Register(new(ops.CustomerOp))
			if connErr != nil {
				fmt.Printf("WARN: main() %s \n", connErr.Error())
				return
			}
			rpcHandler.ServeConn(conn)
		}(conn)
	}
}

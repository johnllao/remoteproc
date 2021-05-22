package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/johnllao/remoteproc/creditcheck/arguments"
	"github.com/johnllao/remoteproc/pkg/hmac"
	"github.com/johnllao/remoteproc/pkg/security"
)

func main() {
	var cmd = os.Args[1]

	if cmd == "gentoken" {
		gentoken(os.Args[2:])
	} else if cmd == "testmsg" {
		testmsg(os.Args[2:])
	} else {
		fmt.Printf("ERR: main() invalid command argument")
	}
}

func gentoken(args []string) {
	var err error

	var key, name string
	var expiry time.Duration

	var flagset = flag.NewFlagSet("gentoken", flag.ContinueOnError)
	flagset.StringVar(&key, "key", "", "secret key")
	flagset.StringVar(&name, "name", "", "name of the token")
	flagset.DurationVar(&expiry, "expiry", 5, "no of days to expire")
	err = flagset.Parse(args)
	if err != nil {
		fmt.Printf("ERR: gentoken() %s \n", err.Error())
		return
	}

	var claim = &hmac.Claim{
		Name:   name,
		Expiry: time.Now().Add(expiry * 24 * time.Hour).UnixNano(),
	}

	var token string
	token, err = hmac.Token(key, claim)
	if err != nil {
		fmt.Printf("ERR: gentoken() %s \n", err.Error())
		return
	}

	var valid bool
	valid, err = hmac.VerifyToken(key, token)
	if err != nil {
		fmt.Printf("ERR: gentoken() %s \n", err.Error())
		return
	}

	if valid {
		fmt.Println("TOKEN:")
		fmt.Println(token)
	}
}

func testmsg(args []string) {
	var token = args[0]

	var err error

	var conn net.Conn
	conn, err = net.Dial("tcp", "localhost:6060")
	if err != nil {
		fmt.Printf("ERR: testmsg() %s \n", err.Error())
		return
	}
	defer conn.Close()

	var cli = rpc.NewClient(conn)
	defer cli.Close()

	err = security.ClientHandshake(conn, conn, token)
	if err != nil {
		fmt.Printf("ERR: testmsg() %s \n", err.Error())
		return
	}

	var a arguments.LoadFileArg
	var r int

	a.Path = "/Users/johnllao/Downloads/nasdaq_comapnies.csv"
	err = cli.Call("CustomerOp.LoadFromFile", &a, &r)
	if err != nil {
		fmt.Printf("ERR: testmsg() %s \n", err.Error())
		return
	}

	var aa arguments.NilArgs
	var rr arguments.CompaniesReply
	err = cli.Call("CustomerOp.Companies", &aa, &rr)
	if err != nil {
		fmt.Printf("ERR: testmsg() %s \n", err.Error())
		return
	}
	for _, co := range rr.Companies {
		fmt.Println(co.Symbol + " " + co.Name)
	}
	fmt.Println("bye!")
}

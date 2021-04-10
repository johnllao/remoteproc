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
)

const (
	DefaultPort = 6060
)

const (
	Certificate = `-----BEGIN CERTIFICATE-----
MIIC5jCCAc4CCQD1whgh8yEwiTANBgkqhkiG9w0BAQsFADA1MQswCQYDVQQGEwJT
RzESMBAGA1UEBwwJU2luZ2Fwb3JlMRIwEAYDVQQDDAlsb2NhbGhvc3QwHhcNMjEw
NDEwMTAyNDU2WhcNMjIwNDEwMTAyNDU2WjA1MQswCQYDVQQGEwJTRzESMBAGA1UE
BwwJU2luZ2Fwb3JlMRIwEAYDVQQDDAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDRvlS1o4W9+eiiT8/MX2PvHPcxST0z6XFt5YJzubuf
WDH9b0r/aJVu9qPi8DQT4R3XkaAJPRXrvqUlunCBzJL9sUgX/N9VHxn/BLL+/MWP
7aOHdW5KcQ1IrDTFMviCSfvFj48B6e81HTyA55B9GgFFPCom+aU7MZIASMeNAZJb
DTus7puUsTubSVRCtmCO4xLrSlOakZorcgVPg3zN7eRf/6enAWclWLTN4NA3N9hQ
eo8DhiWhzqVkoBKULXHLzAdbxi5cgLaRI+MEpucmYTqr9rnhU+YxHysZdvSbG803
qNR2pMCr3jYpFhYcz2st50Y8AiXoq9tO89lN4jbIN73rAgMBAAEwDQYJKoZIhvcN
AQELBQADggEBAArvH4LKbeFG9oZouvFKU65mkHlJTHZtSU8BkA9BujRYSkwAAuFM
B7S4fn9M5V9mTkOmZyKx1SefzW83q/G3DKMLL1rneLbsF0x6dcumRtvclePQbwra
jLy1B7gcQJB55BFlKyA615lFiEUW3CO2DDxVo26krtbCqW8j99vZVnMPP1lWeYJd
O0Y6/i4khmm9eseKmOudqfw1Tui6vKzs37e0OLONJMVSqB3sb7q4m/Sq3jLzA6bM
ujjcX3tHj/1z3lHCRjpcKvbav/zgooBHV1tGjQEqP9wwjUXiHjHO1HcWb/anL/nT
OCFxJAv+B2F5J9Al7ysGw0xgCZTY3uLPzww=
-----END CERTIFICATE-----`

	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0b5UtaOFvfnook/PzF9j7xz3MUk9M+lxbeWCc7m7n1gx/W9K
/2iVbvaj4vA0E+Ed15GgCT0V676lJbpwgcyS/bFIF/zfVR8Z/wSy/vzFj+2jh3Vu
SnENSKw0xTL4gkn7xY+PAenvNR08gOeQfRoBRTwqJvmlOzGSAEjHjQGSWw07rO6b
lLE7m0lUQrZgjuMS60pTmpGaK3IFT4N8ze3kX/+npwFnJVi0zeDQNzfYUHqPA4Yl
oc6lZKASlC1xy8wHW8YuXIC2kSPjBKbnJmE6q/a54VPmMR8rGXb0mxvNN6jUdqTA
q942KRYWHM9rLedGPAIl6KvbTvPZTeI2yDe96wIDAQABAoIBABTyTiFpuESVva7g
C5+ccy2BFgq9Bri1epeAETwfc2Zhd3SY9cN2HV5ckVdrp7fIhqNtrq7dg9/sRS/0
Y6IC3Tcqykli/qbQmVcHkBy4/730/JzdlGsoySvVzttW0MyqONOtF5oYU8RZLB6v
gZWM0E2qyYbk7aRwueT/X7ZsTsJ67IIYkOSr8igLinSjvkLuxP7ugWn5F2iBYwIP
Dr0H9SWZqKnJ8JZ21e6BYYEqwDSCsm4a/i/9TzPCTv3wopIy8LtG+Ci5v1zVzVWe
JH3/NV0VFZ5iaJ6Ezlk5xjY2VD1DxquRcaNYLa7yFJXxqRH2TfPciR5LikSP0JC1
i5JE9+ECgYEA8uZNuPGiGfDpzi1sdC2JJLTJ1di3UmiVxUX5QWXh5HCODlqExDnS
egdxTQwH308rv7nWUPZkMahj3JuWkUvaTUTn3eEoy5DDc1usIl9/s3yzs8nXPVUm
nTu8hMsBe+plsAkDrHd1k3vnXIjpL6ipXE0hjWjSgSJlcJ6+ldfpjVsCgYEA3Q4+
M5fg++HFKThJsUq+Bo2qwoApKHJfQuoM/1z9UznYiT6d8lgqff+bl+CqYmgaUj+o
OzAICMGc6PI6eFidZnSYysN5FUK/FcKUmVrmPyoo1EgQtI6gUY0mI56kSyNA4FFO
bGAlKbtP1WL/V9CFdNHf+H3KE2nEJXdRxaJ8prECgYEAz6Ymp4aaR4b2ubWHU8Jh
zaloKpJ8Fc0mzGDHdyr78+hs6MRlX8L2ti+Koo04ZaUvB1Z9avVYLkOAK2YvT8MC
uq+/cKU91NjK3eFuxGvTpcNjdL2Gbf5PZndc8EED4cU+bUEnjNcLAqwX27mHb6DG
OAwQNO15l7+p7J8o2rycAqUCgYBae9wGLmMPd2jG6J1xjtCdyhtdpiwyvC42K6vK
U3v2NzVlaFYqvuAV1y0PTA0yXr53cEsifxSq0OWzjINWg59aMtvgE4dapomlFJLS
+xxIOq+fxSfhYIhLGWXFKsjBYNrLdzyMrAZKQLv68pzmixo1qTrucj7nF2IMm/zC
0zIG4QKBgGPf7eqytpFOErQhciJxuANOxAdm/OjXkoSvOo5+qc48c8ajIARTw1L1
GDniMj+yCDJocP0ABWLVbun6kRhYuCMY22P83TQeoqrRlwdA5geO7N0mryST7G8q
ZzI2HN28nR3OIgnKWtPh++7JE4HukwuzVc/0ugIlgehUOhhhJ99L
-----END RSA PRIVATE KEY-----`
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
	flag.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flag.Parse()

	fmt.Printf("[start] port: %d \n", port)

	// creates instance of the RPC servers
	var rpcHandler = rpc.NewServer()
	// registers the operations for the RPC
	rpcHandler.Register(new(ServerOp))

	// configure TLS settings here
	var cert tls.Certificate
	cert, err = tls.X509KeyPair([]byte(Certificate), []byte(PrivateKey))
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}
	var ca = x509.NewCertPool()
	ca.AppendCertsFromPEM([]byte(Certificate))

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

	// listener to start accepting RPC calls
	fmt.Printf("[start] service started \n")
	rpcHandler.Accept(l)
}

func send(args []string) {
	var err error

	var port int
	flag.IntVar(&port, "port", DefaultPort, "port number of the RPC service")
	flag.Parse()

	var cert tls.Certificate
	cert, err = tls.X509KeyPair([]byte(Certificate), []byte(PrivateKey))
	if err != nil {
		fmt.Printf("[start] %s \n", err.Error())
		os.Exit(1)
	}
	var ca = x509.NewCertPool()
	ca.AppendCertsFromPEM([]byte(Certificate))

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
	err = c.Call("ServerOp.Ping", new(NilArgs), &r)
	if err != nil {
		fmt.Printf("[send] %s \n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("[send] ping: %d \n", r)
	fmt.Printf("[send] elapsed: %v \n", time.Since(start))
}

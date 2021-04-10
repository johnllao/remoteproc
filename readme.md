# Golang RPC Examples

The following examples demonstrates simple RPC implementations using TCP connections.

## TCP RPC

Simple simplementation of RPC using TCP connection

### Server

From the server side, operations are registered using the Register method of the RPC package

    var rpcHandler = rpc.NewServer()
    rpcHandler.Register(new(ServerOp))

An instance of a TCP listener is then created and pass it to the Accept method to start the service

    l, _ := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
    rpcHandler.Accept(l)

### Client

Connect to servers using the rpc.Dial method

    c, _ = rpc.Dial("tcp", "localhost:"+strconv.Itoa(port))

Perform the standard call to RPC

    var r int
	err = c.Call("ServerOp.Ping", new(NilArgs), &r)


## TLS RPC

The TLS implementation shares the same concept as the simple TCP implementation, except  additional code to setup the certificates are needed

### Server

Setup the TLS certificates for the server

    cert, _ := tls.X509KeyPair([]byte(Certificate), []byte(PrivateKey))

	var ca = x509.NewCertPool()
	ca.AppendCertsFromPEM([]byte(Certificate))

	// configure TLS settings here
	var cfg = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

The tls.Listen will then be invoked to pass the TLS configurations

    l, _ := tls.Listen("tcp", ":"+strconv.Itoa(port), cfg)
	rpcHandler.Accept(l)

### Client

Setup the certificates for the client

    cert, _ := tls.X509KeyPair([]byte(Certificate), []byte(PrivateKey))

	var ca = x509.NewCertPool()
	ca.AppendCertsFromPEM([]byte(Certificate))

	// configure TLS settings here
	var cfg = &tls.Config{
		ServerName:   "localhost",
		RootCAs:      ca,
		Certificates: []tls.Certificate{cert},
	}

Use the tls.Dial to connect to the server and pass the connection to the rpc.NewClient method

    // create a TLS connection with TLS config
	conn, _ := tls.Dial("tcp", ":"+strconv.Itoa(port), cfg)

	// create RPC client with the TLS config
	var c = rpc.NewClient(conn)

Perform the standard call to RPC

	// Calls methos from the registered operations in the RPC server
	var r int
	c.Call("ServerOp.Ping", new(NilArgs), &r)
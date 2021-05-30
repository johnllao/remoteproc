package client

import (
	"net"
	"net/rpc"

	"github.com/johnllao/remoteproc/pkg/security"
)

type Client struct {
	Addr  string
	Token string

	conn net.Conn
	cli  *rpc.Client
}

func (c *Client) Connect() error {
	var err error

	var conn net.Conn
	conn, err = net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	c.conn = conn

	var cli = rpc.NewClient(conn)
	err = security.ClientHandshake(conn, conn, c.Token)
	if err != nil {
		return err
	}
	c.cli = cli

	return nil
}

func (c *Client) Call(method string, args, reply interface{}) error {
	return c.cli.Call(method, args, reply)
}

func (c *Client) Close() error {
	if err := c.cli.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

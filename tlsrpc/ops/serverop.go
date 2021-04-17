package ops

import (
	"os"
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

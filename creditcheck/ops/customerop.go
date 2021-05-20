package ops

type NilArgs struct{}

type CustomerOp struct{}

func (o *CustomerOp) Ping(args *NilArgs, reply *int) error {
	*reply = 1
	return nil
}

package proto

type ConnArg struct {
	Body   []byte
	Server int32
}

type ConnReply struct {
	Key    string
	RoomId int32
}

type SendArg struct {
	Body []byte
}

type SendReply struct {
}

type DisconnArg struct {
	Key    string
	RoomId int32
}

type DisconnReply struct {
	Has bool
}

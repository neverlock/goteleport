package goteleport

import "sync"

type Teleporter struct{
	in        chan interface{}
	out       chan interface{}
	port      int
	lock      sync.Mutex
	client    []string
	master    string
}

const(
	DATA	=	1
	PING	=	2
)

type Message struct {
	MType   uint
	Payload []byte
}

func New(port, size int) (chan interface{}, chan interface{}){
	in := make(chan interface{}, size)
	out := make(chan interface{}, size)
	t := Teleporter{
		in:      	in,
		out:		out,
		port:    	port,
	}

	go t.serverListenForOutboundMessageBuffer()
	go t.serverListenForInboundMessageBuffer()
	return in, out
}

func Connect(address string, port, size int) (chan interface{}, chan interface{}){
	in := make(chan interface{}, size)
	out := make(chan interface{}, size)
	t := Teleporter{
		in:     in,
		out:    out,
		port:   port,
		client: []string{},
		master: address,
	}

	t.clientConnectToMaster()
	go t.clientListenForInboundMessageBuffer()
	go t.clientListenForOutboundMessageBuffer()
	return in, out
}
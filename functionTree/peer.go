package functionTree

import "sync"

type Peer struct {
	mu         sync.Mutex
	id         int
	parent     *Peer
	leftChild  *Peer
	rightChild *Peer
	height     int
	addr       string // VM ip addr
	port       string // VM port
	ready      bool // is ready to be the parent
}

func NewPeer() *Peer {
	return &Peer{
		leftChild:  nil,
		rightChild: nil,
		ready:      false,
		addr:       "",
		port:       "",
	}
}

func (n *Peer) SetAddr(addr string) {
	n.addr = addr
}

func (n *Peer) GetAddr() string {
	return n.addr
}

func (n *Peer) GetHeight() int {
	if n == nil {
		return -1
	}
	return n.height
}

func (n *Peer) GetBalanceFactor() int {
	if n == nil {
		return 0
	}
	return n.leftChild.GetHeight() - n.rightChild.GetHeight()
}

func (n *Peer) GetId() int {
	return n.id
}

func (n *Peer) GetParent() *Peer {
	return n.parent
}

func (n *Peer) SetReady() {
	n.ready = true
}

func (n *Peer) IsReady() bool {
	return n.ready
}

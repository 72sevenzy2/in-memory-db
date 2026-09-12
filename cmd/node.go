package cmd

import (
	"fmt"
	"net"
	"sync"

	"github.com/72sevenzy2/in-memory-database"
)

// a node represents an database instance.
type Node struct {
	lock sync.RWMutex

	ID   string
	Addr string
	role Role

	// Replicas represent a string array of existing nodes addresses.
	Replicas []string

	DB *db.DB
}

func NewNode(id, addr string, role Role, replicas []string) *Node {
	return &Node{
		ID:       id,
		Addr:     addr,
		role:     role,
		Replicas: replicas,
		DB:       db.NewDB(),
	}
}

func (n *Node) Start() error {
	l, err := net.Listen("tcp", n.Addr)
	if err != nil {
		return err
	}

	defer l.Close()

	fmt.Println("node listening on addr:", n.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go n.HandleConnection(conn)
	}
}

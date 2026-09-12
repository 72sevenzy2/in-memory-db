package cmd

import (
	"sync"

	"github.com/72sevenzy2/in-memory-database/db"
)

// a node represents an database instance.

type Role uint8 // small integer which dictates a nodes role.

type Node struct {
	lock sync.RWMutex

	ID   string
	Addr string
	role Role

	// Replicas represent a string array of existing nodes addresses.
	Replicas []string

	DB *db.DB
}

// Represents a nodes role, either an follower or leader.
const (
	Leader Role = iota
	Follower
)

func NewNode(id, addr string, role Role) *Node {
	return &Node{
		ID:   id,
		Addr: addr,
		role: role,
		DB:   db.NewDB(),
	}
}

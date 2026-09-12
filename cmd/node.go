package cmd

import (
	"github.com/72sevenzy2/in-memory-database/db"
)

// a node represents an database instance.

type Node struct {
	ID   string
	Addr string

	// Replicas represent a string array of existing nodes ports.
	Replicas []string

	DB *db.DB
}

type Role uint8

// Represents a nodes role, either an follower or leader.
const (
	Leader Role = iota
	Follower
)

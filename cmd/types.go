package cmd

import "time"

type Role uint8 // small integer which dictates a nodes role.

// Represents a nodes role, either an follower or leader.
const (
	Leader Role = iota
	Follower
)

// Command defines the details of the commands that will be processed by follower nodes after the leader node.

type Command struct {
	Type  string
	Key   string
	Value string

	TTL time.Duration
}

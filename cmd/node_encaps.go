package cmd

import (
	"errors"
)

// Encapsulations for node roles.

func (n *Node) GetRole() Role {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.role
}

func (n *Node) AssignLeader() {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.role = 0
}

func (n *Node) AssignFollower() {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.role = 1
}

var InvalidRoleErr = errors.New("node is not a leader.")

// Commands for which each node will replicate.
func (n *Node) Set(key, value string) error {
	if n.role != Leader {
		return InvalidRoleErr
	}

	cmd := Command{
		Type:  "SET",
		Key:   key,
		Value: value,
	}

	// applying command locally before other nodes
	err := n.DB.SetString(cmd.Type, cmd.Value, cmd.TTL)
	if err != nil {
		return err
	}

	for _, r := range n.Replicas { // replicate for each node
		err := n.Replicate(r, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

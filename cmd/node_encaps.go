package cmd

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

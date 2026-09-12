package cmd

import (
	"encoding/json"
	"fmt"
	"net"
)

func (n *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()

	var cmd Command

	d := json.NewDecoder(conn)
	if err := d.Decode(&cmd); err != nil {
		//TODO: add slog logging
		fmt.Println(err.Error())
		return
	}

	switch cmd.Type {
	case "SET":
		// TODO: later add int/string type cmd.Value parsing
		err := n.DB.SetString(cmd.Type, cmd.Value, cmd.TTL)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

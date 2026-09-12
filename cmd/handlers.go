package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
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
		//	err := n.DB.SetString(cmd.Type, cmd.Value, cmd.TTL)
		err := n.Set(cmd.Key, cmd.Value)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func (n *Node) Replicate(addr string, cmd Command) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*3) // timeout safeguards against unresponsive tcp servers.
	if err != nil {
		return err
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	if err := enc.Encode(cmd); err != nil {
		return err
	}

	return nil
}

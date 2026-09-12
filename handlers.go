package db

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn, node *Node) {
	defer conn.Close()

	//var cmd Command

	//d := json.NewDecoder(conn)
	//if err := d.Decode(&cmd); err != nil {
	//TODO: add slog logging
	//fmt.Println(err.Error())
	//	return
	//}

	buf := make([]byte, 1024) // preallocated buffer

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		input := string(buf[:n])

		parts := strings.Fields(input)

		if len(parts) == 0 { // avoid panic
			// write exit delimeter to avoid hung connections upon empty inputs entered.
			conn.Write(StringToByte(".\n"))
			continue
		}

		// database logic

		UCinput := strings.ToUpper(parts[0]) // normalise to all capital

		switch UCinput {
		case "SET":
			ok := Set(parts, conn, node)
			if !ok {
				conn.Write(StringToByte(".\n"))
			}
		case "GET":
			ok := Get(parts, node, conn)
			if !ok {
				conn.Write(StringToByte(".\n"))
			}
		case "FETCH":
			Fetch(node, conn)
			conn.Write(StringToByte(".\n"))
		case "DEL":
			ok := Del(parts, conn, node)
			if !ok {
				conn.Write(StringToByte(".\n"))
			}
		case "HELP":
			conn.Write(StringToByte("General usage: \n" +
				"SET <KeyName> <value> <TTL> \n" +
				"|| GET <KeyName>\n" +
				"|| DEL <KeyName>\n" +
				"\\n\n" +
				"to exit: run <exit>\n" +
				".\n",
			))
		case "EXIT":
			conn.Write(StringToByte("exited program. \n" +
				",\n",
			))
		default:
			conn.Write(StringToByte("invalid command \n" +
				".\n",
			))
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

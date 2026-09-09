package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/72sevenzy2/in-memory-database/db"
)

func main() {
	// start tcp server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err.Error()) // err acceping connections
			continue
		}
		b := db.NewDB() // db

		// seperate gorounine for each connection
		go func(c net.Conn) {
			defer c.Close() // close connection after reading

			buf := make([]byte, 1024) // preallocated buffer

			for {
				n, err := c.Read(buf)
				if err != nil {
					fmt.Println(err)
					return
				}

				input := string(buf[:n]) // n returns the number of bytes read

				parts := strings.Fields(input)

				if len(parts) == 0 { // avoid panic
					// write exit delimeter to avoid hung connections upon empty inputs entered.
					conn.Write(db.StringToByte(".\n"))
					continue
				}

				// database logic

				UCinput := strings.ToUpper(parts[0]) // normalise to all capital

				switch UCinput {
				case "SET":
					ok := Set(parts, conn, b)
					if !ok {
						conn.Write(db.StringToByte(".\n"))
					}
				case "GET":
					ok := Get(parts, b, conn)
					if !ok {
						conn.Write(db.StringToByte(".\n"))
					}
				case "FETCH":
					ok := Fetch(b, conn)
					if !ok {
						conn.Write(db.StringToByte(".\n"))
					}
					conn.Write(db.StringToByte(".\n"))
				case "DEL":
					ok := Del(parts, conn, b)
					if !ok {
						conn.Write(db.StringToByte(".\n"))
					}
				case "HELP":
					conn.Write(db.StringToByte("General usage: \n" +
						"SET <KeyName> <value> <TTL> \n" +
						"|| GET <KeyName>\n" +
						"|| DEL <KeyName>\n" +
						"\\n\n" +
						"to exit: run <exit>\n" +
						".\n",
					))
				case "EXIT":
					conn.Write(db.StringToByte("exited program. \n" +
						",\n",
					))
				default:
					conn.Write(db.StringToByte("invalid command \n" +
						".\n",
					))
				}

			}
		}(conn)

	}

}

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
					continue
				}

				// database logic

				UCinput := strings.ToUpper(parts[0]) // normalise to all capital

				switch UCinput {
				case "SET":
					ok := Set(parts, conn, b)
					if !ok {
						continue
					}
				case "GET":
					ok := Get(parts, b, conn)
					if !ok {
						continue
					}
				case "FETCH":
					Fetch(b, conn)
				case "DEL":
					ok := Del(parts, conn, b)
					if !ok {
						continue
					}
				case "HELP":
					conn.Write([]byte("General usage:\n"))
					conn.Write([]byte(`
					SET <KeyName> <value>
					GET <KeyName>
					DEL <KeyName>
					\n`))

					conn.Write([]byte("\nto exit: run <exit\n"))
				case "EXIT":
					conn.Write([]byte("exited program.\n"))
					return

				default:
					conn.Write([]byte("invalid command.\n"))
				}

			}
		}(conn)

	}

}

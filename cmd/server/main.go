package main

import (
	"fmt"
	"net"

	"github.com/72sevenzy2/in-memory-database"
)

func main() {
	// start tcp server
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	//b := db.NewDB() // connections share same db instance.
	node1 := db.NewNode("node-1", ":8080", 0, []string{":8082", ":8081"})
	if err := node1.Start(); err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err.Error()) // err acceping connections
			continue
		}

		go db.HandleConnection(conn, node1)
	}

}

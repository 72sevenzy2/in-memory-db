package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// cli usage

func main() {
	conn, err := net.Dial("tcp", ":8080") // connect with tcp server
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Print("> ")
		safe := scanner.Scan()
		if !safe {
			err := scanner.Err()
			if err != nil {
				fmt.Println(err)
			}
			break
		}

		input := scanner.Text()

		// send input to tcp server
		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for {
			resp, err2 := reader.ReadString('\n')
			if resp == ".\n" {
				break
			}
			// exit delimeter
			if resp == ",\n" {
				return
			}
			if err2 != nil {
				fmt.Println(err2.Error())
				continue
			}

			fmt.Println(resp)
		}
	}

}

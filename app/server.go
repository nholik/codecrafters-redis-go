package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Print("Starting server...\n\n")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		fmt.Println("Accepted connection: ", conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		for {
			buf := make([]byte, 1024*5)

			n, err := bufio.NewReader(conn).Read(buf)

			if err != nil {
				fmt.Println("Error reading from connection: ", err.Error())
				conn.Close()
				break
			}

			parts := strings.Split(string(buf[:n]), "\r\n")

			sent := parts[len(parts)-2]

			fmt.Println("Read from connection: ", sent)

			cmd := strings.ToUpper(strings.TrimSpace(sent))
			switch cmd {
			case "PING":
				conn.Write([]byte("+PONG\r\n"))
				break
			case "COMMAND":
				conn.Write([]byte("+OK\r\n"))
				break
			default:
				conn.Write([]byte("-ERR unknown command '" + cmd + "'\r\n"))
			}
		}
	}
}

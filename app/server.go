package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	// "strings"
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
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Printf("Connection Recieved from %s\n", conn.RemoteAddr())
		// reader := bufio.NewReader(conn)
		// msg_type, _, _ := reader.ReadRune()
		// commands := ProcessCommand(reader, msg_type)
		//
		// for _, cmds := range commands {
		// 	if strings.ToUpper(cmds) == "PING" {
		// 		fmt.Println("Found PING")
		// 		conn.Write([]byte("+PONG\r\n"))
		// 		conn.Close()
		// 	}
		// }
		// conn.Close()
		conn.Write([]byte("+PONG\r\n"))
		conn.Close()

	}
}

func GetCommandLength(reader *bufio.Reader) int {
	arr_len_str, _, _ := reader.ReadLine()
	arr_len, _ := strconv.Atoi(string(arr_len_str))
	return arr_len
}

func ProcessCommand(reader *bufio.Reader, cmd rune) []string {
	commands := make([]string, 10)
	command_idx := 0
	switch cmd {
	case '*':
		arr_len := GetCommandLength(reader)
		// fmt.Printf("Array message recieved of length %d\n", arr_len)
		for arr_len > 0 {
			arr_cmd, _, _ := reader.ReadRune()
			arr_cmds := ProcessCommand(reader, arr_cmd)
			arr_len--
			for _, s := range arr_cmds {
				if s != "" {
					fmt.Printf("Command is %s\n", s)
					commands[command_idx] = s
					command_idx++
				}
			}
			// fmt.Printf("%c\n", arr_cmd)
		}
		break
	case '-':
		// fmt.Println("Error message recieved")
		break
	case ':':
		// fmt.Println("Integer message recieved")
		break
	case '$':
		arr_len := GetCommandLength(reader)
		// fmt.Printf("Bulk message recieved of length %d\n", arr_len)
		s := make([]byte, arr_len+2)
		for i := range s {
			b, _ := reader.ReadByte()
			s[i] = b
		}
		// fmt.Printf("Bulk message is %s", string(s[:]))
		commands[command_idx] = string(s[:])
		command_idx++
		break
	case '+':
		// fmt.Println("Simple string message recieved")
		break
	default:
		// fmt.Printf("Unknown message type recieved %c\n", cmd)
		break
	}
	return commands
}

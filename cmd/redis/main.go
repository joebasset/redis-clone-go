package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"strings"
)

func handleConn(conn net.Conn, store *Store) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("connection closed:", err)
			break
		}

		line = strings.TrimSpace(line)

		err = store.handleRequests(line, conn)
		if err != nil {
			fmt.Println("error:", err)
			break
		}
	}
}

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listening on :6379...")
	store := NewStore()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConn(conn, store)
	}

}

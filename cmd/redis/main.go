package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/joebasset/redis-clone-go/internal/server"
)

func handleConn(conn net.Conn, store *server.Store) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("connection closed:", err)
			break
		}

		line = strings.TrimSpace(line)

		err = store.HandleRequests(line, conn)
		if err != nil {
			_, err := conn.Write([]byte("Error: " + err.Error()))
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
	store := server.NewStore()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConn(conn, store)
	}

}

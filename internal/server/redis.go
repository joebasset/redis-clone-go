package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	SET = "SET"
	GET = "GET"
	DEL = "DEL"
)

type Value struct {
	TTL  time.Time
	Data string
}

type Store struct {
	store map[string]Value
	mu    *sync.RWMutex
}

func NewStore() *Store {
	return &Store{store: make(map[string]Value), mu: &sync.RWMutex{}}
}

func GetCommandArgs(request string) []string {
	requestWithoutNewLine := strings.ReplaceAll(request, "\n", "")
	args := strings.Split(requestWithoutNewLine, " ")
	return args

}

func CheckRequestType(request string) string {
	reader := bufio.NewReader(strings.NewReader(request))

	_type, err := reader.Peek(3)
	if err != nil {
		fmt.Println("Error reading first bytes")
	}
	return string(_type)

}

// func handleGet(s *Store,) error {
// 	return

// 	return nil
// }

func (s *Store) handleRequests(request string, conn net.Conn) error {

	args := GetCommandArgs(request)
	command := args[0]
	fmt.Printf("Type: %s\n", command)
	fmt.Printf("Args: %s\n", args)
	switch strings.ToUpper(command) {
	case GET:
		_, err := conn.Write([]byte("Server received: " + command))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
		}
		return nil

	case SET:
		return nil
	case DEL:
		return nil
	default:
		return errors.New("Invalid command")
	}

}

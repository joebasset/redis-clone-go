package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	SET = "SET"
	GET = "GET"
	DEL = "DEL"
)

type Value struct {
	TTL  int16
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

func readFromStore(store *Store, args []string) (Value, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	key := args[1]
	val, ok := store.store[key]
	if !ok {
		return Value{}, errors.New("Key doesn't exist")
	}

	return val, nil
}
func addToStore(store *Store, args []string) Value {
	store.mu.Lock()
	defer store.mu.Unlock()

	key := args[1]
	value := args[2]
	newValue := Value{Data: value, TTL: 300}
	store.store[key] = newValue
	return newValue

}

func deleteFromStore(store *Store, args []string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := args[1]
	delete(store.store, key)

}
func (s *Store) HandleRequests(request string, conn net.Conn) error {

	args := GetCommandArgs(request)
	command := args[0]

	switch command {
	case GET:
		if len(args) != 2 {
			return errors.New("ERR wrong number of arguments for 'GET' command")
		}
		val, err := readFromStore(s, args)
		if err != nil {
			return err
		}
		_, err = conn.Write([]byte("Get value: " + val.Data))
		return nil

	case SET:
		if len(args) != 3 {
			return errors.New("ERR wrong number of arguments for 'SET' command")
		}
		value := addToStore(s, args)

		_, err := conn.Write([]byte("Value has been set: " + value.Data + " with TTL" + string(value.TTL)))
		if err != nil {
			return errors.New("Error sending back message")
		}
		return nil
	case DEL:
		if len(args) != 2 {
			return errors.New("ERR wrong number of arguments for 'DEL' command")
		}
		deleteFromStore(s, args)
		return nil
	default:
		return errors.New("Invalid command")
	}

}

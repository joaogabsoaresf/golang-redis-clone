package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/joaogabsoaresf/golang-redis-clone/aof"
	handler "github.com/joaogabsoaresf/golang-redis-clone/handler"
)

func main() {
	fmt.Println("Listening on port :6371")

	l, err := net.Listen("tcp", ":6371")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value handler.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := handler.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := handler.NewWriter(conn)

		h, ok := handler.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(handler.Value{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := h(args)
		writer.Write(result)
	}
}

package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// Create a server
	listner, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Listening on port :6379")

	// Open aof file
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println("error reading the aof: ", err.Error())
	}
	defer aof.Close()

	// Listen for the connection
	conn, err := listner.Accept()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	err = aof.Read()
	if err != nil {
		fmt.Println("error reading")
		return
	}

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("error reading request: ", err.Error())
			return
		}

		if value.dataType != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		writer := NewWriter(conn)

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{dataType: "string", stringValue: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		writer.Write(handler(args))
	}
}

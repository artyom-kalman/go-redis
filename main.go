package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	listner, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Listening on port :6379")

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println("error reading the aof: ", err.Error())
	}
	defer aof.Close()

	conn, err := listner.Accept()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	if err = aof.Read(); err != nil {
		fmt.Println("error reading aof file: ", err)
		return
	}

	for {
		resp := NewResp(conn)
		writer := NewWriter(conn)

		value, err := resp.Read()
		if err != nil {
			writer.Write(Value{
				dataType:    "error",
				stringValue: fmt.Sprintln("error reading request: ", err.Error()),
			})
			continue
		}

		if value.dataType != "array" {
			writer.Write(Value{
				dataType:    "error",
				stringValue: fmt.Sprintln("Invalid request, expected array"),
			})
			continue
		}

		if len(value.array) == 0 {
			writer.Write(Value{
				dataType:    "error",
				stringValue: fmt.Sprintln("Invalid request, expected array length > 0"),
			})
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{dataType: "string", stringValue: ""})
			continue
		}

		writer.Write(handler(args))

		if command == "SET" || command == "HSET" || command == "DELETE" || command == "HDELETE" {
			aof.Write(value)
		}
	}
}

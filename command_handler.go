package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
}

var sets = map[string]string{}
var setsMutex = sync.RWMutex{}

var hsets = map[string]map[string]string{}
var hsetsMutex = sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{
			dataType:    "string",
			stringValue: "PONG",
		}
	}

	return Value{
		dataType: "bulk",
		bulk:     args[0].bulk,
	}
}

func command(_ []Value) Value {
	return Value{
		dataType:    "string",
		stringValue: "CONNECTED",
	}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{
			dataType:    "error",
			stringValue: "ERR: wrong number of arguments for command `SET`, expected 2",
		}
	}

	key := args[0].bulk
	value := args[1].bulk

	setsMutex.Lock()
	sets[key] = value
	setsMutex.Unlock()

	return Value{
		dataType:    "string",
		stringValue: "OK",
	}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{
			dataType:    "error",
			stringValue: "ERR: wrong number of arguments for command `GET`, expected 1",
		}
	}

	key := args[0].bulk

	setsMutex.RLock()
	value, ok := sets[key]
	setsMutex.RUnlock()

	if !ok {
		return Value{
			dataType:    "error",
			stringValue: "ERR: invalid key",
		}
	}

	return Value{
		dataType: "bulk",
		bulk:     value,
	}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{
			dataType:    "error",
			stringValue: "ERR: wrong nubmer of arguments for command `HSET`, exepected 3",
		}
	}

	key1 := args[0].bulk
	key2 := args[1].bulk
	value := args[2].bulk

	hsetsMutex.Lock()
	_, ok := hsets[key1]
	if !ok {
		hsets[key1] = make(map[string]string)
	}
	hsets[key1][key2] = value
	hsetsMutex.Unlock()

	return Value{
		dataType:    "string",
		stringValue: "OK",
	}
}

func hget(args []Value) Value {
	if len(args) < 1 || len(args) > 2 {
		return Value{
			dataType:    "error",
			stringValue: "ERR: wrong number of arguments for command `HGET`, expected 1 or 2",
		}
	}

	if len(args) == 2 {
		key1 := args[0].bulk
		key2 := args[1].bulk

		hsetsMutex.RLock()
		value, ok := hsets[key1][key2]
		hsetsMutex.RUnlock()

		if !ok {
			return Value{
				dataType:    "error",
				stringValue: "ERR: invalid arguments",
			}
		}

		return Value{
			dataType: "bulk",
			bulk:     value,
		}
	}

	hash := args[0].bulk

	hsetsMutex.RLock()
	values, ok := hsets[hash]
	hsetsMutex.RUnlock()

	if !ok {
		return Value{
			dataType:    "error",
			stringValue: "ERR: invalid arguments",
		}
	}

	array := make([]Value, 0)
	for _, value := range values {
		array = append(array, Value{
			dataType:    "string",
			stringValue: fmt.Sprintf("%s", value),
		})
	}

	return Value{
		dataType: "array",
		array:    array,
	}
}

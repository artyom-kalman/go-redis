package main

import "strconv"

type Value struct {
	dataType    string
	stringValue string
	intValue    int
	bulk        string
	array       []Value
}

func (v *Value) Marshal() []byte {
	switch v.dataType {
	case "string":
		return v.marshalString()
	case "bulk":
		return v.marshalBulk()
	case "array":
		return v.marshalArray()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

func (v *Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.stringValue...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalArray() []byte {
	arrayLen := len(v.array)

	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(arrayLen)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < arrayLen; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v *Value) marshalNull() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(-1)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.stringValue...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

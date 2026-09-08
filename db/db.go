package db

import (
	"encoding/binary"
	"errors"
)

// metadata for value types
type Entity struct {
	Data  []byte
	Value string // will be kind of like: "string", "int", will manually assign these types after manually checking value type in main.go
}

// core logic
type DB struct {
	data map[string]Entity
}

func NewDB() *DB { // initialise a new map to hold data
	return &DB{
		data: make(map[string]Entity),
	}
}

var InvalidInputErr = errors.New("invalid input.")

func (v *DB) SetInt(key string, value uint32) error {
	if value == 0 {
		return InvalidInputErr
	}

	buf := make([]byte, 4) // uint32's has a fixed byte size of 4

	binary.LittleEndian.PutUint32(buf, value)

	v.data[key] = Entity{
		Value: "int", // will then be compared in GetAllInt() func to make sure its int
		Data:  buf,
	}
	return nil
}

func (v *DB) GetInt(key string) (uint32, bool) {
	data, ok := v.data[key]         // data will be of type of the Entity struct
	if !ok || len(data.Data) != 4 { // check if key exists and Data has exactly 4 byte (type uint32 is fixed at 4 bytes)
		return 0, false
	}

	var value uint32         // uint32 response (will be typecated to int when displaying to user)
	if data.Value == "int" { // check if int
		value = binary.LittleEndian.Uint32(data.Data)
	}

	// value := binary.LittleEndian.Uint32(data.Data)
	return value, true
}

// func to retrieve all values at once
func (v *DB) GetAllInt() (map[string]uint32, bool) {
	result := make(map[string]uint32)

	for k, v := range v.data {
		if v.Value == "int" { // make sure type is int before serialization
			result[k] = binary.LittleEndian.Uint32(v.Data)
		}
	}

	if len(result) == 0 { // no existing data case
		return nil, false
	}

	return result, true
}

// string serialization and handlers

func (v *DB) SetString(key string, value string) error {
	if value != "" {
		byteStr := []byte(value) // serialize string to type []byte, because db holds values of type []bte
		v.data[key] = Entity{
			Value: "string",
			Data:  byteStr,
		}
		return nil
	}
	return InvalidInputErr
}

// get method for string

func (v *DB) GetString(key string) (string, bool) {
	val, ok := v.data[key] // val is if type Entity struct

	var resp string            // response string
	if val.Value == "string" { // check if string
		resp = string(val.Data)
	}

	return resp, ok
}

// display all string value data from db

func (v *DB) GetAllString() (map[string]string, bool) {
	results := make(map[string]string)

	for k, val := range v.data {
		if val.Value == "string" {
			results[k] = string(val.Data)
		}
	}

	if len(results) == 0 { // handling no existing data case
		return nil, false
	}

	return results, true
}

// this will stay fixed as key will always be type string
func (v *DB) Del(key string) {
	delete(v.data, key) // built in delete() func to delete a particular key being held in the db map.
}

package service

import (
	"bytes"
	"encoding/json"
	"log"
)

type CreateOrderRequestType struct {
	User_id     string `json:"user_id"`
	Theme       string `json:"theme"`
	Description string `json:"description"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Serialize(msg any) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}

func Deserialize(b []byte) (CreateOrderRequestType, error) {
	var msg CreateOrderRequestType
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
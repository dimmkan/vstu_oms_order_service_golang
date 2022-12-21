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

type ChangeOrderStatusType struct {
	Order_id string `json:"order_id"`
	Status   string `json:"status"`
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

func DeserializeChangeStatus(b []byte) (ChangeOrderStatusType, error) {
	var msg ChangeOrderStatusType
	err := json.Unmarshal(b, &msg)
	return msg, err
}

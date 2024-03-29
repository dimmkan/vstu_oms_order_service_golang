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
	Order_id uint64 `json:"order_id"`
	Status   string `json:"status"`
}

type ChangeOrderDescriptionType struct {
	Order_id    uint64 `json:"order_id"`
	Description string `json:"description"`
}

type DeleteOrderType struct {
	Order_id uint64 `json:"order_id"`
}

type GetUserOrdersType struct {
	User_id uint64 `json:"user_id"`
}

type DeserializeType interface {
	CreateOrderRequestType | ChangeOrderStatusType | ChangeOrderDescriptionType | DeleteOrderType | GetUserOrdersType
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

func Deserialize[T DeserializeType](b []byte) (T, error) {
	var msg T
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

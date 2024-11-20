package models

import "encoding/json"

type Order struct {
	ID       string `json:"id" db:"id"`
	Item     string `json:"item" db:"item"`
	Quantity int32  `json:"quantity" db:"quantity"`
}

func (o *Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

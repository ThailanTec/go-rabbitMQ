package entities

import "github.com/google/uuid"

type Order struct {
	UUID         uuid.UUID `json:"order"`
	Destionation string    `json:"destination"`
}

type Destionation struct {
	Order string `json:"order"`
	Lat   string `json:"lat"`
	Long  string `json:"long"`
}

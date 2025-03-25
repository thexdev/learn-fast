package models

import "time"

type Product struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
}

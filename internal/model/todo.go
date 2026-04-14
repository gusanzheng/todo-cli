package model

import "time"

const DateFormat = "2006-01-02"

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	Date      string    `json:"date"`
}

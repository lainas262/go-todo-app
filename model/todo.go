package model

import "time"

type Todo struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Author        int64     `json:"author" db:"author_id"`
	Collaborators []User    `json:"collaborators" db:"-"`
}

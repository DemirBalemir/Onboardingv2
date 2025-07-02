package entities

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	AuthorID    int       `json:"author_id"`
	Price       float64   `json:"price"`
}

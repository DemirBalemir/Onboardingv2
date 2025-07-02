package entities

import "time"

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	BirthDate time.Time `json:"birthdate"`
}

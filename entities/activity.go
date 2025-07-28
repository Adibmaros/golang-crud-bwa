package entities

import "time"

type Activity struct {
	ID           int       `json:"id"`
	Title        string    `json:"title" validate:"required"`
	Category     string    `json:"category" validate:"required,oneof=TASK EVENT"`
	Description  string    `json:"description" validate:"required"`
	ActivityDate string    `json:"activity_date" validate:"required"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

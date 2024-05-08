package data

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("Item not found")
var ErrInvalidId = errors.New("Invalid ID")
var ErrMissingField = errors.New("Missing field")

type Event interface {
	ETag
	Update(*EventData) error
	PartialUpdate(*EventData) error
	GetID() int64
}

type Data interface {
	AddEvent(*EventData) (Event, error)
	GetEvents() []Event
	GetEvent(int64) Event
	DeleteEvent(int64) error
}

type EventData struct {
	ID          *int64     `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}

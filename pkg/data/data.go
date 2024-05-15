package data

import (
	"errors"
	"time"

	"net/http"
)

var ErrNotFound = errors.New("Item not found")
var ErrInvalidId = errors.New("Invalid ID")
var ErrMissingField = errors.New("Missing field")
var ErrConflict = errors.New("Resource conflict")

func ErrToHttpStatus(err error) int {
	switch {
		case errors.Is(err, ErrNotFound):
			return http.StatusNotFound
		case errors.Is(err, ErrInvalidId):
			return http.StatusNotFound
		case errors.Is(err, ErrMissingField):
			return http.StatusBadRequest
		case errors.Is(err, ErrConflict):
			return http.StatusConflict
	}

	return -1
}

type Event interface {
	ETag
	PartialUpdate(*EventData) error
	GetID() int64
}

type Data interface {
	AddEvent(*EventData) (*EventData, error)
	GetEvents() []EventData
	GetEvent(int64) *EventData
	DeleteEvent(int64, []byte) error

	UpdateEvent(int64, *EventData, []byte) (*EventData, error)
}

type EventData struct {
	ID          *int64     `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
	ETag        []byte     `json:"-"`
}

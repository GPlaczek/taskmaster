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

type Data interface {
	AddEvent(*EventData) (*EventData, error)
	GetEvents() []EventData
	GetEvent(int64) *EventData
	DeleteEvent(int64, []byte) error
	UpdateEvent(int64, *EventData, []byte) (*EventData, error)

	AddAttachment(*AttachmentData) (*AttachmentData, error)
	GetAttachments() []AttachmentData
	GetAttachment(int64) *AttachmentData
	DeleteAttachment(int64, []byte) error
	UpdateAttachment(int64, *AttachmentData, []byte) (*AttachmentData, error)
}

type EventData struct {
	ID          *int64     `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
	ETag        []byte     `json:"-"`
}

type AttachmentData struct {
	ID   *int64  `json:"id"`
	Data *string `json:"data"`
	ETag []byte  `json:"-"`
}

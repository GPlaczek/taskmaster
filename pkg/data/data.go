package data

import (
	"errors"
	"time"

	"net/http"
)

var (
	ErrNotFound     = errors.New("Item not found")
	ErrInvalidId    = errors.New("Invalid ID")
	ErrMissingField = errors.New("Missing field")
	ErrConflict     = errors.New("Resource conflict")
	ErrInvalidEtag  = errors.New("Invalid etag")
)

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
	case errors.Is(err, ErrInvalidEtag):
		return http.StatusPreconditionFailed
	}

	return http.StatusInternalServerError
}

type Data interface {
	AddEvent() (*EventData, error)
	GetEvents() []EventData
	GetEvent(int64) *EventData
	DeleteEvent(int64, []byte) error
	UpdateEvent(int64, *EventData, []byte) (*EventData, error)

	AddAttachment() (*AttachmentData, error)
	GetAttachments() []AttachmentData
	GetAttachment(int64) *AttachmentData
	DeleteAttachment(int64, []byte) error
	UpdateAttachment(int64, *AttachmentData, []byte) (*AttachmentData, error)

	BindAttachment(int64, int64) error
	GetBoundAttachments(int64) ([]AttachmentData, error)

	MergeEvents(*MergeData) (*EventData, *MergeData, error)
	GetMerges() []MergeData
	GetMerge(int64) *MergeData
}

type EventData struct {
	ID          *int64     `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
	ETag        []byte     `json:"-"`
}

type AttachmentData struct {
	ID    *int64  `json:"id"`
	Data  *string `json:"data"`
	ETag  []byte  `json:"-"`
	Event *int64  `json:"-"`
}

type MergeData struct {
	ID    *int64 `json:"id"`
	ID1   *int64 `json:"id1"`
	ID2   *int64 `json:"id2"`
	NewID *int64 `json:"new_id"`
}

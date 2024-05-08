package mem

import (
	"crypto/sha1"
	"encoding/json"
	"time"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	ETag        [20]byte  `json:"-"`
}

func (e *Event) ETagUpdate() error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	t := sha1.Sum(data)
	e.ETag = t

	return nil
}

func (e *Event) ETagGet() []byte {
	return e.ETag[:]
}

func (e *Event) ETagCompare(tag []byte) bool {
	if len(tag) != 20 {
		return false
	}

	for i, b := range e.ETag {
		if tag[i] != b {
			return false
		}
	}

	return true
}

func (e *Event) Update(ed *data.EventData) error {
	if ed.ID != nil && e.ID != *ed.ID {
		return data.ErrInvalidId
	}

	if ed.Name == nil || ed.Description == nil || ed.Date == nil {
		return data.ErrMissingField
	}

	e.Name = *ed.Name
	e.Description = *ed.Description
	e.Date = *ed.Date

	return nil
}

func (e *Event) PartialUpdate(ed *data.EventData) error {
	if ed.ID != nil && e.ID != *ed.ID {
		return data.ErrInvalidId
	}

	if ed.Name != nil {
		e.Name = *ed.Name
	}

	if ed.Description != nil {
		e.Description = *ed.Description
	}

	if ed.Date != nil {
		e.Date = *ed.Date
	}

	return nil
}

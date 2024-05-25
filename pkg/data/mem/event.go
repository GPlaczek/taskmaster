package mem

import (
	"sync"
	"time"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Event struct {
	data.ETag
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Date        time.Time     `json:"date"`
	Attachments []*Attachment `json:"attachments"`
	lock        sync.RWMutex  `json:"-"`
}

func NewEvent(id int64) *Event {
	return &Event{
		ID:   id,
		lock: sync.RWMutex{},
	}
}

func (e *Event) update(ed *data.EventData, tag []byte) (*data.EventData, error) {
	if !e.ETagCompare(tag) {
		return nil, data.ErrConflict
	}

	if ed.ID != nil && e.ID != *ed.ID {
		return nil, data.ErrInvalidId
	}

	if ed.Name == nil || ed.Description == nil || ed.Date == nil {
		return nil, data.ErrMissingField
	}

	e.Name = *ed.Name
	e.Description = *ed.Description
	e.Date = *ed.Date

	e.ETagUpdate()

	return NewEventData(e), nil
}

func (e *Event) Update(ed *data.EventData, tag []byte) (*data.EventData, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.update(ed, tag)
}

func (e *Event) partialUpdate(ed *data.EventData) error {
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

func (e *Event) PartialUpdate(ed *data.EventData) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.partialUpdate(ed)
}

func (e *Event) bindAttachment(at *Attachment) error {
	e.Attachments = append(e.Attachments, at)

	return nil
}

func (e *Event) getBoundAttachments() []data.AttachmentData {
	ad := make([]data.AttachmentData, len(e.Attachments))
	for i, at := range e.Attachments {
		ad[i] = *NewAttachmentData(at)
	}

	return ad
}

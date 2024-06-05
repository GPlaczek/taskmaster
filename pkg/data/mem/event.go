package mem

import (
	"sync"
	"time"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Event struct {
	data.ETag
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Date        time.Time           `json:"date"`
	attachments map[*Attachment]any `json:"-"`
	lock        sync.RWMutex        `json:"-"`
}

func NewEvent(id int64) *Event {
	return &Event{
		ID:   id,
		lock: sync.RWMutex{},
		attachments: make(map[*Attachment]any),
	}
}

func (e *Event) update(ed *data.EventData, tag []byte) (*data.EventData, error) {
	if !e.ETagCompare(tag) {
		return nil, data.ErrInvalidEtag
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
	if at.event != nil {
		return data.ErrConflict
	}

	e.attachments[at] = struct{}{}
	at.event = e

	return nil
}

func (e *Event) unbindAttachment(at *Attachment) error {
	_, ok := e.attachments[at]
	if !ok {
		return data.ErrNotFound
	}

	delete(e.attachments, at)
	at.event = nil

	return nil
}

func (e *Event) getBoundAttachments() []data.AttachmentData {
	ad := make([]data.AttachmentData, 0, len(e.attachments))
	for at := range e.attachments {
		ad = append(ad, *NewAttachmentData(at))
	}

	return ad
}

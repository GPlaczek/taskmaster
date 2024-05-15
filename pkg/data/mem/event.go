package mem

import (
	"crypto/sha1"
	"encoding/json"
	"sync"
	"time"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Event struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Date        time.Time    `json:"date"`
	eTag        []byte       `json:"-"`
	lock        sync.RWMutex `json:"-"`

}

func NewEvent(id int64) *Event {
	return &Event{
		ID: id,
		lock: sync.RWMutex{},
	}
}

func (e *Event) eTagUpdate() error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	t := sha1.Sum(data)
	e.eTag = t[:]

	return nil
}

func (e *Event)ETagUpdate() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	return e.eTagUpdate()
}

func (e *Event) ETagGet() []byte {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.eTag[:]
}

func (e *Event) eTagCompare(tag []byte) bool {
	if e.eTag == nil {
		return true
	}

	if len(tag) != 20 {
		return false
	}

	for i, b := range e.eTag {
		if tag[i] != b {
			return false
		}
	}

	return true
}

func (e *Event) ETagCompare(tag []byte) bool {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.eTagCompare(tag) 
}

func (e *Event) update(ed *data.EventData, tag []byte) (*data.EventData, error) {
	if !e.eTagCompare(tag) {
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

	e.eTagUpdate()

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

func (e *Event)GetID() int64 {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.ID
}

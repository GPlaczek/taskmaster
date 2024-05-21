package mem

import (
	"crypto/sha1"
	"encoding/json"
	"sync"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Attachment struct {
	ID   int64        `json:"id"`
	Data string       `json:"data"`
	eTag []byte       `json:"-"`
	lock sync.RWMutex `json:"-"`
}

func NewAttachment(id int64) *Attachment {
	return &Attachment{
		ID:   id,
		lock: sync.RWMutex{},
	}
}

func (a *Attachment) eTagUpdate() error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	t := sha1.Sum(data)
	a.eTag = t[:]

	return nil
}

func (a *Attachment) ETagUpdate() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.eTagUpdate()
}

func (a *Attachment) ETagGet() []byte {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.eTag
}

func (a *Attachment) eTagCompare(tag []byte) bool {
	if a.eTag == nil {
		return true
	}

	if len(tag) != 20 {
		return false
	}

	for i, b := range a.eTag {
		if tag[i] != b {
			return false
		}
	}

	return true
}

func (a *Attachment) ETagCompare(tag []byte) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()

	return a.eTagCompare(tag)
}

func (a *Attachment) update(ad *data.AttachmentData, tag []byte) (*data.AttachmentData, error) {
	if !a.eTagCompare(tag) {
		return nil, data.ErrConflict
	}

	if ad.ID != nil && a.ID != *ad.ID {
		return nil, data.ErrInvalidId
	}

	if ad.Data == nil {
		return nil, data.ErrMissingField
	}

	a.Data = *ad.Data

	a.eTagUpdate()

	return NewAttachmentData(a), nil
}

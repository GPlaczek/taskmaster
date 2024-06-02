package mem

import (
	"sync"

	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Attachment struct {
	data.ETag
	ID    int64        `json:"id"`
	Data  string       `json:"data"`
	lock  sync.RWMutex `json:"-"`
	event *Event       `json:"-"`
}

func NewAttachment(id int64) *Attachment {
	return &Attachment{
		ID:   id,
		lock: sync.RWMutex{},
		event: nil,
	}
}

func (a *Attachment) update(ad *data.AttachmentData, tag []byte) (*data.AttachmentData, error) {
	if !a.ETagCompare(tag) {
		return nil, data.ErrInvalidEtag
	}

	if ad.ID != nil && a.ID != *ad.ID {
		return nil, data.ErrInvalidId
	}

	if ad.Data == nil {
		return nil, data.ErrMissingField
	}

	a.Data = *ad.Data

	a.ETagUpdate()

	return NewAttachmentData(a), nil
}

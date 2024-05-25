package mem

import (
	"strings"
	"time"


	"github.com/GPlaczek/taskmaster/pkg/data"

	omap "github.com/wk8/go-ordered-map"
)

type Data struct {
	events      *omap.OrderedMap
	evId        int64
	attachments *omap.OrderedMap
	atId        int64
}

func NewData() *Data {
	return &Data{
		omap.New(), 0,
		omap.New(), 0,
	}
}

func NewEventData(e *Event) *data.EventData {
	id := e.ID
	name := strings.Clone(e.Name)
	description := strings.Clone(e.Description)
	date := time.Date(e.Date.Year(), e.Date.Month(), e.Date.Day(),
		e.Date.Hour(), e.Date.Minute(), e.Date.Second(),
		e.Date.Nanosecond(), e.Date.Location())
	tg := e.ETagGet()
	eTag := make([]byte, len(tg))
	copy(eTag, tg)
	return &data.EventData{
		ID:          &id,
		Name:        &name,
		Description: &description,
		Date:        &date,
		ETag:        eTag,
	}
}

func NewAttachmentData(a *Attachment) *data.AttachmentData {
	id := a.ID
	dt := strings.Clone(a.Data)
	tg := a.ETagGet()
	eTag := make([]byte, len(tg))
	copy(eTag, tg)

	return &data.AttachmentData{
		ID:   &id,
		Data: &dt,
		ETag: eTag,
	}
}

func (d *Data) AddEvent(ed *data.EventData) (*data.EventData, error) {
	ev := NewEvent(d.evId)
	ev.lock.Lock()
	defer ev.lock.Unlock()

	_, err := ev.update(ed, nil)
	if err != nil {
		return nil, err
	}

	err = ev.ETagUpdate()
	if err != nil {
		return nil, err
	}

	d.events.Set(ev.ID, ev)
	d.evId++

	return NewEventData(ev), nil
}

func (d *Data) GetEvents() []data.EventData {
	arr := make([]data.EventData, d.events.Len())
	p := d.events.Oldest()
	i := 0

	for p != nil {
		arr[i] = *NewEventData(p.Value.(*Event))
		p = p.Next()
		i = 1
	}

	return arr
}

func (d *Data) GetEvent(id int64) *data.EventData {
	p, ok := d.events.Get(id)
	if !ok {
		return nil
	}

	ev := p.(*Event)
	return NewEventData(ev)
}

func (d *Data) DeleteEvent(id int64, tag []byte) error {
	e, ok := d.events.Get(id)
	if !ok {
		return data.ErrNotFound
	}
	ev := e.(*Event)
	ev.lock.Lock()
	defer ev.lock.Unlock()

	if !ev.ETagCompare(tag) {
		return data.ErrConflict
	}

	_, ok = d.events.Delete(id)
	if !ok {
		return data.ErrNotFound
	}

	return nil
}

func (d *Data) UpdateEvent(id int64, ed *data.EventData, tag []byte) (*data.EventData, error) {
	p, ok := d.events.Get(id)
	if !ok {
		return nil, data.ErrNotFound
	}
	ev := p.(*Event)

	return ev.update(ed, tag)
}

func (d *Data) AddAttachment(ad *data.AttachmentData) (*data.AttachmentData, error) {
	at := NewAttachment(d.atId)
	at.lock.Lock()
	defer at.lock.Unlock()

	_, err := at.update(ad, nil)
	if err != nil {
		return nil, err
	}

	err = at.ETagUpdate()
	if err != nil {
		return nil, err
	}

	d.attachments.Set(at.ID, at)
	d.evId++

	return NewAttachmentData(at), nil
}

func (d *Data) GetAttachments() []data.AttachmentData {
	arr := make([]data.AttachmentData, d.events.Len())
	p := d.events.Oldest()
	i := 0

	for p != nil {
		arr[i] = *NewAttachmentData(p.Value.(*Attachment))
		p = p.Next()
		i = 1
	}

	return arr
}

func (d *Data) GetAttachment(id int64) *data.AttachmentData {
	p, ok := d.attachments.Get(id)
	if !ok {
		return nil
	}
	at := p.(*Attachment)

	at.lock.RLock()
	defer at.lock.RUnlock()

	return NewAttachmentData(at)
}

func (d *Data) UpdateAttachment(id int64, ad *data.AttachmentData, tag []byte) (*data.AttachmentData, error) {
	p, ok := d.attachments.Get(id)
	if !ok {
		return nil, data.ErrNotFound
	}
	at := p.(*Attachment)

	return at.update(ad, tag)
}

func (d *Data) DeleteAttachment(id int64, tag []byte) error {
	a, ok := d.attachments.Get(id)
	if !ok {
		return data.ErrNotFound
	}
	at := a.(*Attachment)
	at.lock.Lock()
	defer at.lock.Unlock()

	if !at.ETagCompare(tag) {
		return data.ErrConflict
	}

	_, ok = d.attachments.Delete(id)
	if !ok {
		return data.ErrNotFound
	}

	return nil
}

func (d *Data) BindAttachment(eid int64, aid int64) error {
	ev_, ok := d.events.Get(eid)
	if !ok {
		return data.ErrNotFound
	}
	at_, ok := d.attachments.Get(aid)
	if !ok {
		return data.ErrNotFound
	}
	ev := ev_.(*Event)
	at := at_.(*Attachment)

	ev.lock.Lock()
	defer ev.lock.Unlock()

	for _, a := range ev.Attachments {
		if a == at {
			return data.ErrConflict
		}
	}
	ev.Attachments = append(ev.Attachments, at)

	return nil
}

func (d *Data) GetBoundAttachments(eid int64) ([]data.AttachmentData, error) {
	ev_, ok := d.events.Get(eid)
	if !ok {
		return nil, data.ErrNotFound
	}

	ev := ev_.(*Event)

	ev.lock.RLock()
	defer ev.lock.RUnlock()

	return ev.getBoundAttachments(), nil
}

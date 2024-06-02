package mem

import (
	"strings"
	"sync"
	"time"

	"github.com/GPlaczek/taskmaster/pkg/data"

	omap "github.com/wk8/go-ordered-map"
)

type Data struct {
	events      *omap.OrderedMap
	evId        int64
	evLock      sync.RWMutex
	attachments *omap.OrderedMap
	atId        int64
	atLock      sync.RWMutex
	merges      []Merge
	meId        int64
	meLock      sync.RWMutex
}

func NewData() *Data {
	return &Data{
		omap.New(), 0, sync.RWMutex{},
		omap.New(), 0, sync.RWMutex{},
		make([]Merge, 0), 0, sync.RWMutex{},
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

func NewMergeData(me *Merge) *data.MergeData {
	id := me.ID
	id1 := me.ID1
	id2 := me.ID2
	newId := me.NewID

	return &data.MergeData{
		ID: &id,
		ID1: &id1,
		ID2: &id2,
		NewID: &newId,
	}
}

func (d *Data) AddEvent() (*data.EventData, error) {
	d.evLock.Lock()
	defer d.evLock.Unlock()

	eid := d.evId
	ev := NewEvent(eid)
	ev.lock.Lock()
	defer ev.lock.Unlock()

	if err := ev.ETagUpdate(); err != nil {
		return nil, err
	}

	ev.ID = eid
	d.events.Set(eid, ev)
	d.evId++

	return NewEventData(ev) , nil
}

func (d *Data) GetEvents() []data.EventData {
	d.evLock.RLock()
	defer d.evLock.RUnlock()

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
	d.evLock.RLock()
	defer d.evLock.RUnlock()

	p, ok := d.events.Get(id)
	if !ok {
		return nil
	}

	ev := p.(*Event)
	return NewEventData(ev)
}

func (d *Data) DeleteEvent(id int64, tag []byte) error {
	d.evLock.Lock()
	defer d.evLock.Unlock()

	e, ok := d.events.Get(id)
	if !ok {
		return data.ErrNotFound
	}
	ev := e.(*Event)
	ev.lock.Lock()
	defer ev.lock.Unlock()

	if !ev.ETagCompare(tag) {
		return data.ErrInvalidEtag
	}

	_, ok = d.events.Delete(id)
	if !ok {
		return data.ErrNotFound
	}

	return nil
}

func (d *Data) UpdateEvent(id int64, ed *data.EventData, tag []byte) (*data.EventData, error) {
	d.evLock.RLock()
	defer d.evLock.RUnlock()

	p, ok := d.events.Get(id)
	if !ok {
		return nil, data.ErrNotFound
	}
	ev := p.(*Event)

	return ev.update(ed, tag)
}

func (d *Data) AddAttachment() (*data.AttachmentData, error) {
	d.atLock.Lock()
	defer d.atLock.Unlock()

	aid := d.atId
	at := NewAttachment(aid)
	at.lock.Lock()
	defer at.lock.Unlock()

	if err := at.ETagUpdate(); err != nil {
		return nil, err
	}

	at.ID = aid
	d.attachments.Set(aid, at)
	d.atId++

	return NewAttachmentData(at), nil
}

func (d *Data) GetAttachments() []data.AttachmentData {
	d.atLock.RLock()
	defer d.atLock.RUnlock()

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
	d.atLock.RLock()
	defer d.atLock.RUnlock()

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
	d.atLock.RLock()
	defer d.atLock.RUnlock()

	p, ok := d.attachments.Get(id)
	if !ok {
		return nil, data.ErrNotFound
	}
	at := p.(*Attachment)

	return at.update(ad, tag)
}

func (d *Data) DeleteAttachment(id int64, tag []byte) error {
	d.atLock.RLock()
	defer d.atLock.RUnlock()

	a, ok := d.attachments.Get(id)
	if !ok {
		return data.ErrNotFound
	}
	at := a.(*Attachment)
	at.lock.Lock()
	defer at.lock.Unlock()

	if !at.ETagCompare(tag) {
		return data.ErrInvalidEtag
	}

	_, ok = d.attachments.Delete(id)
	if !ok {
		return data.ErrNotFound
	}

	return nil
}

func (d *Data) BindAttachment(eid int64, aid int64) error {
	d.evLock.RLock()
	defer d.evLock.RUnlock()
	ev_, ok := d.events.Get(eid)
	if !ok {
		return data.ErrNotFound
	}

	d.atLock.RLock()
	defer d.atLock.RUnlock()
	at_, ok := d.attachments.Get(aid)
	if !ok {
		return data.ErrNotFound
	}
	ev := ev_.(*Event)
	at := at_.(*Attachment)

	ev.lock.Lock()
	defer ev.lock.Unlock()
	at.lock.RLock()
	defer at.lock.RUnlock()

	if err := ev.bindAttachment(at); err != nil {
		return err
	}

	return nil
}

func (d *Data) UnbindAttachment(eid int64, aid int64) error {
	d.evLock.RLock()
	defer d.evLock.RUnlock()
	ev_, ok := d.events.Get(eid)
	if !ok {
		return data.ErrNotFound
	}

	d.atLock.RLock()
	defer d.atLock.RUnlock()
	at_, ok := d.attachments.Get(aid)
	if !ok {
		return data.ErrNotFound
	}
	ev := ev_.(*Event)
	at := at_.(*Attachment)

	ev.lock.Lock()
	defer ev.lock.Unlock()
	at.lock.RLock()
	defer at.lock.RUnlock()

	if err := ev.unbindAttachment(at); err != nil {
		return err
	}

	return nil
}

func (d *Data) GetBoundAttachments(eid int64) ([]data.AttachmentData, error) {
	d.evLock.RLock()
	defer d.evLock.RUnlock()

	ev_, ok := d.events.Get(eid)
	if !ok {
		return nil, data.ErrNotFound
	}

	ev := ev_.(*Event)

	ev.lock.RLock()
	defer ev.lock.RUnlock()

	return ev.getBoundAttachments(), nil
}

func (d *Data) MergeEvents(md *data.MergeData) (*data.EventData, *data.MergeData, error) {
	d.evLock.Lock()
	defer d.evLock.Unlock()
	d.meLock.Lock()
	defer d.meLock.Unlock()

	if md.ID1 == nil {
		return nil, nil, data.ErrMissingField
	}
	e1 := *md.ID1

	if md.ID2 == nil {
		return nil, nil, data.ErrMissingField
	}
	e2 := *md.ID2

	if e1 == e2 {
		return nil, nil, data.ErrConflict
	}

	ev1_, ok := d.events.Get(e1)
	if !ok {
		return nil, nil, data.ErrNotFound
	}
	ev1 := ev1_.(*Event)
	ev1.lock.Lock()
	defer ev1.lock.Unlock()

	ev2_, ok := d.events.Get(e2)
	if !ok {
		return nil, nil, data.ErrNotFound
	}
	ev2 := ev2_.(*Event)
	ev2.lock.Lock()
	defer ev2.lock.Unlock()

	ev3 := NewEvent(d.evId)
	ev3.lock.Lock()
	defer ev3.lock.Unlock()

	me := NewMerge(d.meId, ev1, ev2, ev3)

	d.merges = append(d.merges, *me)
	d.meId++

	med := NewMergeData(me)
	e3d := NewEventData(ev3)

	d.events.Set(d.evId, ev3)
	d.events.Delete(ev1.ID)
	d.events.Delete(ev2.ID)

	return e3d, med, nil
}

func (d *Data) GetMerges() []data.MergeData {
	d.meLock.RLock()
	defer d.meLock.RUnlock()

	arr := make([]data.MergeData, len(d.merges))
	for i := range d.merges {
		arr[i] = *NewMergeData(&d.merges[i])
	}

	return arr
}

func (d *Data) GetMerge(id int64) *data.MergeData {
	d.meLock.RLock()
	defer d.meLock.RUnlock()

	if int64(len(d.merges)) < id {
		return nil
	}

	return NewMergeData(&d.merges[id])
}

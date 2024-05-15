package mem

import (
	"strings"
	"time"
	"github.com/GPlaczek/taskmaster/pkg/data"

	omap "github.com/wk8/go-ordered-map"
)

type Data struct {
	events *omap.OrderedMap
	evId int64
}

func NewData() *Data {
	return &Data {
		omap.New(),
		0,
	}
}

func NewEventData(e *Event) *data.EventData {
	id := e.ID
	name := strings.Clone(e.Name)
	description := strings.Clone(e.Description)
	date := time.Date(e.Date.Year(), e.Date.Month(), e.Date.Day(),
        e.Date.Hour(), e.Date.Minute(), e.Date.Second(),
        e.Date.Nanosecond(), e.Date.Location())
    eTag := make([]byte, len(e.eTag))
    copy(eTag, e.eTag)
	return &data.EventData{
		ID: &id,
		Name: &name,
		Description: &description,
		Date: &date,
		ETag: eTag,
	}
}

func (d *Data)AddEvent(ed *data.EventData) (*data.EventData, error) {
	ev := NewEvent(d.evId)
	ev.lock.Lock()
	defer ev.lock.Unlock()

	_, err := ev.update(ed, nil)
	if err != nil {
		return nil, err
	}

	err = ev.eTagUpdate()
	if err != nil {
		return nil, err
	}

	d.events.Set(ev.ID, ev) 
	d.evId++

	return NewEventData(ev), nil
}

func (d *Data)GetEvents() []data.EventData {
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

func (d *Data)GetEvent(id int64) *data.EventData {
	p, ok := d.events.Get(id)
	if !ok {
		return nil
	}

	ev := p.(*Event)
	return NewEventData(ev) 
}

func (d *Data)DeleteEvent(id int64, tag []byte) error {
	e, ok := d.events.Get(id)
	if !ok {
		return data.ErrNotFound
	}
	ev := e.(*Event)
	ev.lock.Lock()
	defer ev.lock.Unlock()

	if !ev.eTagCompare(tag) {
		return data.ErrConflict
	}

	_, ok = d.events.Delete(id)
	if !ok {
		return data.ErrNotFound
	}

	return nil
}

func (d *Data)UpdateEvent(id int64, ed *data.EventData, tag []byte) (*data.EventData, error) {
	p, ok := d.events.Get(id)
	if !ok {
		return nil, data.ErrNotFound
	}
	ev := p.(*Event)

	return ev.update(ed, tag)
}

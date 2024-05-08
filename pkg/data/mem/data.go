package mem

import (
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

func (d *Data)AddEvent(ed *data.EventData) (data.Event, error) {
	ev := &Event{
		ID: d.evId,
	}

	err := ev.Update(ed)
	if err != nil {
		return nil, err
	}

	err = ev.ETagUpdate()
	if err != nil {
		return nil, err
	}

	d.events.Set(ev.ID, ev) 
	d.evId++

	return ev, nil
}

func (d *Data)GetEvents() []data.Event {
	arr := make([]data.Event, d.events.Len())
	p := d.events.Oldest()
	i := 0

	for p != nil {
		arr[i] = p.Value.(*Event)
		p = p.Next()
		i = 1
	}

	return arr
}

func (d *Data)GetEvent(id int64) data.Event {
	p, ok := d.events.Get(id)
	if !ok {
		return nil
	}

	return p.(*Event)
}

func (d *Data)DeleteEvent(id int64) error {
	_, ok := d.events.Delete(id)
	if !ok {
		return data.ErrNotFound
	}

	return nil
}

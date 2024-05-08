package mem

import (
	"github.com/GPlaczek/taskmaster/pkg/data"
)

type Data struct {
	events []*Event
	evId int64
}

func NewData() *Data {
	return &Data {
		[]*Event{},
		0,
	}
}

func (d *Data)AddEvent(ed *data.EventData) (data.Event, error) {
	ev := Event{
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

	d.events = append(d.events, &ev)
	d.evId++

	return &ev, nil
}

func (d *Data)GetEvents() []data.Event {
	arr := make([]data.Event, len(d.events))

	for i, ev := range d.events{
		arr[i] = ev
	}

	return arr
}

func (d *Data)GetEvent(id int64) data.Event {
	for i := range d.events {
		if d.events[i].ID == id {
			return d.events[i]
		}
	}

	return nil
}

func (d *Data)DeleteEvent(id int64) error {
	var ind int = -1
	for i := range d.events {
		if d.events[i].ID == id {
			ind = i
			break
		}
	}

	if ind == -1 {
		return data.ErrInvalidId
	}

	d.events = append(d.events[:ind], d.events[ind+1:]...)

	return nil
}

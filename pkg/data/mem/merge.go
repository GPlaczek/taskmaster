package mem

type Merge struct {
	ID    int64
	ID1   int64
	ID2   int64
	NewID int64
}

func NewMerge(id int64, ev1 *Event, ev2 *Event, ev3 *Event) *Merge {
	ev3.Name = ev1.Name
	ev3.Description = ev1.Description
	ev3.Date = ev1.Date
	ev3.Attachments = append(ev1.Attachments, ev2.Attachments...)

	return &Merge{
		ID:    id,
		ID1:   ev1.ID,
		ID2:   ev2.ID,
		NewID: ev3.ID,
	}
}

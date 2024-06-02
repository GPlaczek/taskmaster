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
	ev3.attachments = ev1.attachments
	for at := range ev2.attachments {
		ev3.attachments[at] = struct{}{}
	}

	return &Merge{
		ID:    id,
		ID1:   ev1.ID,
		ID2:   ev2.ID,
		NewID: ev3.ID,
	}
}

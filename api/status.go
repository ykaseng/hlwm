package api

import "bytes"

type State int

const (
	Focused State = iota
	Occupied
	Urgent
	Viewed
	Empty
)

var stM = map[State]string{
	Focused:  "focused",
	Occupied: "occupied",
	Urgent:   "urgent",
	Viewed:   "viewed",
	Empty:    "empty",
}

func (s State) String() string {
	return stM[s]
}

func (s State) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`"`)
	b.WriteString(s.String())
	b.WriteString(`"`)
	return b.Bytes(), nil
}

type Status struct {
	Tags  []int   `json:"tags"`
	Views []State `json:"views"`
}

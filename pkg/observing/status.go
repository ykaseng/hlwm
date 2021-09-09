package observing

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
	Tag  int   `json:"tag"`
	View State `json:"view"`
}

type StatusMap map[int]State

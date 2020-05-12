package nodedb

import (
	"time"
)

type Node struct {
	ID            int64
	TypeID        int64
	OwnerID       int64
	Flags         Flag
	Created       time.Time
	Start         time.Time
	End           time.Time
	Deleted       time.Time
	Up            int64
	Left          int64
	Right         int64
	Text          string // should be short < 2 KB
	Content       []byte
	contentLoaded bool
}

func (n *Node) IsContentLoaded() bool {
	return n.contentLoaded
}

func (n *Node) SetContentLoaded(value bool) {
	n.contentLoaded = value
}

func (n *Node) IsFirst() bool {
	return n.Flags.IsFirst()
}

func (n *Node) SetFirst(value bool) {
	n.Flags.Set(First, value)
}

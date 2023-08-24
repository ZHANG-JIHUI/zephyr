package actor

import (
	"strings"

	"github.com/zeebo/xxh3"
)

var pidSeparator = "/"

func NewPID(address, id string, tags ...string) *PID {
	pid := &PID{
		Address: address,
		ID:      id,
	}
	if len(tags) > 0 {
		pid.ID = pid.ID + pidSeparator + strings.Join(tags, pidSeparator)
	}
	return pid
}

func (slf *PID) String() string {
	return slf.Address + pidSeparator + slf.ID
}

func (slf *PID) Equals(other *PID) bool {
	return slf.Address == other.Address && slf.ID == other.ID
}

func (slf *PID) Child(id string, tags ...string) *PID {
	childID := slf.ID + pidSeparator + id
	if len(tags) == 0 {
		return NewPID(slf.Address, childID)
	}
	return NewPID(slf.Address, childID+pidSeparator+strings.Join(tags, pidSeparator))
}

func (slf *PID) HasTag(tag string) bool {
	panic("TODO")
}

func (slf *PID) LookupKey() uint64 {
	key := append([]byte(slf.Address), slf.ID...)
	return xxh3.Hash(key)
}

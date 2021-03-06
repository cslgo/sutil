package center

import (
	"github.com/ZhengHe-MD/agollo/v4"
	"github.com/pkg/errors"
)

type ChangeType int

const (
	ADD ChangeType = iota
	MODIFY
	DELETE
)

type ChangeEventSource int

const (
	Etcd ChangeEventSource = iota
	Apollo
)

var (
	invalidAgolloChangeTypeErr = errors.New("invalid agollo change type")
)

func (c ChangeType) String() string {
	switch c {
	case ADD:
		return "ADD"
	case MODIFY:
		return "MODIFY"
	case DELETE:
		return "DELETE"
	}

	return "UNKOWN"
}

type ChangeEvent struct {
	Source ChangeEventSource
	// Namespace means
	//   Namespace in Apollo
	//   Path in Etcd
	Namespace string
	Changes   map[string]*Change
}

func fromAgolloChangeEvent(ace *agollo.ChangeEvent) *ChangeEvent {
	var changes = map[string]*Change{}
	for k, ac := range ace.Changes {
		if c, err := fromAgolloChange(ac); err == nil {
			changes[k] = c
		}
	}
	return &ChangeEvent{
		Source:    Apollo,
		Namespace: ace.Namespace,
		Changes:   changes,
	}
}

type Change struct {
	OldValue   string
	NewValue   string
	ChangeType ChangeType
}

func fromAgolloChange(ac *agollo.Change) (change *Change, err error) {
	ct, err := fromAgolloChangeType(ac.ChangeType)
	if err != nil {
		return
	}

	change = &Change{
		OldValue:   ac.OldValue,
		NewValue:   ac.NewValue,
		ChangeType: ct,
	}
	return
}

func fromAgolloChangeType(act agollo.ChangeType) (ct ChangeType, err error) {
	switch act {
	case agollo.ADD:
		ct = ADD
	case agollo.MODIFY:
		ct = MODIFY
	case agollo.DELETE:
		ct = DELETE
	default:
		err = invalidAgolloChangeTypeErr
	}

	return
}

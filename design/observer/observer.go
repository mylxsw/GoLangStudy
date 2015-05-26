package observer

import (
	"container/list"
)

type Subject struct {
	observers *list.List
}

type Observer interface {
	Update(subject *Subject)
}

func (subject *Subject) Attach(observer *Observer) {
	subject.observers.PushBack(observer)
}

func (subject *Subject) Detach(observer *Observer) {
	for e := subject.observers.Front(); e != nil; e = e.Next() {
		_observer := e.Value.(*Observer)
		if _observer == observer {
			subject.observers.Remove(e)
			break
		}
	}
}

func (subject *Subject) Notify() {
	for e := subject.observers.Front(); e != nil; e = e.Next() {
		observer := e.Value.(*Observer)
		(*observer).Update(subject)
	}
}

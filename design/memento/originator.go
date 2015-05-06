package memento

import (
	"fmt"
)

type Originator struct {
	state1 string
	state2 string
	state3 string
}

func (ori *Originator) CreateMemento() *Memento {
	states := make(map[string]interface{})
	states["state1"] = ori.GetState1()
	states["state2"] = ori.GetState2()
	states["state3"] = ori.GetState3()

	memento := &Memento{}
	memento.SetState(states)

	return memento
}

func (ori *Originator) RestoreMemento(mem Memento) {
	state := mem.GetState()
	ori.state1 = state["state1"].(string)
	ori.state2 = state["state2"].(string)
	ori.state3 = state["state3"].(string)
}

func (ori Originator) String() string {
	return fmt.Sprintf("state1 = %v, state2 = %v, state3 = %v",
		ori.GetState1(),
		ori.GetState2(),
		ori.GetState3(),
	)
}

func (ori *Originator) SetState1(state string) {
	ori.state1 = state
}

func (ori *Originator) SetState2(state string) {
	ori.state2 = state
}

func (ori *Originator) SetState3(state string) {
	ori.state3 = state
}

func (ori *Originator) GetState1() string {
	return ori.state1
}

func (ori *Originator) GetState2() string {
	return ori.state2
}

func (ori *Originator) GetState3() string {
	return ori.state3
}

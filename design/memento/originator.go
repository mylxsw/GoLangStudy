package memento

import (
	"fmt"
)

type Originator struct {
	state1 string
	state2 string
}

type MementoImpl struct {
	state map[string]string
}

func (mem *MementoImpl) GetState() interface{} {
	return mem.state
}

func (mem *MementoImpl) SetState(state interface{}) {
	mem.state = state.(map[string]string)
}

func (ori *Originator) CreateMemento() *MementoImpl {
	states := make(map[string]string)
	states["state1"] = ori.GetState1()
	states["state2"] = ori.GetState2()

	memento := &MementoImpl{}
	memento.SetState(states)

	return memento
}

func (ori *Originator) RestoreMemento(mem Memento) {
	state := mem.GetState().(map[string]string)
	ori.state1 = state["state1"]
	ori.state2 = state["state2"]
}

func (ori Originator) String() string {
	return fmt.Sprintf("state1 = %v, state2 = %v",
		ori.GetState1(),
		ori.GetState2(),
	)
}

func (ori *Originator) SetState1(state string) {
	ori.state1 = state
}

func (ori *Originator) SetState2(state string) {
	ori.state2 = state
}

func (ori *Originator) GetState1() string {
	return ori.state1
}

func (ori *Originator) GetState2() string {
	return ori.state2
}

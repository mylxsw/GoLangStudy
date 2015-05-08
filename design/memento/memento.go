package memento

type Memento interface {
	SetState(state interface{})
	GetState() interface{}
}

//type Memento struct {
//	state map[string]interface{}
//}

//func (mem *Memento) SetState(state map[string]interface{}) {
//	mem.state = state
//}

//func (mem *Memento) GetState() map[string]interface{} {
//	return mem.state
//}

package memento

type Caretaker struct {
	memento Memento
}

func (care *Caretaker) SetMemento(mem Memento) {
	care.memento = mem
}

func (care Caretaker) GetMemento() Memento {
	return care.memento
}

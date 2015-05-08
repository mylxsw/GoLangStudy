package facade

type Facade struct {
	systemA *SystemA
	systemB *SystemB
	systemC *SystemC
}

func New() (facade *Facade) {

	facade = &Facade{}

	facade.systemA = &SystemA{}
	facade.systemB = &SystemB{}
	facade.systemC = &SystemC{}

	return
}

func (facade *Facade) Operation() {
	facade.systemA.operationA()
	facade.systemB.operationB()
	facade.systemC.operationC()
}

package composite

import "fmt"
import "container/list"

type Element struct {
	Value GraphicComponent
	next  *Element
	prev  *Element
}

func (ele *Element) Next() *Element {
	return ele.next
}

func (ele *Element) Prev() *Element {
	return ele.prev
}

type GraphicComponent interface {
	Draw()
	Add(comp GraphicComponent)
	Remove(comp GraphicComponent)
	Iterator() *Element
}

type GraphLeaf struct{}

func (leaf *GraphLeaf) Draw() {
	fmt.Println("叶子绘图")
}

func (leaf *GraphLeaf) Add(comp GraphicComponent) {
	panic("不支持该操作")
}

func (leaf *GraphLeaf) Remove(comp GraphicComponent) {
	panic("不支持该操作")
}

func (leaf *GraphLeaf) Iterator() *Element {
	panic("不支持该操作")
}

type GraphComposite struct {
	children *list.List
}

func NewComposite() *GraphComposite {
	comp := &GraphComposite{}
	comp.children = list.New()

	return comp
}

func (composite *GraphComposite) Draw() {
	fmt.Println("容器绘图")
}

func (composite *GraphComposite) Add(comp GraphicComponent) {
	composite.children.PushBack(comp)
}

func (composite *GraphComposite) Remove(comp GraphicComponent) {
	ele := &list.Element{}
	ele.Value = comp
	composite.children.Remove(ele)
}

func (composite *GraphComposite) Iterator() *Element {
	root := &Element{}
	root.next = nil
	root.prev = nil

	for e := composite.children.Front(); e != nil; e = e.Next() {
		element := &Element{}
		element.Value = e.Value.(GraphicComponent)
		element.next = root.next
		element.prev = root
		if root.next != nil {
			root.next.prev = element
		}
		root.next = element
	}

	return root.next
}

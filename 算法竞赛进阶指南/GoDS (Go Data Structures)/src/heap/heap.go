package heap

type Heap interface {
	Push(value interface{})
	Pop() (value interface{})
	Peek() (value interface{})
	Len() int
}

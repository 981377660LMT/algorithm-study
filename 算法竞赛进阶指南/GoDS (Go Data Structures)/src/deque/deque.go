package deque

type Deque interface {
	Append(value interface{})
	AppendLeft(value interface{})
	Pop() (value interface{})
	PopLeft() (value interface{})
	At(index int) (value interface{})
	ForEach(func(value interface{}, index int))
	Len() int
}

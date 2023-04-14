// https://scrapbox.io/data-structures/Realtime_Queue
package main

import "fmt"

func main() {
	s := NewStream()
	s2 := s.Push(1).Push(2).Push(3)
	fmt.Println(s2.Top())
	rev := s2.Reverse()
	fmt.Println(rev.Top())
	fmt.Println(rev)
}

type S int
type Cell struct {
	resolved S
	next     *Stream
}
type Stream struct {
	*Suspension
}

// 惰性求值的流.
func NewStream() *Stream {
	return &Stream{Suspension: NewSuspension()}
}

// 连接两个流。
func Concat(x, y *Stream) *Stream {
	return &Stream{Suspension: NewSuspensionWith(func() interface{} {
		if x.Empty() {
			return y.Resolve()
		}
		return &Cell{resolved: x.Top(), next: Concat(x.Pop(), y)}
	})}
}

func (s *Stream) Empty() bool {
	return s.Resolve() == nil
}

func (s *Stream) Top() S {
	return s.Resolve().(*Cell).resolved
}

func (s *Stream) Pop() *Stream {
	return s.Resolve().(*Cell).next
}

func (s *Stream) Push(x S) *Stream {
	return &Stream{Suspension: NewSuspensionWith(&Cell{resolved: x, next: s})}
}

func (s *Stream) Reverse() *Stream {
	return &Stream{Suspension: NewSuspensionWith(func() interface{} {
		x := s
		res := NewStream()
		for !x.Empty() {
			res = res.Push(x.Top())
			x = x.Pop()
		}
		return res.Resolve()
	})}
}

func (s *Stream) String() string {
	x := s
	res := []S{}
	for !x.Empty() {
		res = append(res, x.Top())
		x = x.Pop()
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return fmt.Sprintf("Stream%v", res)
}

// 惰性求值.
type Suspension struct {
	x        interface{}
	resolved interface{}
}

func NewSuspension() *Suspension {
	return &Suspension{}
}

func NewSuspensionWith(x interface{}) *Suspension {
	return &Suspension{x: x}
}

func (s *Suspension) Resolve() interface{} {
	if s.resolved == nil {
		if f, ok := s.x.(func() interface{}); ok {
			s.resolved = f()
		} else {
			s.resolved = s.x
		}
	}
	return s.resolved
}

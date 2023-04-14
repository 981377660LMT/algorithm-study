package main

import "fmt"

func main() {
	s := NewSuspensionWith(1)
	fmt.Println(s.Resolve())
	s = NewSuspensionWith(func() interface{} { return 2 })
	fmt.Println(s.Resolve())
	s = NewSuspension()
	fmt.Println(s.Resolve())
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

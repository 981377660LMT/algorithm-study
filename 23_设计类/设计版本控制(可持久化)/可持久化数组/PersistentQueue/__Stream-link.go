// // https://scrapbox.io/data-structures/Realtime_Queue
// package main

// import "fmt"

// type S int
// type Stream struct {
// 	value S
// 	next  *Stream
// }

// func NewStream() *Stream {
// 	return &Stream{}
// }

// func Concat(x, y *Stream) *Stream {
// 	if x.next == nil {
// 		return &Stream{value: y.value, next: y.next}
// 	}
// 	return &Stream{value: x.value, next: Concat(x.next, y)}
// }

// func (s *Stream) Empty() bool {
// 	return s.next == nil
// }

// func (s *Stream) Top() S {
// 	return s.value
// }

// func (s *Stream) Pop() *Stream {
// 	return s.next
// }

// func (s *Stream) Push(x S) *Stream {
// 	return &Stream{value: x, next: s}
// }

// func (s *Stream) Reverse() *Stream {
// 	x := s
// 	res := NewStream()
// 	for x.next != nil {
// 		res = res.Push(x.value)
// 		x = x.next
// 	}
// 	return &Stream{value: res.value, next: res.next}
// }

// func (s *Stream) String() string {
// 	x := s
// 	res := []S{}
// 	for x.next != nil {
// 		res = append(res, x.value)
// 		x = x.next
// 	}
// 	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
// 		res[i], res[j] = res[j], res[i]
// 	}
// 	return fmt.Sprintf("Link{%v}", res)
// }

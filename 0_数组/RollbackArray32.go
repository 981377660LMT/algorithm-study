/* eslint-disable @typescript-eslint/no-non-null-assertion */

package main

import "fmt"

func main() {
	arr := NewRollbackArray32(10, func(index int32) int32 { return 0 })
	arr.Set(0, 1)
	arr.Set(1, 2)
	fmt.Println(arr.GetAll())
	arr.Undo()
	fmt.Println(arr.GetAll())
}

const mask int = 1<<32 - 1

type RollbackArray32 struct {
	n       int32
	data    []int32
	history []int // (index, value)
}

func NewRollbackArray32(n int32, f func(index int32) int32) *RollbackArray32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray32{
		n:    n,
		data: data,
	}
}

func NewRollbackArray32From(data []int32) *RollbackArray32 {
	return &RollbackArray32{n: int32(len(data)), data: data}
}

func (r *RollbackArray32) GetTime() int32 {
	return int32(len(r.history))
}

func (r *RollbackArray32) Rollback(time int32) {
	for i := int32(len(r.history)) - 1; i >= time; i-- {
		pair := r.history[i]
		r.data[pair>>32] = int32(pair & mask)
	}
	r.history = r.history[:time]
}

func (r *RollbackArray32) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair>>32] = int32(pair & mask)
	return true
}

func (r *RollbackArray32) Get(index int32) int32 {
	return r.data[index]
}

func (r *RollbackArray32) Set(index int32, value int32) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, int(index)<<32|int(r.data[index]))
	r.data[index] = value
	return true
}

func (r *RollbackArray32) GetAll() []int32 {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray32) Len() int32 {
	return r.n
}

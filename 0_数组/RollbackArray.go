/* eslint-disable @typescript-eslint/no-non-null-assertion */

package main

import "fmt"

func main() {
	arr := NewRollbackArray(10, func(index int) V { return 0 })
	arr.Set(0, 1)
	arr.Set(1, 2)
	fmt.Println(arr.GetAll())
	arr.Undo()
	fmt.Println(arr.GetAll())
}

type V = interface{}
type HistoryItem struct {
	index int
	value V
}

type RollbackArray struct {
	n       int
	data    []V
	history []HistoryItem
}

func NewRollbackArray(n int, f func(index int) V) *RollbackArray {
	data := make([]V, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray{
		n:    n,
		data: data,
	}
}

func (r *RollbackArray) GetTime() int {
	return len(r.history)
}

func (r *RollbackArray) Rollback(time int) {
	for len(r.history) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *RollbackArray) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray) Get(index int) V {
	return r.data[index]
}

func (r *RollbackArray) Set(index int, value V) {
	r.history = append(r.history, HistoryItem{index: index, value: r.data[index]})
	r.data[index] = value
}

func (r *RollbackArray) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray) Len() int {
	return r.n
}

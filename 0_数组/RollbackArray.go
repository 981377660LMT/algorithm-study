/* eslint-disable @typescript-eslint/no-non-null-assertion */

package main

import "fmt"

func main() {
	arr := NewRollbackArray(10, func(index int32) int32 { return 0 })
	arr.Set(0, 1)
	arr.Set(1, 2)
	fmt.Println(arr.GetAll())
	arr.Undo()
	fmt.Println(arr.GetAll())
}

type HistoryItem[V comparable] struct {
	index int32
	value V
}

type RollbackArray[V comparable] struct {
	n       int32
	data    []V
	history []HistoryItem[V]
}

func NewRollbackArray[V comparable](n int32, f func(index int32) V) *RollbackArray[V] {
	data := make([]V, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray[V]{
		n:    n,
		data: data,
	}
}

func NewRollbackArrayFrom[V comparable](data []V) *RollbackArray[V] {
	return &RollbackArray[V]{n: int32(len(data)), data: data}
}

func (r *RollbackArray[V]) GetTime() int32 {
	return int32(len(r.history))
}

func (r *RollbackArray[V]) Rollback(time int32) {
	for i := int32(len(r.history)) - 1; i >= time; i-- {
		pair := r.history[i]
		r.data[pair.index] = pair.value
	}
	r.history = r.history[:time]
}

func (r *RollbackArray[V]) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray[V]) Get(index int32) V {
	return r.data[index]
}

func (r *RollbackArray[V]) Set(index int32, value V) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, HistoryItem[V]{index: index, value: r.data[index]})
	r.data[index] = value
	return true
}

func (r *RollbackArray[V]) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray[V]) Len() int32 {
	return r.n
}

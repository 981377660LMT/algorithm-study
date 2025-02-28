package main

type Iterator[T any] interface {
	Next() T
	HasNext() bool
}

type Peekable[T any] struct {
	iter    Iterator[T]
	hasNext bool
	next    T
}

func NewPeekable[T any](iter Iterator[T]) *Peekable[T] {
	return &Peekable[T]{iter, iter.HasNext(), iter.Next()}
}

func (p *Peekable[T]) Peek() T {
	return p.next
}

func (p *Peekable[T]) Next() T {
	res := p.next
	p.hasNext = p.iter.HasNext()
	if p.hasNext {
		p.next = p.iter.Next()
	}
	return res
}

func (p *Peekable[T]) HasNext() bool {
	return p.hasNext
}

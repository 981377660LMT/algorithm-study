// BinaryGrouping-二进制分组

package main

type IPreprocessor[V any] interface {
	Add(value V)
	Build()
	Clear()
}

type BinaryGrouping[T IPreprocessor[V], V any] struct {
	groups             [][]V
	preprocessors      []IPreprocessor[V]
	createPreprocessor func() IPreprocessor[V]
}

func NewBinaryGrouping[T IPreprocessor[V], V any](createPreprocessor func() IPreprocessor[V]) *BinaryGrouping[T, V] {
	return &BinaryGrouping[T, V]{
		createPreprocessor: createPreprocessor,
	}
}

func (b *BinaryGrouping[T, V]) Add(value V) {
	k := 0
	for k < len(b.groups) && len(b.groups[k]) > 0 {
		k++
	}
	if k == len(b.groups) {
		b.groups = append(b.groups, []V{})
		b.preprocessors = append(b.preprocessors, b.createPreprocessor())
	}
	b.groups[k] = append(b.groups[k], value)
	b.preprocessors[k].Add(value)
	for i := 0; i < k; i++ {
		for _, v := range b.groups[i] {
			b.preprocessors[k].Add(v)
		}
		b.groups[k] = append(b.groups[k], b.groups[i]...)
		b.preprocessors[i].Clear()
		b.groups[i] = b.groups[i][:0]
	}
	b.preprocessors[k].Build()
}

func (b *BinaryGrouping[T, V]) Query(onQuery func(p IPreprocessor[V]) (shouldBreak bool), ignoreEmpty bool) {
	for i := 0; i < len(b.preprocessors); i++ {
		if ignoreEmpty && len(b.groups[i]) == 0 {
			continue
		}
		if onQuery(b.preprocessors[i]) {
			break
		}
	}
}

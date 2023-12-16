package main

func main() {

}

type V = int

type IPreprocessor interface {
	Add(value V)
	Build()
	Clear()
}

type BinaryGrouping struct {
	groups             [][]V
	preprocessors      []IPreprocessor
	createPreprocessor func() IPreprocessor
}

func NewBinaryGrouping(createPreprocessor func() IPreprocessor) *BinaryGrouping {
	return &BinaryGrouping{
		createPreprocessor: createPreprocessor,
	}
}

func (b *BinaryGrouping) Add(value V) {
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
}

func (b *BinaryGrouping) Query(onQuery func(p IPreprocessor)) {
	for i := 0; i < len(b.preprocessors); i++ {
		onQuery(b.preprocessors[i])
	}
}

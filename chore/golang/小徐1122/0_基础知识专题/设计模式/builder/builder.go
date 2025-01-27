package main

func main() {
	bigClass := NewBigClass(WithName("小红"), WithAge(18), WithSex("女"), WithWeight(50), WithFieldC("C"))
	_ = bigClass
}

type BigClass struct {
	Options
}

type Options struct {
	name   string
	age    int
	sex    string
	weight float64
	height float64
	width  float64
	fieldA string
	fieldB string
	fieldC string
}

type Option func(opts *Options)

func WithName(name string) Option {
	return func(opts *Options) {
		opts.name = name
	}
}

func WithAge(age int) Option {
	return func(opts *Options) {
		opts.age = age
	}
}

func WithSex(sex string) Option {
	return func(opts *Options) {
		opts.sex = sex
	}
}

func WithWeight(weight float64) Option {
	return func(opts *Options) {
		opts.weight = weight
	}
}

func WithHeight(height float64) Option {
	return func(opts *Options) {
		opts.height = height
	}
}

func WithWidth(width float64) Option {
	return func(opts *Options) {
		opts.width = width
	}
}

func WithFieldA(fieldA string) Option {
	return func(opts *Options) {
		opts.fieldA = fieldA
	}
}

func WithFieldB(fieldB string) Option {
	return func(opts *Options) {
		opts.fieldB = fieldB
	}
}

func WithFieldC(fieldC string) Option {
	return func(opts *Options) {
		opts.fieldC = fieldC
	}
}

func repair(opts *Options) {
	if opts.name == "" {
		opts.name = "小明"
	}
	if opts.age == 0 {
		opts.age = 20
	}
}

func NewBigClass(opts ...Option) *BigClass {
	bigClass := BigClass{}
	for _, opt := range opts {
		opt(&bigClass.Options)
	}

	repair(&bigClass.Options)
	return &bigClass
}

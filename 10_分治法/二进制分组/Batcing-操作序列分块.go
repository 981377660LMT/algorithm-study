package main

func main() {

}

type V = int

type IPreprocessor interface {
	Add(value V)
	Build()
	Clear()
}

// 操作分块.
type Blocking struct {
}

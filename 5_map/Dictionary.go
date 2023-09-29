package main

import "fmt"

func main() {
	dict := NewDictionary()
	fmt.Println(dict.Size(), dict.Id("a"), dict.Id("b"), dict.Id("c"))
	fmt.Println(dict.Value(0))
	fmt.Println(dict.Value(10))
}

type V = interface{}
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}

func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}

func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}

func (d *Dictionary) Size() int {
	return len(d._idToValue)
}

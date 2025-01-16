// https://github.com/Jsewill/morton/blob/master/morton.go

/*
Morton implements Z-Order Curve encoding and decoding for N-dimensions, using lookup tables and magic bits respectively.

In order to supply for N-dimensions, this library generates the magic bits used in decoding.  While this library does supply for N-dimensions, because this type of ordering uses bit interleaving for encoding it is limited by the width of the uint64 type divided by the number of dimensions (i.e., uint64/3 for 3 dimensions).
*/
package main

import (
	"errors"
	"fmt"
	"sort"
)

func main() {
	//Create a new Morton
	m := new(Morton)
	//Generate Tables and Magic bits
	m.Create(4, 512)

	//Create arbitrary coordinates
	c := []uint32{511, 472, 103, 7}

	//Encode aforementioned coordinates
	e, err := m.Encode(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Coordinates: %v\n", c)
	fmt.Printf("Encoded Coordinates: %v\n", e)
	fmt.Printf("Decoded Coordinates: %v\n", m.Decode(e))
}

type Table struct {
	Index  uint8
	Length uint32
	Encode []Bit
}

func (t Table) String() string {
	var bits string
	for _, b := range t.Encode {
		bits = fmt.Sprintf("%v%v\n", bits, b)
	}
	return fmt.Sprintf("Index: %v\nLength: %v\n%v", t.Index, t.Length, bits)
}

// Sortable Table slice type to satisfy the sort package interface
type ByTable []Table

func (t ByTable) Len() int {
	return len(t)
}

func (t ByTable) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByTable) Less(i, j int) bool {
	return t[i].Index < t[j].Index
}

type Bit struct {
	Index uint32
	Value uint64
}

func (b Bit) String() string {
	return fmt.Sprintf("[%v]%08b", b.Index, b.Value)
}

// Sortable Table slice type to satisfy the sort package interface
type ByBit []Bit

func (b ByBit) Len() int {
	return len(b)
}

func (b ByBit) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByBit) Less(i, j int) bool {
	return b[i].Index < b[j].Index
}

type Morton struct {
	Dimensions uint8    // 维度数，如 2、3、4
	Tables     []Table  // 针对每个维度，存储了一个 Table 用于快速编码
	Magic      []uint64 // 用来协助解码（从 Morton code 反推多维坐标），是一系列的掩码或拆位信息
}

// Convenience function
func New(dimensions uint8, size uint32) *Morton {
	m := new(Morton)
	m.Create(dimensions, size)
	return m
}

func (m *Morton) Create(dimensions uint8, size uint32) {
	done := make(chan struct{})
	mch := make(chan []uint64)
	go func() {
		m.CreateTables(dimensions, size)
		done <- struct{}{}
	}()
	go func() {
		mch <- MakeMagic(dimensions)
		m.Magic = MakeMagic(dimensions)
	}()
	m.Magic = <-mch
	close(mch)
	<-done
	close(done)
}

func (m *Morton) CreateTables(dimensions uint8, length uint32) {
	ch := make(chan Table)

	m.Dimensions = dimensions
	for i := uint8(0); i < dimensions; i++ {
		go func(i uint8) {
			ch <- CreateTable(i, dimensions, length)
		}(i)
	}
	for i := uint8(0); i < dimensions; i++ {
		t := <-ch
		m.Tables = append(m.Tables, t)
	}
	close(ch)

	sort.Sort(ByTable(m.Tables))
}

func MakeMagic(dimensions uint8) []uint64 {
	// Generate nth and ith bits variables
	d := uint64(dimensions)
	limit := 64/d + 1
	nth := []uint64{0, 0, 0, 0, 0, 0}
	for i := uint64(0); i < limit; i++ {

		switch {
		case i <= 32:
			//32
			nth[0] |= 1 << (i * d)
			fallthrough
		case i <= 16:
			//16
			nth[1] |= 3 << (i * (d << 1))
			fallthrough
		case i <= 8:
			//8
			nth[2] |= 0xf << (i * (d << 2))
			fallthrough
		case i <= 4:
			//4
			nth[3] |= 0xff << (i * (d << 3))
			fallthrough
		case i <= 2:
			//2
			nth[4] |= 0xffff << (i * (d << 4))
			fallthrough
		case i <= 1:
			//1
			nth[5] |= 0xffffff << (i * (d << 5))
		}
	}

	return nth
}

func (m *Morton) Decode(code uint64) (result []uint32) {
	if m.Dimensions == 0 {
		return
	}

	d := uint64(m.Dimensions)
	r := make([]uint64, d)

	// Process each dimension
	for i := uint64(0); i < d; i++ {
		r[i] = (code >> i) & m.Magic[0]
		for j := uint64(0); int(j) < len(m.Magic)-1; j++ {
			r[i] = (r[i] ^ (r[i] >> ((d - 1) * (1 << j)))) & m.Magic[j+1]
		}

		result = append(result, uint32(r[i]))
	}

	return
}

func (m *Morton) Encode(vector []uint32) (result uint64, err error) {
	length := len(m.Tables)
	if length == 0 {
		err = errors.New("No lookup tables.  Please generate them via CreateTables().")
		return
	}

	if len(vector) > length {
		err = errors.New("Input vector slice length exceeds the number of lookup tables.  Please regenerate them via CreateTables()")
		return
	}

	//sort.Sort(sort.Reverse(ByUint32Index(vector)))

	for k, v := range vector {
		if v > uint32(len(m.Tables[k].Encode)-1) {
			err = errors.New(fmt.Sprint("Input vector component, ", k, " length exceeds the corresponding lookup table's size.  Please regenerate them via CreateTables() and specify the appropriate table length"))
			return
		}

		result |= m.Tables[k].Encode[v].Value
	}

	return
}

func CreateTable(index, dimensions uint8, length uint32) Table {
	t := Table{Index: index, Length: length}
	bch := make(chan Bit)

	// Build interleave queue
	for i := uint32(0); i < length; i++ {
		go func(i uint32) {
			bch <- InterleaveBits(i, uint32(index), uint32(dimensions-1))
		}(i)
	}
	// Pull from interleave queue
	for i := uint32(0); i < length; i++ {
		ib := <-bch
		t.Encode = append(t.Encode, ib)
	}
	close(bch)

	sort.Sort(ByBit(t.Encode))
	return t
}

// Interleave bits of a uint32.
func InterleaveBits(value, offset, spread uint32) Bit {
	ib := Bit{value, 0}

	// Determine the minimum number of single shifts required. There's likely a better, and more efficient, way to do this.
	n := value
	limit := uint64(0)
	for i := uint32(0); n != 0; i++ {
		n = n >> 1
		limit++
	}

	// Offset value for interleaving and reconcile types
	v, o, s := uint64(value), uint64(offset), uint64(spread)
	for i := uint64(0); i < limit; i++ {
		// Interleave bits, bit by bit.
		ib.Value |= (v & (1 << i)) << (i * s)
	}
	ib.Value = ib.Value << o

	return ib
}

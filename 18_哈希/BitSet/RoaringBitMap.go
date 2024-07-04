// RoaringBitMap (咆哮位图,RBM,高效压缩位图，高效求交集并集，是对bitmap的改进)
// https://github.com/RoaringBitmap/roaring 功能齐全，优化好
// https://github.com/fzandona/goroar 简单实现
// Use Roaring for bitmap compression whenever possible. Do not use other bitmap compression methods (Wang et al., SIGMOD 2017)
// goroar is an implementation of Roaring Bitmaps in Golang.
// Roaring bitmaps is a new form of compressed bitmaps, proposed by Daniel Lemire et. al.,
// which often offers better compression and fast access than other compressed bitmap approaches.
//
// https://cloud.tencent.com/developer/article/2207347
// container最多只能存65536bit的数据
// 当整数个数小于等于4096时，使用array container，否则使用bitmap
//
// https://www.jianshu.com/p/b09bb3e7652e
//
// RBM实现源理
// !1.分块：
// 将32位无符号整数按照高16位分桶，即最多可能有2^16=65536个桶，论文内称为container。
// 存储数据时，按照数据的高16位找到container（找不到就会新建一个），再将低16位放入container中。
// !只为用到的容器分配空间，解决了稀疏数据浪费空间的问题
// 也就是说，一个RBM就是很多container的集合。
// !2.container:
// Container只需要处理低16位的数据.
// 4096 是ArrayContainer和BitmapContainer转换的阈值.
// 还有一种RunContainer：使用Run-Length Encoding方式压缩存储的元素，对连续数据的压缩效果特别好，但如果数据比较散列，反而会更占用空间，长度没有限制
// !ArrayContainer使用2字节的short类型来存储每个元素，4096*2byte=8kb；BitmapContainer是定长2^16个bit，即bitmap固定大小8k

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"strconv"
)

func main() {
	abc350_g()
}

func demo() {
	rb1 := NewBitmapOf(1, 2, 3, 4, 5)
	rb2 := NewBitmapOf(2, 3, 4)
	rb3 := NewBitmap()

	fmt.Println("Cardinality: ", rb1.Cardinality())

	fmt.Println("Contains 3? ", rb1.Contains(3))

	rb1.And(rb2)

	rb3.Add(1)
	rb3.Add(5)

	rb3.Or(rb1)

	// prints 1, 2, 3, 4, 5
	rb3.ForEach(func(value uint32) bool {
		fmt.Println(value)
		return false
	})

	Example_stats()
}

func Example_stats() {
	rb1 := NewBitmap()
	rb2 := NewBitmap()

	for i := 0; i < 1000000; i += 2 {
		rb1.Add(uint32(i))
		rb2.Add(uint32(i + 1))
	}

	rb1.Or(rb2)
	rb1.Stats()
}

// [ABC350G] Mediator (Roaring Bitmaps)
// https://www.luogu.com.cn/problem/AT_abc350_g
// 初始时，有 N 个点，编号为 0 到 N-1,没有边存在.
// 有 Q 次操作，每次操作有三个整数 a, b, c.
// 1 u v: 在 u 和 v 之间添加一条边,保证 u 和 v 不在同一个连通块中.
// 2 u v: 询问是否存在和 u,v 都相邻的点，若存在输出编号，若不存在输出0.
func abc350_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	const MOD int = 998244353

	var n, q int
	fmt.Fscan(in, &n, &q)

	G := make([]*RoaringBitmap, n)
	for i := 0; i < n; i++ {
		G[i] = NewBitmap()
	}

	// 添加边.
	addEdge := func(a, b int32) {
		G[a].Add(uint32(b))
		G[b].Add(uint32(a))
	}

	// 查询a和b是否有共同的邻居.
	query := func(a, b int32) (int32, bool) {
		intersection := G[a].Clone()
		intersection.And(G[b])
		if intersection.Cardinality() == 0 {
			return 0, false
		}
		var res int32
		intersection.ForEach(func(u uint32) bool {
			res = int32(u)
			return true
		})
		return res, true
	}

	preRes := 0
	for i := 0; i < q; i++ {
		var A, B, C int
		fmt.Fscan(in, &A, &B, &C)
		A = 1 + (((A * (1 + preRes)) % MOD) % 2)
		B = 1 + (((B * (1 + preRes)) % MOD) % n)
		C = 1 + (((C * (1 + preRes)) % MOD) % n)
		B--
		C--
		if A == 1 {
			addEdge(int32(B), int32(C))
		} else {
			res, ok := query(int32(B), int32(C))
			if !ok {
				preRes = 0
				fmt.Fprintln(out, preRes)
			} else {
				preRes = int(res + 1)
				fmt.Fprintln(out, preRes)
			}
		}
	}
}

const (
	arrayContainerInitSize = 4
	// arrayContainerMaxSize  = 4096
	arrayContainerMaxSize = 1 << 10
)

const (
	bitmapContainerMaxCapacity = uint32(1 << 16)
	one                        = uint64(1)
)

type container interface {
	add(x uint16) container
	contains(x uint16) bool

	and(x container) container
	or(x container) container
	andNot(x container) container
	xor(x container) container

	clone() container

	getCardinality() int
	sizeInBytes() int
}

type entry struct {
	key       uint16
	container container
}

type RoaringBitmap struct {
	containers []entry
}

// NewBitmap creates a new RoaringBitmap
func NewBitmap() *RoaringBitmap {
	containers := make([]entry, 0, 4)
	return &RoaringBitmap{containers}
}

// NewBitmapOf generates a new bitmap with the specified values set to true.
// The provided values don't have to be in sorted order, but it may be preferable to sort them from a performance point of view.
func NewBitmapOf(values ...uint32) *RoaringBitmap {
	rb := NewBitmap()
	for _, value := range values {
		rb.Add(value)
	}
	return rb
}

// Add adds a uint32 value to the RoaringBitmap
func (rb *RoaringBitmap) Add(x uint32) {
	hb, lb := highlowbits(x)

	pos := rb.containerIndex(hb)
	if pos >= 0 {
		container := rb.containers[pos].container
		rb.containers[pos].container = container.add(lb)
	} else {
		ac := newArrayContainer()
		ac.add(lb)
		rb.increaseCapacity()

		loc := -pos - 1

		// insertion : shift the elements > x by one position to
		// the right and put x in it's appropriate place
		rb.containers = rb.containers[:len(rb.containers)+1]
		copy(rb.containers[loc+1:], rb.containers[loc:])
		e := entry{hb, ac}
		rb.containers[loc] = e
	}
}

// Contains checks whether the value in included, which is equivalent to checking
// if the corresponding bit is set (get in BitSet class).
func (rb *RoaringBitmap) Contains(i uint32) bool {
	pos := rb.containerIndex(highbits(i))
	if pos < 0 {
		return false
	}
	return rb.containers[pos].container.contains(lowbits(i))
}

// Cardinality returns the number of distinct integers (uint32) in the bitmap.
func (rb *RoaringBitmap) Cardinality() int {
	var cardinality int
	for _, entry := range rb.containers {
		cardinality = cardinality + entry.container.getCardinality()
	}
	return cardinality
}

// And computes the bitwise AND operation.
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) And(other *RoaringBitmap) {
	pos1 := 0
	pos2 := 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				rb.removeAtIndex(pos1)
				length1 = length1 - 1
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				pos2 = pos2 + 1
				if pos2 == length2 {
					break main
				}
				s2 = other.keyAtIndex(pos2)
			} else {
				c := rb.containers[pos1].container.and(other.containers[pos2].container)

				if c.getCardinality() > 0 {
					rb.containers[pos1].container = c
					pos1 = pos1 + 1
				} else {
					rb.removeAtIndex(pos1)
					length1 = length1 - 1
				}
				pos2 = pos2 + 1
				if (pos1 == length1) || (pos2 == length2) {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
				s2 = other.keyAtIndex(pos2)
			}
		}
	}
	rb.resize(pos1)
}

// Or computes the bitwise OR operation.
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) Or(other *RoaringBitmap) {
	pos1, pos2 := 0, 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				pos1++
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				rb.insertAt(pos1, s2, other.containers[pos2].container)
				pos1++
				length1++
				pos2++
				if pos2 == length2 {
					break main
				}
				s2 = other.containers[pos2].key
			} else {
				rb.containers[pos1].container = rb.containers[pos1].container.or(other.containers[pos2].container)
				pos1++
				pos2++
				if pos1 == length1 || pos2 == length2 {
					break main
				}
				s1 = rb.containers[pos1].key
				s2 = other.containers[pos2].key
			}
		}
	}
	if pos1 == length1 {
		rb.containers = append(rb.containers, other.containers[pos2:length2]...)
	}
}

// Xor computes the bitwise XOR operation.
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) Xor(other *RoaringBitmap) {
	pos1, pos2 := 0, 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				pos1++
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				rb.insertAt(pos1, s2, other.containers[pos2].container)
				pos1++
				length1++
				pos2++
				if pos2 == length2 {
					break main
				}
				s2 = other.containers[pos2].key
			} else {
				c := rb.containers[pos1].container.xor(other.containers[pos2].container)
				if c.getCardinality() > 0 {
					rb.containers[pos1].container = c
					pos1++
				} else {
					rb.removeAtIndex(pos1)
					length1--
				}
				pos2++
				if pos1 == length1 || pos2 == length2 {
					break main
				}
				s1 = rb.containers[pos1].key
				s2 = other.containers[pos2].key
			}
		}
	}
	if pos1 == length1 {
		rb.containers = append(rb.containers, other.containers[pos2:length2]...)
	}
}

// AndNot computes the bitwise andNot operation (difference)
// The receiving RoaringBitmap is modified - the input one is not.
func (rb *RoaringBitmap) AndNot(other *RoaringBitmap) {
	pos1, pos2 := 0, 0
	length1 := len(rb.containers)
	length2 := len(other.containers)

main:
	for pos1 < length1 && pos2 < length2 {
		s1 := rb.keyAtIndex(pos1)
		s2 := other.keyAtIndex(pos2)
		for {
			if s1 < s2 {
				pos1++
				if pos1 == length1 {
					break main
				}
				s1 = rb.keyAtIndex(pos1)
			} else if s1 > s2 {
				pos2++
				if pos2 == length2 {
					break main
				}
				s2 = other.containers[pos2].key
			} else {
				c := rb.containers[pos1].container.andNot(other.containers[pos2].container)
				if c.getCardinality() > 0 {
					rb.containers[pos1].container = c
					pos1++
				} else {
					rb.removeAtIndex(pos1)
					length1--
				}
				pos2++
				if pos1 == length1 || pos2 == length2 {
					break main
				}
				s1 = rb.containers[pos1].key
				s2 = other.containers[pos2].key
			}
		}
	}
}

// Iterator returns an iterator over the RoaringBitmap which can be used with "for range".
func (rb *RoaringBitmap) ForEach(f func(uint32) bool) {
	// iterate over data
	for _, entry := range rb.containers {
		hs := uint32(entry.key) << 16
		switch typedContainer := entry.container.(type) {
		case *arrayContainer:
			pos := 0
			for pos < typedContainer.cardinality {
				ls := typedContainer.content[pos]
				pos++
				if f(hs | uint32(ls)) {
					return
				}
			}
		case *bitmapContainer:
			i := typedContainer.nextSetBit(0)
			for i >= 0 {
				if f(hs | uint32(i)) {
					return
				}
				i = typedContainer.nextSetBit(i + 1)
			}
		}
	}
}

// Clone returns a copy of the original RoaringBitmap
func (rb *RoaringBitmap) Clone() *RoaringBitmap {
	containers := make([]entry, len(rb.containers))
	for i, value := range rb.containers {
		containers[i] = entry{value.key, value.container.clone()}
	}
	// copy(containers, rb.containers[0:])
	return &RoaringBitmap{containers}
}

func (rb *RoaringBitmap) String() string {
	var buffer bytes.Buffer
	name := []byte("RoaringBitmap[")

	buffer.Write(name)
	rb.ForEach(func(u uint32) bool {
		buffer.WriteString(strconv.Itoa(int(u)))
		buffer.WriteString(", ")
		return false
	})
	if buffer.Len() > len(name) {
		buffer.Truncate(buffer.Len() - 2) // removes the last ", "
	}
	buffer.WriteString("]")
	return buffer.String()
}

// Stats prints statistics about the Roaring Bitmap's internals.
func (rb *RoaringBitmap) Stats() {
	const output = `* Roaring Bitmap Stats *
Cardinality: {{.Cardinality}}
Size uncompressed: {{.UncompressedSize}} bytes
Size compressed: {{.CompressedSize}} bytes ({{.CompressionRate}}%)
Number of containers: {{.TotalContainers}}
    {{.TotalAC}} ArrayContainers
    {{.TotalBC}} BitmapContainers
Average entries per ArrayContainer: {{.AverageAC}}
Max entries per ArrayContainer: {{.MaxAC}}
`
	type stats struct {
		Cardinality, TotalContainers, TotalAC, TotalBC int
		AverageAC, MaxAC                               string
		CompressedSize, UncompressedSize               int
		CompressionRate                                string
	}

	var totalAC, totalBC, totalCardinalityAC int
	var maxAC int

	for _, c := range rb.containers {
		switch typedContainer := c.container.(type) {
		case *arrayContainer:
			if typedContainer.cardinality > maxAC {
				maxAC = typedContainer.cardinality
			}
			totalCardinalityAC += typedContainer.cardinality
			totalAC++
		case *bitmapContainer:
			totalBC++
		default:
		}
	}

	s := new(stats)
	s.Cardinality = rb.Cardinality()
	s.TotalContainers = len(rb.containers)
	s.TotalAC = totalAC
	s.TotalBC = totalBC
	s.CompressedSize = rb.SizeInBytes()
	s.UncompressedSize = rb.Cardinality() * 4
	s.CompressionRate = fmt.Sprintf("%3.1f",
		float32(s.CompressedSize)/float32(s.UncompressedSize)*100.0)

	if totalCardinalityAC > 0 {
		s.AverageAC = fmt.Sprintf("%6.2f", float32(totalCardinalityAC)/float32(totalAC))
		s.MaxAC = fmt.Sprintf("%d", maxAC)
	} else {
		s.AverageAC = "--"
		s.MaxAC = "--"
	}

	t := template.Must(template.New("stats").Parse(output))
	if err := t.Execute(os.Stdout, s); err != nil {
		log.Println("RoaringBitmap stats: ", err)
	}
}

func (rb *RoaringBitmap) SizeInBytes() int {
	size := 12 // size of RoaringBitmap struct
	for _, c := range rb.containers {
		size += 12 + c.container.sizeInBytes()
	}
	return size
}

func (rb *RoaringBitmap) resize(newLength int) {
	for i := newLength; i < len(rb.containers); i++ {
		rb.containers[i] = entry{}
	}
	rb.containers = rb.containers[:newLength]
}

func (rb *RoaringBitmap) keyAtIndex(pos int) uint16 {
	return rb.containers[pos].key
}

func (rb *RoaringBitmap) removeAtIndex(i int) {
	copy(rb.containers[i:], rb.containers[i+1:])
	rb.containers[len(rb.containers)-1] = entry{}
	rb.containers = rb.containers[:len(rb.containers)-1]
}

func (rb *RoaringBitmap) insertAt(i int, key uint16, c container) {
	rb.containers = append(rb.containers, entry{})
	copy(rb.containers[i+1:], rb.containers[i:])
	rb.containers[i] = entry{key, c}
}

func (rb *RoaringBitmap) containerIndex(key uint16) int {
	length := len(rb.containers)

	if length == 0 || rb.containers[length-1].key == key {
		return length - 1
	}

	return searchContainer(rb.containers, length, key)
}

func searchContainer(containers []entry, length int, key uint16) int {
	low := 0
	high := length - 1

	for low <= high {
		middleIndex := (low + high) >> 1
		middleValue := containers[middleIndex].key

		switch {
		case middleValue < key:
			low = middleIndex + 1
		case middleValue > key:
			high = middleIndex - 1
		default:
			return middleIndex
		}
	}
	return -(low + 1)
}

// increaseCapacity increases the slice capacity keeping the same length.
func (rb *RoaringBitmap) increaseCapacity() {
	length := len(rb.containers)
	if length+1 > cap(rb.containers) {
		var newCapacity int
		if length < 1024 {
			newCapacity = 2 * (length + 1)
		} else {
			newCapacity = 5 * (length + 1) / 4
		}

		newSlice := make([]entry, length, newCapacity)
		copy(newSlice, rb.containers)

		// increasing the length by 1
		rb.containers = newSlice
	}
}

type arrayContainer struct {
	cardinality int
	content     []uint16
}

func newArrayContainer() *arrayContainer {
	content := make([]uint16, arrayContainerInitSize)
	return &arrayContainer{0, content}
}

func newArrayContainerWithCapacity(capacity int) *arrayContainer {
	content := make([]uint16, capacity)
	return &arrayContainer{0, content}
}

func newArrayContainerRunOfOnes(firstOfRun, lastOfRun int) *arrayContainer {
	valuesInRange := lastOfRun - firstOfRun + 1
	content := make([]uint16, valuesInRange)
	for i := 0; i < valuesInRange; i++ {
		content[i] = uint16(firstOfRun + i)
	}
	return &arrayContainer{int(valuesInRange), content}
}

func (ac *arrayContainer) add(x uint16) container {
	if ac.cardinality >= arrayContainerMaxSize {
		bc := ac.toBitmapContainer()
		bc.add(x)
		return bc
	}

	if ac.cardinality == 0 || x > ac.content[ac.cardinality-1] {
		if ac.cardinality >= len(ac.content) {
			ac.increaseCapacity()
		}
		ac.content[ac.cardinality] = x
		ac.cardinality++
		return ac
	}

	loc := binarySearch(ac.content, ac.cardinality, x)
	if loc < 0 {
		if ac.cardinality >= len(ac.content) {
			ac.increaseCapacity()
		}
		loc = -loc - 1
		// insertion : shift the elements > x by one position to
		// the right and put x in it's appropriate place
		copy(ac.content[loc+1:], ac.content[loc:])
		ac.content[loc] = x
		ac.cardinality++
	}
	return ac
}

func (ac *arrayContainer) and(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return ac.andArray(oc)
	case *bitmapContainer:
		return ac.andBitmap(oc)
	}
	return nil
}

func (ac *arrayContainer) andArray(value2 *arrayContainer) *arrayContainer {
	value1 := ac

	cardinality, content := intersect2by2(value1.content,
		value1.cardinality, value2.content,
		value2.cardinality)

	return &arrayContainer{cardinality, content}
}

func (ac *arrayContainer) andBitmap(bc *bitmapContainer) *arrayContainer {
	return bc.andArray(ac)
}

func (ac *arrayContainer) or(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return ac.orArray(oc)
	case *bitmapContainer:
		return ac.orBitmap(oc)
	}
	return nil
}

func (ac *arrayContainer) orArray(other *arrayContainer) container {
	totalCardinality := ac.cardinality + other.cardinality
	if totalCardinality > arrayContainerMaxSize {
		bc := newBitmapContainer()
		for i := 0; i < other.cardinality; i++ {
			bc.add(other.content[i])
		}
		for i := 0; i < ac.cardinality; i++ {
			bc.add(ac.content[i])
		}
		if bc.cardinality <= arrayContainerMaxSize {
			return bc.toArrayContainer()
		}
		return bc
	}
	answer := arrayContainer{}
	pos, content := union2by2(ac.content, ac.cardinality, other.content, other.cardinality, totalCardinality)
	answer.cardinality = pos
	answer.content = content
	return &answer
}

func (ac *arrayContainer) orBitmap(bc *bitmapContainer) container {
	return bc.or(ac)
}

func (ac *arrayContainer) xor(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return ac.xorArray(oc)
	case *bitmapContainer:
		return ac.xorBitmap(oc)
	}
	return nil
}

func (ac *arrayContainer) xorArray(other *arrayContainer) container {
	totalCardinality := ac.cardinality + other.cardinality
	if totalCardinality > arrayContainerMaxSize {
		bc := newBitmapContainer()
		for i := 0; i < other.cardinality; i++ {
			index := other.content[i] >> 6
			bc.bitmap[index] ^= one << other.content[i]
		}
		for i := 0; i < ac.cardinality; i++ {
			index := ac.content[i] >> 6
			bc.bitmap[index] ^= one << ac.content[i]
		}
		for _, bitmap := range bc.bitmap {
			bc.cardinality += countBits(bitmap)
		}
		if bc.cardinality <= arrayContainerMaxSize {
			return bc.toArrayContainer()
		}
		return bc
	}
	answer := arrayContainer{}
	pos, content := exclusiveUnion2by2(ac.content, ac.cardinality, other.content, other.cardinality, totalCardinality)
	answer.cardinality = pos
	answer.content = content
	return &answer
}

func (ac *arrayContainer) xorBitmap(bc *bitmapContainer) container {
	return bc.xor(ac)
}

func (ac *arrayContainer) andNot(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return ac.andNotArray(oc)
	case *bitmapContainer:
		return ac.andNotBitmap(oc)
	}
	return nil
}

func (ac *arrayContainer) andNotArray(value2 *arrayContainer) *arrayContainer {
	cardinality, content := difference(ac.content, ac.cardinality,
		value2.content, value2.cardinality)

	return &arrayContainer{cardinality, content}
}

func (ac *arrayContainer) andNotBitmap(value2 *bitmapContainer) *arrayContainer {
	pos := 0
	for k := 0; k < ac.cardinality; k++ {
		if !value2.contains(ac.content[k]) {
			ac.content[pos] = ac.content[k]
			pos++
		}
	}
	ac.cardinality = pos

	return ac
}

func (ac *arrayContainer) contains(x uint16) bool {
	return binarySearch(ac.content, ac.cardinality, x) >= 0
}

func (ac *arrayContainer) clear() {
	ac.content = make([]uint16, arrayContainerInitSize)
	ac.cardinality = 0
}

func (ac *arrayContainer) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	bc.loadData(ac)
	return bc
}

func (ac *arrayContainer) getCardinality() int {
	return ac.cardinality
}

func (ac *arrayContainer) arraySizeInBytes() int {
	return ac.cardinality * 2
}

func (ac *arrayContainer) increaseCapacity() {
	length := len(ac.content)
	var newLength int
	switch {
	case length < 64:
		newLength = length * 2
	case length < 1024:
		newLength = length * 3 / 2
	default:
		newLength = length * 5 / 4
	}
	if newLength > arrayContainerMaxSize {
		newLength = arrayContainerMaxSize
	}
	newSlice := make([]uint16, newLength)
	copy(newSlice, ac.content)
	ac.content = newSlice
}

func (ac *arrayContainer) sizeInBytes() int {
	return ac.cardinality*2 + 16
}

func (ac *arrayContainer) clone() container {
	newContent := make([]uint16, len(ac.content), cap(ac.content))
	copy(newContent, ac.content)
	return &arrayContainer{ac.cardinality, newContent}
}

type bitmapContainer struct {
	cardinality int
	bitmap      []uint64
}

var _ container = (*bitmapContainer)(nil)

func newBitmapContainer() *bitmapContainer {
	return &bitmapContainer{0, make([]uint64, bitmapContainerMaxCapacity>>6)}
}

func (bc *bitmapContainer) loadData(ac *arrayContainer) {
	bc.cardinality = ac.cardinality
	for i := 0; i < ac.cardinality; i++ {
		bc.bitmap[uint32(ac.content[i])>>6] |= one << (ac.content[i] & 63)
	}
}

func (bc *bitmapContainer) add(i uint16) container {
	x := uint32(i)
	index := x >> 6
	mod := x & 63
	previous := bc.bitmap[index]
	bc.bitmap[index] |= one << mod
	bc.cardinality += int((previous ^ bc.bitmap[index]) >> mod)
	return bc
}

func (bc *bitmapContainer) and(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.andArray(oc)
	case *bitmapContainer:
		return bc.andBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) andArray(value2 *arrayContainer) *arrayContainer {
	answer := make([]uint16, value2.cardinality)

	cardinality := 0
	for k := 0; k < value2.cardinality; k++ {
		if bc.contains(value2.content[k]) {
			answer[cardinality] = value2.content[k]
			cardinality++
		}
	}

	return &arrayContainer{cardinality, answer[:cardinality]}
}

func (bc *bitmapContainer) andBitmap(value2 *bitmapContainer) container {
	newCardinality := 0
	for k, v := range bc.bitmap {
		newCardinality += countBits(v & value2.bitmap[k])
	}

	if newCardinality > arrayContainerMaxSize {
		answer := newBitmapContainer()
		for k, v := range bc.bitmap {
			answer.bitmap[k] = v & value2.bitmap[k]
		}
		answer.cardinality = newCardinality
		return answer

	}
	content := fillArrayAND(bc.bitmap, value2.bitmap, newCardinality)
	return &arrayContainer{newCardinality, content}
}

func (bc *bitmapContainer) or(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.orArray(oc)
	case *bitmapContainer:
		return bc.orBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) orArray(ac *arrayContainer) *bitmapContainer {
	answer := bc.clone().(*bitmapContainer)
	for i := 0; i < ac.cardinality; i++ {
		answer.add(ac.content[i])
	}
	return answer
}

func (bc *bitmapContainer) orBitmap(other *bitmapContainer) container {
	answer := newBitmapContainer()

	for i := 0; i < len(bc.bitmap); i++ {
		answer.bitmap[i] = bc.bitmap[i] | other.bitmap[i]
		answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
	}
	return answer
}
func (bc *bitmapContainer) xor(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.xorArray(oc)
	case *bitmapContainer:
		return bc.xorBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) xorArray(ac *arrayContainer) container {
	answer := bc.clone().(*bitmapContainer)
	for i := 0; i < ac.cardinality; i++ {
		v := ac.content[i]
		mod := v & 63
		index := v >> 6
		shift := one << v
		answer.cardinality += 1 - 2*int((answer.bitmap[index]&shift)>>mod)
		answer.bitmap[index] ^= shift

	}
	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}

func (bc *bitmapContainer) xorBitmap(other *bitmapContainer) container {
	answer := newBitmapContainer()

	for i := 0; i < len(bc.bitmap); i++ {
		answer.bitmap[i] = bc.bitmap[i] ^ other.bitmap[i]
		answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
	}

	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}
func (bc *bitmapContainer) andNot(other container) container {
	switch oc := other.(type) {
	case *arrayContainer:
		return bc.andNotArray(oc)
	case *bitmapContainer:
		return bc.andNotBitmap(oc)
	}
	return nil
}

func (bc *bitmapContainer) andNotArray(ac *arrayContainer) container {
	answer := bc.clone().(*bitmapContainer)
	for i := 0; i < ac.cardinality; i++ {
		v := ac.content[i]
		mod := v & 63
		index := v >> 6
		shift := one << v
		answer.bitmap[index] = answer.bitmap[index] & (^shift)
		answer.cardinality -= int((answer.bitmap[index] ^ bc.bitmap[index]) >> mod)
	}
	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}

func (bc *bitmapContainer) andNotBitmap(other *bitmapContainer) container {
	answer := newBitmapContainer()

	for i := 0; i < len(bc.bitmap); i++ {
		answer.bitmap[i] = bc.bitmap[i] & (^other.bitmap[i])
		answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
	}

	if answer.cardinality <= arrayContainerMaxSize {
		return answer.toArrayContainer()
	}
	return answer
}

func (bc *bitmapContainer) clone() container { //*bitmapContainer {
	bitmap := make([]uint64, len(bc.bitmap))
	copy(bitmap, bc.bitmap)
	return &bitmapContainer{bc.cardinality, bitmap}
}

func (bc *bitmapContainer) contains(x uint16) bool {
	return bc.bitmap[uint32(x)>>6]&(one<<(x&63)) != 0
}

// nextSetBit finds the index of the next set bit greater or equal to i.
// It returns -1 if none is found.
func (bc *bitmapContainer) nextSetBit(i int) int {
	x := i >> 6
	if x >= len(bc.bitmap) {
		return -1
	}

	w := bc.bitmap[x]
	w = w >> (uint(i) & 63)
	if w != 0 {
		return i + trailingZeros(w)
	}

	x = x + 1
	for ; x < len(bc.bitmap); x++ {
		if bc.bitmap[x] != 0 {
			return x<<6 + trailingZeros(bc.bitmap[x])
		}
	}

	return -1
}

func (bc *bitmapContainer) getCardinality() int {
	return bc.cardinality
}

func (bc *bitmapContainer) toArrayContainer() *arrayContainer {
	container := make([]uint16, bc.cardinality)
	pos := 0
	for k := 0; k < len(bc.bitmap); k++ {
		bitset := bc.bitmap[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k<<6 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return &arrayContainer{bc.cardinality, container}
}

func (bc *bitmapContainer) sizeInBytes() int {
	return 16 + len(bc.bitmap)*8
}

func binarySearch(array []uint16, length int, k uint16) int {
	low := 0
	high := length - 1

	for low <= high {
		middleIndex := (low + high) >> 1
		middleValue := array[middleIndex]

		switch {
		case middleValue < k:
			low = middleIndex + 1
		case middleValue > k:
			high = middleIndex - 1
		default:
			return middleIndex
		}
	}
	return -(low + 1)
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func intersect2by2(set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	if length1<<6 < length2 {
		return oneSidedGallopingIntersect2by2(set1, length1, set2, length2)
	}

	if length2<<6 < length1 {
		return oneSidedGallopingIntersect2by2(set2, length2, set1, length1)
	}

	return localIntersect2by2(set1, length1, set2, length2)
}

func localIntersect2by2(set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	if 0 == length1 || 0 == length2 {
		return 0, make([]uint16, 0)
	}

	finalLength := min(length1, length2)
	buffer := make([]uint16, finalLength)
	k1, k2, pos := 0, 0, 0

Mainwhile:
	for {
		if set2[k2] < set1[k1] {
			for {
				k2++
				if k2 == length2 {
					break Mainwhile
				}
				if set2[k2] >= set1[k1] {
					break
				}
			}
		}
		if set1[k1] < set2[k2] {
			for {
				k1++
				if k1 == length1 {
					break Mainwhile
				}
				if set1[k1] >= set2[k2] {
					break
				}
			}
		} else {
			buffer[pos] = set1[k1]
			pos++
			k1++
			if k1 == length1 {
				break
			}
			k2++
			if k2 == length2 {
				break
			}
		}
	}
	return pos, buffer[:pos]
}

func oneSidedGallopingIntersect2by2(
	smallSet []uint16, smallLength int,
	largeSet []uint16, largeLength int) (int, []uint16) {

	if 0 == smallLength {
		return 0, make([]uint16, 0)
	}

	buffer := make([]uint16, smallLength)
	k1, k2, pos := 0, 0, 0

	for {
		if largeSet[k1] < smallSet[k2] {
			k1 = advanceUntil(largeSet, k1, largeLength, smallSet[k2])
			if k1 == largeLength {
				break
			}
		}
		if smallSet[k2] < largeSet[k1] {
			k2++
			if k2 == smallLength {
				break
			}
		} else { // (set2[k2] == set1[k1])
			buffer[pos] = smallSet[k2]
			pos++
			k2++
			if k2 == smallLength {
				break
			}
			k1 = advanceUntil(largeSet, k1, largeLength, smallSet[k2])
			if k1 == largeLength {
				break
			}
		}

	}
	return pos, buffer[:pos]
}

// Find the smallest integer larger than pos such that array[pos]>= min.
// If none can be found, return length. Based on code by O. Kaser.
func advanceUntil(array []uint16, pos, length int, min uint16) int {
	lower := pos + 1

	// special handling for a possibly common sequential case
	if lower >= length || array[lower] >= min {
		return lower
	}

	spansize := 1 // could set larger  bootstrap an upper limit

	for (lower+spansize) < length && array[lower+spansize] < min {
		spansize *= 2
	}
	var upper int
	if lower+spansize < length {
		upper = lower + spansize
	} else {
		upper = length - 1
	}

	// maybe we are lucky (could be common case when the seek ahead
	// expected to be small and sequential will otherwise make us look bad)
	if array[upper] == min {
		return upper
	}

	if array[upper] < min { // means array has no item >= min
		return length
	}

	// we know that the next-smallest span was too small
	lower += (spansize / 2)

	// else begin binary search
	// invariant: array[lower]<min && array[upper]>min
	for lower+1 != upper {
		mid := (lower + upper) / 2
		if array[mid] == min {
			return mid
		} else if array[mid] < min {
			lower = mid
		} else {
			upper = mid
		}
	}
	return upper
}

func difference(
	set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	k1, k2, pos := 0, 0, 0

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		return 0, make([]uint16, 0)
	}

	buffer := make([]uint16, length1)

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos++
			k1++
			if k1 >= length1 {
				break
			}
		} else if set1[k1] == set2[k2] {
			k1++
			k2++
			if k1 >= length1 {
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos++
				}
				break
			}
		} else { // if (val1>val2)
			k2++
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos++
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

// http://en.wikipedia.org/wiki/Hamming_weight
func countBits(i uint64) int {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	result := (((i + (i >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56
	return int(result)
}

func highbits(x uint32) uint16 {
	return uint16(x >> 16)
}

func lowbits(x uint32) uint16 {
	return uint16(x & 0xFFFF)
}

func highlowbits(x uint32) (uint16, uint16) {
	return highbits(x), lowbits(x)
}

func fillArrayAND(bitmap1, bitmap2 []uint64, newCardinality int) []uint16 {
	pos := 0

	if len(bitmap1) != len(bitmap2) {
		panic("Bitmaps have different length - not supported.")
	}

	container := make([]uint16, newCardinality)
	for k := 0; k < len(bitmap1); k++ {
		bitset := bitmap1[k] & bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k<<6 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return container
}

func fillArrayXOR(bitmap1, bitmap2 []uint64, newCardinality int) []uint16 {
	pos := 0

	if len(bitmap1) != len(bitmap2) {
		panic("Bitmaps have different length - not supported.")
	}

	container := make([]uint16, newCardinality)
	for k := 0; k < len(bitmap1); k++ {
		bitset := bitmap1[k] ^ bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k<<6 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return container
}

// http://graphics.stanford.edu/~seander/bithacks.html#ZerosOnRightBinSearch
func trailingZeros(v uint64) int {
	if v&0x1 == 1 {
		return 0
	}

	c := 1

	if (v & 0xFFFFFFFF) == 0 {
		v = v >> 32
		c = c + 32
	}

	if (v & 0xFFFF) == 0 {
		v = v >> 16
		c = c + 16
	}
	if (v & 0xFF) == 0 {
		v = v >> 8
		c = c + 8
	}
	if (v & 0xF) == 0 {
		v = v >> 4
		c = c + 4
	}
	if (v & 0x3) == 0 {
		v = v >> 2
		c = c + 2
	}

	return c - int(v&0x1)
}

// Unite two sorted lists
func union2by2(set1 []uint16, length1 int,
	set2 []uint16, length2, bufferSize int) (int, []uint16) {

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		buffer := make([]uint16, length2)
		copy(buffer, set2)
		return length2, buffer
	}

	buffer := make([]uint16, bufferSize)

	k1, k2, pos := 0, 0, 0

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
		} else if set1[k1] == set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			k2 = k2 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		} else {
			buffer[pos] = set2[k2]
			pos = pos + 1
			k2 = k2 + 1
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

// Compute the exclusive union of two sorted lists
func exclusiveUnion2by2(set1 []uint16, length1 int,
	set2 []uint16, length2, bufferSize int) (int, []uint16) {

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		buffer := make([]uint16, length2)
		copy(buffer, set2)
		return length2, buffer
	}

	buffer := make([]uint16, bufferSize)

	k1, k2, pos := 0, 0, 0

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
		} else if set1[k1] == set2[k2] {
			buffer[pos] = set1[k1]
			k1 = k1 + 1
			k2 = k2 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		} else {
			buffer[pos] = set2[k2]
			pos = pos + 1
			k2 = k2 + 1
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

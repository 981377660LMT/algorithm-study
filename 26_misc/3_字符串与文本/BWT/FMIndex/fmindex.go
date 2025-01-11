// TODO: 大致结构同 https://github.com/rossmerr/fm-index/tree/77e6c665a79e
// 实现需要改进。不可用。

package main

import (
	"fmt"
	"math"
	"math/bits"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// #region BWT
	index, err := NewFMIndex("The quick brown fox jumps over the lazy dog", WithCompression(2))
	if err != nil {
		panic(err)
	}

	fmt.Println(index.Count("the"))
	fmt.Println(index.Locate("the"))
	fmt.Println(index.Extract(0, 3))
}

// #region FMIndex
type FMIndex struct {
	// first column of the BWT matrix
	f *WaveletTree
	// last column of the BWT matrix
	l *WaveletTree
	// suffix array
	suffix Suffix
	// prefix tree
	prefix          *Prefix
	caseinsensitive bool
	compression     *int
}

type FMIndexOption func(f *FMIndex)

func WithCaseInsensitive(caseinsensitive bool) FMIndexOption {
	return func(f *FMIndex) {
		f.caseinsensitive = caseinsensitive
	}
}

func WithPrefixTree(prefix Prefix) FMIndexOption {
	return func(f *FMIndex) {
		f.prefix = &prefix
	}
}

func WithCompression(compression int) FMIndexOption {
	return func(s *FMIndex) {
		if compression >= 2 {
			s.compression = &compression
		}
	}
}

func NewFMIndex(text string, opts ...FMIndexOption) (*FMIndex, error) {
	index := &FMIndex{}

	for _, opt := range opts {
		opt(index)
	}

	if index.caseinsensitive {
		text = strings.ToUpper(text)
	}

	if index.compression == nil {
		two := 2
		index.compression = &two
	}

	first, last, sa, err := FirstLastSuffix(text, WithSampleRate(*index.compression))
	if err != nil {
		return nil, err
	}

	if index.prefix == nil {
		index.prefix = NewHuffmanCodeTree(first)
	}

	index.suffix = sa
	index.f = NewWaveletTree(first, index.prefix)
	index.l = NewWaveletTree(last, index.prefix)

	return index, nil
}

func (s *FMIndex) Extract(offset, length int) string {
	result := []rune{}

	mappedSuffix := map[int]int{}
	iterator := s.suffix.Enumerate()
	for iterator.HasNext() {
		k, i := iterator.Next()
		mappedSuffix[k] = i
	}

	index := 0
	ok := false
	for i := offset; i < offset+length; i++ {
		if index, ok = mappedSuffix[i]; !ok {
			index = i - 1
			hops := 1

			for {
				if i, ok := mappedSuffix[index]; ok {
					index = i
					break
				}
				index--
				hops++
			}

			for i := 0; i < hops; i++ {
				r := s.f.Access(index)
				rank, _ := s.f.Rank(r, index)
				index = s.l.Select(r, rank)
			}
		}

		r := s.f.Access(index)
		result = append(result, r)

	}
	return string(result)
}

func (s *FMIndex) Count(pattern string) int {
	f, l := s.query(pattern)
	return l - f
}

func (s *FMIndex) Locate(pattern string) []int {
	f, l := s.query(pattern)
	result := []int{}
	for i := f; i < l; i++ {
		index := s.walkToNearest(i, 0)
		r := s.suffix.Get(index)
		result = append(result, r)
	}
	if f == l {
		index := s.walkToNearest(f, 0)

		r := s.suffix.Get(index)
		result = append(result, r)
	}
	return result
}

func (s *FMIndex) walkToNearest(index, count int) int {
	b := s.suffix.Has(index)
	if b {
		return index + count
	}
	count++
	a := s.l.Access(index)
	r, _ := s.l.Rank(a, index)
	nextIndex := s.f.Select(a, r)
	return s.walkToNearest(nextIndex, count)
}

func (s *FMIndex) query(pattern string) (top, bottom int) {
	if s.caseinsensitive {
		pattern = strings.ToUpper(pattern)
	}

	length := len(pattern)

	// // look at the pattern in reverse order
	next := rune(pattern[length-1])

	n1, _ := s.f.Rank(next, 0)
	top = s.f.Select(next, n1)
	n2, _ := s.f.Rank(next, s.l.Length())
	bottom = s.f.Select(next, n2+1)

	i := length - 2
	for i >= 0 && bottom >= top {
		next = rune(pattern[i])
		n1, _ := s.l.Rank(next, top)
		n2, _ := s.l.Rank(next, bottom)
		skip := s.f.Select(next, 0)
		top = (n1 + skip)
		bottom = (n2 + skip)
		i--
	}

	return
}

func (s *FMIndex) CaseInsensitive() bool {
	return s.caseinsensitive
}

func (s *FMIndex) PrefixTree() *Prefix {
	return s.prefix

}

func (s *FMIndex) Compression() int {
	return *s.compression

}

// #endregion

// #region Suffixarray
type SuffixConstraints interface {
	SuffixArray | SampleSuffixArray
}

type Suffix interface {
	Has(index int) bool
	Get(index int) int
	Set(index, value int)
	Enumerate() SuffixIterator
}

type SuffixIterator interface {
	HasNext() bool
	Next() (int, int)
}

type SuffixArray struct {
	sa      []int
	version int
}

func NewSuffixArray(size int) Suffix {
	suffix := &SuffixArray{
		sa: make([]int, size),
	}
	return suffix
}

func (s *SuffixArray) Get(index int) int {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}
	return s.sa[index]
}

func (s *SuffixArray) Set(index, value int) {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}
	s.sa[index] = value
	s.version++
}

func (s *SuffixArray) Has(index int) bool {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}
	return true
}

func (s *SuffixArray) Length() int {
	return len(s.sa)
}

func (s *SuffixArray) Enumerate() SuffixIterator {
	return NewSuffixArrayIterator(s)
}

func NewSuffixArrayIterator(suffix *SuffixArray) *SuffixArrayIterator {
	return &SuffixArrayIterator{
		suffix:     suffix,
		indexStart: 0,
		indexEnd:   suffix.Length(),
		version:    suffix.version,
	}
}

type SuffixArrayIterator struct {
	suffix     *SuffixArray
	version    int
	indexStart int
	indexEnd   int
}

func (s *SuffixArrayIterator) HasNext() bool {
	return s.indexStart < s.indexEnd
}

func (s *SuffixArrayIterator) Next() (int, int) {
	if s.version != s.suffix.version {
		panic("version failed")
	}
	if s.indexStart < s.suffix.Length() {
		index := s.indexStart
		currentElement := s.suffix.Get(index)
		s.indexStart++
		return currentElement, index
	}
	s.indexStart = s.suffix.Length()
	return 0, s.indexStart
}

type SampleSuffixArray struct {
	sa         []int
	sampleRate int
	size       int
	length     int
	version    int
}

func NewSampleSuffixArray(size, sampleRate int) Suffix {
	l := int(math.Ceil(float64(size) / float64(sampleRate)))
	suffix := &SampleSuffixArray{
		sampleRate: sampleRate,
		sa:         make([]int, l),
		size:       l,
		length:     size,
	}
	return suffix
}

func (s *SampleSuffixArray) Has(index int) bool {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}
	return index%s.sampleRate == 0
}

func (s *SampleSuffixArray) Get(index int) int {
	return s.walk(index, 0)
}

func (s *SampleSuffixArray) get(index, count int) int {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}
	if index%s.sampleRate == 0 {
		i := index / s.sampleRate
		return s.sa[i] + count
	}
	return -1
}

func (s *SampleSuffixArray) walk(i, count int) int {
	if s.Has(i) {
		return s.get(i, count)
	} else {
		return s.walk(i-1, count+1)
	}
}

func (s *SampleSuffixArray) Set(index, value int) {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}
	if index%s.sampleRate == 0 {
		i := index / s.sampleRate
		s.sa[i] = value
	}
	s.version++
}

func (s *SampleSuffixArray) Length() int {
	return s.length
}

func (s *SampleSuffixArray) Enumerate() SuffixIterator {
	return NewSampleSuffixArrayIterator(s)
}

func NewSampleSuffixArrayIterator(suffix *SampleSuffixArray) *SampleSuffixArrayIterator {
	return &SampleSuffixArrayIterator{
		suffix:       suffix,
		indexStart:   0,
		currentIndex: 0,
		indexEnd:     suffix.size,
		version:      suffix.version,
	}
}

type SampleSuffixArrayIterator struct {
	suffix       *SampleSuffixArray
	version      int
	indexStart   int
	indexEnd     int
	currentIndex int
}

func (s *SampleSuffixArrayIterator) HasNext() bool {
	return s.indexStart < s.indexEnd
}

func (s *SampleSuffixArrayIterator) Next() (int, int) {
	if s.version != s.suffix.version {
		panic("version failed")
	}
	if s.indexStart < s.indexEnd {
		currentIndex := s.currentIndex
		currentElement := s.suffix.Get(currentIndex)
		s.indexStart++
		s.currentIndex += s.suffix.sampleRate
		return currentElement, currentIndex
	}
	s.indexStart = s.indexEnd
	return 0, s.indexStart
}

// #endregion

// #region BWT
// Matrix of the BWT
func Matrix(str string, options ...func(*OptionsBwt)) ([][]rune, error) {
	opts := buildOptions(options)

	if strings.Contains(str, opts.ext) {
		err := fmt.Errorf("input string cannot contain EXT character")
		return [][]rune{}, err
	}

	str = str + opts.ext
	size := len(str)

	matrix := make(matrix, size)

	for i := size - 1; i >= 0; i-- {
		matrix[i] = append([]rune(str[i:]), []rune(str[:i])...)
	}
	sort.Sort(&matrix)

	return matrix, nil
}

// Last column of the BWT matrix
// Optimized to only do the last rotation
func Last(str string, options ...func(*OptionsBwt)) ([]rune, error) {
	opts := buildOptions(options)

	appendFirst := func(i int, r rune) {
	}

	set := func(s, o int) {
	}
	last, err := firstLastSuffix(str, appendFirst, set, opts)
	return last, err
}

// First and Last column of the BWT matrix
// Optimized to only do the last rotation
func FirstLast(str string, options ...func(*OptionsBwt)) ([]rune, []rune, error) {
	opts := buildOptions(options)

	size := len(str + opts.ext)
	first := make([]rune, size)

	appendFirst := func(i int, r rune) {
		first[i] = r
	}

	set := func(s, o int) {
	}
	last, err := firstLastSuffix(str, appendFirst, set, opts)
	return first, last, err
}

// First and Last column of the BWT matrix with a SuffixArray
// The SuffixArray returns the offset of the original string relative to the first column of the BWT matrix
// Optimized to only do the last rotation
func FirstLastSuffix(str string, options ...func(*OptionsBwt)) ([]rune, []rune, Suffix, error) {
	opts := buildOptions(options)

	size := len(str + opts.ext)

	sa := NewSampleSuffixArray(size, opts.SampleRate())

	first := make([]rune, size)

	appendFirst := func(i int, r rune) {
		first[i] = r
	}

	last, err := firstLastSuffix(str, appendFirst, sa.Set, opts)
	return first, last, sa, err
}

func firstLastSuffix(str string, appendFirst func(i int, r rune), set func(index, value int), opts *OptionsBwt) ([]rune, error) {
	if strings.Contains(str, opts.ext) {
		err := fmt.Errorf("input string cannot contain EXT character")
		return []rune{}, err
	}

	str = str + opts.ext
	size := len(str)

	suffixes := make([]string, size)
	for i := 0; i < size; i++ {
		suffixes[i] = str[i:]
	}

	sort.Strings(suffixes)

	last := make([]rune, size)
	for i := 0; i < size; i++ {
		appendFirst(i, rune(suffixes[i][0]))
		s := size - len(suffixes[i])
		mod := (s + size - 1) % size
		last[i] = rune(str[mod])
		set(i, s)
	}

	return last, nil
}

// Reverse the BWT transformation, last column of the BWT matrix back to the original text
func Reverse(str string, options ...func(*OptionsBwt)) string {
	opts := buildOptions(options)

	size := len(str)
	table := make([]string, size)
	for range table {
		for i := 0; i < size; i++ {
			table[i] = str[i:i+1] + table[i]
		}
		sort.Strings(table)
	}
	for _, row := range table {
		if strings.HasPrefix(row, opts.ext) {
			return row[1:]
		}
	}
	return ""
}

type matrix [][]rune

func (m matrix) Len() int { return len(m) }
func (m matrix) Less(i, j int) bool {
	for x := range m[i] {
		if m[i][x] == m[j][x] {
			continue
		}
		return m[i][x] < m[j][x]
	}
	return false
}

func (m matrix) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type OptionsBwt struct {
	// End-of-Text code
	ext string
	// Sample rate of suffix index (compresses the array)
	sampleRate int
}

// End-of-Text code
func WithEndOfText(ext rune) func(*OptionsBwt) {
	return func(s *OptionsBwt) {
		s.ext = string(ext)
	}
}

// Sample rate of suffix index (compresses the array)
func WithSampleRate(mod int) func(*OptionsBwt) {
	return func(s *OptionsBwt) {
		s.sampleRate = mod
	}
}

func (s *OptionsBwt) SampleRate() int {
	return s.sampleRate
}

func buildOptions(options []func(*OptionsBwt)) *OptionsBwt {
	opts := &OptionsBwt{}
	for _, o := range options {
		o(opts)
	}

	if opts.sampleRate <= 0 {
		opts.sampleRate = 1
	}

	if opts.ext == "" {
		opts.ext = "\003"
	}

	return opts
}

// #endregion

// #region waveletTree

// WaveletTree is a succinct data structure to store strings in compressed space.
type WaveletTree struct {
	root   *node
	prefix *Prefix
	n      int // Length of the root bitvector
}

func NewWaveletTree(value []rune, prefix *Prefix) *WaveletTree {
	root := buildNode(value, prefix)
	tree := &WaveletTree{
		root:   root,
		prefix: prefix,
		n:      root.Length(),
	}
	return tree
}

func NewBalancedWaveletTree(value []rune) *WaveletTree {
	prefix := NewBinaryTree(value)
	root := buildNode(value, prefix)
	tree := &WaveletTree{
		root:   root,
		prefix: prefix,
		n:      root.Length(),
	}
	return tree
}

func NewHuffmanCodeWaveletTree(value []rune) *WaveletTree {
	prefix := NewHuffmanCodeTree(value)
	root := buildNode(value, prefix)
	tree := &WaveletTree{
		root:   root,
		prefix: prefix,
		n:      root.Length(),
	}
	return tree
}

// Access gets the run at the index.
func (wt *WaveletTree) Access(index int) rune {
	return wt.root.Access(index)
}

// Rank counts the number of times the rune occurs up to but not including the offset.
func (wt *WaveletTree) Rank(c rune, offset int) (int, error) {
	prefix := wt.prefix.Get(c)
	if prefix == nil {
		return 0, fmt.Errorf("rune '%v' code %v not found in prefix", string(c), c)
	}
	return wt.root.Rank(prefix, offset), nil
}

// Select returns the index of the rune with the given rank
func (wt *WaveletTree) Select(c rune, rank int) int {
	prefix := wt.prefix.Get(c)
	start := wt.root.Walk(prefix)
	return start.Select(prefix, rank)
}

func (wt *WaveletTree) Length() int {
	return wt.n
}

func (wt WaveletTree) String() string {
	str := ""
	str += fmt.Sprintf(" length: %v", wt.n)

	if wt.root != nil {
		str += fmt.Sprintf(", root: %s", wt.root)
	}

	if wt.prefix != nil {
		str += fmt.Sprintf(", prefix: %+v", wt.prefix)
	}

	return fmt.Sprintf("{%s }", str)
}

type node struct {
	parent *node
	left   *node
	right  *node
	value  *rune
	vector *BitVector
}

func buildNode(data []rune, prefix *Prefix) *node {
	return buildChildNode(data, prefix, nil, 0)
}
func buildChildNode(data []rune, prefix *Prefix, parent *node, depth int) *node {
	vector := NewBitVector(len(data))
	left, right := []rune{}, []rune{}

	for i, entry := range data {

		partitions := prefix.Get(entry)

		if depth >= partitions.Length() {
			return nil
		}

		c := partitions.Get(depth)
		vector.Set(i, c)
		if c {
			right = append(right, entry)
		} else {
			left = append(left, entry)
		}
	}

	t := &node{
		vector: vector,
		parent: parent,
	}

	if len(left) > 0 {
		n := buildChildNode(left, prefix, t, depth+1)

		if n != nil {
			t.left = n
		} else {
			t.left = &node{
				value:  &left[0],
				parent: t,
			}
		}

	}
	if len(right) > 0 {
		n := buildChildNode(right, prefix, t, depth+1)

		if n != nil {
			t.right = n
		} else {
			t.right = &node{
				value:  &right[0],
				parent: t,
			}
		}
	}
	return t
}

func (t *node) Length() int {
	return t.vector.Length()
}

func (t *node) isLeaf() bool {
	return t.vector == nil
}

func (t *node) Access(i int) rune {
	if t.isLeaf() {
		return rune(*t.value)
	}
	c := t.vector.Get(i)
	rank := t.vector.Rank(c, i)
	if c {
		return t.right.Access(rank)
	} else {
		return t.left.Access(rank)
	}
}

func (t *node) Rank(prefix *BitVector, offset int) int {

	c := prefix.Get(0)

	rank := t.vector.Rank(c, offset)

	vector := NewBitVector(prefix.Length() - 1)
	if prefix.Length() > 1 {
		prefix.Copy(vector, 1, prefix.Length())
	} else {
		return rank
	}

	if c {
		return t.right.Rank(vector, rank)
	} else {
		return t.left.Rank(vector, rank)
	}
}

func (t *node) Walk(prefix *BitVector) *node {

	if t.isLeaf() {
		return t
	}

	c := prefix.Get(0)

	vector := NewBitVector(prefix.Length() - 1)
	if prefix.Length() > 1 {
		prefix.Copy(vector, 1, prefix.Length())
	}

	if c {
		return t.right.Walk(vector)
	} else {
		return t.left.Walk(vector)
	}
}

func (t *node) Select(prefix *BitVector, rank int) int {

	if t.isLeaf() {
		return t.parent.Select(prefix, rank)
	}
	i := prefix.Get(prefix.Length() - 1)

	r := t.vector.Select(i, rank)

	if t.parent != nil {

		vector := NewBitVector(prefix.Length() - 1)
		prefix.Copy(vector, 0, prefix.Length()-1)

		return t.parent.Select(vector, r)
	}
	return r

}

func (t node) String() string {
	str := ""
	if t.left != nil {
		str += fmt.Sprintf(" left: %s", t.left)
	}
	if t.right != nil {
		str += fmt.Sprintf(" right: %s", t.right)
	}

	if t.value != nil {
		str += fmt.Sprintf(" value: %s", string(*t.value))
	}

	return fmt.Sprintf("{%s }", str)
}

// #endregion

// #region prefixTree
const ext = '\003'

type BinaryTree struct {
	Left  *BinaryTree
	Right *BinaryTree
	Value *rune
}

func NewBinaryTree(value []rune) *Prefix {
	runeFrequencies, keys := binaryCount(value)
	binaryList := rankByBinaryCount(runeFrequencies, keys)
	tree := buildBinaryTree(binaryList)
	return tree.prefix()
}

func binaryCount(value []rune) (map[rune]int, []rune) {
	runeFrequencies := make(map[rune]int)
	keys := make([]rune, 0)

	for _, r := range value {
		if _, ok := runeFrequencies[r]; !ok {
			runeFrequencies[r] = len(runeFrequencies)
			keys = append(keys, r)
		}
	}

	return runeFrequencies, keys
}

type binaryList []*BinaryTree

func rankByBinaryCount(runeFrequencies map[rune]int, keys []rune) binaryList {
	list := make(binaryList, len(runeFrequencies))

	for i, r := range keys {
		v := r
		list[i] = &BinaryTree{
			Value: &v,
		}
	}

	return list
}

func buildBinaryTree(list binaryList) *BinaryTree {

	for {
		first := list[0]
		list = list[1:]

		if len(list) == 0 {
			return first
		}

		second := list[0]
		list = list[1:]

		t := &BinaryTree{
			Left:  first,
			Right: second,
		}

		if len(list) == 0 {
			return t
		}

		list = append(list, t)
	}

}

func (s *BinaryTree) isLeaf() bool {
	return s.Value != nil
}

func (s *BinaryTree) prefix() *Prefix {
	prefix := NewPrefix()
	left := s.Left
	if left.isLeaf() {
		vector := NewBitVectorFromBool([]bool{false})
		prefix.Append(*left.Value, vector)

	} else {
		m := left.prefix()
		iterator := m.Enumerate()
		for iterator.HasNext() {
			k, v, _ := iterator.Next()
			vector := NewBitVectorFromVectorPadStart(v, 1)
			vector.Set(0, false)
			prefix.Append(k, vector)
		}
	}

	right := s.Right
	if right.isLeaf() {
		vector := NewBitVectorFromBool([]bool{true})
		prefix.Append(*right.Value, vector)

	} else {
		m := right.prefix()

		iterator := m.Enumerate()
		for iterator.HasNext() {
			k, v, _ := iterator.Next()
			vector := NewBitVectorFromVectorPadStart(v, 1)
			vector.Set(0, true)
			prefix.Append(k, vector)
		}
	}

	return prefix
}

type RuneSlice []rune

func (p RuneSlice) Len() int           { return len(p) }
func (p RuneSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p RuneSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Prefix struct {
	arr     []rune
	prefix  map[rune]*BitVector
	version int
}

func NewPrefix() *Prefix {
	return &Prefix{
		arr:    []rune{},
		prefix: map[rune]*BitVector{},
	}
}

func NewPrefixFromMap(prefixMap map[rune]*BitVector) *Prefix {
	prefix := NewPrefix()

	for k, v := range prefixMap {
		prefix.Append(k, v)
	}

	return prefix
}

func (s *Prefix) Append(r rune, vector *BitVector) {
	s.prefix[r] = vector
	s.arr = append(s.arr, r)
	s.version++
}

func (s *Prefix) Element(index int) (rune, *BitVector) {
	r := s.arr[index]
	return r, s.prefix[r]
}

func (s *Prefix) Get(r rune) *BitVector {
	return s.prefix[r]
}

func (s *Prefix) Has(r rune) bool {
	_, ok := s.prefix[r]
	return ok
}

func (s *Prefix) Length() int {
	return len(s.arr)
}

func (s *Prefix) Enumerate() *PrefixIterator {
	return NewPrefixIteratorWithOffset(s, 0, s.Length())
}

type PrefixIterator struct {
	prefix     *Prefix
	version    int
	indexStart int
	indexEnd   int
}

func NewPrefixIteratorWithOffset(prefix *Prefix, indexStart, indexEnd int) *PrefixIterator {
	if indexStart > prefix.Length() {
		panic("indexStart grater or equal to length")
	}
	if indexStart < 0 {
		panic("indexStart must be non negative number")
	}

	if indexStart > indexEnd {
		panic("indexEnd must be greater then indexStart")
	}

	if prefix.Length()-indexStart < 0 {
		panic("invalid indexStart length")
	}

	if indexEnd > prefix.Length() {
		panic("indexEnd must be greater then prefix length")
	}

	return &PrefixIterator{
		prefix:     prefix,
		indexStart: indexStart,
		indexEnd:   indexEnd,
		version:    prefix.version,
	}
}

func (s *PrefixIterator) Reset() {
	if s.version != s.prefix.version {
		panic("version failed")
	}
	s.indexStart = 0
}

func (s *PrefixIterator) HasNext() bool {
	return s.indexStart < s.indexEnd
}

func (s *PrefixIterator) Next() (rune, *BitVector, int) {
	if s.version != s.prefix.version {
		panic("version failed")
	}

	if s.indexStart < s.indexEnd {
		index := s.indexStart
		currentElement, currentVector := s.prefix.Element(index)
		s.indexStart++
		return currentElement, currentVector, index
	}

	s.indexStart = s.prefix.Length()

	return rune('\x10'), nil, s.indexStart
}

type HuffmanCodeTree struct {
	Left      *HuffmanCodeTree
	Right     *HuffmanCodeTree
	Value     *rune
	Frequency int
}

func NewHuffmanCodeTree(value []rune) *Prefix {
	runeFrequencies, keys := frequencyCount(value)
	huffmanList := rankByRuneCount(runeFrequencies, keys)
	tree := buildHuffmanTree(huffmanList)
	return tree.prefix()
}

func NewHuffmanCodeTreeFromFrequencies(runeFrequencies map[rune]int, keys []rune) *Prefix {
	huffmanList := rankByRuneCount(runeFrequencies, keys)
	tree := buildHuffmanTree(huffmanList)
	return tree.prefix()
}

func (s *HuffmanCodeTree) isLeaf() bool {
	return s.Value != nil
}

func (s *HuffmanCodeTree) prefix() *Prefix {
	prefix := NewPrefix()
	left := s.Left
	if left.isLeaf() {
		vector := NewBitVectorFromBool([]bool{false})
		prefix.Append(*left.Value, vector)
	} else {
		m := left.prefix()
		iterator := m.Enumerate()
		for iterator.HasNext() {
			k, v, _ := iterator.Next()
			vector := NewBitVectorFromVectorPadStart(v, 1)
			vector.Set(0, false)
			prefix.Append(k, vector)
		}
	}

	right := s.Right
	if right.isLeaf() {
		vector := NewBitVectorFromBool([]bool{true})
		prefix.Append(*right.Value, vector)
	} else {
		m := right.prefix()
		iterator := m.Enumerate()
		for iterator.HasNext() {
			k, v, _ := iterator.Next()
			vector := NewBitVectorFromVectorPadStart(v, 1)
			vector.Set(0, true)
			prefix.Append(k, vector)
		}
	}

	return prefix
}

func buildHuffmanTree(list huffmanList) *HuffmanCodeTree {
	for {
		first := list[0]
		list = list[1:]

		if len(list) == 0 {
			return first
		}

		second := list[0]
		list = list[1:]

		sum := first.Frequency + second.Frequency

		t := &HuffmanCodeTree{
			Frequency: sum,
			Left:      first,
			Right:     second,
		}

		if len(list) == 0 {
			return t
		}

		for _, pair := range list {
			if pair.Frequency >= sum {
				list = append([]*HuffmanCodeTree{t}, list...)
			} else {
				list = append(list, t)

			}
			break
		}

	}
}

func frequencyCount(value []rune) (map[rune]int, []rune) {
	runeFrequencies := make(map[rune]int)
	keys := make([]rune, 0)

	for _, r := range value {
		if _, ok := runeFrequencies[r]; ok {
			runeFrequencies[r] = runeFrequencies[r] + 1
		} else {
			runeFrequencies[r] = 1
			keys = append(keys, r)

		}
	}

	return runeFrequencies, keys
}

func rankByRuneCount(runeFrequencies map[rune]int, keys []rune) huffmanList {
	list := make(huffmanList, len(runeFrequencies))

	for i, r := range keys {
		v := r
		list[i] = &HuffmanCodeTree{Value: &v, Frequency: runeFrequencies[r]}
	}
	sort.Sort(list)
	return list
}

type huffmanList []*HuffmanCodeTree

func (p huffmanList) Len() int           { return len(p) }
func (p huffmanList) Less(i, j int) bool { return p[i].Frequency < p[j].Frequency }
func (p huffmanList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func English() *Prefix {
	runeFrequencies := map[rune]int{
		'A': 280937,
		'B': 169474,
		'C': 229363,
		'D': 129632,
		'E': 138443,
		'F': 100751,
		'G': 93212,
		'H': 123632,
		'I': 223312,
		'J': 78706,
		'K': 46580,
		'L': 106984,
		'M': 259474,
		'N': 205409,
		'O': 105700,
		'P': 144239,
		'Q': 11659,
		'R': 146448,
		'S': 304971,
		'T': 325462,
		'U': 57488,
		'V': 31053,
		'W': 107195,
		'X': 7578,
		'Y': 94297,
		'Z': 5610,
		'a': 5263779,
		'b': 866156,
		'c': 1960412,
		'd': 2369820,
		'e': 7741842,
		'f': 1296925,
		'g': 1206747,
		'h': 2955858,
		'i': 4527332,
		'j': 65856,
		'k': 460788,
		'l': 2553152,
		'm': 1467376,
		'n': 4535545,
		'o': 4729266,
		'p': 1255579,
		'q': 54221,
		'r': 4137949,
		's': 4186210,
		't': 5507692,
		'u': 1613323,
		'v': 653370,
		'w': 1015656,
		'x': 123577,
		'y': 1062040,
		'z': 66423,
		' ': 5263778,
		'!': 2178,
		'“': 284671,
		'#': 10,
		'$': 51572,
		'%': 1993,
		'&': 6523,
		'‘': 204497,
		'(': 53398,
		')': 53735,
		'*': 20716,
		'+': 309,
		',': 984969,
		'-': 252302,
		'.': 946136,
		'/': 8161,
		'0': 546233,
		'1': 460946,
		'2': 333499,
		'3': 187606,
		'4': 192528,
		'5': 374413,
		'6': 153865,
		'7': 120094,
		'8': 182627,
		'9': 282364,
		':': 54036,
		';': 36727,
		'<': 82,
		'=': 22,
		'>': 83,
		'?': 12357,
		'@': 1,
		ext: 0,
	}

	keys := RuneSlice{}

	for i := range runeFrequencies {
		keys = append(keys, i)
	}

	sort.Sort(keys)

	return NewHuffmanCodeTreeFromFrequencies(runeFrequencies, keys)
}

// #endregion

// #region bitvector
const bitsPerInt32 = 32

type BitVector struct {
	array   []uint32
	length  int
	version int
}

// Allocates space to hold the length of bit. All of the values in the BitVector are set to false.
func NewBitVector(length int) *BitVector {
	return NewBitVectorOfLength(length, false)
}

// Allocates space to hold the length of bit. All of the values in the BitVector are set to defaultBit.
func NewBitVectorOfLength(length int, defaultBit bool) *BitVector {
	arrayLength, err := getArrayLength(length, bitsPerInt32)
	if err != nil {
		panic(err)
	}
	array := make([]uint32, arrayLength)

	fillValue := uint32(0)
	if defaultBit {
		fillValue = 0xffffffff
	}

	for i := 0; i < arrayLength; i++ {
		array[i] = fillValue
	}

	return &BitVector{
		array:   array,
		length:  length,
		version: 0,
	}
}

// Allocates space to hold the values from the booleans.
func NewBitVectorFromBool(values []bool) *BitVector {
	arrayLength, err := getArrayLength(len(values), bitsPerInt32)
	if err != nil {
		panic(err)
	}
	array := make([]uint32, arrayLength)
	for i, value := range values {
		if value {
			array[i/bitsPerInt32] |= (1 << (i % bitsPerInt32))
		} else {
			array[i/bitsPerInt32] &= ^(1 << (i % bitsPerInt32))
		}
	}

	return &BitVector{
		array:   array,
		length:  len(values),
		version: 0,
	}
}

// Allocates a new BitVector with the same length and bit values as vector.
func NewBitVectorFromVector(vector BitVector) *BitVector {
	array := make([]uint32, len(vector.array))

	copy(array, vector.array)

	return &BitVector{
		array:   array,
		length:  vector.length,
		version: 0,
	}
}

// Allocates a new BitVector padded with the same length and values as the vector but left shifted by the padding.
func NewBitVectorFromVectorPadStart(vector *BitVector, padding int) *BitVector {
	length, err := getArrayLength(vector.Length()+padding, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	array := make([]uint32, length)

	index, err := getArrayLength(padding+1, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	index--

	offset := padding % bitsPerInt32

	arrayLength, err := getArrayLength(vector.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := index; i < length; i++ {
		for y := 0; y < arrayLength; y++ {
			array[i] = (vector.array[y] << offset)
		}
	}

	return &BitVector{
		array:   array,
		length:  vector.length + padding,
		version: 0,
	}
}

// Returns the bit value at position index.
func (s BitVector) Get(index int) bool {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}

	return (s.array[index/bitsPerInt32] & (1 << (index % bitsPerInt32))) != 0
}

// Sets the bit value at position index to value.
func (s BitVector) Set(index int, bit bool) {
	if index < 0 || index >= s.Length() {
		panic(fmt.Sprintf("index %v out of range", index))
	}

	if bit {
		s.array[index/bitsPerInt32] |= (1 << (index % bitsPerInt32))
	} else {
		s.array[index/bitsPerInt32] &= ^(1 << (index % bitsPerInt32))
	}

	s.version++
}

// Sets all the bit values to value.
func (s *BitVector) SetAll(bit bool) {
	fillValue := uint32(0)
	if bit {
		fillValue = 0xffffffff
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] = fillValue
	}

	s.version++
}

func (s *BitVector) Copy(vector *BitVector, indexStart, indexEnd int) {
	if indexStart < 0 {
		panic("indexStart must be non negative number")
	}

	if indexEnd > s.Length() {
		panic("indexEnd must be equal to or less than bitvector")
	}

	if vector.Length() < (s.Length()-indexStart)-indexEnd {
		panic("invalid vector length is to small")
	}

	var err error

	arrayEnd := 0
	if indexEnd > 0 {
		arrayEnd, err = getArrayLength(indexEnd+1, bitsPerInt32)
		if err != nil {
			panic(err)
		}
		arrayEnd--
	}

	arrayStart := 0
	if indexStart > 0 {
		arrayStart, err = getArrayLength(indexStart+1, bitsPerInt32)
		if err != nil {
			panic(err)
		}
		arrayStart--
	}

	index := 0
	offset := indexStart % bitsPerInt32

	for i := arrayStart; i < arrayEnd; i++ {
		vector.array[index] = (s.array[i] >> offset) ^ (s.array[i+1] << (bitsPerInt32 - offset))
		index++
	}

	vector.array[index] = s.array[arrayEnd] >> offset
	vector.version++
}

func (s *BitVector) Length() int {
	return s.length
}

func (s *BitVector) Resize(length int) {
	if length < 0 {
		panic(fmt.Errorf("need non-negative number"))
	}

	arrayLength, err := getArrayLength(length, bitsPerInt32)

	if err != nil {
		panic(err)
	}

	if length != s.length {
		newarray := make([]uint32, arrayLength)
		if len(s.array) != 0 {
			copy(newarray, s.array[:arrayLength])
		}
		s.array = newarray
	}

	if length > s.length {
		last, err := getArrayLength(s.length, bitsPerInt32)
		if err != nil {
			panic(err)
		}
		last--

		bits := s.length % bitsPerInt32
		if bits > 0 {
			s.array[last] &= (1 << bits) - 1
		}
	}

	s.length = length
	s.version++
}

func getArrayLength(n int, div int) (int, error) {
	if div < 0 {
		return 0, fmt.Errorf("div arg must be greater than 0")
	}
	if n > 0 {
		return ((n - 1) / div) + 1, nil
	}
	return 0, nil
}

// Rank counts the number of true or false (depending on what the bit is set to)
// in the bitvector but not including the offset
func (s *BitVector) Rank(bit bool, offset int) int {
	rank := 0

	iterator := s.EnumerateFromOffset(0, offset)
	for iterator.HasNext() {
		v, _ := iterator.Next()
		if v == bit {
			rank++
		}
	}

	return rank
}

// find the offset of true or false (depending on what the bit is set to) from the rank
// (number of times the bit occurs)
func (s *BitVector) Select(bit bool, rank int) int {
	offset := -1
	match := -1
	iterator := s.EnumerateFromOffset(0, s.Length())

	for iterator.HasNext() {
		v, index := iterator.Next()

		if v == bit {
			offset++
			match = index
		}
		if offset == rank {
			break
		}
	}

	return match
}

func (s *BitVector) Concat(vectors []*BitVector) *BitVector {

	length := s.Length()
	for _, v := range vectors {
		length += v.Length()
	}

	vector := NewBitVector(length)

	iterator := s.Enumerate()

	for iterator.HasNext() {
		value, i := iterator.Next()

		vector.Set(i, value)
	}

	index := s.Length()
	for i, v := range vectors {
		index += i

		iterator := v.Enumerate()

		for iterator.HasNext() {
			value, i := iterator.Next()

			index += i
			vector.Set(index, value)
		}

	}
	return vector
}

func (s *BitVector) TrueBits() int {
	output := 0

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		output += bits.OnesCount32(s.array[i])
	}

	return output
}

// ANDed with vector.
func (s *BitVector) And(vector *BitVector) {
	if vector == nil {
		panic(fmt.Errorf("vector is null"))
	}

	if s.Length() != vector.Length() {
		panic(fmt.Errorf("vector length is different"))
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] &= vector.array[i]
	}

	s.version++
}

// ORed with vector.
func (s *BitVector) Or(vector *BitVector) {
	if vector == nil {
		panic(fmt.Errorf("vector is null"))
	}

	if s.Length() != vector.Length() {
		panic(fmt.Errorf("vector length is different"))
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] |= vector.array[i]
	}

	s.version++
}

// XORed with vector.
func (s *BitVector) Xor(vector *BitVector) {
	if vector == nil {
		panic(fmt.Errorf("vector is null"))
	}

	if s.Length() != vector.Length() {
		panic(fmt.Errorf("vector length is different"))
	}

	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] ^= vector.array[i]
	}

	s.version++
}

// Inverts all the bit values. On/true bit values are converted to off/false. Off/false bit values are turned on/true.
func (s *BitVector) Not() {
	arrayLength, err := getArrayLength(s.length, bitsPerInt32)
	if err != nil {
		panic(err)
	}

	for i := 0; i < arrayLength; i++ {
		s.array[i] = ^s.array[i]
	}

	s.version++
}

func (s BitVector) String() string {
	str := []string{}
	iterator := s.Enumerate()

	for iterator.HasNext() {
		value, _ := iterator.Next()

		str = append(str, strconv.FormatBool(value))
	}
	return fmt.Sprintf("{ %s }\n", strings.Join(str, ", "))
}

func (s *BitVector) Enumerate() *BitVectorIterator {
	return NewBitVectorIteratorWithOffset(s, 0, s.Length())
}

func (s *BitVector) EnumerateFromOffset(indexStart, indexEnd int) *BitVectorIterator {
	return NewBitVectorIteratorWithOffset(s, indexStart, indexEnd)
}

type BitVectorIterator struct {
	vector         *BitVector
	version        int
	indexStart     int
	indexEnd       int
	currentElement bool
}

func NewBitVectorIteratorWithOffset(vector *BitVector, indexStart, indexEnd int) *BitVectorIterator {
	if indexStart > vector.Length() {
		panic("indexStart grater or equal to length")
	}
	if indexStart < 0 {
		panic("indexStart must be non negative number")
	}

	if indexStart > indexEnd {
		panic("indexEnd must be greater then indexStart")
	}

	if vector.Length()-indexStart < 0 {
		panic("invalid indexStart length")
	}

	if indexEnd > vector.Length() {
		panic("indexEnd must be greater then vector length")
	}

	return &BitVectorIterator{
		vector:     vector,
		indexStart: indexStart,
		indexEnd:   indexEnd,
		version:    vector.version,
	}
}

func (s *BitVectorIterator) Reset() {
	if s.version != s.vector.version {
		panic("version failed")
	}
	s.indexStart = 0
}

func (s *BitVectorIterator) HasNext() bool {
	return s.indexStart < s.indexEnd
}

func (s *BitVectorIterator) Next() (bool, int) {
	if s.version != s.vector.version {
		panic("version failed")
	}

	if s.indexStart < s.vector.Length() {
		index := s.indexStart
		currentElement := s.vector.Get(s.indexStart)
		s.currentElement = currentElement
		s.indexStart++
		return currentElement, index
	}

	s.indexStart = s.vector.Length()

	return false, s.indexStart
}

// #endregion

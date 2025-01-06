// https://github.com/siongui/go-succinct-data-structure-trie?tab=readme-ov-file
//
// 字典树结构非常适合快速查找字典单词，但如果字典的词汇量很大，构建的字典树可能会占用大量空间。
// !因此，简洁数据结构被应用于字典树结构，使我们能够同时实现快速查找和较小的空间需求。

package main

import (
	"math"
	"strings"
	"unicode/utf8"
)

func main() {
	insertNotInAlphabeticalOrder := func(te *Trie) {
		te.Insert("apple")
		te.Insert("orange")
		te.Insert("alphapha")
		te.Insert("lamp")
		te.Insert("hello")
		te.Insert("jello")
		te.Insert("quiz")
	}

	// optional: set alphabet of words
	// bits.SetAllowedCharacters("abcdeghijklmnoprstuvyāīūṁṃŋṇṅñṭḍḷ…'’° -")
	// Note that you must include space " " in your alphabet if you do not use the
	// default alphabet.
	// default alphabet is [a-z ], i.e.,
	// bits.SetAllowedCharacters("abcdefghijklmnopqrstuvwxyz ")

	// encode: build succinct trie
	te := Trie{}
	te.Init()
	// encode: insert words
	insertNotInAlphabeticalOrder(&te)
	// encode: trie encoding
	teData := te.Encode()
	println(teData)
	println(te.GetNodeCount())
	// encode: build cache for quick lookup
	rd := CreateRankDirectory(teData, te.GetNodeCount()*2+1, L1, L2)
	println(rd.GetData())

	// decode: build frozen succinct trie
	ft := FrozenTrie{}
	ft.Init(teData, rd.GetData(), te.GetNodeCount())

	// decode: look up words
	println(ft.Lookup("apple"))
	println(ft.Lookup("appl"))
	println(ft.Lookup("applee"))

	// decode: words suggestion (find words that start with "prefix")
	// find words starts with "a", max number of returned words is 10
	for _, word := range ft.GetSuggestedWords("a", 10) {
		println(word)
	}
}

// #region trie

// https://blog.golang.org/strings
// https://golang.org/pkg/unicode/utf8/

/*
*

	A Trie node, for use in building the encoding trie. This is not needed for
	the decoder.
*/
type TrieNode struct {
	letter   string
	final    bool
	children []*TrieNode
}

type Trie struct {
	previousWord string
	root         *TrieNode
	cache        []*TrieNode
	nodeCount    uint
}

func (t *Trie) Init() {
	t.previousWord = ""
	t.root = &TrieNode{
		letter: " ",
		final:  false,
	}
	t.cache = append(t.cache, t.root)
	t.nodeCount = 1
}

/*
*

	Returns the number of nodes in the trie
*/
func (t *Trie) GetNodeCount() uint {
	return t.nodeCount
}

/*
*

	Inserts a word into the trie. This function is fastest if the words are
	inserted in alphabetical order.
*/
func (t *Trie) Insert(word string) {

	commonPrefixWidth := 0
	commonRuneCount := 0

	minRuneCount := utf8.RuneCountInString(word)
	if minRuneCount > utf8.RuneCountInString(t.previousWord) {
		minRuneCount = utf8.RuneCountInString(t.previousWord)
	}

	for ; commonRuneCount < minRuneCount; commonRuneCount++ {
		runeValue1, width1 := utf8.DecodeRuneInString(word[commonPrefixWidth:])
		runeValue2, _ := utf8.DecodeRuneInString(t.previousWord[commonPrefixWidth:])
		if runeValue1 != runeValue2 {
			break
		}
		commonPrefixWidth += width1
	}

	t.cache = t.cache[:commonRuneCount+1]
	node := t.cache[commonRuneCount]

	for i, w := commonPrefixWidth, 0; i < len(word); i += w {
		// fix the bug if words not inserted in alphabetical order
		isLetterExist := false
		runeValue, width := utf8.DecodeRuneInString(word[i:])
		w = width
		for _, cld := range node.children {
			if cld.letter == string(runeValue) {
				t.cache = append(t.cache, cld)
				node = cld
				isLetterExist = true
				break
			}
		}
		if isLetterExist {
			continue
		}

		next := &TrieNode{
			letter: string(runeValue),
			final:  false,
		}
		t.nodeCount++
		node.children = append(node.children, next)
		t.cache = append(t.cache, next)
		node = next
	}

	node.final = true
	t.previousWord = word
}

/*
*

	Apply a function to each node, traversing the trie in level order.
*/
func (t *Trie) Apply(fn func(*TrieNode)) {
	var level []*TrieNode
	level = append(level, t.root)
	for len(level) > 0 {
		node := level[0]
		level = level[1:]
		for i := 0; i < len(node.children); i++ {
			level = append(level, node.children[i])
		}
		fn(node)
	}
}

/*
*

	Encode the trie and all of its nodes. Returns a string representing the
	encoded data.
*/
func (t *Trie) Encode() string {
	// Write the unary encoding of the tree in level order.
	bits := BitWriter{}
	bits.Write(0x02, 2)
	t.Apply(func(node *TrieNode) {
		for i := 0; i < len(node.children); i++ {
			bits.Write(1, 1)
		}
		bits.Write(0, 1)
	})

	// Write the data for each node, using (dataBits) bits for one node.
	// 1 bit stores the "final" indicator. The other (dataBits-1) bits store
	// one of the characters of the alphabet.
	t.Apply(func(node *TrieNode) {
		value, ok := mapCharToUint[node.letter]
		if !ok {
			panic("illegal character:" + node.letter)
		}
		if node.final {
			value |= (1 << (dataBits - 1))
		}

		bits.Write(uint(value), dataBits)
	})

	return bits.GetData()
}

// #endregion

// #region alphabet

// var allowedCharacters = "abcdeghijklmnoprstuvyāīūṁṃŋṇṅñṭḍḷ…'’° -"
var allowedCharacters = "abcdefghijklmnopqrstuvwxyz "
var mapCharToUint = getCharToUintMap(allowedCharacters)
var mapUintToChar = getUintToCharMap(mapCharToUint)

/**
 * Write the data for each node, call getDataBits() to calculate how many bits
 * for one node.
 * 1 bit stores the "final" indicator. The other bits store one of the
 * characters of the alphabet.
 */
var dataBits = getDataBits(allowedCharacters)

func SetAllowedCharacters(alphabet string) {
	allowedCharacters = alphabet
	mapCharToUint = getCharToUintMap(alphabet)
	mapUintToChar = getUintToCharMap(mapCharToUint)
	dataBits = getDataBits(alphabet)
}

func getCharToUintMap(alphabet string) map[string]uint {
	result := map[string]uint{}

	var i uint = 0
	chars := strings.Split(alphabet, "")
	for _, char := range chars {
		result[char] = i
		i++
	}

	return result
}

func getUintToCharMap(c2ui map[string]uint) map[uint]string {
	result := map[uint]string{}
	for k, v := range c2ui {
		result[v] = k
	}
	return result
}

func getDataBits(alphabet string) uint {
	numOfChars := len(strings.Split(alphabet, ""))
	var i uint = 0

	for (1 << i) < numOfChars {
		i++
	}

	// one more bit for the "final" indicator
	return (i + 1)
}

// #endregion

// #region base64
// Configure the bit writing and reading functions to work natively in BASE-64
// encoding. That way, we don't have to convert back and forth to bytes.

var BASE64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

/*
*

	The width of each unit of the encoding, in bits. Here we use 6, for base-64
	encoding.
*/
var W uint = 6

/*
*

	Returns the character unit that represents the given value. If this were
	binary data, we would simply return id.
*/
func CHR(id uint) string {
	return BASE64[id : id+1]
}

/*
*

	Returns the decimal value of the given character unit.
*/
var BASE64_CACHE = map[string]uint{
	"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7,
	"I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13, "O": 14,
	"P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20, "V": 21,
	"W": 22, "X": 23, "Y": 24, "Z": 25, "a": 26, "b": 27, "c": 28,
	"d": 29, "e": 30, "f": 31, "g": 32, "h": 33, "i": 34, "j": 35,
	"k": 36, "l": 37, "m": 38, "n": 39, "o": 40, "p": 41, "q": 42,
	"r": 43, "s": 44, "t": 45, "u": 46, "v": 47, "w": 48, "x": 49,
	"y": 50, "z": 51, "0": 52, "1": 53, "2": 54, "3": 55, "4": 56,
	"5": 57, "6": 58, "7": 59, "8": 60, "9": 61, "-": 62, "_": 63,
}

func ORD(ch string) uint {
	// Used to be: return BASE64.indexOf(ch);
	return BASE64_CACHE[ch]
}

// #endregion

// #region bitstring
/**
  Given a string of data (eg, in BASE-64), the BitString class supports
  reading or counting a number of bits from an arbitrary position in the
  string.
*/
type BitString struct {
	base64DataString string
	length           uint
}

var MaskTop = [7]uint{
	0x3f, 0x1f, 0x0f, 0x07, 0x03, 0x01, 0x00,
}

var BitsInByte = [256]uint{
	0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2,
	3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3,
	3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3,
	4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4,
	3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5,
	6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4,
	4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5,
	6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3,
	4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6,
	6, 7, 6, 7, 7, 8,
}

func (bs *BitString) Init(data string) {
	bs.base64DataString = data
	bs.length = uint(len(bs.base64DataString)) * W
}

/*
*

	Returns the internal string of bytes
*/
func (bs *BitString) GetData() string {
	return bs.base64DataString
}

/*
*

	Returns a decimal number, consisting of a certain number, n, of bits
	starting at a certain position, p.
*/
func (bs *BitString) Get(p, n uint) uint {

	// case 1: bits lie within the given byte
	if (p%W)+n <= W {
		idx := p/W | 0
		return (ORD(bs.base64DataString[idx:idx+1]) & MaskTop[p%W]) >>
			(W - p%W - n)

		// case 2: bits lie incompletely in the given byte
	} else {
		idx := p/W | 0
		result := (ORD(bs.base64DataString[idx:idx+1]) & MaskTop[p%W])

		l := W - p%W
		p += l
		n -= l

		for n >= W {
			idx := p/W | 0
			result = (result << W) | ORD(bs.base64DataString[idx:idx+1])
			p += W
			n -= W
		}

		if n > 0 {
			idx := p/W | 0
			result = (result << n) | (ORD(bs.base64DataString[idx:idx+1]) >>
				(W - n))
		}

		return result
	}
}

/*
*

	Counts the number of bits set to 1 starting at position p and
	ending at position p + n
*/
func (bs *BitString) Count(p, n uint) uint {

	var count uint = 0
	for n >= 8 {
		count += BitsInByte[bs.Get(p, 8)]
		p += 8
		n -= 8
	}

	return count + BitsInByte[bs.Get(p, n)]
}

/*
*

	Returns the number of bits set to 1 up to and including position x.
	This is the slow implementation used for testing.
*/
func (bs *BitString) Rank(x uint) uint {
	var rank uint = 0
	var i uint = 0
	for i = 0; i <= x; i++ {
		// FIXME: the above line should be the following???
		//for i = 0; i < x; i++ {
		if bs.Get(i, 1) != 0 {
			rank++
		}
	}

	return rank
}

// #endregion

// #region bitwriter

/*
*

	The BitWriter will create a stream of bytes, letting you write a certain
	number of bits at a time. This is part of the encoder, so it is not
	optimized for memory or speed.
*/
type BitWriter struct {
	bits []uint
}

/*
*

	Write some data to the bit string. The number of bits must be 32 or
	fewer.
*/
func (bw *BitWriter) Write(data, numBits uint) {
	//for i := (numBits-1); i >= 0; i-- {
	// @siongui: the above commented line will cause infinite loop, why???
	// answer from @xphoenix:
	// Because i becomes uint, let's check iteration when i == 0, at the end
	// of loop, i-- takes place but as i is uint, it leads to 2^32-1 instead
	// of -1, loop condition is still true...
	for i := numBits; i > 0; i-- {
		j := i - 1
		if (data & (1 << j)) != 0 {
			bw.bits = append(bw.bits, 1)
		} else {
			bw.bits = append(bw.bits, 0)
		}
	}
}

/*
*

	Get the bitstring represented as a javascript string of bytes
*/
func (bw *BitWriter) GetData() string {
	var chars []string
	var b, i uint = 0, 0

	for j := 0; j < len(bw.bits); j++ {
		b = (b << 1) | bw.bits[j]
		i += 1
		if i == W {
			chars = append(chars, CHR(b))
			i = 0
			b = 0
		}
	}

	if i != 0 {
		chars = append(chars, CHR(b<<(W-i)))
	}

	return strings.Join(chars, "")
}

/*
*

	Returns the bits as a human readable binary string for debugging
*/
func (bw *BitWriter) GetDebugString(group uint) string {
	var chars []string
	var i uint = 0

	for j := 0; j < len(bw.bits); j++ {
		if bw.bits[j] == 1 {
			chars = append(chars, "1")
		} else {
			chars = append(chars, "0")
		}
		i++
		if i == group {
			chars = append(chars, " ")
			i = 0
		}
	}

	return strings.Join(chars, "")
}

// #endregion

// #region frozentrie
/**
  This class is used for traversing the succinctly encoded trie.
*/
type FrozenTrieNode struct {
	trie       *FrozenTrie
	index      uint
	letter     string
	final      bool
	firstChild uint
	childCount uint
}

/*
*

	Returns the number of children.
*/
func (f *FrozenTrieNode) GetChildCount() uint {
	return f.childCount
}

/*
*

	Returns the FrozenTrieNode for the given child.

	@param index The 0-based index of the child of this node. For example, if
	the node has 5 children, and you wanted the 0th one, pass in 0.
*/
func (f *FrozenTrieNode) GetChild(index uint) FrozenTrieNode {
	return f.trie.GetNodeByIndex(f.firstChild + index)
}

/*
*

	The FrozenTrie is used for looking up words in the encoded trie.

	@param data A string representing the encoded trie.

	@param directoryData A string representing the RankDirectory. The global L1
	and L2 constants are used to determine the L1Size and L2size.

	@param nodeCount The number of nodes in the trie.
*/
type FrozenTrie struct {
	data        BitString
	directory   RankDirectory
	letterStart uint
}

func (f *FrozenTrie) Init(data, directoryData string, nodeCount uint) {
	f.data.Init(data)
	f.directory.Init(directoryData, data, nodeCount*2+1, L1, L2)

	// The position of the first bit of the data in 0th node. In non-root
	// nodes, this would contain 6-bit letters.
	f.letterStart = nodeCount*2 + 1
}

/*
*

	Retrieve the FrozenTrieNode of the trie, given its index in level-order.
	This is a private function that you don't have to use.
*/
func (f *FrozenTrie) GetNodeByIndex(index uint) FrozenTrieNode {
	// retrieve the (dataBits)-bit letter.
	final := (f.data.Get(f.letterStart+index*dataBits, 1) == 1)
	letter, ok := mapUintToChar[f.data.Get(f.letterStart+index*dataBits+1, (dataBits-1))]
	if !ok {
		panic("illegal: bits -> char failed")
	}
	firstChild := f.directory.Select(0, index+1) - index

	// Since the nodes are in level order, this nodes children must go up
	// until the next node's children start.
	childOfNextNode := f.directory.Select(0, index+2) - index - 1

	return FrozenTrieNode{
		trie:       f,
		index:      index,
		letter:     letter,
		final:      final,
		firstChild: firstChild,
		childCount: (childOfNextNode - firstChild),
	}
}

/*
*

	Retrieve the root node. You can use this node to obtain all of the other
	nodes in the trie.
*/
func (f *FrozenTrie) GetRoot() FrozenTrieNode {
	return f.GetNodeByIndex(0)
}

/*
*

	Look-up a word in the trie. Returns true if and only if the word exists
	in the trie.
*/
func (f *FrozenTrie) Lookup(word string) bool {
	node := f.GetRoot()
	for i, w := 0, 0; i < len(word); i += w {
		runeValue, width := utf8.DecodeRuneInString(word[i:])
		w = width
		var child FrozenTrieNode
		var j uint = 0
		for ; j < node.GetChildCount(); j++ {
			child = node.GetChild(j)
			if child.letter == string(runeValue) {
				break
			}
		}

		if j == node.GetChildCount() {
			return false
		}
		node = child
	}

	return node.final
}

// #endregion

// #region rankdirectory

/*
*

	Fixed values for the L1 and L2 table sizes in the Rank Directory
*/
var L1 uint = 32 * 32
var L2 uint = 32

/*
*

	The rank directory allows you to build an index to quickly compute the
	rank() and select() functions. The index can itself be encoded as a binary
	string.
*/
type RankDirectory struct {
	directory   BitString
	data        BitString // data of succinct trie
	l1Size      uint
	l2Size      uint
	l1Bits      uint
	l2Bits      uint
	sectionBits uint
	numBits     uint
}

/*
*

	Used to build a rank directory from the given input string.

	@param data A javascript string containing the data, as readable using the
	BitString object.

	@param numBits The number of bits to index.

	@param l1Size The number of bits that each entry in the Level 1 table
	summarizes. This should be a multiple of l2Size.

	@param l2Size The number of bits that each entry in the Level 2 table
	summarizes.
*/
func CreateRankDirectory(data string, numBits, l1Size, l2Size uint) RankDirectory {
	bits := BitString{}
	bits.Init(data)
	var p, i uint = 0, 0
	var count1, count2 uint = 0, 0
	l1bits := uint(math.Ceil(math.Log2(float64(numBits))))
	l2bits := uint(math.Ceil(math.Log2(float64(l1Size))))

	directory := BitWriter{}

	for p+l2Size <= numBits {
		count2 += bits.Count(p, l2Size)
		i += l2Size
		p += l2Size
		if i == l1Size {
			count1 += count2
			directory.Write(count1, l1bits)
			count2 = 0
			i = 0
		} else {
			directory.Write(count2, l2bits)
		}
	}

	rd := RankDirectory{}
	rd.Init(directory.GetData(), data, numBits, l1Size, l2Size)
	return rd
}

func (rd *RankDirectory) Init(directoryData, bitData string, numBits, l1Size, l2Size uint) {
	rd.directory.Init(directoryData)
	rd.data.Init(bitData)
	rd.l1Size = l1Size
	rd.l2Size = l2Size
	rd.l1Bits = uint(math.Ceil(math.Log2(float64(numBits))))
	rd.l2Bits = uint(math.Ceil(math.Log2(float64(l1Size))))
	rd.sectionBits = (l1Size/l2Size-1)*rd.l2Bits + rd.l1Bits
	rd.numBits = numBits
}

/*
*

	Returns the string representation of the directory.
*/
func (rd *RankDirectory) GetData() string {
	return rd.directory.GetData()
}

/*
*

	Returns the number of 1 or 0 bits (depending on the "which" parameter) to
	to and including position x.
*/
func (rd *RankDirectory) Rank(which, x uint) uint {

	if which == 0 {
		return x - rd.Rank(1, x) + 1
	}

	var rank uint = 0
	o := x
	var sectionPos uint = 0

	if o >= rd.l1Size {
		sectionPos = (o/rd.l1Size | 0) * rd.sectionBits
		rank = rd.directory.Get(sectionPos-rd.l1Bits, rd.l1Bits)
		o = o % rd.l1Size
	}

	if o >= rd.l2Size {
		sectionPos += (o/rd.l2Size | 0) * rd.l2Bits
		rank += rd.directory.Get(sectionPos-rd.l2Bits, rd.l2Bits)
	}

	rank += rd.data.Count(x-x%rd.l2Size, x%rd.l2Size+1)

	return rank
}

/*
*

	Returns the position of the y'th 0 or 1 bit, depending on the "which"
	parameter.
*/
func (rd *RankDirectory) Select(which, y uint) uint {
	high := int(rd.numBits)
	low := -1
	val := -1

	for high-low > 1 {
		probe := (high+low)/2 | 0
		r := rd.Rank(which, uint(probe))

		if r == y {
			// We have to continue searching after we have found it,
			// because we want the _first_ occurrence.
			val = probe
			high = probe
		} else if r < y {
			low = probe
		} else {
			high = probe
		}
	}

	return uint(val)
}

// #endregion

// #region search
/**
 * Given a word, returns array of words, prefix of which is word
 */
func (f *FrozenTrie) GetSuggestedWords(word string, limit int) []string {
	var result []string

	node := f.GetRoot()

	// find the node corresponding to the last char of input
	for _, runeValue := range word {
		var child FrozenTrieNode
		var j uint = 0
		for ; j < node.GetChildCount(); j++ {
			child = node.GetChild(j)
			if child.letter == string(runeValue) {
				break
			}
		}

		// not found, return.
		if j == node.GetChildCount() {
			return result
		}

		node = child
	}

	// The node corresponding to the last letter of word is found.
	// Use this node as root. traversing the trie in level order.
	return f.traverseSubTrie(node, word, limit)
}

func (f *FrozenTrie) traverseSubTrie(node FrozenTrieNode, prefix string, limit int) []string {
	var result []string

	var level []FrozenTrieNode
	level = append(level, node)
	var prefixLevel []string
	prefixLevel = append(prefixLevel, prefix)

	for len(level) > 0 {
		nodeNow := level[0]
		level = level[1:]
		prefixNow := prefixLevel[0]
		prefixLevel = prefixLevel[1:]

		// if the prefix is a legal word.
		if nodeNow.final {
			result = append(result, prefixNow)
			if len(result) > limit {
				return result
			}
		}

		var i uint = 0
		for ; i < nodeNow.GetChildCount(); i++ {
			child := nodeNow.GetChild(i)
			level = append(level, child)
			prefixLevel = append(prefixLevel, prefixNow+child.letter)
		}
	}

	return result
}

// #endregion

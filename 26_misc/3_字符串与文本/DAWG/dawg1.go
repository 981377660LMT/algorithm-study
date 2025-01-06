// https://github.com/smhanov/dawg

/*
Package dawg is an implemention of a Directed Acyclic Word Graph, as described
on my blog at http://stevehanov.ca/blog/?id=115

A DAWG provides fast lookup of all possible prefixes of words in a dictionary, as well
as the ability to get the index number of any word.

This particular implementation may be different from others because it is very memory
efficient, and it also works fast with large character sets. It can deal with
thousands of branches out of a single node without needing to go through each one.

The storage format is as small as possible. Bits are used instead of bytes so that
no space is wasted as padding, and there are no practical limitations to the number of
nodes or characters. A summary of the data format is found at the top of disk.go.

In general, to use it you first create a builder using dawg.New(). You can then
add words to the Dawg. The two restrictions are that you cannot repeat a word, and
they must be in strictly increasing alphabetical order.

After all the words are added, call Finish() which returns a dawg.Finder interface.
You can perform queries with this interface, such as finding all prefixes of a given string
which are also words, or looking up a word's index that you have previously added.

After you have called Finish() on a Builder, you may choose to write it to disk using the
Save() function. The DAWG can then be opened again later using the Load() function.
When opened from disk, no memory is used. The structure is accessed in-place on disk.
*/
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
	"strconv"
)

func main() {

}

// #region dawg
// FindResult is the result of a lookup in the d. It
// contains both the word found, and it's index based on the
// order it was added.
type FindResult struct {
	Word  string
	Index int
}

type edgeStart struct {
	node int
	ch   rune
}

func (edge edgeStart) String() string {
	return fmt.Sprintf("(%d, '%c')", edge.node, edge.ch)
}

type edgeEnd struct {
	node  int
	count int
}

type uncheckedNode struct {
	parent int
	ch     rune
	child  int
}

// EnumFn is a method that you implement. It will be called with
// all prefixes stored in the DAWG. If final is true, the prefix
// represents a complete word that has been stored.
type EnumFn = func(index int, word []rune, final bool) EnumerationResult

// EnumerationResult is returned by the enumeration function to indicate whether
// indication should continue below this depth or to stop altogether
type EnumerationResult = int

const (
	// Continue enumerating all words with this prefix
	Continue EnumerationResult = iota

	// Skip will skip all words with this prefix
	Skip

	// Stop will immediately stop enumerating words
	Stop
)

// Finder is the interface for querying a dawg. Use either
// Builder.Finish() or Load() to obtain one.
type Finder interface {
	// Find all prefixes of the given string
	FindAllPrefixesOf(input string) []FindResult

	// Find the index of the given string
	IndexOf(input string) int

	AtIndex(index int) (string, error)

	// Enumerate all prefixes stored in the dawg.
	Enumerate(fn EnumFn)

	// Returns the number of words
	NumAdded() int

	// Returns the number of edges
	NumEdges() int

	// Returns the number of nodes
	NumNodes() int

	// Output a human-readable description of the dawg to stdout
	Print()

	// Close the dawg that was opened with Load(). After this, it is no longer
	// accessible.
	Close() error

	// Save to a writer
	Write(w io.Writer) (int64, error)

	// Save to a file
	Save(filename string) (int64, error)
}

// Builder is the interface for creating a new Dawg. Use New() to create it.
type Builder interface {
	// Add the word to the dawg
	Add(wordIn string)

	// Returns true if the word can be added.
	CanAdd(word string) bool

	// Complete the dawg and return a Finder.
	Finish() Finder
}

const rootNode = 0

type node struct {
	final bool
	count int
	edges []edgeStart
}

// dawg represents a Directed Acyclic Word Graph
type dawg struct {
	// these are erased after we finish building
	lastWord       []rune
	nextID         int
	uncheckedNodes []uncheckedNode
	minimizedNodes map[string]int
	nodes          map[int]*node

	// if read from a file, this is set
	r    io.ReaderAt
	size int64 // size of the readerAt

	// these are kept
	finished        bool
	numAdded        int
	numNodes        int
	numEdges        int
	cbits           int64 // bits to represent character value
	abits           int64 // bits to represent node address
	wbits           int64 // bits to represent number of words / counts
	firstNodeOffset int64 // first node offset in bits in the file
	hasEmptyWord    bool
}

// New creates a new dawg
func New() Builder {
	return &dawg{
		nextID:         1,
		minimizedNodes: make(map[string]int),
		nodes: map[int]*node{
			0: {count: -1},
		},
	}
}

// CanAdd will return true if the word can be added to the d.
// Words must be added in alphabetical order.
func (d *dawg) CanAdd(word string) bool {
	return !d.finished &&
		(d.numAdded == 0 || word > string(d.lastWord))
}

// Add adds a word to the structure.
// Adding a word not in alphaetical order, or to a finished dawg will panic.
func (d *dawg) Add(wordIn string) {
	if d.numAdded > 0 && wordIn <= string(d.lastWord) {
		log.Printf("Last word=%s newword=%s", string(d.lastWord), wordIn)
		panic(errors.New("d.AddWord(): Words not in alphabetical order"))
	} else if d.finished {
		panic(errors.New("d.AddWord(): Tried to add to a finished dawg"))
	}

	word := []rune(wordIn)

	// find common prefix between word and previous word
	commonPrefix := 0
	for i := 0; i < min(len(word), len(d.lastWord)); i++ {
		if word[i] != d.lastWord[i] {
			break
		}
		commonPrefix++
	}

	// Check the uncheckedNodes for redundant nodes, proceeding from last
	// one down to the common prefix size. Then truncate the list at that
	// point.
	d.minimize(commonPrefix)

	// add the suffix, starting from the correct node mid-way through the
	// graph
	var node int
	if len(d.uncheckedNodes) == 0 {
		node = rootNode
	} else {
		node = d.uncheckedNodes[len(d.uncheckedNodes)-1].child
	}

	for _, letter := range word[commonPrefix:] {
		nextNode := d.newNode()
		d.addChild(node, letter, nextNode)
		d.uncheckedNodes = append(d.uncheckedNodes, uncheckedNode{node, letter, nextNode})
		node = nextNode
	}

	d.setFinal(node)
	d.lastWord = word
	d.numAdded++
}

// Finish will mark the dawg as complete. The dawg cannot be used for lookups
// until Finish has been called.
func (d *dawg) Finish() Finder {
	if !d.finished {
		d.finished = true

		d.minimize(0)

		d.numNodes = len(d.minimizedNodes) + 1

		// Fill in the counts
		d.calculateSkipped(rootNode)

		// no longer need the names.
		d.uncheckedNodes = nil
		d.minimizedNodes = nil
		d.lastWord = nil

		d.renumber()

		var buffer bytes.Buffer
		d.size, _ = d.Write(&buffer)
		d.r = bytes.NewReader(buffer.Bytes())
		d.nodes = nil
	}

	finder, _ := Read(d.r, 0)

	return finder
}

func (d *dawg) renumber() {
	// after minimization, nodes have been removed so there are gaps in the node IDs.
	// Renumber them all to be consecutive.
	// process them in a depth-first order so that runs of characters
	// will appear in consecutive nodes, which is more efficient for encoding.

	remap := make(map[int]int)

	var process func(id int)

	process = func(id int) {
		if _, ok := remap[id]; ok {
			return
		}

		remap[id] = len(remap)
		node := d.nodes[id]
		for _, edge := range node.edges {
			process(edge.node)
		}
	}

	process(rootNode)

	nodes := make(map[int]*node)
	for id, node := range d.nodes {
		nodes[remap[id]] = node
		for i := range node.edges {
			node.edges[i].node = remap[node.edges[i].node]
		}
	}
	d.nodes = nodes
}

// Print will print all edges to the standard output
func (d *dawg) Print() {
	DumpFile(d.r)
}

// FindAllPrefixesOf returns all items in the dawg that are a prefix of the input string.
// It will panic if the dawg is not finished.
func (d *dawg) FindAllPrefixesOf(input string) []FindResult {

	d.checkFinished()

	var results []FindResult
	skipped := 0
	final := d.hasEmptyWord
	node := rootNode
	var edgeEnd edgeEnd
	var ok bool

	r := newBitSeeker(d.r)

	// for each character of the input
	for pos, letter := range input {
		// if the node is final, add a result
		if final {
			results = append(results, FindResult{
				Word:  input[:pos],
				Index: skipped,
			})
		}

		// check if there is an outgoing edge for the letter
		edgeEnd, final, ok = d.getEdge(&r, edgeStart{node: node, ch: letter})
		if !ok {
			return results
		}

		// we found an edge.
		node = edgeEnd.node
		skipped += edgeEnd.count
	}

	if final {
		results = append(results, FindResult{
			Word:  input,
			Index: skipped,
		})
	}

	return results
}

// IndexOf returns the index, which is the order the item was inserted.
// If the item was never inserted, it returns -1
// It will panic if the dawg is not finished.
func (d *dawg) IndexOf(input string) int {
	skipped := 0
	node := rootNode
	final := d.hasEmptyWord
	var ok bool
	var edgeEnd edgeEnd
	r := newBitSeeker(d.r)

	// for each character of the input
	for _, letter := range input {
		// check if there is an outgoing edge for the letter
		edgeEnd, final, ok = d.getEdge(&r, edgeStart{node: node, ch: letter})
		//log.Printf("Follow %v:%v=>%v (ok=%v)", node, string(letter), edgeEnd.node, ok)
		if !ok {
			// not found
			return -1
		}

		// we found an edge.
		node = edgeEnd.node
		skipped += edgeEnd.count
	}

	//log.Printf("IsFinal %d: %v", node, final)
	if final {
		return skipped
	}
	return -1
}

// NumAdded returns the number of words added
func (d *dawg) NumAdded() int {
	return d.numAdded
}

// NumNodes returns the number of nodes in the d.
func (d *dawg) NumNodes() int {
	return d.numNodes
}

// NumEdges returns the number of edges in the d. This includes transitions to
// the "final" node after each word.
func (d *dawg) NumEdges() int {
	return d.numEdges
}

func (d *dawg) checkFinished() {
	if !d.finished {
		panic(errors.New("dawg was not Finished()"))
	}
}

func (d *dawg) minimize(downTo int) {
	// proceed from the leaf up to a certain point
	for i := len(d.uncheckedNodes) - 1; i >= downTo; i-- {
		u := d.uncheckedNodes[i]
		name := d.nameOf(u.child)
		if node, ok := d.minimizedNodes[name]; ok {
			// replace the child with the previously encountered one
			d.replaceChild(u.parent, u.ch, node)
		} else {
			// add the state to the minimized nodes.
			d.minimizedNodes[name] = u.child
		}
	}

	d.uncheckedNodes = d.uncheckedNodes[:downTo]
}

func (d *dawg) newNode() int {
	d.nextID++
	return d.nextID - 1
}

func (d *dawg) nameOf(nodeid int) string {
	node := d.nodes[nodeid]

	// node name is id_ch:id... for each child
	buff := bytes.Buffer{}
	for _, edge := range node.edges {
		buff.WriteByte('_')
		buff.WriteRune(edge.ch)
		buff.WriteByte(':')
		buff.WriteString(strconv.Itoa(edge.node))
	}

	if node.final {
		buff.WriteByte('!')
	}

	return buff.String()
}

func (d *dawg) setFinal(node int) {
	d.nodes[node].final = true
	if node == rootNode {
		d.hasEmptyWord = true
	}
}

func (d *dawg) addChild(parent int, ch rune, child int) {
	//log.Printf("Addchild %v(%v)->%v", parent, string(ch), child)
	d.numEdges++
	if d.nodes[child] == nil {
		d.nodes[child] = &node{
			count: -1,
		}
	}
	node := d.nodes[parent]
	if len(node.edges) > 0 && ch <= node.edges[len(node.edges)-1].ch {
		log.Panic("Not strictly increasing")
	}
	node.edges = append(node.edges, edgeStart{child, ch})
}

func (d *dawg) replaceChild(parent int, ch rune, child int) {
	pnode := d.nodes[parent]
	//TODO: should be bsearch
	i := bsearch(len(pnode.edges), func(i int) int {
		return int(pnode.edges[i].ch - ch)
	})

	if pnode.edges[i].ch != ch {
		//for _, edge := range pnode.edges {
		//	log.Printf("Edge %c %d", rune(edge.ch), edge.node)
		//}
		log.Panicf("Not found: %c", ch)
	}

	//log.Printf("ReplaceChild(%v:%v=>%v, %v:%v=>%v)",
	//	parent, string(ch), pnode.edges[i].node,
	//	parent, string(ch), child)

	delete(d.nodes, pnode.edges[i].node)
	pnode.edges[i].node = child

}

func (d *dawg) calculateSkipped(nodeid int) int {
	// for each child of the node, calculate now many nodes
	// are skipped over by following that child. This is the
	// sum of all skipped-over counts of its previous siblings.

	// returns the number of leaves reachable from the node.
	node := d.nodes[nodeid]
	if node.count >= 0 {
		return node.count
	}

	numReachable := 0

	if node.final {
		numReachable++
	}

	for _, edge := range node.edges {
		numReachable += d.calculateSkipped(edge.node)
	}

	node.count = numReachable

	return numReachable
}

// Enumerate will call the given method, passing it every possible prefix of words in the index.
// Return Continue to continue enumeration, Skip to skip this branch, or Stop to stop enumeration.
func (d *dawg) Enumerate(fn EnumFn) {
	r := newBitSeeker(d.r)
	d.enumerate(&r, 0, rootNode, nil, fn)
}

func (d *dawg) enumerate(r *bitSeeker, index int, address int, runes []rune, fn EnumFn) EnumerationResult {
	// get the node and whether its final
	node := d.getNode(r, address)

	// call the enum function on the runes
	result := fn(index, runes, node.final)

	// if the function didn't say to continue, then return.
	if result != Continue {
		return result
	}

	l := len(runes)
	runes = append(runes, 0)

	// for each edge
	for _, edge := range node.edges {
		// add ch to the runes
		runes[l] = edge.ch
		// recurse
		result = d.enumerate(r, index+edge.count, edge.node, runes, fn)
		if result == Stop {
			break
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *dawg) AtIndex(index int) (string, error) {
	if index < 0 || index >= d.NumAdded() {
		return "", errors.New("invalid index")
	}

	r := newBitSeeker(d.r)
	// start at first node and empty string
	result, _ := d.atIndex(&r, rootNode, 0, index, nil)
	return result, nil
}

func (d *dawg) atIndex(r *bitSeeker, nodeNumber, atIndex, targetIndex int, runes []rune) (string, bool) {
	node := d.getNode(r, nodeNumber)
	// if node is final and index matches, return it
	if node.final && atIndex == targetIndex {
		return string(runes), true
	}

	next := bsearch(len(node.edges), func(i int) int {
		return atIndex + node.edges[i].count - targetIndex
	})

	if next == len(node.edges) || atIndex+node.edges[next].count > targetIndex {
		next--
	}

	//log.Printf("Follow edge %v %c skip=%d", node.edges[next], node.edges[next].ch, node.edges[next].count)
	runes = append(runes, 0)
	for i := next; i < len(node.edges); i++ {
		runes[len(runes)-1] = node.edges[i].ch
		if result, ok := d.atIndex(r, node.edges[i].node, atIndex+node.edges[i].count, targetIndex, runes); ok {
			return result, ok
		}
	}
	return "", false

}

// #endregion

// #region bits
type bitWriter struct {
	io.Writer
	cache uint8
	used  int
}

// NewBitWriter creates a new BitWriter from an io writer.
func newBitWriter(w io.Writer) *bitWriter {
	return &bitWriter{w, 0, 0}
}

func (w *bitWriter) WriteBits(data uint64, n int) error {
	var mask uint8
	for n > 0 {
		written := n
		if written+w.used > 8 {
			written = 8 - w.used
		}

		mask = uint8(uint16(1<<(written)) - 1)
		w.used += written
		w.cache = (w.cache << written) | byte(data>>(n-written))&mask

		if w.used == 8 {
			_, err := w.Write([]byte{w.cache})
			if err != nil {
				return err
			}
			w.used = 0
		}

		n -= written
	}
	return nil
}

func (w *bitWriter) Flush() error {
	if w.used > 0 {
		_, err := w.Write([]byte{w.cache << (8 - w.used)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *bitWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}

	if closer, ok := w.Writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

var maskTop = [64]uint64{
	0xffffffffffffffff,
	0x7fffffffffffffff,
	0x3fffffffffffffff,
	0x1fffffffffffffff,
	0x0fffffffffffffff,
	0x07ffffffffffffff,
	0x03ffffffffffffff,
	0x01ffffffffffffff,
	0x00ffffffffffffff,
	0x007fffffffffffff,
	0x003fffffffffffff,
	0x001fffffffffffff,
	0x000fffffffffffff,
	0x0007ffffffffffff,
	0x0003ffffffffffff,
	0x0001ffffffffffff,
	0x0000ffffffffffff,
	0x00007fffffffffff,
	0x00003fffffffffff,
	0x00001fffffffffff,
	0x00000fffffffffff,
	0x000007ffffffffff,
	0x000003ffffffffff,
	0x000001ffffffffff,
	0x000000ffffffffff,
	0x0000007fffffffff,
	0x0000003fffffffff,
	0x0000001fffffffff,
	0x0000000fffffffff,
	0x00000007ffffffff,
	0x00000003ffffffff,
	0x00000001ffffffff,
	0x00000000ffffffff,
	0x000000007fffffff,
	0x000000003fffffff,
	0x000000001fffffff,
	0x000000000fffffff,
	0x0000000007ffffff,
	0x0000000003ffffff,
	0x0000000001ffffff,
	0x0000000000ffffff,
	0x00000000007fffff,
	0x00000000003fffff,
	0x00000000001fffff,
	0x00000000000fffff,
	0x000000000007ffff,
	0x000000000003ffff,
	0x000000000001ffff,
	0x000000000000ffff,
	0x0000000000007fff,
	0x0000000000003fff,
	0x0000000000001fff,
	0x0000000000000fff,
	0x00000000000007ff,
	0x00000000000003ff,
	0x00000000000001ff,
	0x00000000000000ff,
	0x000000000000007f,
	0x000000000000003f,
	0x000000000000001f,
	0x000000000000000f,
	0x0000000000000007,
	0x0000000000000003,
	0x0000000000000001,
}

// BitSeeker reads bits from a given offset in bits
type bitSeeker struct {
	io.ReaderAt
	p      int64
	have   int64
	buffer [8]byte
	slice  []byte
	cache  uint64
}

// NewBitSeeker creates a new bitreaderat
func newBitSeeker(r io.ReaderAt) bitSeeker {
	bs := bitSeeker{ReaderAt: r, have: -1}
	// avoids re-creating the slice over and over.
	bs.slice = bs.buffer[:]
	return bs
}

func (r *bitSeeker) nextWord(at int64) uint64 {
	at = at >> 6
	if at != r.have {
		r.ReadAt(r.slice, at<<3)
		r.have = at
		r.cache = binary.BigEndian.Uint64(r.slice)
	}
	return r.cache
}

func (r *bitSeeker) ReadBits(n int64) uint64 {
	var result uint64

	p := r.p & 63
	//mask := uint64((1 << (64 - uint8(p))) - 1)
	mask := maskTop[p]
	if p+n <= 64 {
		result = (r.nextWord(r.p) & mask) >> (64 - p - n)
		r.p += n
		return result
	}

	// case 2: bits lie incompletely in the given byte
	result = r.nextWord(r.p) & mask

	l := 64 - p
	r.p += l
	n -= l

	if n > 0 {
		r.p += n
		result = (result << n) | r.nextWord(r.p)>>(64-n)
	}

	return result

}

func (r *bitSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.p = offset
	case io.SeekCurrent:
		r.p += offset
	default:
		log.Panicf("Seek whence=%d not supported", whence)
	}
	return r.p, nil
}

func (r *bitSeeker) Skip(offset int64) {
	r.p += offset
}

func (r *bitSeeker) Tell() int64 {
	return r.p
}

// #endregion

// #region disk
/* FILE FORMAT
- 4 bytes - total size of file
- 1 byte: cbits
- 1 byte: abits
- 7code - number of words
- 7code - number of nodes
- 7code - number of edges
- let wbits be the number of bits to represent the total number of words in the file.
- for each node:
	- 1 bit: is node final?
	- 1 bit: fallthrough?

	- if fallthrough
		cbits: character
	else:
		1 bit: single edge?
		- if !single edge:
			7code: number of edges
			log(wbits): nskip (number of bits in skip field)
		- for each edge:
			cbits: character
			if this is not the first edge:
				nskip: count
			abits: location in bits of the node to jump to from start of file.

We define 7code to be an unsigned that can be read the following way:

result = 0
for {
	data = next 8 bits
	result = result << 7 | data & 0x7f
	if data & 0x80 == 0 break
}

*/

// Save writes the dawg to disk. Returns the number of bytes written
func (d *dawg) Save(filename string) (int64, error) {
	d.checkFinished()

	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}

	defer f.Close()
	return d.Write(f)
}

func readUint32(r io.ReaderAt, at int64) uint32 {
	data := make([]byte, 4, 4)
	_, err := r.ReadAt(data, at)
	if err != nil {
		log.Panic(err)
	}
	return (uint32(data[0]) << 24) |
		(uint32(data[1]) << 16) |
		(uint32(data[2]) << 8) |
		(uint32(data[3]) << 0)
}

func (n *node) isFallthrough(id int) bool {
	return len(n.edges) == 1 && n.edges[0].node == id+1
}

// Save writes the dawg to an io.Writer. Returns the number of bytes written
func (d *dawg) Write(wIn io.Writer) (int64, error) {
	if d.r != nil {
		return io.Copy(wIn, io.NewSectionReader(d.r, 0, d.size))
	}

	if !d.finished {
		return 0, errors.New("dawg not finished")
	}

	w := newBitWriter(wIn)

	// get maximum character and calculate cbits
	// record node addresses, calculate counts and number of edges
	addresses := make([]uint64, d.NumNodes(), d.NumNodes())
	var maxChar rune
	for _, node := range d.nodes {
		for _, edge := range node.edges {
			if edge.ch > maxChar {
				maxChar = edge.ch
			}
		}
	}

	cbits := uint64(bits.Len(uint(maxChar)))
	wbits := uint64(bits.Len(uint(d.NumAdded())))
	nskiplen := uint64(bits.Len(uint(wbits)))

	// let abits = 1
	abits := uint64(1)
	var pos uint64
	for {
		// position = 32 + 8 + 8 + encoded length of number of words, nodes, and edges
		pos = 32 + 8 + 8
		pos += unsignedLength(uint64(d.NumAdded())) * 8
		pos += unsignedLength(uint64(d.NumNodes())) * 8
		pos += unsignedLength(uint64(d.NumEdges())) * 8

		// for each node,
		for i := range addresses {
			node := d.nodes[i]

			// record its position
			addresses[i] = pos

			// final bit
			pos++

			// fallthrough?
			pos++

			if node.isFallthrough((i)) {
				pos += cbits
			} else {
				// add number of edges
				pos++ // singleEdge?

				numEdges := uint64(len(node.edges))

				// find maximum value of skip

				skip := 0
				if node.final {
					skip = 1
				}

				for _, edge := range node.edges {
					skip += d.nodes[edge.node].count
				}

				nskipbits := uint64(bits.Len(uint(skip)))

				if numEdges != 1 {
					pos += unsignedLength(numEdges) * 8
					pos += nskiplen
				}

				// add #edges * (cbits + wbits + abits)
				if numEdges > 0 {
					pos += numEdges*(cbits+nskipbits+abits) - nskipbits
				}
			}
		}

		// if file position fits into abits, then break out.
		if uint64(bits.Len(uint(pos))) <= abits {
			break
		}
		abits = uint64(bits.Len(uint(pos)))
	}

	size := (pos + 7) / 8

	// write file size, cbits, abits
	w.WriteBits(size, 32)
	w.WriteBits(cbits, 8)
	w.WriteBits(abits, 8)

	// write number of words, nodes, and edges.
	writeUnsigned(w, uint64(d.NumAdded()))
	writeUnsigned(w, uint64(d.NumNodes()))
	writeUnsigned(w, uint64(d.NumEdges()))

	// for each edge,
	for i := range addresses {
		node := d.nodes[i]
		count := 0
		if node.final {
			count++
			w.WriteBits(1, 1)
		} else {
			w.WriteBits(0, 1)
		}

		if node.isFallthrough(i) {
			w.WriteBits(1, 1)
			w.WriteBits(uint64(node.edges[0].ch), int(cbits))
		} else {
			w.WriteBits(0, 1)
			skip := 0
			if node.final {
				skip = 1
			}

			for _, edge := range node.edges {
				skip += d.nodes[edge.node].count
			}

			nskipbits := uint64(bits.Len(uint(skip)))

			if len(node.edges) == 1 {
				w.WriteBits(1, 1)
			} else {
				w.WriteBits(0, 1)
				writeUnsigned(w, uint64(len(node.edges)))
				w.WriteBits(nskipbits, int(nskiplen))
			}

			for index, edge := range node.edges {
				// write character, address
				w.WriteBits(uint64(edge.ch), int(cbits))
				if index > 0 {
					w.WriteBits(uint64(count), int(nskipbits))
				}
				w.WriteBits(addresses[edge.node], int(abits))
				count += d.nodes[edge.node].count
			}
		}
	}

	w.Flush()

	return int64(size), nil
}

const edgesOffset = (32*4 + 8 + 8)

// Read returns a finder that accesses the dawg in-place using the
// given io.ReaderAt
func Read(f io.ReaderAt, offset int64) (Finder, error) {
	size := readUint32(f, offset)
	if offset != 0 {
		f = io.NewSectionReader(f, offset, int64(size))
	}

	r := newBitSeeker(f)

	r.Seek(32, 0)
	cbits := r.ReadBits(8)
	abits := r.ReadBits(8)
	numAdded := int(readUnsigned(&r))
	numNodes := int(readUnsigned(&r))
	numEdges := int(readUnsigned(&r))
	firstNodeOffset := r.Tell()
	hasEmpty := r.ReadBits(1) == 1
	wbits := int64(bits.Len(uint(numAdded)))
	dawg := &dawg{
		finished:        true,
		numAdded:        numAdded,
		numNodes:        numNodes,
		numEdges:        numEdges,
		abits:           int64(abits),
		cbits:           int64(cbits),
		wbits:           wbits,
		hasEmptyWord:    hasEmpty,
		firstNodeOffset: firstNodeOffset,
		r:               f,
		size:            int64(size),
	}

	return dawg, nil
}

// Close ...
func (d *dawg) Close() error {
	if closer, ok := d.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (d *dawg) getEdge(r *bitSeeker, eStart edgeStart) (edgeEnd, bool, bool) {
	var edgeEnd edgeEnd
	var final, ok bool
	if d.numEdges > 0 {
		pos := int64(eStart.node)
		if pos == 0 {
			// its the first node
			pos = d.firstNodeOffset
		}

		r.Seek(pos, 0)
		nodeFinal := int(r.ReadBits(1))
		fallthr := int(r.ReadBits(1))

		if fallthr == 1 {
			ch := rune(r.ReadBits(d.cbits))
			if ch == eStart.ch {
				edgeEnd.count = nodeFinal
				edgeEnd.node = int(r.Tell())
				final = r.ReadBits(1) == 1
				ok = true
			}
		} else {
			singleEdge := r.ReadBits(1)
			numEdges := uint64(1)
			nskiplen := int64(bits.Len(uint(d.wbits)))
			nskip := int64(0)
			if singleEdge != 1 {
				numEdges = readUnsigned(r)
				nskip = int64(r.ReadBits(nskiplen))
			}

			pos = r.Tell()
			bsearch(int(numEdges), func(i int) int {
				seekTo := pos + int64(i)*int64(d.cbits+nskip+d.abits)
				if i > 0 {
					seekTo -= nskip
				}

				r.Seek(seekTo, 0)
				ch := rune(r.ReadBits(d.cbits))
				if ch == eStart.ch {
					if i > 0 {
						edgeEnd.count = int(r.ReadBits(nskip))
					} else {
						edgeEnd.count = nodeFinal
					}
					edgeEnd.node = int(r.ReadBits(d.abits))
					r.Seek(int64(edgeEnd.node), 0)
					final = r.ReadBits(1) == 1
					ok = true
				}
				return int(ch - eStart.ch)
			})
		}
	}

	return edgeEnd, final, ok
}

type nodeResult struct {
	node  int
	final bool
	edges []edgeResult
}

type edgeResult struct {
	ch    rune
	count int
	node  int
}

func (d *dawg) getNode(r *bitSeeker, node int) nodeResult {
	var result nodeResult
	pos := int64(node)
	if pos == 0 {
		// its the first node
		pos = d.firstNodeOffset
	}

	r.Seek(pos, 0)
	nodeFinal := r.ReadBits(1)
	fallthr := r.ReadBits(1)

	result.node = node
	result.final = nodeFinal == 1

	if fallthr == 1 {
		result.edges = append(result.edges, edgeResult{
			ch:    rune(r.ReadBits(d.cbits)),
			count: int(nodeFinal),
			node:  int(r.Tell()),
		})
	} else {
		nskiplen := int64(bits.Len(uint(d.wbits)))
		nskip := int64(0)

		singleEdge := r.ReadBits(1)
		numEdges := uint64(1)
		if singleEdge != 1 {
			numEdges = readUnsigned(r)
			nskip = int64(r.ReadBits(nskiplen))
		}

		for i := uint64(0); i < numEdges; i++ {
			ch := r.ReadBits(int64(d.cbits))
			var count uint64
			if i > 0 {
				count = r.ReadBits(int64(nskip))
			} else {
				count = nodeFinal
			}
			address := r.ReadBits(int64(d.abits))
			result.edges = append(result.edges, edgeResult{
				ch:    rune(ch),
				count: int(count),
				node:  int(address),
			})
		}
	}
	return result
}

// DumpFile prints out the file
func DumpFile(f io.ReaderAt) {
	r := newBitSeeker(f)
	size := r.ReadBits(32)
	fmt.Printf("[%08x] Size=%v bytes\n", r.Tell()-32, size)

	cbits := r.ReadBits(8)
	fmt.Printf("[%08x] cbits=%d\n", r.Tell()-8, cbits)

	abits := r.ReadBits(8)
	fmt.Printf("[%08x] abits=%d\n", r.Tell()-8, cbits)

	wordCount := readUnsigned(&r)
	fmt.Printf("[%08x] WordCount=%v\n", r.Tell()-int64(unsignedLength(wordCount)*8), wordCount)

	nodeCount := readUnsigned(&r)
	fmt.Printf("[%08x] NodeCount=%v\n", r.Tell()-int64(unsignedLength(nodeCount)*8), nodeCount)
	wbits := bits.Len(uint(wordCount))

	edgeCount := readUnsigned(&r)
	fmt.Printf("[%08x] EdgeCount=%v\n", r.Tell()-int64(unsignedLength(edgeCount)*8), edgeCount)

	nskiplen := bits.Len(uint(wbits))

	for i := 0; i < int(nodeCount); i++ {
		at := r.Tell()
		final := r.ReadBits(1)
		fallthr := r.ReadBits(1)

		if fallthr == 1 {
			ch := r.ReadBits(int64(cbits))
			fmt.Printf("[%08x] Node final=%d ch='%c' (fallthrough)\n", at, final, rune(ch))
			continue
		}

		singleEdge := r.ReadBits(1)
		edges := uint64(1)
		nskip := uint64(0)
		if singleEdge != 1 {
			edges = readUnsigned(&r)
			nskip = r.ReadBits(int64(nskiplen))
		}

		fmt.Printf("[%08x] Node final=%d has %d edges, skipfieldlen=%d\n", at, final, edges, nskip)

		for j := uint64(0); j < edges; j++ {
			at = r.Tell()
			ch := r.ReadBits(int64(cbits))
			var count uint64
			if j > 0 {
				count = r.ReadBits(int64(nskip))
			} else {
				count = final
			}
			address := r.ReadBits(int64(abits))
			fmt.Printf("[%08x] '%c' goto <%08x> skipping %d\n",
				at, rune(ch), address, count)
		}

	}
}

func writeUnsigned(w *bitWriter, n uint64) {
	if n < 0x7f {
		w.WriteBits(n, 8)
	} else if n < 0x3fff {
		w.WriteBits((n>>7)&0x7f|0x80, 8)
		w.WriteBits(n&0x7f, 8)
	} else if n < 0x1fffff {
		w.WriteBits((n>>14)&0x7f|0x80, 8)
		w.WriteBits((n>>7)&0x7f|0x80, 8)
		w.WriteBits(n&0x7f, 8)
	} else if n < 0xfffffff {
		w.WriteBits((n>>21)&0x7f|0x80, 8)
		w.WriteBits((n>>14)&0x7f|0x80, 8)
		w.WriteBits((n>>7)&0x7f|0x80, 8)
		w.WriteBits(n&0x7f, 8)
	} else {
		// could go further
		log.Panic("Not implemented")
	}
}

func readUnsigned(r *bitSeeker) uint64 {
	var result uint64
	for {
		d := r.ReadBits(8)
		result = (result << 7) | d&0x7f
		if d&0x80 == 0 {
			break
		}
	}
	return result
}

func unsignedLength(n uint64) uint64 {
	if n < 0x7f {
		return 1
	} else if n < 0x3fff {
		return 2
	} else if n < 0x1fffff {
		return 3
	} else if n < 0xfffffff {
		return 4
	}
	log.Panicf("Not implemented: %d", n)
	return 0
}

/** @param cmp returns target - i  or cmp(i, target)*/
func bsearch(count int, cmp func(i int) int) int {
	high := count
	low := -1
	var match, probe int
	for high-low > 1 {
		probe = (high + low) >> 1

		match = cmp(probe)

		if match == 0 {
			return probe
		} else if match < 0 {
			low = probe
		} else {
			high = probe
		}
	}

	return high
}

// #endregion

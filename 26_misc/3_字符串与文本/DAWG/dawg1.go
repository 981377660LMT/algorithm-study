// https://github.com/smhanov/dawg

// **Package dawg** 实现了一个有向无环词图（Directed Acyclic Word Graph），
// 其原理可参见作者在 [http://stevehanov.ca/blog/?id=115](http://stevehanov.ca/blog/?id=115) 上的博客。
// !**DAWG** 能够快速地查找字典（词典）中所有可能的前缀，并且支持根据单词获取其在字典中的索引。
// 与其他实现不同，本实现非常注重**内存利用**，
// 同时在支持超大字符集时依旧保持高效——它可以在一个节点下处理数千条分支，而无需依次遍历每一条。
// 存储格式尽可能地紧凑，使用了**按位（bit-level）**而非字节来记录数据，避免任何填充字节的浪费。
// 同时，对于节点或字符的数量也几乎没有实际限制。有关详细的数据格式说明，可以在 **disk.go** 文件开头找到概要。
//
// 通常情况下，如果你要使用它，先通过 **dawg.New()** 创建一个构造器（builder），然后逐个向它添加单词。需要注意的是：
// !1. 不可重复添加相同单词；
// !2. 所有待添加的单词必须按严格的字典序（alphabetical order）递增。
//
// 待所有单词都添加完后，调用 **Finish()** 方法会返回一个 **dawg.Finder** 接口。
// 你可以通过这个接口执行各种查询，比如找到某个字符串在字典中对应的所有前缀、或者检索先前添加的单词对应的索引值。
//
// 在调用 **Finish()** 之后，你可以选择使用 **Save()** 方法将构造好的 DAWG 写入磁盘。
// !之后可通过 **Load()** 方法重新打开。重新打开时，无需占用额外内存即可访问整个数据结构——一切都在磁盘上以只读方式直接访问。
//
// api:
// 1. **`New()`**: 初始化一个待构建的空 DAWG（Builder）。
// 2. **`Add(word)`**: 按字典序插入单词，内部进行**最小化**合并。
// 3. **`Finish()`**: 标记构建完毕，计算并压缩，最终序列化到内存并作为 Finder 返回。
// 4. **`FindAllPrefixesOf(input)`**: 查找全部前缀对应的单词索引。
// 5. **`IndexOf(input)`**: 返回单词对应的插入顺序（若不存在返回 -1）。
// 6. **`AtIndex(index)`**: 根据插入顺序反向取出单词。
// 7. **`Save(filename)` / `Write(io.Writer)`**: 将 DAWG 的位编码完整写出到文件/流，以便下次重用。
// 8. **`Read(io.ReaderAt, offset)`**: 从外部文件/流中按同样格式读取，生成只读 DAWG。

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
	"slices"
	"strconv"
)

func main() {
	words := []string{"hello", "world", "hello", "world", "h", "he"}
	slices.Sort(words)
	words = slices.Compact(words)

	fmt.Println(words)

	d := NewDawgBuilder()
	for _, w := range words {
		d.Add(w)
	}
	finder := d.Finish()

	// find all prefixes of "hello"
	fmt.Println(finder.FindAllPrefixesOf("hello"))

	// find the index of "world"
	fmt.Println(finder.IndexOf("world"))

	// find the word at index 2
	word, _ := finder.AtIndex(2)
	fmt.Println(word)

	// enumerate all prefixes
	finder.Enumerate(func(index int, word []rune, final bool) EnumerationResult {
		fmt.Println(index, string(word), final)
		return Continue
	})

	// save to a file
	finder.Save("dawg.bin")

	// read from a file
	f, _ := os.Open("dawg.bin")
	finder, _ = Read(f, 0)

	fmt.Println(finder.NumAdded(), finder.NumEdges(), finder.NumNodes())
	finder.Print()
}

// #region dawg
// FindResult is the result of a lookup in the d. It
// contains both the word found, and it's index based on the
// order it was added.
type FindResult struct {
	Word  string
	Index int
}

// 查询时输入“起点+字符”，用于定位下一状态
type edgeStart struct {
	node int
	ch   rune
}

func (edge edgeStart) String() string {
	return fmt.Sprintf("(%d, '%c')", edge.node, edge.ch)
}

// 存储“下一节点 ID + 跳过多少个单词计数”等信息
type edgeEnd struct {
	node  int
	count int
}

// 用于在最小化过程中暂时存储“尚未固定/合并”的路径信息(未冻结的节点)
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
type EnumerationResult int

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
	final bool        // 是否是一个单词的结尾
	count int         // 该节点下可达的单词数量，用于快速 skip
	edges []edgeStart // 从当前节点出发的边，每条边包含的字符按照字典序递增
}

// 既是构建期的 “Builder” 又是完成后的 “Finder”，通过内部标志 `finished` 区分所处阶段
type dawg struct {
	// these are erased after we finish building
	lastWord       []rune // 上一次插入的单词，用于确保有序和找公共前缀
	nextID         int
	uncheckedNodes []uncheckedNode // 存放当前还未最小化的 “路径节点”
	minimizedNodes map[string]int  // 已最小化的状态(子树) 缓存: signature -> nodeID
	nodes          map[int]*node   // 所有节点的临时存储

	// if read from a file, this is set
	r    io.ReaderAt
	size int64 // size of the readerAt

	// these are kept
	finished bool
	numAdded int
	numNodes int
	numEdges int

	// 在写盘时确定的各类“位宽”，用于后续解析
	cbits int64 // bits to represent character value
	abits int64 // bits to represent node address
	wbits int64 // bits to represent number of words / counts

	firstNodeOffset int64 // first node offset in bits in the file
	hasEmptyWord    bool
}

// NewDawgBuilder creates a new dawg
func NewDawgBuilder() Builder {
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

	// !对 [commonPrefix, end) 的 uncheckedNodes 进行 minimize，“把之前多余的后缀进行最小化合并”
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

	// 从公共前缀之后，把字符一个个加进来
	// 新建节点并插入 `edges`。同时将新建节点推到 `uncheckedNodes` 里等待后续合并
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

		// 合并剩余的 `uncheckedNodes`
		d.minimize(0)

		d.numNodes = len(d.minimizedNodes) + 1

		// Fill in the counts
		// 给每个节点计算 `count`，代表“从此节点可达多少终止单词”。这在实现 `IndexOf`、`AtIndex` 时会用到
		d.calculateSkipped(rootNode)

		// no longer need the names.
		d.uncheckedNodes = nil
		d.minimizedNodes = nil
		d.lastWord = nil

		// 重新给节点编号。构建阶段我们可能创建了很多“废弃”节点，被合并后就不再使用，
		// `renumber()` 会把现有节点重新组织成 0,1,2... 的连贯编号，以便在后续序列化中更加紧凑、简洁、可预测
		d.renumber()

		// 预写入到内存并切换到 Finder
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

		// 若中途没有对应边，则停止
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

// 返回某单词在 DAWG 中的插入顺序（0-based）。若不存在，则返回 -1.
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

// 从 `uncheckedNodes` 的末尾往前处理，
// 依次调用 `nameOf(child)` 构建子树的字符串标识（包含所有边及其子节点 ID），
// 若已经在 `minimizedNodes` 映射里，说明有重复结构，就把当前的 child 换成已经存在的那个节点；
// 否则把它插入 `minimizedNodes`。
// 处理完毕后，截断 `uncheckedNodes` 列表到指定长度 `downTo`
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

// 结点哈希.
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

// 在 `parent` 节点里找到对应 `ch` 的边，然后把它的 `node` 替换为 `child`。
// 同时删除原先 `child` 对应的节点信息（因为已被合并）。
func (d *dawg) replaceChild(parent int, ch rune, child int) {
	pnode := d.nodes[parent]
	//TODO: should be bsearch
	i := bsearch(len(pnode.edges), func(i int) int {
		return int(pnode.edges[i].ch - ch)
	})

	if pnode.edges[i].ch != ch {
		log.Panicf("Not found: %c", ch)
	}

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

// 列举 DAWG 中的所有前缀（包括中间节点也会回调）
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

// 给定插入顺序 `index`，找出对应的单词
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
	cache  uint64 // 维护一个 64-bit 缓存 `cache`，只要访问到某个 64-bit 边界，就 `ReadAt()` 读取 8 字节并 BigEndian 存入 `cache`
}

// NewBitSeeker creates a new bitreaderat
func newBitSeeker(r io.ReaderAt) bitSeeker {
	bs := bitSeeker{ReaderAt: r, have: -1}

	// 创建一个引用整个 buffer 切片的新的切片
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

// 在当前 `p`(bit 偏移) 下取 n bits。若超出当前 64-bit，就跨越下一个64-bit缓存
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

// 将整个 DAWG（节点信息）**按位**序列化到 `buffer`
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

// **关键之处**： 读出的 Finder 不会把整个结构加载到内存，而是**懒解析**：
// 每次查询时，通过 `bitSeeker` 在文件中定位到某个节点的 bit offset，然后读其结构、决定下一步走向，
// 因而在巨大文件场景下仍然只读“用到的部分”。
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

// 根据当前 `node` 的位偏移，解析“final bit, fallthrough bit, ...” 等信息找到对应字符 `letter` 的那条边以及 `count` 值
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

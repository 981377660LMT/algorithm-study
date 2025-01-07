// TrieFuzzy/TrieEditDistance/TrieLevenshteinDistance
// https://github.com/shivamMg/trie
// 支持模糊查询的Trie树(Levenshtein Distance)
// 它提供了比普通 Trie 前缀搜索更多的功能，旨在用于自动补全。
// 搜索结果的顺序是确定的。它遵循插入顺序。
//
// Api:
// - NewTrieFuzzy() *Trie
// - Put([]string, interface{}) bool
// - Delete([]string) (interface{}, bool)
// - Search([]string, ...func(*SearchOptions)) *SearchResults
// - Walk([]string, WalkFunc) error
// - Root() *Node
// - Print()
// - Sprint() string
// - PrintWithError() error
// - SprintWithError() (string, error)
// - PrintHrn(*Node)
// - SprintHrn(*Node) string
//
// SearchOptions:
// - WithExactKey
// - WithMaxResults
// - WithMaxEditDistance
//   - WithEditOps
//   - WithTopKLeastEdited
package main

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

func main() {
	trie := NewTrieFuzzy()
	// Put keys ([]string) and values (any)
	trie.Put([]string{"the"}, 1)
	trie.Put([]string{"the", "quick", "brown", "fox"}, 2) // !key 是一个元组
	trie.Put([]string{"the", "quick", "sports", "car"}, 3)
	trie.Put([]string{"the", "green", "tree"}, 4)
	trie.Put([]string{"an", "apple", "tree"}, 5)
	trie.Put([]string{"an", "umbrella"}, 6)

	trie.Root().Print()
	// Output (full trie with terminals ending with ($)):
	// ^
	// ├─ the ($)
	// │  ├─ quick
	// │  │  ├─ brown
	// │  │  │  └─ fox ($)
	// │  │  └─ sports
	// │  │     └─ car ($)
	// │  └─ green
	// │     └─ tree ($)
	// └─ an
	//    ├─ apple
	//    │  └─ tree ($)
	//    └─ umbrella ($)

	results := trie.Search([]string{"the", "quick"})
	for _, res := range results.Results {
		fmt.Println(res.Key, res.Value)
	}
	// Output (prefix-based search):
	// [the quick brown fox] 2
	// [the quick sports car] 3

	key := []string{"the", "tree"}
	results = trie.Search(key, WithMaxEditDistance(2), // An edit can be insert, delete, replace
		WithEditOps())
	for _, res := range results.Results {
		fmt.Println(res.Key, res.EditDistance) // EditDistance is number of edits needed to convert to [the tree]
	}
	// Output (results not more than 2 edits away from [the tree]):
	// [the] 1
	// [the green tree] 1
	// [an apple tree] 2
	// [an umbrella] 2

	result := results.Results[2]
	fmt.Printf("To convert %v to %v:\n", result.Key, key)
	printEditOps(result.EditOps)
	// Output (edit operations needed to covert a result to [the tree]):
	// To convert [an apple tree] to [the tree]:
	// - delete "an"
	// - replace "apple" with "the"
	// - don't edit "tree"

	results = trie.Search(key, WithMaxEditDistance(2), WithTopKLeastEdited(), WithMaxResults(2))
	for _, res := range results.Results {
		fmt.Println(res.Key, res.Value, res.EditDistance)
	}
	// Output (top 2 least edited results):
	// [the] 1 1
	// [the green tree] 4 1
}

// #region trie
const (
	RootKeyPart    = "^"
	terminalSuffix = "($)"
)

// TrieFuzzy is the trie data structure.
type TrieFuzzy struct {
	root *Node
}

// Node is a tree node inside Trie.
type Node struct {
	keyPart     string
	isTerminal  bool
	value       interface{}
	dllNode     *dllNode          // 在childrenDLL中对应的节点，用于删除
	children    map[string]*Node  // 子节点，无序存储
	childrenDLL *doublyLinkedList // 按插入顺序存储子节点，用于遍历
}

func newNode(keyPart string) *Node {
	return &Node{
		keyPart:     keyPart,
		children:    make(map[string]*Node),
		childrenDLL: &doublyLinkedList{},
	}
}

// KeyPart returns the part (string) of the key ([]string) that this Node represents.
func (n *Node) KeyPart() string {
	return n.keyPart
}

// IsTerminal returns a boolean that tells whether a key ends at this Node.
func (n *Node) IsTerminal() bool {
	return n.isTerminal
}

// Value returns the value stored for the key ending at this Node. If Node is not a terminal, it returns nil.
func (n *Node) Value() interface{} {
	return n.value
}

// SetValue sets the value for the key ending at this Node. If Node is not a terminal, value is not set.
func (n *Node) SetValue(value interface{}) {
	if n.isTerminal {
		n.value = value
	}
}

// ChildNodes returns the child-nodes of this Node.
func (n *Node) ChildNodes() []*Node {
	return n.childNodes()
}

// !_Data is used in Print(). Use Value() to get value at this Node.
func (n *Node) _Data() interface{} {
	data := n.keyPart
	if n.isTerminal {
		data += " " + terminalSuffix
	}
	return data
}

// !_Children is used in Print(). Use ChildNodes() to get child-nodes of this Node.
func (n *Node) _Children() []INode {
	children := n.childNodes()
	result := make([]INode, len(children))
	for i, child := range children {
		result[i] = INode(child)
	}
	return result
}

// Print prints the tree rooted at this Node. A Trie's root node is printed as RootKeyPart.
// All the terminal nodes are suffixed with ($).
func (n *Node) Print() {
	PrintHrn(n)
}

func (n *Node) Sprint() string {
	return SprintHrn(n)
}

func (n *Node) childNodes() []*Node {
	children := make([]*Node, 0, len(n.children))
	dllNode := n.childrenDLL.head
	for dllNode != nil {
		children = append(children, dllNode.trieNode)
		dllNode = dllNode.next
	}
	return children
}

// NewTrieFuzzy returns a new instance of Trie.
func NewTrieFuzzy() *TrieFuzzy {
	return &TrieFuzzy{root: newNode(RootKeyPart)}
}

// Root returns the root node of the Trie.
func (t *TrieFuzzy) Root() *Node {
	return t.root
}

// 如果键已经存在，则更新其值，并返回true；否则返回false.
func (t *TrieFuzzy) Put(key []string, value interface{}) (existed bool) {
	node := t.root
	for i, part := range key {
		child, ok := node.children[part]
		if !ok {
			child = newNode(part)
			child.dllNode = newDLLNode(child)
			node.children[part] = child
			node.childrenDLL.append(child.dllNode)
		}
		if i == len(key)-1 {
			existed = child.isTerminal
			child.isTerminal = true
			child.value = value
		}
		node = child
	}
	return existed
}

// 从Trie中删除一个键，并返回其对应的值。如果键不存在，返回false.
func (t *TrieFuzzy) Delete(key []string) (value interface{}, existed bool) {
	node := t.root
	parent := make(map[*Node]*Node)
	for _, keyPart := range key {
		child, ok := node.children[keyPart]
		if !ok {
			return nil, false
		}
		parent[child] = node
		node = child
	}
	if !node.isTerminal {
		return nil, false
	}
	node.isTerminal = false
	value = node.value
	node.value = nil
	for node != nil && !node.isTerminal && len(node.children) == 0 {
		delete(parent[node].children, node.keyPart)
		parent[node].childrenDLL.pop(node.dllNode)
		node = parent[node]
	}
	return value, true
}

// #endregion

// #region dll
type doublyLinkedList struct {
	head, tail *dllNode
}

type dllNode struct {
	trieNode   *Node
	next, prev *dllNode
}

func newDLLNode(trieNode *Node) *dllNode {
	return &dllNode{trieNode: trieNode}
}

func (dll *doublyLinkedList) append(node *dllNode) {
	if dll.head == nil {
		dll.head = node
		dll.tail = node
		return
	}
	dll.tail.next = node
	node.prev = dll.tail
	dll.tail = node
}

func (dll *doublyLinkedList) pop(node *dllNode) {
	if node == dll.head && node == dll.tail {
		dll.head = nil
		dll.tail = nil
		return
	}
	if node == dll.head {
		dll.head = node.next
		dll.head.prev = nil
		node.next = nil
		return
	}
	if node == dll.tail {
		dll.tail = node.prev
		dll.tail.next = nil
		node.prev = nil
		return
	}
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
	node.next = nil
	node.prev = nil
}

// #endregion

// #region heap
type searchResultMaxHeap []*SearchResult

func (s searchResultMaxHeap) Len() int {
	return len(s)
}

func (s searchResultMaxHeap) Less(i, j int) bool {
	if s[i].EditDistance == s[j].EditDistance {
		return s[i].tiebreaker > s[j].tiebreaker
	}
	return s[i].EditDistance > s[j].EditDistance
}

func (s searchResultMaxHeap) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *searchResultMaxHeap) Push(x interface{}) {
	*s = append(*s, x.(*SearchResult))
}

func (s *searchResultMaxHeap) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// #endregion

// #region search
type EditOpType int

const (
	EditOpTypeNoEdit EditOpType = iota
	EditOpTypeInsert
	EditOpTypeDelete
	EditOpTypeReplace
)

// EditOp represents an Edit Operation.

type EditOp struct {
	Type EditOpType
	// KeyPart:
	// - In case of NoEdit, KeyPart is to be retained.
	// - In case of Insert, KeyPart is to be inserted in the key.
	// - In case of Delete/Replace, KeyPart is the part of the key on which delete/replace is performed.
	KeyPart string
	// ReplaceWith is set for Type=EditOpTypeReplace
	ReplaceWith string
}

type SearchResults struct {
	Results         []*SearchResult
	heap            *searchResultMaxHeap
	tiebreakerCount int
}

type SearchResult struct {
	// Key is the key that was Put() into the Trie.
	Key []string
	// Value is the value that was Put() into the Trie.
	Value interface{}
	// EditDistance is the number of edits (insert/delete/replace) needed to convert Key into the Search()-ed key.
	EditDistance int
	// EditOps is the list of edit operations (see EditOpType) needed to convert Key into the Search()-ed key.
	EditOps []*EditOp

	tiebreaker int // 用于在堆中解决编辑距离相同的情况，保证插入顺序
}

// - WithExactKey
// - WithMaxResults
// - WithMaxEditDistance
//   - WithEditOps
//   - WithTopKLeastEdited
type SearchOptions struct {
	exactKey bool // 是否精确匹配

	maxResults      bool // 是否限制最大结果数
	maxResultsCount int

	editDistance    bool
	maxEditDistance int

	editOps         bool // 是否返回编辑操作（插入、删除、替换信息）
	topKLeastEdited bool // 是否仅返回编辑距离最小的前K个结果
}

// WithExactKey can be passed to Search(). When passed, Search() returns just the result with
// Key=Search()-ed key. If the key does not exist, result list will be empty.
func WithExactKey() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.exactKey = true
	}
}

// WithMaxResults can be passed to Search(). When passed, Search() will return at most maxResults
// number of results.
func WithMaxResults(maxResults int) func(*SearchOptions) {
	if maxResults <= 0 {
		panic(errors.New("invalid usage: maxResults must be greater than zero"))
	}
	return func(so *SearchOptions) {
		so.maxResults = true
		so.maxResultsCount = maxResults
	}
}

// WithMaxEditDistance can be passed to Search(). When passed, Search() changes its default behaviour from
// Prefix search to Edit distance search. It can be used to return "Approximate" results instead of strict
// Prefix search results.
//
// maxDistance is the maximum number of edits allowed on Trie keys to consider them as a SearchResult.
// Higher the maxDistance, more lenient and slower the search becomes.
//
// e.g. If a Trie stores English words, then searching for "wheat" with maxDistance=1 might return similar
// looking words like "wheat", "cheat", "heat", "what", etc. With maxDistance=2 it might also return words like
// "beat", "ahead", etc.
//
// Read about Edit distance: https://en.wikipedia.org/wiki/Edit_distance
func WithMaxEditDistance(maxDistance int) func(*SearchOptions) {
	if maxDistance <= 0 {
		panic(errors.New("invalid usage: maxDistance must be greater than zero"))
	}
	return func(so *SearchOptions) {
		so.editDistance = true
		so.maxEditDistance = maxDistance
	}
}

// WithEditOps can be passed to Search() alongside WithMaxEditDistance(). When passed, Search() also returns EditOps
// for each SearchResult. EditOps can be used to determine the minimum number of edit operations needed to convert
// a result Key into the Search()-ed key.
//
// e.g. Searching for "wheat" in a Trie that stores English words might return "beat". EditOps for this result might be:
// 1. insert "w" 2. replace "b" with "h".
//
// There might be multiple ways to edit a key into another. EditOps represents only one.
//
// Computing EditOps makes Search() slower.
func WithEditOps() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.editOps = true
	}
}

// WithTopKLeastEdited can be passed to Search() alongside WithMaxEditDistance() and WithMaxResults(). When passed,
// Search() returns maxResults number of results that have the lowest EditDistances. Results are sorted on EditDistance
// (lowest to highest).
//
// e.g. In a Trie that stores English words searching for "wheat" might return "wheat" (EditDistance=0), "cheat" (EditDistance=1),
// "beat" (EditDistance=2) - in that order.
func WithTopKLeastEdited() func(*SearchOptions) {
	return func(so *SearchOptions) {
		so.topKLeastEdited = true
	}
}

// Search() takes a key and some options to return results (see SearchResult) from the Trie.
// Without any options, it does a Prefix search i.e. result Keys have the same prefix as key.
// Order of the results is deterministic and will follow the order in which Put() was called for the keys.
// See "With*" functions for options accepted by Search().
func (t *TrieFuzzy) Search(key []string, options ...func(*SearchOptions)) *SearchResults {
	opts := &SearchOptions{}
	for _, f := range options {
		f(opts)
	}
	if opts.editOps && !opts.editDistance {
		panic(errors.New("invalid usage: WithEditOps() must not be passed without WithMaxEditDistance()"))
	}
	if opts.topKLeastEdited && !opts.editDistance {
		panic(errors.New("invalid usage: WithTopKLeastEdited() must not be passed without WithMaxEditDistance()"))
	}
	if opts.exactKey && opts.editDistance {
		panic(errors.New("invalid usage: WithExactKey() must not be passed with WithMaxEditDistance()"))
	}
	if opts.exactKey && opts.maxResults {
		panic(errors.New("invalid usage: WithExactKey() must not be passed with WithMaxResults()"))
	}
	if opts.topKLeastEdited && !opts.maxResults {
		panic(errors.New("invalid usage: WithTopKLeastEdited() must not be passed without WithMaxResults()"))
	}

	if opts.editDistance {
		return t.searchWithEditDistance(key, opts)
	}
	return t.search(key, opts)
}

func (t *TrieFuzzy) searchWithEditDistance(key []string, opts *SearchOptions) *SearchResults {
	// https://en.wikipedia.org/wiki/Levenshtein_distance#Iterative_with_full_matrix
	// http://stevehanov.ca/blog/?id=114
	columns := len(key) + 1
	dp := make([]int, columns)
	for i := 0; i < columns; i++ {
		dp[i] = i
	}
	m := len(key)
	if m == 0 {
		m = 1
	}
	rows := make([][]int, 1, m)
	rows[0] = dp
	results := &SearchResults{}
	if opts.topKLeastEdited {
		results.heap = &searchResultMaxHeap{}
	}

	// 优先遍历与查询键首字母匹配的节点
	keyColumn := make([]string, 1, m)
	stop := false
	// prioritize Node that has the same keyPart as key. this results in better results
	// e.g. if key=national, build with Node(keyPart=n) first so that keys like notional, nation, nationally, etc. are prioritized
	// same logic is used inside the recursive buildWithEditDistance() method
	var prioritizedNode *Node
	if len(key) > 0 {
		if prioritizedNode = t.root.children[key[0]]; prioritizedNode != nil {
			keyColumn[0] = prioritizedNode.keyPart
			t.buildWithEditDistance(&stop, results, prioritizedNode, &keyColumn, &rows, key, opts)
		}
	}

	// 遍历其他子节点
	for dllNode := t.root.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		node := dllNode.trieNode
		if node == prioritizedNode {
			continue
		}
		keyColumn[0] = node.keyPart
		t.buildWithEditDistance(&stop, results, node, &keyColumn, &rows, key, opts)
	}
	if opts.topKLeastEdited {
		n := results.heap.Len()
		results.Results = make([]*SearchResult, n)
		for n != 0 {
			result := heap.Pop(results.heap).(*SearchResult)
			result.tiebreaker = 0
			results.Results[n-1] = result
			n--
		}
		results.heap = nil
		results.tiebreakerCount = 0
	}
	return results
}

func (t *TrieFuzzy) buildWithEditDistance(stop *bool, results *SearchResults, node *Node, keyColumn *[]string, rows *[][]int, key []string, opts *SearchOptions) {
	if *stop {
		return
	}
	prevRow := (*rows)[len(*rows)-1]
	columns := len(key) + 1
	newRow := make([]int, columns)
	newRow[0] = prevRow[0] + 1
	for i := 1; i < columns; i++ {
		replaceCost := 1
		if key[i-1] == (*keyColumn)[len(*keyColumn)-1] {
			replaceCost = 0
		}
		newRow[i] = mins(
			newRow[i-1]+1,            // insertion
			prevRow[i]+1,             // deletion
			prevRow[i-1]+replaceCost, // substitution
		)
	}
	*rows = append(*rows, newRow)

	if newRow[columns-1] <= opts.maxEditDistance && node.isTerminal {
		editDistance := newRow[columns-1]
		lazyCreate := func() *SearchResult { // optimization for the case where topKLeastEdited=true and the result should not be pushed to heap
			resultKey := make([]string, len(*keyColumn))
			copy(resultKey, *keyColumn)
			result := &SearchResult{Key: resultKey, Value: node.value, EditDistance: editDistance}
			if opts.editOps {
				result.EditOps = t.getEditOps(rows, keyColumn, key)
			}
			return result
		}
		if opts.topKLeastEdited {
			results.tiebreakerCount++
			if results.heap.Len() < opts.maxResultsCount {
				result := lazyCreate()
				result.tiebreaker = results.tiebreakerCount
				heap.Push(results.heap, result)
			} else if (*results.heap)[0].EditDistance > editDistance {
				result := lazyCreate()
				result.tiebreaker = results.tiebreakerCount
				heap.Pop(results.heap)
				heap.Push(results.heap, result)
			}
		} else {
			result := lazyCreate()
			results.Results = append(results.Results, result)
			if opts.maxResults && len(results.Results) == opts.maxResultsCount {
				*stop = true
				return
			}
		}
	}

	if mins(newRow...) <= opts.maxEditDistance {
		var prioritizedNode *Node
		m := len(*keyColumn)
		if m < len(key) {
			if prioritizedNode = node.children[key[m]]; prioritizedNode != nil {
				*keyColumn = append(*keyColumn, prioritizedNode.keyPart)
				t.buildWithEditDistance(stop, results, prioritizedNode, keyColumn, rows, key, opts)
				*keyColumn = (*keyColumn)[:len(*keyColumn)-1]
			}
		}
		for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
			child := dllNode.trieNode
			if child == prioritizedNode {
				continue
			}
			*keyColumn = append(*keyColumn, child.keyPart)
			t.buildWithEditDistance(stop, results, child, keyColumn, rows, key, opts)
			*keyColumn = (*keyColumn)[:len(*keyColumn)-1]
		}
	}

	// 滚动编辑距离矩阵
	*rows = (*rows)[:len(*rows)-1]
}

func (t *TrieFuzzy) getEditOps(rows *[][]int, keyColumn *[]string, key []string) []*EditOp {
	// https://gist.github.com/jlherren/d97839b1276b9bd7faa930f74711a4b6
	ops := make([]*EditOp, 0, len(key))
	r, c := len(*rows)-1, len((*rows)[0])-1
	for r > 0 || c > 0 {
		insertionCost, deletionCost, substitutionCost := math.MaxInt, math.MaxInt, math.MaxInt
		if c > 0 {
			insertionCost = (*rows)[r][c-1]
		}
		if r > 0 {
			deletionCost = (*rows)[r-1][c]
		}
		if r > 0 && c > 0 {
			substitutionCost = (*rows)[r-1][c-1]
		}
		minCost := mins(insertionCost, deletionCost, substitutionCost)
		if minCost == substitutionCost {
			if (*rows)[r][c] > (*rows)[r-1][c-1] {
				ops = append(ops, &EditOp{Type: EditOpTypeReplace, KeyPart: (*keyColumn)[r-1], ReplaceWith: key[c-1]})
			} else {
				ops = append(ops, &EditOp{Type: EditOpTypeNoEdit, KeyPart: (*keyColumn)[r-1]})
			}
			r -= 1
			c -= 1
		} else if minCost == deletionCost {
			ops = append(ops, &EditOp{Type: EditOpTypeDelete, KeyPart: (*keyColumn)[r-1]})
			r -= 1
		} else if minCost == insertionCost {
			ops = append(ops, &EditOp{Type: EditOpTypeInsert, KeyPart: key[c-1]})
			c -= 1
		}
	}
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

func (t *TrieFuzzy) search(prefixKey []string, opts *SearchOptions) *SearchResults {
	results := &SearchResults{}
	node := t.root
	for _, keyPart := range prefixKey {
		child, ok := node.children[keyPart]
		if !ok {
			return results
		}
		node = child
	}
	if opts.exactKey {
		if node.isTerminal {
			result := &SearchResult{Key: prefixKey, Value: node.value}
			results.Results = append(results.Results, result)
		}
		return results
	}
	t.build(results, node, &prefixKey, opts)
	return results
}

func (t *TrieFuzzy) build(results *SearchResults, node *Node, prefixKey *[]string, opts *SearchOptions) (stop bool) {
	if node.isTerminal {
		key := make([]string, len(*prefixKey))
		copy(key, *prefixKey)
		result := &SearchResult{Key: key, Value: node.value}
		results.Results = append(results.Results, result)
		if opts.maxResults && len(results.Results) == opts.maxResultsCount {
			return true
		}
	}

	for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		child := dllNode.trieNode
		*prefixKey = append(*prefixKey, child.keyPart)
		stop := t.build(results, child, prefixKey, opts)
		*prefixKey = (*prefixKey)[:len(*prefixKey)-1]
		if stop {
			return true
		}
	}
	return false
}

func mins(s ...int) int {
	m := s[0]
	for _, a := range s[1:] {
		if a < m {
			m = a
		}
	}
	return m
}

// #endregion

// #region walk
type WalkFunc func(key []string, node *Node) error

// Walk traverses the Trie and calls walker function. If walker function returns an error, Walk early-returns with that error.
// Traversal follows insertion order.
func (t *TrieFuzzy) Walk(key []string, walker WalkFunc) error {
	node := t.root
	for _, keyPart := range key {
		child, ok := node.children[keyPart]
		if !ok {
			return nil
		}
		node = child
	}
	return t.walk(node, &key, walker)
}

func (t *TrieFuzzy) walk(node *Node, prefixKey *[]string, walker WalkFunc) error {
	if node.isTerminal {
		key := make([]string, len(*prefixKey))
		copy(key, *prefixKey)
		if err := walker(key, node); err != nil {
			return err
		}
	}

	for dllNode := node.childrenDLL.head; dllNode != nil; dllNode = dllNode.next {
		child := dllNode.trieNode
		*prefixKey = append(*prefixKey, child.keyPart)
		err := t.walk(child, prefixKey, walker)
		*prefixKey = (*prefixKey)[:len(*prefixKey)-1]
		if err != nil {
			return err
		}
	}
	return nil
}

// #endregion

// #region PrettyPrintDataStructures

const (
	BoxVer       = "│"
	BoxHor       = "─"
	BoxVerRight  = "├"
	BoxDownLeft  = "┐"
	BoxDownRight = "┌"
	BoxDownHor   = "┬"
	BoxUpRight   = "└"
	// Gutter is number of spaces between two adjacent child nodes.
	Gutter = 2
)

// ErrDuplicateNode indicates that a duplicate Node (node with same hash) was
// encountered while going through the tree. As of now Sprint/Print and
// SprintWithError/PrintWithError cannot operate on such trees.
//
// This error is returned by SprintWithError/PrintWithError. It's also used
// in Sprint/Print as error for panic for the same case.
//
// FIXME: create internal representation of trees that copies data
var ErrDuplicateNode = errors.New("duplicate node")

// INode represents a node in a tree. Type that satisfies INode must be a hashable type.
type INode interface {
	// _Data must return a value representing the node. It is stringified using "%v".
	// If empty, a space is used.
	_Data() interface{}
	// _Children must return a list of all child nodes of the node.
	_Children() []INode
}

type queue struct {
	arr []INode
}

func (q queue) empty() bool {
	return len(q.arr) == 0
}

func (q queue) len() int {
	return len(q.arr)
}

func (q *queue) push(n INode) {
	q.arr = append(q.arr, n)
}

func (q *queue) pop() INode {
	if q.empty() {
		return nil
	}
	ele := q.arr[0]
	q.arr = q.arr[1:]
	return ele
}

func (q *queue) peek() INode {
	if q.empty() {
		return nil
	}
	return q.arr[0]
}

// Print prints the formatted tree to standard output. To handle ErrDuplicateNode use PrintWithError.
func Print(root INode) {
	fmt.Print(Sprint(root))
}

// Sprint returns the formatted tree. To handle ErrDuplicateNode use SprintWithError.
func Sprint(root INode) string {
	parents := map[INode]INode{}
	if err := setParents(parents, root); err != nil {
		panic(err)
	}
	return sprint(parents, root)
}

// PrintWithError prints the formatted tree to standard output.
func PrintWithError(root INode) error {
	s, err := SprintWithError(root)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

// SprintWithError returns the formatted tree.
func SprintWithError(root INode) (string, error) {
	parents := map[INode]INode{}
	if err := setParents(parents, root); err != nil {
		return "", err
	}
	return sprint(parents, root), nil
}

func sprint(parents map[INode]INode, root INode) string {
	isLeftMostChild := func(n INode) bool {
		p, ok := parents[n]
		if !ok {
			// root
			return true
		}
		return p._Children()[0] == n
	}

	paddings := map[INode]int{}
	setPaddings(paddings, map[INode]int{}, 0, root)

	q := queue{}
	q.push(root)
	lines := []string{}
	for !q.empty() {
		// line storing branches, and line storing nodes
		branches, nodes := "", ""
		// runes covered
		covered := 0
		qLen := q.len()
		for i := 0; i < qLen; i++ {
			n := q.pop()
			for _, c := range n._Children() {
				q.push(c)
			}

			spaces := paddings[n] - covered
			data := safeData(n)
			nodes += strings.Repeat(" ", spaces) + data

			w := utf8.RuneCountInString(data)
			covered += spaces + w
			current, next := isLeftMostChild(n), isLeftMostChild(q.peek())
			if current {
				branches += strings.Repeat(" ", spaces)
			} else {
				branches += strings.Repeat(BoxHor, spaces)
			}

			if current && next {
				branches += BoxVer
			} else if current {
				branches += BoxVerRight
			} else if next {
				branches += BoxDownLeft
			} else {
				branches += BoxDownHor
			}

			if next {
				branches += strings.Repeat(" ", w-1)
			} else {
				branches += strings.Repeat(BoxHor, w-1)
			}
		}
		lines = append(lines, branches, nodes)
	}

	s := ""
	// ignore first line since it's the branch above root
	for _, line := range lines[1:] {
		s += strings.TrimRight(line, " ") + "\n"

	}
	return s
}

// safeData always returns non-empty representation of n's data. Empty data
// messes up tree structure, and ignoring such node will return incomplete
// tree output (tree without an entire subtree). So it returns a space.
func safeData(n INode) string {
	data := fmt.Sprintf("%v", n._Data())
	if data == "" {
		return " "
	}
	return data
}

// setPaddings sets left padding (distance of a node from the root)
// for each node in the tree.
func setPaddings(paddings map[INode]int, widths map[INode]int, pad int, root INode) {
	for _, c := range root._Children() {
		paddings[c] = pad
		setPaddings(paddings, widths, pad, c)
		pad += width(widths, c)
	}
}

// setParents sets child-parent relationships for the tree rooted
// at root.
func setParents(parents map[INode]INode, root INode) error {
	for _, c := range root._Children() {
		if _, ok := parents[c]; ok {
			return ErrDuplicateNode
		}
		parents[c] = root
		if err := setParents(parents, c); err != nil {
			return err
		}
	}
	return nil
}

// width returns either the sum of widths of it's children or its own
// data length depending on which one is bigger. widths is used in
// memoization.
func width(widths map[INode]int, n INode) int {
	if w, ok := widths[n]; ok {
		return w
	}

	w := utf8.RuneCountInString(safeData(n)) + Gutter
	widths[n] = w
	if len(n._Children()) == 0 {
		return w
	}

	sum := 0
	for _, c := range n._Children() {
		sum += width(widths, c)
	}
	if sum > w {
		widths[n] = sum
		return sum
	}
	return w
}

// PrintHr prints the horizontal formatted tree to standard output.
func PrintHr(root INode) {
	fmt.Print(SprintHr(root))
}

// SprintHr returns the horizontal formatted tree.
func SprintHr(root INode) (s string) {
	for _, line := range lines(root) {
		// ignore runes before root node
		line = string([]rune(line)[2:])
		s += strings.TrimRight(line, " ") + "\n"
	}
	return
}

func lines(root INode) (s []string) {
	data := fmt.Sprintf("%s %v ", BoxHor, root._Data())
	l := len(root._Children())
	if l == 0 {
		s = append(s, data)
		return
	}

	w := utf8.RuneCountInString(data)
	for i, c := range root._Children() {
		for j, line := range lines(c) {
			if i == 0 && j == 0 {
				if l == 1 {
					s = append(s, data+BoxHor+line)
				} else {
					s = append(s, data+BoxDownHor+line)
				}
				continue
			}

			var box string
			if i == l-1 && j == 0 {
				// first line of the last child
				box = BoxUpRight
			} else if i == l-1 {
				box = " "
			} else if j == 0 {
				box = BoxVerRight
			} else {
				box = BoxVer
			}
			s = append(s, strings.Repeat(" ", w)+box+line)
		}
	}
	return
}

// PrintHrn prints the horizontal-newline formatted tree to standard output.
func PrintHrn(root INode) {
	fmt.Print(SprintHrn(root))
}

// SprintHrn returns the horizontal-newline formatted tree.
func SprintHrn(root INode) (s string) {
	return strings.Join(lines2(root), "\n") + "\n"
}

func lines2(root INode) (s []string) {
	s = append(s, fmt.Sprintf("%v", root._Data()))
	l := len(root._Children())
	if l == 0 {
		return
	}

	for i, c := range root._Children() {
		for j, line := range lines2(c) {
			// first line of the last child
			if i == l-1 && j == 0 {
				s = append(s, BoxUpRight+BoxHor+" "+line)
			} else if j == 0 {
				s = append(s, BoxVerRight+BoxHor+" "+line)
			} else if i == l-1 {
				s = append(s, "   "+line)
			} else {
				s = append(s, BoxVer+"  "+line)
			}
		}
	}
	return
}

func printEditOps(ops []*EditOp) {
	for _, op := range ops {
		switch op.Type {
		case EditOpTypeNoEdit:
			fmt.Printf("- don't edit %q\n", op.KeyPart)
		case EditOpTypeInsert:
			fmt.Printf("- insert %q\n", op.KeyPart)
		case EditOpTypeDelete:
			fmt.Printf("- delete %q\n", op.KeyPart)
		case EditOpTypeReplace:
			fmt.Printf("- replace %q with %q\n", op.KeyPart, op.ReplaceWith)
		}
	}
}

// #endregion

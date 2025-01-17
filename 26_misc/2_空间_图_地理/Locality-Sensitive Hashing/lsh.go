// 该库包含用于 L2 度量空间中近似最近邻搜索问题的各种局部敏感哈希（LSH）算法。L2 的 LSH 函数族是 Mayur Datar 等人的研究成果。
// Currently includes:  当前包括：
//
// Basic LSH  基本局部敏感哈希 (LSH)
// Multi-probe LSH  多探针局部敏感哈希
// LSH Forest  LSH 森林
//
// 给定高维向量数据，通过 LSH 将相似（欧几里得距离小）的点映射到同样或相似的哈希值，以便在近似最近邻搜索中快速筛选候选点

package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 示例：数据维度
	dim := 2
	// 造一些示例点
	dataset := []struct {
		ID  string
		Vec Point
	}{
		{"id1", Point{1.0, 2.0}},
		{"id2", Point{2.5, 2.1}},
		{"id3", Point{9.0, 9.5}},
		{"id4", Point{1.1, 2.4}},
		// ... 可以添加更多
	}

	// 示例：我们想要查找点 {2.0, 2.0} 的近邻
	queryPt := Point{2.0, 2.0}

	// 后面我们演示三种 LSH：Basic、Multiprobe、LSH Forest
	basicExample(dim, dataset, queryPt)
	multiProbeExample(dim, dataset, queryPt)
	forestExample(dim, dataset, queryPt)
}

func basicExample(dim int, dataset []struct {
	ID  string
	Vec Point
}, query Point) {
	// 设定超参数
	// l=5个哈希表, 每个表m=3个哈希函数, w=4.0 表示桶宽
	l, m := 5, 3
	w := 4.0
	index := NewBasicLsh(dim, l, m, w)

	// 将数据集插入 LSH
	for _, data := range dataset {
		index.Insert(data.Vec, data.ID)
	}

	// 执行查询
	candidates := index.Query(query)
	fmt.Println("Basic LSH candidates:", candidates)
}

func multiProbeExample(dim int, dataset []struct {
	ID  string
	Vec Point
}, query Point) {
	// 与 Basic LSH 相同的超参数
	l, m := 5, 3
	w := 4.0
	// 额外参数 t=2, 表示 Multi-probe 每次探针数
	t := 2

	index := NewMultiprobeLsh(dim, l, m, w, t)

	// 插入数据
	for _, data := range dataset {
		index.Insert(data.Vec, data.ID)
	}

	// 查询
	candidates := index.Query(query)
	fmt.Println("Multi-probe LSH candidates:", candidates)
}

func forestExample(dim int, dataset []struct {
	ID  string
	Vec Point
}, query Point) {
	l, m := 5, 3
	w := 4.0
	index := NewLshForest(dim, l, m, w)

	// 插入
	for _, data := range dataset {
		index.Insert(data.Vec, data.ID)
	}

	// 查询, k=2表示我们想要找到最多2个近似最近邻
	k := 2
	candidates := index.Query(query, k)
	fmt.Println("LSH Forest top-k candidates:", candidates)

	// 如果用完要删除索引以回收内存
	index.Delete()
}

// #region lsh
const (
	rand_seed = 1
)

// Key is a way to index into a table.
type hashTableKey []int

// Value is an index into the input dataset.
type hashTableBucket []string

type lshParams struct {
	// Dimensionality of the input data.
	dim int
	// Number of hash tables.
	l int
	// Number of hash functions for each table.
	m int
	// Shared constant for each table.
	w float64

	// Hash function params for each (l, m).
	a [][]Point
	b [][]float64
}

// NewLshParams initializes the LSH settings.
func newLshParams(dim, l, m int, w float64) *lshParams {
	// Initialize hash params.
	a := make([][]Point, l)
	b := make([][]float64, l)
	random := rand.New(rand.NewSource(rand_seed))
	for i := range a {
		a[i] = make([]Point, m)
		b[i] = make([]float64, m)
		for j := range a[i] {
			a[i][j] = make(Point, dim)
			for d := 0; d < dim; d++ {
				a[i][j][d] = random.NormFloat64()
			}
			b[i][j] = random.Float64() * float64(w)
		}
	}
	return &lshParams{
		dim: dim,
		l:   l,
		m:   m,
		a:   a,
		b:   b,
		w:   w,
	}
}

// Hash returns all combined hash values for all hash tables.
func (lsh *lshParams) hash(point Point) []hashTableKey {
	hvs := make([]hashTableKey, lsh.l)
	for i := range hvs {
		s := make(hashTableKey, lsh.m)
		for j := 0; j < lsh.m; j++ {
			hv := (point.Dot(lsh.a[i][j]) + lsh.b[i][j]) / lsh.w
			s[j] = int(math.Floor(hv))
		}
		hvs[i] = s
	}
	return hvs
}

// #endregion

// #region basic_lsh
type basicHashTableKey string

type hashTable map[basicHashTableKey]hashTableBucket

// BasicLsh implements the original LSH algorithm for L2 distance.
type BasicLsh struct {
	*lshParams
	// Hash tables.
	tables []hashTable
}

// NewBasicLsh creates a basic LSH for L2 distance.
// dim is the diminsionality of the data, l is the number of hash
// tables to use, m is the number of hash values to concatenate to
// form the key to the hash tables, w is the slot size for the
// family of LSH functions.
func NewBasicLsh(dim, l, m int, w float64) *BasicLsh {
	tables := make([]hashTable, l)
	for i := range tables {
		tables[i] = make(hashTable)
	}
	return &BasicLsh{
		lshParams: newLshParams(dim, l, m, w),
		tables:    tables,
	}
}

func (index *BasicLsh) toBasicHashTableKeys(keys []hashTableKey) []basicHashTableKey {
	basicKeys := make([]basicHashTableKey, index.l)
	for i, key := range keys {
		s := ""
		for _, hashVal := range key {
			s += fmt.Sprintf("%.16x", hashVal)
		}
		basicKeys[i] = basicHashTableKey(s)
	}
	return basicKeys
}

// Insert adds a new data point to the LSH.
// id is the unique identifier for the data point.
func (index *BasicLsh) Insert(point Point, id string) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(point))
	// Insert key into all hash tables
	var wg sync.WaitGroup
	wg.Add(len(index.tables))
	for i := range index.tables {
		hv := hvs[i]
		table := index.tables[i]
		go func(table hashTable, hv basicHashTableKey) {
			if _, exist := table[hv]; !exist {
				table[hv] = make(hashTableBucket, 0)
			}
			table[hv] = append(table[hv], id)
			wg.Done()
		}(table, hv)
	}
	wg.Wait()
}

// Query finds the ids of approximate nearest neighbour candidates,
// in un-sorted order, given the query point,
func (index *BasicLsh) Query(q Point) []string {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(index.hash(q))
	// Keep track of keys seen
	seen := make(map[string]bool)
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				if _, exist := seen[id]; exist {
					continue
				}
				seen[id] = true
			}
		}
	}
	// Collect results
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// Delete removes a new data point to the LSH.
// id is the unique identifier for the data point.
func (index *BasicLsh) Delete(id string) {
	// Delete key from all hash tables
	var wg sync.WaitGroup
	wg.Add(len(index.tables))
	for i := range index.tables {
		table := index.tables[i]
		go func(table hashTable) {
			for tableIndex, bucket := range table {
				for index, identifier := range bucket {
					if id == identifier {
						table[tableIndex] = remove(bucket, index)
						if len(table[tableIndex]) == 0 {
							delete(table, tableIndex)
						}
					}
				}
			}
			wg.Done()
		}(table)
	}
	wg.Wait()
}

func remove(original []string, index int) []string {
	original[index] = original[len(original)-1]
	original = original[:len(original)-1]
	return original
}

// #endregion

// #region multi_probe_lsh
type perturbSet map[int]bool

func (ps perturbSet) isValid(m int) bool {
	for key := range ps {
		// At most one perturbation on same index.
		if _, ok := ps[2*m+1-key]; ok {
			return false
		}
		// No keys larger than 2m.
		if key > 2*m {
			return false
		}
	}
	return true
}

func (ps perturbSet) shift() perturbSet {
	next := make(perturbSet)
	max := 0
	for k := range ps {
		if k > max {
			max = k
		}
		next[k] = true
	}
	delete(next, max)
	next[max+1] = true
	return next
}

func (ps perturbSet) expand() perturbSet {
	next := make(perturbSet)
	max := 0
	for k := range ps {
		if k > max {
			max = k
		}
		next[k] = true
	}
	next[max+1] = true
	return next
}

// A pair of perturbation set and its score.
type perturbSetPair struct {
	ps    perturbSet
	score float64
}

// perturbSetHeap is a min-heap of perturbSetPairs.
type perturbSetHeap []perturbSetPair

func (h perturbSetHeap) Len() int           { return len(h) }
func (h perturbSetHeap) Less(i, j int) bool { return h[i].score < h[j].score }
func (h perturbSetHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *perturbSetHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(perturbSetPair))
}

func (h *perturbSetHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MultiprobeLsh implements the Multi-probe LSH algorithm by Qin Lv et.al.
// The Multi-probe LSH does not support k-NN query directly.
type MultiprobeLsh struct {
	*BasicLsh
	// The size of our probe sequence.
	t int

	// The scores of perturbation values.
	scores []float64

	perturbSets []perturbSet

	// Each hash table has a list of perturbation vectors
	// each perturbation vector is list of -+ 1 or 0 that will
	// be applied to the hashTableKey of the query hash value
	// t x l x m
	perturbVecs [][][]int
}

// NewMultiprobeLsh creates a new Multi-probe LSH for L2 distance.
// dim is the diminsionality of the data, l is the number of hash
// tables to use, m is the number of hash values to concatenate to
// form the key to the hash tables, and w is the slot size for the
// family of LSH functions.
// t is the number of perturbation vectors that will be applied to
// each query.
// Increasing t increases the running time of the Query function.
func NewMultiprobeLsh(dim, l, m int, w float64, t int) *MultiprobeLsh {
	index := &MultiprobeLsh{
		BasicLsh: NewBasicLsh(dim, l, m, w),
		t:        t,
	}
	index.initProbeSequence()
	return index
}

func (index *MultiprobeLsh) initProbeSequence() {
	m := index.m
	index.scores = make([]float64, 2*m)
	// Use j's starting from 1 to match the paper.
	for j := 1; j <= m; j++ {
		index.scores[j-1] = float64(j*(j+1)) / float64(4*(m+1)*(m+2))
	}
	for j := m + 1; j <= 2*m; j++ {
		index.scores[j-1] = 1 - float64(2*m+1-j)/float64(m+1) + float64((2*m+1-j)*(2*m+2-j))/float64(4*(m+1)*(m+2))
	}
	index.genPerturbSets()
	index.genPerturbVecs()
}

func (index *MultiprobeLsh) getScore(ps *perturbSet) float64 {
	score := 0.0
	for j := range *ps {
		score += index.scores[j-1]
	}
	return score
}

func (index *MultiprobeLsh) genPerturbSets() {
	setHeap := make(perturbSetHeap, 1)
	start := perturbSet{1: true}
	setHeap[0] = perturbSetPair{
		ps:    start,
		score: index.getScore(&start),
	}
	heap.Init(&setHeap)
	index.perturbSets = make([]perturbSet, index.t)
	m := index.m

	for i := 0; i < index.t; i++ {
		for counter := 0; true; counter++ {
			currentTop := heap.Pop(&setHeap).(perturbSetPair)
			nextShift := currentTop.ps.shift()
			heap.Push(&setHeap, perturbSetPair{
				ps:    nextShift,
				score: index.getScore(&nextShift),
			})
			nextExpand := currentTop.ps.expand()
			heap.Push(&setHeap, perturbSetPair{
				ps:    nextExpand,
				score: index.getScore(&nextExpand),
			})

			if currentTop.ps.isValid(m) {
				index.perturbSets[i] = currentTop.ps
				break
			}
			if counter >= 2*m {
				panic("too many iterations, probably infinite loop!")
			}
		}
	}
}

func (index *MultiprobeLsh) genPerturbVecs() {
	// First we need to generate the permutation tables
	// that maps the ids of the unit perturbation in each
	// perturbation set to the index of the unit hash
	// value
	perms := make([][]int, index.l)
	for i := range index.tables {
		random := rand.New(rand.NewSource(int64(i)))
		perm := random.Perm(index.m)
		perms[i] = make([]int, index.m*2)
		for j := 0; j < index.m; j++ {
			perms[i][j] = perm[j]
		}
		for j := 0; j < index.m; j++ {
			perms[i][j+index.m] = perm[index.m-1-j]
		}
	}

	// Generate the vectors
	index.perturbVecs = make([][][]int, len(index.perturbSets))
	for i, ps := range index.perturbSets {
		perTableVecs := make([][]int, index.l)
		for j := range perTableVecs {
			vec := make([]int, index.m)
			for k := range ps {
				mapped_ind := perms[j][k-1]
				if k > index.m {
					// If it is -1
					vec[mapped_ind] = -1
				} else {
					// if it is +1
					vec[mapped_ind] = 1
				}
			}
			perTableVecs[j] = vec
		}
		index.perturbVecs[i] = perTableVecs
	}
}

func (index *MultiprobeLsh) queryHelper(tableKeys []hashTableKey, out chan<- string) {
	// Apply hash functions
	hvs := index.toBasicHashTableKeys(tableKeys)

	// Lookup in each table.
	for i, table := range index.tables {
		if candidates, exist := table[hvs[i]]; exist {
			for _, id := range candidates {
				out <- id
			}
		}
	}
}

// perturb returns the result of applying perturbation on each baseKey.
func (index *MultiprobeLsh) perturb(baseKey []hashTableKey, perturbation [][]int) []hashTableKey {
	if len(baseKey) != len(perturbation) {
		panic("Number tables does not match with number of perturb vecs")
	}
	perturbedTableKeys := make([]hashTableKey, len(baseKey))
	for i, p := range perturbation {
		perturbedTableKeys[i] = make(hashTableKey, index.m)
		for j, h := range baseKey[i] {
			perturbedTableKeys[i][j] = h + p[j]
		}
	}
	return perturbedTableKeys
}

// Query finds the ids of nearest neighbour candidates,
// given the query point
func (index *MultiprobeLsh) Query(q Point) []string {
	// Hash
	baseKey := index.hash(q)
	// Query
	results := make(chan string)
	go func() {
		defer close(results)
		for i := 0; i < len(index.perturbVecs)+1; i++ {
			perturbedTableKeys := baseKey
			if i != 0 {
				// Generate new hash key based on perturbation.
				perturbedTableKeys = index.perturb(baseKey, index.perturbVecs[i-1])
			}
			// Perform lookup.
			index.queryHelper(perturbedTableKeys, results)
		}
	}()
	seen := make(map[string]bool)
	for id := range results {
		if _, exist := seen[id]; exist {
			continue
		}
		seen[id] = true
	}
	// Collect results
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// #endregion

// #region lsh_forest
type treeNode struct {
	// Hash key for this intermediate node. nil/empty for root nodes.
	hashKey int
	// A list of ids to the source dataset, only leaf nodes have non-empty ids.
	ids []string
	// Child nodes, keyed by partial hash value.
	children map[int]*treeNode
}

func (node *treeNode) recursiveDelete() {
	for _, child := range node.children {
		if len((child).children) > 0 {
			(child).recursiveDelete()
		}
		if len(child.ids) > 0 {
			node.ids = nil
		}
	}
	node.ids = nil
	node.children = nil
}

// recursiveAdd recurses down the tree to find the correct location to insert id.
// Returns whether a new hash value was added.
func (node *treeNode) recursiveAdd(level int, id string, tableKey hashTableKey) bool {
	if level == len(tableKey) {
		node.ids = append(node.ids, id)
		return false
	}
	// Check if next hash exists in children map. If not, create.
	var next *treeNode
	hasNewHash := false
	if nextNode, ok := node.children[tableKey[level]]; !ok {
		next = &treeNode{
			hashKey:  tableKey[level],
			ids:      make([]string, 0),
			children: make(map[int]*treeNode),
		}
		node.children[tableKey[level]] = next
		hasNewHash = true
	} else {
		next = nextNode
	}
	// Recurse using next level's hash value.
	recursive := next.recursiveAdd(level+1, id, tableKey)
	return hasNewHash || recursive
}

func tab(times int) {
	for i := 0; i < times; i++ {
		fmt.Print("    ")
	}
}

func (node *treeNode) dump(level int) {
	tab(level)
	fmt.Printf("{ (%v): ids %v ", node.hashKey, node.ids)
	if len(node.children) > 0 {
		fmt.Printf("[\n")
		for _, v := range node.children {
			v.dump(level + 1)
		}
		tab(level)
		fmt.Print("] }\n")
	} else {
		fmt.Print("}\n")
	}
}

type prefixTree struct {
	// Number of distinct elements in the tree.
	count int
	// Pointer to the root node.
	root *treeNode
}

func (tree *prefixTree) insertIntoTree(id string, tableKey hashTableKey) {
	if tree.root.recursiveAdd(0, id, tableKey) {
		tree.count++
	}
}

// lookup find ids and write them to out channel
func (tree *prefixTree) lookup(maxLevel int, tableKey hashTableKey,
	done <-chan struct{}, out chan<- string) {
	currentNode := tree.root
	for level := 0; level < len(tableKey) && level < maxLevel; level++ {
		if next, ok := currentNode.children[tableKey[level]]; ok {
			currentNode = next
		} else {
			return
		}
	}

	// Grab all ids of nodes descendent from the current node.
	queue := []*treeNode{currentNode}
	for len(queue) > 0 {
		// Add node's ids to main list.
		for _, id := range queue[0].ids {
			select {
			case out <- id:
			case <-done:
				return
			}
		}

		// Add children.
		for _, child := range queue[0].children {
			queue = append(queue, child)
		}

		// Done with head.
		queue = queue[1:]
	}
}

// LshForest implements the LSH Forest algorithm by Mayank Bawa et.al.
// It supports both nearest neighbour candidate query and k-NN query.
type LshForest struct {
	// Embedded type
	*lshParams
	// Trees.
	trees []prefixTree
}

// NewLshForest creates a new LSH Forest for L2 distance.
// dim is the diminsionality of the data, l is the number of hash
// tables to use, m is the number of hash values to concatenate to
// form the key to the hash tables, w is the slot size for the
// family of LSH functions.
func NewLshForest(dim, l, m int, w float64) *LshForest {
	trees := make([]prefixTree, l)
	for i := range trees {
		trees[i].count = 0
		trees[i].root = &treeNode{
			hashKey:  0,
			ids:      make([]string, 0),
			children: make(map[int]*treeNode),
		}
	}
	return &LshForest{
		lshParams: newLshParams(dim, l, m, w),
		trees:     trees,
	}
}

// Delete releases the memory used by this index.
func (index *LshForest) Delete() {
	for _, tree := range index.trees {
		(*tree.root).recursiveDelete()
	}
}

// Insert adds a new data point to the LSH Forest.
// id is the unique identifier for the data point.
func (index *LshForest) Insert(point Point, id string) {
	// Apply hash functions.
	hvs := index.hash(point)
	// Parallel insert
	var wg sync.WaitGroup
	wg.Add(len(index.trees))
	for i := range index.trees {
		hv := hvs[i]
		tree := &(index.trees[i])
		go func(tree *prefixTree, hv hashTableKey) {
			tree.insertIntoTree(id, hv)
			wg.Done()
		}(tree, hv)
	}
	wg.Wait()
}

// Helper that queries all trees and returns an channel ids.
func (index *LshForest) queryHelper(maxLevel int, tableKeys []hashTableKey, done <-chan struct{}, out chan<- string) {
	var wg sync.WaitGroup
	wg.Add(len(index.trees))
	for i := range index.trees {
		key := tableKeys[i]
		tree := index.trees[i]
		go func() {
			tree.lookup(maxLevel, key, done, out)
			wg.Done()
		}()
	}
	wg.Wait()
}

// Query finds at top-k ids of approximate nearest neighbour candidates,
// in unsorted order, given the query point.
func (index *LshForest) Query(q Point, k int) []string {
	// Apply hash functions
	hvs := index.hash(q)
	// Query
	results := make(chan string)
	done := make(chan struct{})
	go func() {
		for maxLevels := index.m; maxLevels >= 0; maxLevels-- {
			select {
			case <-done:
				return
			default:
				index.queryHelper(maxLevels, hvs, done, results)
			}
		}
		close(results)
	}()
	seen := make(map[string]bool)
	for id := range results {
		if len(seen) >= k {
			break
		}
		if _, exist := seen[id]; exist {
			continue
		}
		seen[id] = true
	}
	close(done)
	// Collect results
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// Dump prints out the index for debugging
func (index *LshForest) dump() {
	for i, tree := range index.trees {
		fmt.Printf("Tree %d (%d hash values):\n", i, tree.count)
		tree.root.dump(0)
	}
}

// #endregion

// #region metrics
// Point is a vector in the L2 metric space.
type Point []float64

// Dot returns the dot product of two points.
func (p Point) Dot(q Point) float64 {
	s := 0.0
	for i := 0; i < len(p); i++ {
		s += p[i] * q[i]
	}
	return s
}

// L2 returns the L2 distance of two points.
func (p Point) L2(q Point) float64 {
	s := 0.0
	for i := 0; i < len(p); i++ {
		d := p[i] - q[i]
		s += d * d
	}
	return math.Sqrt(s)
}

// #endregion

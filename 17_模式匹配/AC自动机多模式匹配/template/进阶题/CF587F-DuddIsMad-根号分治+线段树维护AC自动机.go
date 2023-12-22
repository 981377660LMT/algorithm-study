package main

import (
	"bufio"
	"fmt"
	"os"
)

// Duff is Mad
// !给定n个字符串，q次询问[start:end)中的字符串在s[index]中出现次数之和.
// 字符串总长度<=1e5.
// https://www.luogu.com.cn/problem/CF587F
// https://www.luogu.com.cn/blog/abruce-home/solution-cf587f
//
// CF587F-DuddIsMad-根号分治+线段树维护AC自动机
// !对于短串，在线段树的每个节点维护一个AC自动机.
// !对于长串，这样的串不超过sqrt(n)个，直接预处理.
// 为了卡空间，把根号分治的大小开大一点，减少一点预处理
// 时间复杂度O(nsqrt(sum)logn)，空间复杂度O(nsqrt(sum)+26nlogn)
// !MLE
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	THRESHOLD := 800

	var n, q int
	fmt.Fscan(in, &n, &q)

	seg := NewDI(n)
	acms := make([]*AC, seg.Size())
	for i := range acms {
		acms[i] = NewAC(26, 97)
	}
	bigId := make([]int32, n)
	for i := 0; i < n; i++ {
		bigId[i] = -1
	}
	bigCount := int32(0)
	words := make([]string, n)
	// 长串需要预处理，短串需要插入AC自动机
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
		if len(words[i]) <= THRESHOLD {
			seg.EnumeratePoint(i, func(segmentId int) {
				acms[segmentId].AddString(words[i])
			})
		} else {
			bigId[i] = bigCount
			bigCount++
		}
	}
	for i := range acms {
		acms[i].BuildSuffixLink()
	}

	preSum := make([][]int16, bigCount) // 每个串在bigId[i]中出现次数的前缀和

	for i := 0; i < n; i++ {
		id := bigId[i]
		if id == -1 {
			continue
		}
		preSum[id] = make([]int16, n+1)
		longText := words[i]
		for j := 0; j < n; j++ {
			preSum[id][j+1] = preSum[id][j] + int16(CountIndexOfAll(longText, words[j], 0, nil))
		}
	}

	for i := 0; i < q; i++ {
		var start, end, index int
		fmt.Fscan(in, &start, &end, &index)
		start--
		index--

		res := int32(0)
		if id := bigId[index]; id == -1 {
			seg.EnumerateSegment(start, end, func(segmentId int) {
				acm := acms[segmentId]
				counter := acm.Counter
				pos := 0
				text := words[index]
				for _, c := range text {
					pos = acm.Move(pos, int(c))
					res += counter[pos]
				}
			})
		} else {
			res = int32(preSum[id][end] - preSum[id][start])
		}

		fmt.Println(res)
	}
}

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个字符串.
type AC struct {
	Counter  []int32
	sigma    int32     // 字符集大小.
	offset   int32     // 字符集的偏移量.
	children [][]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
}

func NewAC(sigma, offset int) *AC {
	res := &AC{sigma: int32(sigma), offset: int32(offset)}
	res.newNode()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *AC) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := int32(0)
	for _, s := range str {
		ord := int32(s) - trie.offset
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
		}
		pos = trie.children[pos][ord]
	}
	trie.Counter[pos]++
	return int(pos)
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *AC) Move(pos int, ord int) int {
	ord -= int(trie.offset)
	return int(trie.children[pos][ord])
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *AC) Size() int {
	return len(trie.children)
}

func (trie *AC) Empty() bool {
	return len(trie.children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *AC) BuildSuffixLink() {
	suffixLink := make([]int32, len(trie.children))
	for i := range suffixLink {
		suffixLink[i] = -1
	}
	bfsOrder := make([]int32, len(trie.children))
	head, tail := 0, 0
	bfsOrder[tail] = 0
	tail++
	for head < tail {
		v := bfsOrder[head]
		head++
		for i, next := range trie.children[v] {
			if next == -1 {
				continue
			}
			bfsOrder[tail] = next
			tail++
			f := suffixLink[v]
			for f != -1 && trie.children[f][i] == -1 {
				f = suffixLink[f]
			}
			suffixLink[next] = f
			if f == -1 {
				suffixLink[next] = 0
			} else {
				suffixLink[next] = trie.children[f][i]
			}
		}
	}

	for _, v := range bfsOrder {
		for i, next := range trie.children[v] {
			if next == -1 {
				f := suffixLink[v]
				if f == -1 {
					trie.children[v][i] = 0
				} else {
					trie.children[v][i] = trie.children[f][i]
				}
			}
		}
		// dp
		if v != 0 {
			trie.Counter[v] += trie.Counter[suffixLink[v]]
		}
	}
}

func (trie *AC) newNode() int32 {
	nexts := make([]int32, trie.sigma)
	for i := range nexts {
		nexts[i] = -1
	}
	trie.children = append(trie.children, nexts)
	trie.Counter = append(trie.Counter, 0)
	return int32(len(trie.children) - 1)
}

// 线段树每个节点维护一个ac自动机

type DI struct {
	Offset int // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	n      int
}

// 线段树分割区间.
// 将长度为n的序列搬到长度为offset+n的线段树上, 以实现快速的区间操作.
func NewDI(n int) *DI {
	offset := 1
	for offset < n {
		offset <<= 1
	}
	return &DI{Offset: offset, n: n}
}

// 获取原下标为i的元素在树中的(叶子)编号.
func (d *DI) Id(rawIndex int) int {
	return rawIndex + d.Offset
}

// O(logn) 顺序遍历`[start,end)`区间对应的线段树节点.
func (d *DI) EnumerateSegment(start, end int, f func(segmentId int)) {
	if !(0 <= start && start <= end && end <= d.n) {
		panic("invalid range")
	}
	for _, i := range d.getSegmentIds(start, end) {
		f(i)
	}
}

func (d *DI) EnumeratePoint(index int, f func(segmentId int)) {
	if index < 0 || index >= d.n {
		return
	}
	index += d.Offset
	for index > 0 {
		f(index)
		index >>= 1
	}
}

// O(n) 从根向叶子方向push.
func (d *DI) PushDown(f func(parent, child int)) {
	for p := 1; p < d.Offset; p++ {
		f(p, p<<1)
		f(p, p<<1|1)
	}
}

// O(n) 从叶子向根方向update.
func (d *DI) PushUp(f func(parent, child1, child2 int)) {
	for p := d.Offset - 1; p > 0; p-- {
		f(p, p<<1, p<<1|1)
	}
}

// 线段树的节点个数.
func (d *DI) Size() int {
	return d.Offset + d.n
}

func (d *DI) IsLeaf(segmentId int) bool {
	return segmentId >= d.Offset
}

func (d *DI) getSegmentIds(start, end int) []int {
	if !(0 <= start && start <= end && end <= d.n) {
		panic("invalid range")
	}
	leftRes, rightRes := []int{}, []int{}
	for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
		if start&1 == 1 {
			leftRes = append(leftRes, start)
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = append(rightRes, end)
		}
	}
	for i := len(rightRes) - 1; i >= 0; i-- {
		leftRes = append(leftRes, rightRes[i])
	}
	return leftRes
}

type Str = string

func GetNext(pattern Str) []int32 {
	n := int32(len(pattern))
	next := make([]int32, n)
	j := int32(0)
	for i := int32(1); i < n; i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = int32(next[j-1])
		}
		if pattern[i] == pattern[j] {
			j++
		}
		next[i] = j
	}
	return next
}

func CountIndexOfAll(longer Str, shorter Str, position int32, nexts []int32) int {
	if len(shorter) == 0 {
		return 0
	}
	if len(longer) < len(shorter) {
		return 0
	}
	res := 0
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	n := int32(len(longer))
	m := int32(len(shorter))
	hitJ := int32(0)
	for i := position; i < n; i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == m {
			res++
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc257_g()
	// fmt.Println(ZAlgo("ababab"))
}

const INF int = 1e18

// [ABC257G] Prefix Concatenation
// https://www.luogu.com.cn/problem/AT_abc257_g
// 给定两个字符串s和t，将t划分为k个子串，使得每一个子串都是字符串s的前缀.
// 求满足条件的最小的k，不存在则输出-1.
// |t|,|s|<=5e5
//
// 1.求出t的每个后缀与s的最长公共前缀长度.
// 2.后缀优化建图+bfs01求最短路.
func abc257_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)

	n2 := len(t)
	lcp := ZAlgoTwoString(s, t)
	adjList := make([][][2]int, n2+1)

	// 前后缀优化建图 (i 向范围`[i+1, i+lcp[i]]`连边)
	for i, v := range lcp {
		if v > 0 {
			adjList[i] = append(adjList[i], [2]int{i + v, 1})
		}
	}
	for i := n2 - 1; i >= 0; i-- {
		adjList[i+1] = append(adjList[i+1], [2]int{i, 0})
	}

	dist := bfs01(adjList, 0)
	if res := dist[n2]; res == INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, res)
	}
}

func sumScores(s string) int64 {
	n := len(s)
	res := int64(0)
	z := ZAlgo(s)
	for i := 0; i < n; i++ {
		res += int64(z[i])
	}
	return res + int64(n)
}

// z算法求字符串每个后缀与原串的最长公共前缀长度
//
// z[0]=0
// z[i]是后缀s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
func ZAlgo(s string) []int {
	n := len(s)
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

// 对t的每个后缀，求与s的最长公共前缀长度.
// z[i]是后缀t[i:]与s的最长公共前缀(LCP)的长度 (i>=0)
func ZAlgoTwoString(s, t string) []int {
	n1, n2 := len(s), len(t)
	z := ZAlgo(s + t)
	for i := n1; i < n1+n2; i++ {
		z[i] = min(z[i], n1)
	}
	return z[n1:]
}

func ZAlgoNums(nums []int) []int {
	n := len(nums)
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && nums[z[i]] == nums[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

// 对nums2的每个后缀，求与nums1的最长公共前缀长度.
// z[i]是后缀nums2[i:]与nums1的最长公共前缀(LCP)的长度 (i>=0)
func ZAlgoTwoNums(nums1, nums2 []int) []int {
	n1, n2 := len(nums1), len(nums2)
	nums := append(nums1, nums2...)
	z := ZAlgoNums(nums)
	for i := n1; i < n1+n2; i++ {
		z[i] = min(z[i], n1)
	}
	return z[n1:]
}

type DiffArray struct {
	diff  []int
	dirty bool
}

func NewDiffArray(n int) *DiffArray {
	return &DiffArray{
		diff: make([]int, n+1),
	}
}

func (d *DiffArray) Add(start, end, delta int) {
	if start < 0 {
		start = 0
	}
	if end >= len(d.diff) {
		end = len(d.diff) - 1
	}
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffArray) Build() {
	if d.dirty {
		preSum := make([]int, len(d.diff))
		for i := 1; i < len(d.diff); i++ {
			preSum[i] = preSum[i-1] + d.diff[i]
		}
		d.diff = preSum
		d.dirty = false
	}
}

func (d *DiffArray) Get(pos int) int {
	d.Build()
	return d.diff[pos]
}

func (d *DiffArray) GetAll() []int {
	d.Build()
	return d.diff[:len(d.diff)-1]
}

func bfs01(adjList [][][2]int, start int) (dist []int) {
	n := len(adjList)
	dist = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
	}
	queue := NewDeque(n)
	queue.Append(start)
	dist[start] = 0
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

type D = int
type Deque struct{ l, r []D }

func NewDeque(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
	"strconv"
	"strings"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	H, W, N := io.NextInt(), io.NextInt(), io.NextInt()
	xs, ys, ws := make([]int32, N), make([]int32, N), make([]E, N)
	for i := 0; i < N; i++ {
		x, y := io.NextInt(), io.NextInt()
		xs[i], ys[i], ws[i] = int32(x), int32(y), E{id: int32(i), max: -INF32}
	}
	order := make([]int, N)
	for i := 0; i < N; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		if xs[order[i]] == xs[order[j]] {
			return ys[order[i]] < ys[order[j]]
		}
		return xs[order[i]] < xs[order[j]]
	})

	dp := NewSegmentTree2DSparse32FastWithWeights(xs, ys, ws, true)
	pre := make([]int32, N)
	for i := 0; i < N; i++ {
		pre[i] = -1
	}

	for _, o := range order {
		curX, curY := xs[o], ys[o]
		preMax := dp.Query(1, curX+1, 1, curY+1)
		dp.Set(int32(o), E{id: int32(o), max: preMax.max + 1})
		if preMax.id != -1 {
			pre[o] = preMax.id
		}
	}

	best := dp.Query(int32(1), int32(H+1), int32(1), int32(W+1))
	path := []int32{best.id}
	curId := best.id
	for pre[curId] != -1 {
		curId = pre[curId]
		path = append(path, curId)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	points := make([][2]int, 0, len(path)+2)
	points = append(points, [2]int{1, 1})
	for _, o := range path {
		points = append(points, [2]int{int(xs[o]), int(ys[o])})
	}
	points = append(points, [2]int{H, W})

	res := strings.Builder{}
	for i := 1; i < len(points); i++ {
		x1, y1 := points[i-1][0], points[i-1][1]
		x2, y2 := points[i][0], points[i][1]
		res.WriteString(strings.Repeat("R", y2-y1))
		res.WriteString(strings.Repeat("D", x2-x1))
	}

	io.Println(best.max)
	io.Println(res.String())
}

const INF32 int32 = 1e9 + 10

type E = struct {
	id  int32
	max int32
}

func e() E { return E{id: -1, max: 0} }
func op(a, b E) E {
	if a.max > b.max {
		return a
	}
	return b
}

type SegmentTree2DSparse32Fast struct {
	n          int32
	keyX       []int32
	keyY       []int32
	minX       int32
	allY       []int32
	pos        []int32
	indptr     []int32
	size       int32
	data       []E
	discretize bool
	unit       E
	toLeft     []int32
}

func NewSegmentTree2DSparse32Fast(xs, ys []int32, discretize bool) *SegmentTree2DSparse32Fast {
	res := &SegmentTree2DSparse32Fast{discretize: discretize, unit: e()}
	ws := make([]E, len(xs))
	for i := range ws {
		ws[i] = res.unit
	}
	res._build(xs, ys, ws)
	return res
}

func NewSegmentTree2DSparse32FastWithWeights(xs, ys []int32, ws []E, discretize bool) *SegmentTree2DSparse32Fast {
	res := &SegmentTree2DSparse32Fast{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

func (t *SegmentTree2DSparse32Fast) Update(rawIndex int32, value E) {
	i := int32(1)
	p := t.pos[rawIndex]
	indPtr, toLeft := t.indptr, t.toLeft
	for {
		t._update(i, p-indPtr[i], value)
		if i >= t.size {
			break
		}
		lc := toLeft[p] - toLeft[indPtr[i]]
		rc := p - indPtr[i] - lc
		if toLeft[p+1] > toLeft[p] {
			p = indPtr[i<<1] + lc
			i <<= 1
		} else {
			p = indPtr[i<<1|1] + rc
			i = i<<1 | 1
		}
	}
}

func (t *SegmentTree2DSparse32Fast) Set(rawIndex int32, value E) {
	i := int32(1)
	p := t.pos[rawIndex]
	indPtr, toLeft := t.indptr, t.toLeft
	for {
		t._set(i, p-indPtr[i], value)
		if i >= t.size {
			break
		}
		lc := toLeft[p] - toLeft[indPtr[i]]
		rc := p - indPtr[i] - lc
		if toLeft[p+1] > toLeft[p] {
			p = indPtr[i<<1] + lc
			i <<= 1
		} else {
			p = indPtr[i<<1|1] + rc
			i = i<<1 | 1
		}
	}
}

func (t *SegmentTree2DSparse32Fast) Query(lx, rx, ly, ry int32) E {
	L := t._xtoi(lx)
	R := t._xtoi(rx)
	res := t.unit
	indPtr, toLeft := t.indptr, t.toLeft
	var dfs func(i, l, r, a, b int32)
	dfs = func(i, l, r, a, b int32) {
		if a == b || R <= l || r <= L {
			return
		}
		if L <= l && r <= R {
			res = op(res, t._query(i, a, b))
			return
		}
		la := toLeft[indPtr[i]+a] - toLeft[indPtr[i]]
		ra := a - la
		lb := toLeft[indPtr[i]+b] - toLeft[indPtr[i]]
		rb := b - lb
		m := (l + r) >> 1
		dfs(i<<1, l, m, la, lb)
		dfs(i<<1|1, m, r, ra, rb)
	}
	dfs(1, 0, t.size, bisectLeft(t.allY, ly, 0, int32(len(t.allY)-1)), bisectLeft(t.allY, ry, 0, int32(len(t.allY)-1)))
	return res
}

func (seg *SegmentTree2DSparse32Fast) Count(lx, rx, ly, ry int32) int32 {
	L := seg._xtoi(lx)
	R := seg._xtoi(rx)
	res := int32(0)
	indPtr, toLeft := seg.indptr, seg.toLeft
	var dfs func(i, l, r, a, b int32)
	dfs = func(i, l, r, a, b int32) {
		if a == b || R <= l || r <= L {
			return
		}
		if L <= l && r <= R {
			res += b - a
			return
		}
		la := toLeft[indPtr[i]+a] - toLeft[indPtr[i]]
		ra := a - la
		lb := toLeft[indPtr[i]+b] - toLeft[indPtr[i]]
		rb := b - lb
		m := (l + r) >> 1
		dfs(i<<1, l, m, la, lb)
		dfs(i<<1|1, m, r, ra, rb)
	}
	dfs(1, 0, seg.size, bisectLeft(seg.allY, ly, 0, int32(len(seg.allY)-1)), bisectLeft(seg.allY, ry, 0, int32(len(seg.allY)-1)))
	return res
}

func (t *SegmentTree2DSparse32Fast) _update(i int32, y int32, val E) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	offset := lid << 1
	y += n
	for y > 0 {
		t.data[offset+y] = op(t.data[offset+y], val)
		y >>= 1
	}
}

func (seg *SegmentTree2DSparse32Fast) _set(i, y int32, val E) {
	lid := seg.indptr[i]
	n := seg.indptr[i+1] - seg.indptr[i]
	off := lid << 1
	y += n
	seg.data[off+y] = val
	for y > 1 {
		y >>= 1
		seg.data[off+y] = op(seg.data[off+y<<1], seg.data[off+y<<1|1])
	}
}

func (t *SegmentTree2DSparse32Fast) _query(i int32, ly, ry int32) E {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	offset := lid << 1
	left, right := n+ly, n+ry
	val := t.unit
	for left < right {
		if left&1 == 1 {
			val = op(val, t.data[offset+left])
			left++
		}
		if right&1 == 1 {
			right--
			val = op(t.data[offset+right], val)
		}
		left >>= 1
		right >>= 1
	}
	return val
}

func (seg *SegmentTree2DSparse32Fast) _build(X, Y []int32, wt []E) {
	if len(X) != len(Y) || len(X) != len(wt) {
		panic("Lengths of X, Y, and wt must be equal.")
	}

	if seg.discretize {
		seg.keyX = unique(X)
		seg.n = int32(len(seg.keyX))
	} else {
		if len(X) > 0 {
			min_, max_ := int32(0), int32(0)
			for _, x := range X {
				if x < min_ {
					min_ = x
				}
				if x > max_ {
					max_ = x
				}
			}
			seg.minX = min_
			seg.n = max_ - min_ + 1
		}
	}

	log := int32(0)
	for 1<<log < seg.n {
		log++
	}
	size := int32(1 << log)
	seg.size = size

	orderX := make([]int32, len(X))
	for i := range orderX {
		orderX[i] = seg._xtoi(X[i])
	}
	seg.indptr = make([]int32, 2*size+1)
	for _, i := range orderX {
		i += size
		for i > 0 {
			seg.indptr[i+1]++
			i >>= 1
		}
	}
	for i := int32(1); i <= 2*size; i++ {
		seg.indptr[i] += seg.indptr[i-1]
	}
	seg.data = make([]E, 2*seg.indptr[2*size])
	for i := range seg.data {
		seg.data[i] = seg.unit
	}

	seg.toLeft = make([]int32, seg.indptr[size]+1)
	ptr := append([]int32(nil), seg.indptr...)
	order := argSort(Y)
	seg.pos = make([]int32, len(X))
	for i, v := range order {
		seg.pos[v] = int32(i)
	}
	for _, rawIdx := range order {
		i := orderX[rawIdx] + size
		j := int32(-1)
		for i > 0 {
			p := ptr[i]
			ptr[i]++
			seg.data[seg.indptr[i+1]+p] = wt[rawIdx]
			if j != -1 && j&1 == 0 {
				seg.toLeft[p+1] = 1
			}
			j = i
			i >>= 1
		}
	}
	for i := int32(1); i < int32(len(seg.toLeft)); i++ {
		seg.toLeft[i] += seg.toLeft[i-1]
	}

	for i := int32(0); i < 2*size; i++ {
		off := 2 * seg.indptr[i]
		n := seg.indptr[i+1] - seg.indptr[i]
		for j := n - 1; j >= 1; j-- {
			seg.data[off+j] = op(seg.data[off+j<<1], seg.data[off+j<<1|1])
		}
	}

	allY := append([]int32(nil), Y...)
	sort.Slice(allY, func(i, j int) bool { return allY[i] < allY[j] })
	seg.allY = allY
}

func (seg *SegmentTree2DSparse32Fast) _xtoi(x int32) int32 {
	if seg.discretize {
		return bisectLeft(seg.keyX, x, 0, int32(len(seg.keyX)-1))
	}
	tmp := x - seg.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > seg.n {
		tmp = seg.n
	}
	return tmp
}

func bisectLeft(nums []int32, x int32, left, right int32) int32 {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func unique(nums []int32) (sorted []int32) {
	set := make(map[int32]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int32, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return
}

func argSort(nums []int32) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

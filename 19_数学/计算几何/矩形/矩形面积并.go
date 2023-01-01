package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	rectangles := make([][]int, n)
	for i := 0; i < n; i++ {
		var x1, y1, x2, y2 int
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		rectangles[i] = []int{x1, y1, x2, y2}
	}

	fmt.Fprintln(out, rectangleArea(rectangles))
}

// 求矩形的面积并
//  rectangle[i] = [x1, y1, x2, y2]
//  0<=x1<x2<=10^9
//  0<=y1<y2<=10^9
//  1<=rectangle.length<=1e5
// https://leetcode.cn/problems/rectangle-area-ii/solution/ju-xing-mian-ji-ii-by-leetcode-solution-ulqz/
func rectangleArea(rectangles [][]int) (res int) {
	n := len(rectangles) * 2
	heights := make([]int, 0, n)
	for _, r := range rectangles {
		heights = append(heights, r[1], r[3])
	}

	// 排序+去重
	sort.Ints(heights)
	m := 0
	for _, h := range heights[1:] {
		if heights[m] != h {
			m++
			heights[m] = h
		}
	}
	heights = heights[:m+1]

	tree := make(data, m*4)
	tree.build(heights, 1, 1, m)

	type event struct{ x, i, d int }
	events := make([]event, 0, n)
	for i, r := range rectangles {
		events = append(events, event{r[0], i, 1}, event{r[2], i, -1})
	}
	sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })

	for i := 0; i < n; i++ {
		j := i
		for j+1 < n && events[j+1].x == events[i].x {
			j++
		}
		if j+1 == n {
			break
		}
		// 一次性地处理掉一批横坐标相同的左右边界
		for k := i; k <= j; k++ {
			index, diff := events[k].i, events[k].d
			// 使用二分查找得到完整覆盖的线段的编号范围
			left := sort.SearchInts(heights, rectangles[index][1]) + 1
			right := sort.SearchInts(heights, rectangles[index][3])
			tree.update(1, 1, m, left, right, diff)
		}
		res += tree[1].len * (events[j+1].x - events[j].x)
		i = j
	}
	return res
}

type data []struct{ cover, len, maxLen int }

func (t data) build(heights []int, idx, l, r int) {
	if l == r {
		t[idx].maxLen = heights[l] - heights[l-1]
		return
	}
	mid := (l + r) / 2
	t.build(heights, idx*2, l, mid)
	t.build(heights, idx*2+1, mid+1, r)
	t[idx].maxLen = t[idx*2].maxLen + t[idx*2+1].maxLen
}

func (t data) update(idx, l, r, ul, ur, diff int) {
	if l > ur || r < ul {
		return
	}
	if ul <= l && r <= ur {
		t[idx].cover += diff
		t.pushUp(idx, l, r)
		return
	}
	mid := (l + r) / 2
	t.update(idx*2, l, mid, ul, ur, diff)
	t.update(idx*2+1, mid+1, r, ul, ur, diff)
	t.pushUp(idx, l, r)
}

func (t data) pushUp(idx, l, r int) {
	if t[idx].cover > 0 {
		t[idx].len = t[idx].maxLen
	} else if l == r {
		t[idx].len = 0
	} else {
		t[idx].len = t[idx*2].len + t[idx*2+1].len
	}
}

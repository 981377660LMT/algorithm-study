package main

import "sort"

func main() {

}

// 判断 6 个矩形是否为长方体的 6 个面
// https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/misc.go#L735
// NEERC04 https://www.luogu.com.cn/problem/UVA1587
func isCuboid(rectangles [][2]int) bool {
	for i, r := range rectangles {
		if r[0] > r[1] {
			rectangles[i] = [2]int{r[1], r[0]} // swap
		}
	}

	sort.Slice(rectangles, func(i, j int) bool {
		a, b := rectangles[i], rectangles[j]
		return a[0] < b[0] || a[0] == b[0] && a[1] < b[1]
	})

	for i := 0; i < 6; i += 2 {
		if rectangles[i] != rectangles[i+1] { // !相对面
			return false
		}
	}

	minH, maxH := rectangles[0][1], rectangles[2][1]
	return rectangles[2][0] == rectangles[0][0] && (rectangles[4] == [2]int{minH, maxH} || rectangles[4] == [2]int{maxH, minH})
}

// abc365G - AtCoderOffice(在线查询区间交集长度，区间相交)
// !交集考虑根号分治
// https://atcoder.jp/contests/abc365/tasks/abc365_g

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	n, m := int32(NextInt()), int32(NextInt())
	times, personIds := make([]int32, m), make([]int32, m)
	inTime := make([]int32, n)
	for i := range inTime {
		inTime[i] = -1
	}
	intervals := make([][][2]int32, n)
	for i := int32(0); i < m; i++ {
		t, p := int32(NextInt()), int32(NextInt())
		p--
		times[i], personIds[i] = t, p
		if inTime[p] == -1 {
			inTime[p] = t
		} else {
			intervals[p] = append(intervals[p], [2]int32{inTime[p], t})
			inTime[p] = -1
		}
	}
	for _, vs := range intervals {
		sort.Slice(vs, func(i, j int) bool { return vs[i][0] < vs[j][0] })
	}

	const threshold int32 = 1000 // sqrt(nm/q)
	var big []int32
	isBig := func(id int32) bool { return int32(len(intervals[id])) > threshold }
	for i := int32(0); i < n; i++ {
		if isBig(i) {
			big = append(big, i)
		}
	}
	bigRes := make(map[int32][]int32, len(big))
	for i := 0; i < len(big); i++ {
		bigRes[big[i]] = make([]int32, n)
	}

	// !扫描线+前缀和处理所有区间与大区间的交集范围
	type event = struct {
		id   int32 // in: id, out: -id-1
		time int32
	}
	events := make([]event, 0, m)
	for i := int32(0); i < n; i++ {
		vs := intervals[i]
		for j := 0; j < len(vs); j++ {
			events = append(events, event{id: i, time: vs[j][0]})
			events = append(events, event{id: ^i, time: vs[j][1]})
		}
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].time != events[j].time {
			return events[i].time < events[j].time
		}
		return events[i].id > events[j].id // in first
	})
	dp := make([]int32, n)
	for _, bigId := range big {
		for i := range dp {
			dp[i] = 0
		}
		pre, curSum, inside := int32(0), int32(0), int32(0)
		for _, e := range events {
			id, time := e.id, e.time
			if id < 0 {
				id = ^id
			}
			if id == bigId {
				curSum += inside * (time - pre)
				pre = time
				inside ^= 1
			}
			dp[id] = curSum + inside*(time-pre) - dp[id] // !更新区间交集长度
		}
		for i := int32(0); i < n; i++ {
			bigRes[bigId][i] = dp[i]
		}
	}

	memo := make(map[int]int32)
	const mask32 int = 1<<32 - 1
	query := func(a, b int32) int32 {
		if a > b {
			a, b = b, a
		}
		hash := int(a)<<32 | int(b)
		if res, has := memo[hash]; has {
			return res
		}
		var res int32
		if !isBig(a) && !isBig(b) {
			res = IntervalsIntersection(intervals[a], intervals[b])
		} else if isBig(a) {
			res = bigRes[a][b]
		} else {
			res = bigRes[b][a]
		}
		memo[hash] = res
		return res
	}

	q := int32(NextInt())
	for i := int32(0); i < q; i++ {
		a, b := int32(NextInt()), int32(NextInt())
		a--
		b--
		fmt.Fprintln(out, query(a, b))
	}
}

// 有序区间交集长度.
func IntervalsIntersection(intervals1, intervals2 [][2]int32) int32 {
	n1, n2 := len(intervals1), len(intervals2)
	res := int32(0)
	left, right := 0, 0
	for left < n1 && right < n2 {
		s1, e1, s2, e2 := intervals1[left][0], intervals1[left][1], intervals2[right][0], intervals2[right][1]
		if (s1 <= e2 && e2 <= e1) || (s2 <= e1 && e1 <= e2) {
			res += min32(e1, e2) - max32(s1, s2)
		}
		if e1 < e2 {
			left++
		} else {
			right++
		}
	}
	return res
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

func min32(a, b int32) int32 {
	if a < b {
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

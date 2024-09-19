// 接口与bfs相同.
// 分为求最短路的双向bfs、带路径还原的双向bfs.
//
// 如果状态不是可哈希的，需要实现encode和decode函数.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc361_d()
	// P1379()
}

// D - Go Stone Puzzle
// https://atcoder.jp/contests/abc361/tasks/abc361_d
// n+2个格子，其中前 n个格子有石头，石头有黑有白。每次操作。
// 将相邻两个石头移动到无石头的位置，俩石头相对顺序不变。
// 给定初始局面和最终局面，问操作次数的最小值。
// n<=14
func abc361_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	var S, T string
	fmt.Fscan(in, &N)
	fmt.Fscan(in, &S, &T)

	// 0: 空格子, 1: 黑, 2: 白.
	toBytes := func(s string) []byte {
		res := make([]byte, len(s)+2)
		for i := 0; i < len(s); i++ {
			if s[i] == 'B' {
				res[i] = 1
			} else if s[i] == 'W' {
				res[i] = 2
			}
		}
		return res
	}
	start := toBytes(S)
	target := toBytes(T)

	encode := func(nums []byte) int {
		state := 0
		for _, v := range nums {
			state = state<<2 + int(v)
		}
		return state
	}
	decode := func(state int) []byte {
		bytes := make([]byte, N+2)
		for i := N + 1; i >= 0; i-- {
			bytes[i] = byte(state & 3)
			state >>= 2
		}
		return bytes
	}

	getNextStates := func(state int) (nexts []int) {
		nums := decode(state)
		emptyIndex := -1
		for i := 0; i < N+2; i++ {
			if nums[i] == 0 {
				emptyIndex = i
				break
			}
		}
		for i := 0; i < N+1; i++ {
			if nums[i] != 0 && nums[i+1] != 0 {
				nums[i], nums[emptyIndex] = nums[emptyIndex], nums[i]
				nums[i+1], nums[emptyIndex+1] = nums[emptyIndex+1], nums[i+1]
				nexts = append(nexts, encode(nums))
				nums[i], nums[emptyIndex] = nums[emptyIndex], nums[i]
				nums[i+1], nums[emptyIndex+1] = nums[emptyIndex+1], nums[i+1]
			}
		}
		return
	}

	res := BiBfs(encode(start), encode(target), getNextStates)
	fmt.Fprintln(out, res)
}

// P1379 八数码难题
// https://www.luogu.com.cn/problem/P1379
// 目标是 123804765
func P1379() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	encode := func(nums []byte) int {
		state := 0
		for _, v := range nums {
			state = state*16 + int(v-'0')
		}
		return state
	}
	decode := func(state int) []byte {
		nums := make([]byte, 9)
		for i := 8; i >= 0; i-- {
			nums[i] = byte(state&15) + '0'
			state >>= 4
		}
		return nums
	}

	nums := []byte(s)
	start, target := encode(nums), encode([]byte("123804765"))
	getNextStates := func(state int) (nexts []int) {
		nums := decode(state)
		pos := 0
		for i, v := range nums {
			if v == '0' {
				pos = i
				break
			}
		}
		dirs := []int{-3, 3, -1, 1}
		for _, d := range dirs {
			if d == -1 && pos%3 == 0 || d == 1 && pos%3 == 2 {
				continue
			}
			if pos+d < 0 || pos+d >= 9 {
				continue
			}
			nums[pos], nums[pos+d] = nums[pos+d], nums[pos]
			nexts = append(nexts, encode(nums))
			nums[pos], nums[pos+d] = nums[pos+d], nums[pos]
		}
		return
	}

	res := BiBfs(start, target, getNextStates)
	fmt.Fprintln(out, res)
}

// 127. 单词接龙
// https://leetcode.cn/problems/word-ladder/description/
func ladderLength(beginWord string, endWord string, wordList []string) int {
	wordSet := map[string]struct{}{}
	for _, word := range wordList {
		wordSet[word] = struct{}{}
	}
	if _, has := wordSet[endWord]; !has {
		return 0
	}

	getNextStates := func(cur string) (nexts []string) {
		sb := []byte(cur)
		for i := 0; i < len(sb); i++ {
			c := sb[i]
			for j := 0; j < 26; j++ {
				sb[i] = byte('a' + j)
				next := string(sb)
				if _, has := wordSet[next]; has {
					nexts = append(nexts, next)
				}
			}
			sb[i] = c
		}
		return
	}

	// path := BiBfsPath(beginWord, endWord, getNextStates)
	// return len(path)
	res := BiBfs(beginWord, endWord, getNextStates)
	if res == -1 {
		return 0
	}
	return int(res) + 1
}

const INF32 int32 = 1e9 + 10

// 双向bfs求最短路.
// 如果不存在，返回-1.
func BiBfs[S comparable](start, target S, getNextStates func(cur S) (nexts []S)) int32 {
	queue := [2]map[S]struct{}{make(map[S]struct{}), make(map[S]struct{})}
	queue[0][start], queue[1][target] = struct{}{}, struct{}{}
	visited := map[S]struct{}{}
	steps := int32(0)
	for len(queue[0]) > 0 && len(queue[1]) > 0 {
		qi := int32(0)
		if len(queue[0]) > len(queue[1]) {
			qi = 1
		}
		nextQueue := map[S]struct{}{}
		for cur := range queue[qi] {
			if _, has := queue[qi^1][cur]; has {
				return steps
			}
			if _, has := visited[cur]; has {
				continue
			}
			visited[cur] = struct{}{}
			for _, next := range getNextStates(cur) {
				nextQueue[next] = struct{}{}
			}
		}
		steps++
		queue[qi] = nextQueue
	}
	return -1
}

func BiBfsPath[S comparable](start, target S, getNextStates func(cur S) (nexts []S)) []S {
	queue := [2]map[S]struct{}{make(map[S]struct{}), make(map[S]struct{})}
	queue[0][start], queue[1][target] = struct{}{}, struct{}{}
	pre := [2]map[S]S{make(map[S]S), make(map[S]S)}
	visited := [2]map[S]struct{}{make(map[S]struct{}), make(map[S]struct{})}
	visited[0][start], visited[1][target] = struct{}{}, struct{}{}

	// start -> mid -> target.
	restorePath := func(mid S) []S {
		pre1, path1, cur1 := pre[0], []S{mid}, mid
		for {
			p, has := pre1[cur1]
			if !has {
				break
			}
			cur1 = p
			path1 = append(path1, cur1)
		}
		pre2, path2, cur2 := pre[1], []S{}, mid
		for {
			p, has := pre2[cur2]
			if !has {
				break
			}
			cur2 = p
			path2 = append(path2, cur2)
		}
		for i, j := 0, len(path1)-1; i < j; i, j = i+1, j-1 {
			path1[i], path1[j] = path1[j], path1[i]
		}
		return append(path1, path2...)
	}

	for len(queue[0]) > 0 && len(queue[1]) > 0 {
		qi := int32(0)
		if len(queue[0]) > len(queue[1]) {
			qi = 1
		}

		nextQueue := map[S]struct{}{}
		curVisited, curPre, otherQueue := visited[qi], pre[qi], queue[qi^1]
		for cur := range queue[qi] {
			if _, has := otherQueue[cur]; has {
				return restorePath(cur)
			}
			for _, next := range getNextStates(cur) {
				if _, has := curVisited[next]; has {
					continue
				}
				curVisited[next] = struct{}{}
				nextQueue[next] = struct{}{}
				curPre[next] = cur
			}
		}
		queue[qi] = nextQueue
	}

	return nil
}

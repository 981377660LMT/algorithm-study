// 接口与bfs相同.
// 分为求最短路的双向bfs、带路径还原的双向bfs.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P1379()
}

// P1379 八数码难题
// https://www.luogu.com.cn/problem/P1379
func P1379() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
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

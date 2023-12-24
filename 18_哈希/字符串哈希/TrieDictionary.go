// !TrieHash(TrieDictionary)：用字典树把每个字符串转换成一个整数编号.
// 实际案例：React 的 cache 就是基于字典树哈希实现的.
//
// api:
//  Add(seq) -> (id, pos)
//  Has(pos) -> bool
//  GetId(pos) -> id
//  GetValue(id) -> seq
//  Size() -> int
//  Move(pos, ord) -> int

package main

import "fmt"

const INF int = 1e18

func main() {
	fmt.Println(minimumCost("abcd", "acbe", []string{"a", "b", "c", "c", "e", "d"}, []string{"b", "c", "b", "e", "b", "e"}, []int{2, 5, 5, 1, 2, 20}))
}

// 100158. 转换字符串的最小成本 II
// https://leetcode.cn/problems/minimum-cost-to-convert-string-ii/
func minimumCost(source string, target string, original []string, changed []string, cost []int) int64 {
	D := NewDictionaryTrie()
	originalId := make([]int, len(original))
	for i, s := range original {
		id, _ := D.Add(s)
		originalId[i] = id
	}
	changedId := make([]int, len(changed))
	for i, s := range changed {
		id, _ := D.Add(s)
		changedId[i] = id
	}

	// floyd
	dist := make([][]int, D.Size())
	for i := range dist {
		dist[i] = make([]int, D.Size())
		for j := range dist[i] {
			dist[i][j] = INF
		}
		dist[i][i] = 0
	}
	for i := range original {
		u, v, w := originalId[i], changedId[i], cost[i]
		dist[u][v] = min(dist[u][v], w)
	}
	for k := range dist {
		for i := range dist {
			for j := range dist {
				dist[i][j] = min(dist[i][j], dist[i][k]+dist[k][j])
			}
		}
	}

	// dp
	n := len(target)
	memo := make([]int, n+1)
	for i := range memo {
		memo[i] = -1
	}
	var dfs func(int) int
	dfs = func(index int) int {
		if index == n {
			return 0
		}
		if memo[index] != -1 {
			return memo[index]
		}
		res := INF
		if source[index] == target[index] {
			res = dfs(index + 1)
		}
		pos1, pos2 := 0, 0 // 在trie中的位置
		for j := index; j < n; j++ {
			pos1 = D.Move(pos1, int(source[j]))
			pos2 = D.Move(pos2, int(target[j]))
			if pos1 == -1 || pos2 == -1 {
				break
			}
			id1, id2 := D.GetId(pos1), D.GetId(pos2)
			if id1 != -1 && id2 != -1 {
				res = min(res, dfs(j+1)+dist[id1][id2])
			}
		}

		memo[index] = res
		return res
	}

	res := dfs(0)
	if res == INF {
		return -1
	}
	return int64(res)
}

const SIGMA int = 26
const OFFSET int = 'a'

type Sequence = string

type DictionaryTrie struct {
	Children  [][SIGMA]int32
	PosToId   map[int32]int32
	IdToValue []Sequence
}

// TrieHash.
// A trie that maps sequence to unique ids.
func NewDictionaryTrie() *DictionaryTrie {
	res := &DictionaryTrie{PosToId: make(map[int32]int32)}
	res.newNode()
	return res
}

// 添加字符串，返回对应的(id, pos).
func (trie *DictionaryTrie) Add(seq Sequence) (id, pos int) {
	node := int32(0)
	for _, s := range seq {
		ord := int(s) - OFFSET
		if trie.Children[node][ord] == -1 {
			trie.Children[node][ord] = trie.newNode()
		}
		node = trie.Children[node][ord]
	}
	if v, ok := trie.PosToId[node]; ok {
		return int(v), int(node)
	} else {
		len_ := len(trie.PosToId)
		trie.IdToValue = append(trie.IdToValue, seq)
		trie.PosToId[node] = int32(len_)
		return len_, int(node)
	}
}

// 判断pos对应的字符串是否存在.
func (trie *DictionaryTrie) Has(pos int) (ok bool) {
	if pos < 0 || pos >= len(trie.Children) {
		return
	}
	_, ok = trie.PosToId[int32(pos)]
	return
}

// 返回pos对应的字符串的id.如果不存在，返回-1.
func (trie *DictionaryTrie) GetId(pos int) int {
	if v, ok := trie.PosToId[int32(pos)]; ok {
		return int(v)
	} else {
		return -1
	}
}

func (trie *DictionaryTrie) GetValue(id int) Sequence {
	return trie.IdToValue[id]
}

func (trie *DictionaryTrie) Size() int {
	return len(trie.PosToId)
}

func (trie *DictionaryTrie) Move(pos int, ord int) int {
	ord -= OFFSET
	return int(trie.Children[pos][ord])
}

func (trie *DictionaryTrie) newNode() int32 {
	nexts := [SIGMA]int32{}
	for i := range nexts {
		nexts[i] = -1
	}
	trie.Children = append(trie.Children, nexts)
	return int32(len(trie.Children) - 1)
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

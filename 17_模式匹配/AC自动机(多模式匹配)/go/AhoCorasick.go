// https://nyaannyaan.github.io/library/string/aho-corasick.hpp

// #pragma once

// #include "trie.hpp"

// template <size_t X = 26, char margin = 'a'>
// struct AhoCorasick : Trie<X + 1, margin> {
//   using TRIE = Trie<X + 1, margin>;
//   using TRIE::next;
//   using TRIE::st;
//   using TRIE::TRIE;
//   vector<int> cnt;

//   void build(int heavy = true) {
//     int n = st.size();
//     cnt.resize(n);
//     for (int i = 0; i < n; i++) {
//       if (heavy) sort(st[i].idxs.begin(), st[i].idxs.end());
//       cnt[i] = st[i].idxs.size();
//     }

//     queue<int> que;
//     for (int i = 0; i < (int)X; i++) {
//       if (~next(0, i)) {
//         next(next(0, i), X) = 0;
//         que.emplace(next(0, i));
//       } else {
//         next(0, i) = 0;
//       }
//     }

//     while (!que.empty()) {
//       auto &x = st[que.front()];
//       int fail = x.nxt[X];

//       cnt[que.front()] += cnt[fail];
//       que.pop();

//       for (int i = 0; i < (int)X; i++) {
//         int &nx = x.nxt[i];
//         if (nx < 0) {
//           nx = next(fail, i);
//           continue;
//         }
//         que.emplace(nx);
//         next(nx, X) = next(fail, i);
//         if (heavy) {
//           auto &idx = st[nx].idxs;
//           auto &idy = st[next(fail, i)].idxs;
//           vector<int> idz;
//           set_union(idx.begin(), idx.end(), idy.begin(), idy.end(),
//                     back_inserter(idz));
//           idx = idz;
//         }
//       }
//     }
//   }

//   vector<int> match(string s, int heavy = true) {
//     vector<int> res(heavy ? TRIE::size() : 1);
//     int pos = 0;
//     for (auto &c : s) {
//       pos = next(pos, c - margin);
//       if (heavy)
//         for (auto &x : st[pos].idxs) res[x]++;
//       else
//         res[0] += cnt[pos];
//     }
//     return res;
//   }

//   int count(int pos) { return cnt[pos]; }
// };

// API
// Add(s, id)
// Build(heavy)
// Match(s, heavy)
// MatchFrom(s, pos, heavy)
// Move(pos, char)
// Next(pos, j)
// Count(pos)
// Size()
package main

import (
	"fmt"
	"sort"
)

func main() {
	aho := NewAhoCorasick(26, 'a')
	aho.Add("abc", 1)
	aho.Add("abc", 1)
	aho.Add("bcd", 2)
	aho.Build(true)
	fmt.Println(aho.Match("abcabc", true))
	fmt.Println(aho.Count(3))
	fmt.Println(aho.Move(0, 'g'))
}

type AhoCorasick struct {
	*Trie
	count []int
}

// size: 字符集大小
// margin: 字符集起始字符
func NewAhoCorasick(size int, margin byte) *AhoCorasick {
	res := &AhoCorasick{Trie: NewTrie(size+1, margin)}
	return res
}

func (ac *AhoCorasick) Build(heavy bool) {
	n := len(ac.stack)
	ac.count = make([]int, n)
	for i := 0; i < n; i++ {
		if heavy {
			sort.Ints(ac.stack[i].indexes)
		}
		ac.count[i] = len(ac.stack[i].indexes)
	}

	var que []int
	for i := 0; i < ac.size-1; i++ {
		if *ac.Next(0, i) != -1 {
			*ac.Next(*ac.Next(0, i), ac.size-1) = 0
			que = append(que, *ac.Next(0, i))
		} else {
			*ac.Next(0, i) = 0
		}
	}

	for len(que) > 0 {
		x := ac.stack[que[0]]
		fail := x.next[ac.size-1]
		ac.count[que[0]] += ac.count[fail]
		que = que[1:]

		for i := 0; i < ac.size-1; i++ {
			nx := &x.next[i]
			if *nx < 0 {
				*nx = *ac.Next(fail, i)
				continue
			}
			que = append(que, *nx)
			*ac.Next(*nx, ac.size-1) = *ac.Next(fail, i)
			if heavy {
				idx := ac.stack[*nx].indexes
				idy := ac.stack[*ac.Next(fail, i)].indexes
				idz := make([]int, 0, len(idx)+len(idy))

				// set union
				i, j := 0, 0
				for i < len(idx) && j < len(idy) {
					if idx[i] < idy[j] {
						idz = append(idz, idx[i])
						i++
					} else if idx[i] > idy[j] {
						idz = append(idz, idy[j])
						j++
					} else {
						idz = append(idz, idx[i])
						i++
						j++
					}
				}
				for i < len(idx) {
					idz = append(idz, idx[i])
					i++
				}
				for j < len(idy) {
					idz = append(idz, idy[j])
					j++
				}

				ac.stack[*nx].indexes = idz
			}
		}

	}
}

func (ac *AhoCorasick) Match(s string, heavy bool) []int {
	size := 1
	if heavy {
		size = ac.Size()
	}
	res := make([]int, size)
	pos := 0
	for i := 0; i < len(s); i++ {
		pos = *ac.Next(pos, int(s[i]-ac.margin))
		if heavy {
			for _, x := range ac.stack[pos].indexes {
				res[x]++
			}
		} else {
			res[0] += ac.count[pos]
		}
	}
	return res
}

func (ac *AhoCorasick) MatchFrom(s string, pos int, heavy bool) []int {
	size := 1
	if heavy {
		size = ac.Size()
	}
	res := make([]int, size)
	for i := 0; i < len(s); i++ {
		pos = *ac.Next(pos, int(s[i]-ac.margin))
		if heavy {
			for _, x := range ac.stack[pos].indexes {
				res[x]++
			}
		} else {
			res[0] += ac.count[pos]
		}
	}
	return res
}

func (ac *AhoCorasick) Count(pos int) int {
	return ac.count[pos]
}

//
//
//
//
type Trie struct {
	size   int
	margin byte
	stack  []*trieNode
}

type trieNode struct {
	index   int // 最后一次被更新的字符串的索引
	key     byte
	indexes []int // 存储了哪些字符串的索引
	next    []int // children position
}

// size: 字符集大小
// margin: 字符集起始字符
func NewTrie(size int, margin byte) *Trie {
	res := &Trie{size: size, margin: margin}
	root := res.newNode('$')
	res.stack = append(res.stack, root)
	return res
}

func (t *Trie) Add(s string, index int) {
	pos := 0
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.Next(pos, k) != -1 {
			pos = *t.Next(pos, k)
			continue
		}
		nextPos := len(t.stack)
		*t.Next(pos, k) = nextPos
		node := t.newNode(s[i])
		t.stack = append(t.stack, node)
		pos = nextPos
	}
	t.stack[pos].index = index
	t.stack[pos].indexes = append(t.stack[pos].indexes, index)
}

func (t *Trie) Find(s string) int {
	pos := 0
	for i := 0; i < len(s); i++ {
		k := int(s[i] - t.margin)
		if *t.Next(pos, k) == -1 {
			return -1
		}
		pos = *t.Next(pos, k)
	}
	return pos
}

func (t *Trie) Move(pos int, c byte) int {
	if pos < 0 || pos >= len(t.stack) {
		return -1
	}
	return *t.Next(pos, int(c-t.margin))
}

func (t *Trie) Index(pos int) int {
	if pos < 0 {
		return -1
	}
	return t.stack[pos].index
}

func (t *Trie) IndexAll(pos int) []int {
	if pos < 0 {
		return []int{}
	}
	return t.stack[pos].indexes
}

func (t *Trie) Size() int { return len(t.stack) }

func (t *Trie) Next(pos, j int) *int {
	return &t.stack[pos].next[j]
}

func (t *Trie) newNode(c byte) *trieNode {
	next := make([]int, t.size)
	for i := range next {
		next[i] = -1
	}
	return &trieNode{
		index: -1,
		key:   c,
		next:  next,
	}
}

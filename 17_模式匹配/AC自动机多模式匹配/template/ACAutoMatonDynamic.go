// https://sotanishy.github.io/cp-library-cpp/string/aho_corasick.cpp
// 静的データ構造で動的に処理する (1) - えびちゃんの日記
// https://rsk0315.hatenablog.com/entry/2019/06/19/124528

// API
// !NewAhoCorasick() *AhoCorasick
// Insert(s string)
// Build()
// Move(state int,c byte) int  // 转移到下一个状态.
// Count(s string) int  // s中有多少个模式串.
// Match(s string) [][2]int // 每一项是(匹配到的模式串id, 匹配到的结尾的index).

// !NewDynamicAhoCorasick() *DynamicAhoCorasick
// Insert(s string)
// Count(s string) int  // s中有多少个模式串.

package main

import "fmt"

func main() {
	ac := NewAhoCorasick()
	ac.Insert("he", 1)
	ac.Insert("she", 2)
	ac.Insert("his", 3)
	ac.Build()
	fmt.Println(ac.Match("ushershis"))
	fmt.Println(ac.Count("ushershis"))

	s0 := 0
	s1 := ac.Move(s0, 's')
	fmt.Println(s1)
	s2 := ac.Move(s1, 'h')
	fmt.Println(s2)
	s3 := ac.Move(s2, 'e')
	s4 := ac.Move(s2, 'q')
	fmt.Println(s3, s4)

	dac := NewDynamicAhoCorasick()
	dac.Insert("he")
	dac.Insert("she")
	fmt.Println(dac.Count("heheshe")) // 4
}

// 面试题 17.17. 多次搜索
func multiSearch(big string, smalls []string) [][]int {
	ac := NewAhoCorasick()
	for i, s := range smalls {
		ac.Insert(s, i)
	}
	ac.Build()
	res := make([][]int, len(smalls))
	matched := ac.Match(big)
	for _, m := range matched {
		id, last := m[0], m[1]
		res[id] = append(res[id], last-len(smalls[id])+1)
	}
	return res

}

type ANode struct {
	ch     map[byte]int
	accept []int
	link   int
	cnt    int
}

func NewANode() *ANode {
	return &ANode{
		ch:   make(map[byte]int),
		link: -1,
	}
}

// AC自动机.
type AhoCorasick struct {
	states      []*ANode
	acceptState map[int]int
}

func NewAhoCorasick() *AhoCorasick {
	return &AhoCorasick{
		states:      []*ANode{NewANode()},
		acceptState: make(map[int]int),
	}
}

// 插入字符串s, id为该字符串的标识符.
//  id为-1时, 表示不需要标识符.
func (ac *AhoCorasick) Insert(s string, id int) {
	if len(s) == 0 {
		return
	}
	i := 0
	for j := 0; j < len(s); j++ {
		c := s[j]
		if _, ok := ac.states[i].ch[c]; !ok {
			ac.states[i].ch[c] = len(ac.states)
			ac.states = append(ac.states, NewANode())
		}
		i = ac.states[i].ch[c]
	}
	ac.states[i].cnt++
	ac.states[i].accept = append(ac.states[i].accept, id)
	ac.acceptState[id] = i
}

func (ac *AhoCorasick) Build() {
	que := []int{0}
	for len(que) > 0 {
		i := que[0]
		que = que[1:]
		for c, j := range ac.states[i].ch {
			ac.states[j].link = ac.Move(ac.states[i].link, c)
			ac.states[j].cnt += ac.states[ac.states[j].link].cnt
			a := ac.states[j].accept
			b := ac.states[ac.states[j].link].accept
			accept := make([]int, 0, len(a)+len(b))
			p1, p2 := 0, 0
			for p1 < len(a) && p2 < len(b) {
				if a[p1] < b[p2] {
					accept = append(accept, a[p1])
					p1++
				} else {
					accept = append(accept, b[p2])
					p2++
				}
			}
			accept = append(accept, a[p1:]...)
			accept = append(accept, b[p2:]...)
			ac.states[j].accept = accept
			que = append(que, j)
		}
	}

}

// 当前状态为state, 输入字符c, 转移到下一个状态.
//  初始/失败状态为0.
func (ac *AhoCorasick) Move(state int, c byte) int {
	for state != -1 {
		if _, ok := ac.states[state].ch[c]; !ok {
			state = ac.states[state].link
		} else {
			break
		}
	}
	if state == -1 {
		return 0
	}
	return ac.states[state].ch[c]
}

// 匹配字符串s, 返回匹配的位置.
//  (id, index) 数组.
func (ac *AhoCorasick) Match(s string) [][2]int {
	res := make([][2]int, 0)
	i := 0
	for k := 0; k < len(s); k++ {
		c := s[k]
		i = ac.Move(i, c)
		for _, id := range ac.states[i].accept {
			res = append(res, [2]int{id, k})
		}
	}
	return res
}

// 匹配字符串s, 返回匹配的个数.
func (ac *AhoCorasick) Count(s string) int {
	res := 0
	i := 0
	for k := 0; k < len(s); k++ {
		c := s[k]
		i = ac.Move(i, c)
		res += ac.states[i].cnt
	}
	return res
}

// clear
func (ac *AhoCorasick) Clear() {
	ac.states = ac.states[:0]
	ac.states = append(ac.states, NewANode())
}

// 动态AC自动机.
//
type DynamicAhoCorasick struct {
	ac   []*AhoCorasick
	dict [][]string
}

func NewDynamicAhoCorasick() *DynamicAhoCorasick {
	return &DynamicAhoCorasick{}
}

func (dac *DynamicAhoCorasick) Insert(s string) {
	if len(s) == 0 {
		return
	}
	k := 0
	for k < len(dac.dict) && len(dac.dict[k]) != 0 {
		k++
	}
	if k == len(dac.dict) {
		dac.dict = append(dac.dict, []string{})
		dac.ac = append(dac.ac, NewAhoCorasick())
	}

	dac.dict[k] = append(dac.dict[k], s)
	dac.ac[k].Insert(s, -1)
	for i := 0; i < k; i++ {
		for _, t := range dac.dict[i] {
			dac.ac[k].Insert(t, -1)
		}
		dac.dict[k] = append(dac.dict[k], dac.dict[i]...)
		dac.ac[i].Clear()
		dac.dict[i] = dac.dict[i][:0]
	}

	dac.ac[k].Build()
}

// 统计字符串s中出现的模式串的个数.
func (dac *DynamicAhoCorasick) Count(s string) int {
	res := 0
	for i := 0; i < len(dac.ac); i++ {
		res += dac.ac[i].Count(s)
	}
	return res
}

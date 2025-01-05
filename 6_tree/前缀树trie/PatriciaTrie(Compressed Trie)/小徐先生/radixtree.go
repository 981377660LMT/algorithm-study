// 这里实现较为简单.

package main

import "strings"

// https://leetcode.cn/problems/implement-trie-prefix-tree/
type Trie struct {
	radixTree *Radix
}

func Constructor() Trie {
	return Trie{
		radixTree: NewRadix(),
	}
}

func (this *Trie) Insert(word string) {
	this.radixTree.Insert(word)
}

func (this *Trie) Search(word string) bool {
	return this.radixTree.Search(word)
}

func (this *Trie) StartsWith(prefix string) bool {
	return this.radixTree.StartsWith(prefix)
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */

type Radix struct {
	root *radixNode
}

func NewRadix() *Radix {
	return &Radix{
		root: &radixNode{},
	}
}

func (r *Radix) Insert(word string) {
	if r.Search(word) {
		return
	}
	r.root.insert(word)
}

func (r *Radix) Search(word string) bool {
	node := r.root.search(word)
	return node != nil && node.fullPath == word && node.end
}

func (r *Radix) StartsWith(prefix string) bool {
	node := r.root.search(prefix)
	return node != nil && strings.HasPrefix(node.fullPath, prefix)
}

func (r *Radix) CountPre(prefix string) int {
	node := r.root.search(prefix)
	if node == nil || !strings.HasPrefix(node.fullPath, prefix) {
		return 0
	}
	return node.preCnt
}

// 删除调一个字典
func (r *Radix) Discard(word string) bool {
	if !r.Search(word) {
		return false
	}

	// root 直接精准命中了
	if r.root.fullPath == word {
		// 如果一个孩子都没有
		if len(r.root.indices) == 0 {
			r.root.path = ""
			r.root.fullPath = ""
			r.root.end = false
			r.root.preCnt = 0
			return true
		}

		// 如果只有一个孩子
		if len(r.root.indices) == 1 {
			child := r.root.children[0]
			child.path = r.root.path + child.path
			r.root = child
			return true
		}

		// 如果有多个孩子
		for i := 0; i < len(r.root.indices); i++ {
			r.root.children[i].path = r.root.path + r.root.children[0].path
		}

		newRoot := radixNode{
			indices:  r.root.indices,
			children: r.root.children,
			preCnt:   r.root.preCnt - 1,
		}
		r.root = &newRoot
		return true
	}

	// 确定 word 存在的情况下
	move := r.root
	// root 单独作为一个分支处理
	// 其他情况下，需要对孩子进行处理
walk:
	for {
		move.preCnt--
		prefix := move.path
		word = word[len(prefix):]
		c := word[0]
		for i := 0; i < len(move.indices); i++ {
			if move.indices[i] != c {
				continue
			}

			// 精准命中但是他仍有后继节点
			if move.children[i].path == word && move.children[i].preCnt > 1 {
				move.children[i].end = false
				move.children[i].preCnt--
				return true
			}

			// 找到对应的 child 了
			// 如果说后继节点的 preCnt = 1，直接干掉
			if move.children[i].preCnt > 1 {
				move = move.children[i]
				continue walk
			}

			// 删除move.children[i]
			move.children = append(move.children[:i], move.children[i+1:]...)
			move.indices = move.indices[:i] + move.indices[i+1:]
			// 如果干掉一个孩子后，发现只有一个孩子了，并且自身 end 为 false 则需要进行合并
			if !move.end && len(move.indices) == 1 {
				// 合并自己与唯一的孩子
				child := move.children[0]
				move.path += child.path
				move.fullPath = child.fullPath
				move.end = child.end
				move.indices = child.indices
				move.children = child.children
			}

			return true
		}
	}
}

type radixNode struct {
	path     string       // 当前节点的相对路径
	fullPath string       // 完整路径
	indices  string       // 每个 indice 字符对应一个孩子节点的 path 首字母
	children []*radixNode // 后继节点
	end      bool         // 是否有路径以当前节点为终点
	preCnt   int          // 记录有多少路径途径当前节点
}

// 传入相对路径和完整路径，补充一个新生成的节点信息
func (rn *radixNode) fill(path, fullPath string) {
	rn.path, rn.fullPath = path, fullPath
	rn.preCnt = 1
	rn.end = true
}

// 不断比较、消耗前缀，直到找到合适的位置插入新节点或者更新已有节点.
func (rn *radixNode) insert(word string) {
	fullWord := word

	// 如果当前节点为 root，此之前没有注册过子节点，则直接插入并返回
	if rn.path == "" && len(rn.children) == 0 {
		rn.fill(word, word)
		return
	}

walk:
	for {
		// 获取到 word 和当前节点 path 的公共前缀长度
		i := lcp(word, rn.path)
		// 只要公共前缀大于 0，则一定经过当前节点，需要累加 preCnt
		if i > 0 {
			rn.preCnt++
		}

		// 公共前缀小于当前节点的相对路径，需要对节点进行分解
		if i < len(rn.path) {
			// 需要进行节点切割
			child := radixNode{
				// 进行相对路径切分
				path: rn.path[i:],
				// 继承完整路径
				fullPath: rn.fullPath,
				// 当前节点的后继节点进行委托
				children: rn.children,
				indices:  rn.indices,
				end:      rn.end,
				// 传承给孩子节点时，需要把之前累加上的 preCnt 计数扣除
				preCnt: rn.preCnt - 1,
			}

			// 续接上孩子节点
			rn.indices = string(rn.path[i])
			rn.children = []*radixNode{&child}
			// 调整原节点的 full path
			rn.fullPath = rn.fullPath[:len(rn.fullPath)-(len(rn.path)-i)]
			// 调整原节点的 path
			rn.path = rn.path[:i]
			// 原节点是新拆分出来的，目前不可能有单词以该节点结尾
			rn.end = false
		}

		// 公共前缀小于插入 word 的长度
		if i < len(word) {
			// 对 word 扣除公共前缀部分
			word = word[i:]
			// 获取 word 剩余部分的首字母
			c := word[0]
			for i := 0; i < len(rn.indices); i++ {
				// 如果与后继节点还有公共前缀，则将 rn 指向子节点，然后递归执行流程
				if rn.indices[i] == c {
					rn = rn.children[i]
					continue walk
				}
			}

			// 到了这里，意味着 word 剩余部分与后继节点没有公共前缀了
			// 此时直接构造新的节点进行插入
			rn.indices += string(c)
			child := radixNode{}
			child.fill(word, fullWord)
			rn.children = append(rn.children, &child)
			return
		}

		// 倘若公共前缀恰好是 path，需要将 end 置为 true
		rn.end = true
		return
	}
}

func (rn *radixNode) search(word string) *radixNode {
walk:
	for {
		prefix := rn.path
		// word 长于 path
		if len(word) > len(prefix) {
			// 没匹配上，直接返回 nil
			if word[:len(prefix)] != prefix {
				return nil
			}
			// word 扣除公共前缀后的剩余部分
			word = word[len(prefix):]
			c := word[0]
			for i := 0; i < len(rn.indices); i++ {
				// 后继节点还有公共前缀，继续匹配
				if c == rn.indices[i] {
					rn = rn.children[i]
					continue walk
				}
			}
			// word 还有剩余部分，但是 prefix 不存在后继节点和 word 剩余部分有公共前缀了
			// 必然不存在
			return nil
		}

		// 和当前节点精准匹配上了
		if word == prefix {
			return rn
		}

		// !走到这里意味着 len(word) <= len(prefix) && word != prefix
		return rn
	}
}

func lcp(wordA, wordB string) int {
	var res int
	for res < len(wordA) && res < len(wordB) && wordA[res] == wordB[res] {
		res++
	}
	return res
}

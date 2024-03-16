// 745. 前缀和后缀搜索
// https://leetcode.cn/problems/prefix-and-suffix-search/description/

package main

const BASE uint = 131 // 自然溢出哈希
const THRESHOLD int32 = 20

type WordFilter struct {
	words                  []string          // 词典
	prefixHash, suffixHash []map[int32]uint  // 每个串每个前缀/后缀的哈希值
	shortHash              map[[2]uint]int32 // 长度<=20的前后缀哈希值对应的下标最大值.
	longWordIndex          []int32           // 长度>20的串的下标
}

// 使用词典中的单词 words 初始化对象。
func Constructor(words []string) WordFilter {
	n := int32(len(words))
	prefixHash := make([]map[int32]uint, n) // 每个串每个前缀的哈希值
	suffixHash := make([]map[int32]uint, n) // 每个串每个后缀的哈希值
	for i := int32(0); i < n; i++ {
		prefixHash[i] = make(map[int32]uint)
		h1 := uint(0)
		for j, c := range words[i] {
			h1 = (h1*BASE + uint(c))
			prefixHash[i][int32(j+1)] = h1
		}
		suffixHash[i] = make(map[int32]uint)
		h2 := uint(0)
		for j := len(words[i]) - 1; j >= 0; j-- {
			h2 = (h2*BASE + uint(words[i][j]))
			suffixHash[i][int32(len(words[i])-j)] = h2
		}
	}

	shortHash := make(map[[2]uint]int32) // 长度<=20的前后缀哈希值对应的下标最大值.
	for i, w := range words {
		for preLen := int32(1); preLen <= min32(THRESHOLD, int32(len(w))); preLen++ {
			for sufLen := int32(1); sufLen <= min32(THRESHOLD, int32(len(w))); sufLen++ {
				h1, h2 := prefixHash[i][preLen], suffixHash[i][sufLen]
				shortHash[[2]uint{h1, h2}] = int32(i)
			}
		}
	}

	longWordIndex := []int32{} // 长度>20的串的下标
	for i, w := range words {
		if int32(len(w)) > THRESHOLD {
			longWordIndex = append(longWordIndex, int32(i))
		}
	}

	return WordFilter{words: words, prefixHash: prefixHash, suffixHash: suffixHash, shortHash: shortHash, longWordIndex: longWordIndex}
}

// 返回词典中具有前缀 prefix 和后缀 suff 的单词的下标。
// 如果存在不止一个满足要求的下标，返回其中 最大的下标 。
// 如果不存在这样的单词，返回 -1 。
func (this *WordFilter) F(prefix string, suffix string) int {
	n1, n2 := int32(len(prefix)), int32(len(suffix))
	h1 := uint(0)
	for j := int32(0); j < n1; j++ {
		h1 = (h1*BASE + uint(prefix[j]))
	}
	h2 := uint(0)
	for j := n2 - 1; j >= 0; j-- {
		h2 = (h2*BASE + uint(suffix[j]))
	}

	if n1 <= THRESHOLD && n2 <= THRESHOLD {
		res, ok := this.shortHash[[2]uint{h1, h2}]
		if ok {
			return int(res)
		} else {
			return -1
		}
	} else {
		for i := len(this.longWordIndex) - 1; i >= 0; i-- {
			wi := this.longWordIndex[i]
			n3 := int32(len(this.words[wi]))
			if n1 <= n3 && n2 <= n3 && this.prefixHash[wi][n1] == h1 && this.suffixHash[wi][n2] == h2 {
				return int(wi)
			}
		}
		return -1
	}
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

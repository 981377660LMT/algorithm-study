// 720. 词典中最长的单词
// https://leetcode.cn/problems/longest-word-in-dictionary/description/
//
// 给出一个字符串数组 words 组成的一本英语词典。返回 words 中最长的一个单词，该单词是由 words 词典中其他单词逐步添加一个字母组成。
// 若其中有多个可行的答案，则返回答案中字典序最小的单词。若无答案，则返回空字符串。
// 请注意，单词应该从左到右构建，每个额外的字符都添加到前一个单词的结尾。

package main

type Trie struct {
	children [26]*Trie
	isEnd    bool
}

func (t *Trie) Insert(word string) {
	node := t
	for _, c := range word {
		c -= 'a'
		if node.children[c] == nil {
			node.children[c] = &Trie{}
		}
		node = node.children[c]
	}
	node.isEnd = true
}

// 要求插入的单词的前缀必须存在.
func (t *Trie) Search(word string) bool {
	node := t
	for _, c := range word {
		c -= 'a'
		if node.children[c] == nil || !node.children[c].isEnd {
			return false
		}
		node = node.children[c]
	}
	return true
}

func longestWord(words []string) (res string) {
	root := &Trie{}
	for _, word := range words {
		root.Insert(word)
	}
	for _, word := range words {
		if !root.Search(word) {
			continue
		}
		if len(word) > len(res) || (len(word) == len(res) && word < res) {
			res = word
		}
	}
	return
}

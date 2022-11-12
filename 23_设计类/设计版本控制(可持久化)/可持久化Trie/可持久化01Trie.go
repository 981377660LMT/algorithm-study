// https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/trie.go#L344
package main

import "runtime/debug"

// 注：由于用的是指针写法，必要时禁止 GC，能加速不少
func init() { debug.SetGCPercent(-1) }

type TrieNode struct {
	wordCount, prefixCount int
	chidlren               [26]*TrieNode
}

// 可持久化字典树
// 注意为了拷贝一份 trieNode，这里的接收器不是指针
// https://oi-wiki.org/ds/persistent-trie/
// usage:
//  git := make([]*TrieNode, maxVersion+1)  // restore all versions
//  git[0] = &TrieNode{}  // init version 0
//  newTrie = git[0].Insert([]byte("bcd"))  // version 1
//  git[1] = newTrie
func (o TrieNode) Insert(s []byte) *TrieNode {
	o.prefixCount++
	if len(s) == 0 {
		o.wordCount++
		return &o
	}

	ord_ := ord(s[0])
	if o.chidlren[ord_] == nil {
		o.chidlren[ord_] = &TrieNode{}
	}
	o.chidlren[ord_] = o.chidlren[ord_].Insert(s[1:])
	return &o
}

// 查找字符串s
func (o *TrieNode) Search(s []byte) *TrieNode {
	root := o
	for _, b := range s {
		root = root.chidlren[ord(b)]
		if root == nil {
			return nil
		}
	}

	if root.wordCount == 0 {
		return nil
	}
	return root
}

func ord(c byte) byte { return c - 'a' }
func chr(v byte) byte { return v + 'a' }

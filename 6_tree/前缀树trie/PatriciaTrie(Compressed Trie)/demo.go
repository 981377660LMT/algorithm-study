package main

import (
	"fmt"
	"strings"
)

// RadixNode 表示 Radix Trie 的节点
type RadixNode struct {
	prefix   string // 该节点代表的“边”前缀
	children map[string]*RadixNode
	isEnd    bool // 是否存在字符串在此结束
}

// NewRadixNode 创建一个新的 RadixNode
func NewRadixNode(prefix string) *RadixNode {
	return &RadixNode{
		prefix:   prefix,
		children: make(map[string]*RadixNode),
		isEnd:    false,
	}
}

// Insert 往当前节点插入字符串 key
func (node *RadixNode) Insert(key string) {
	// 如果 key 为空，则表示应在当前节点结束
	if len(key) == 0 {
		node.isEnd = true
		return
	}

	// 遍历子节点，找是否有公共前缀
	for childPrefix, childNode := range node.children {
		// 计算当前 key 与子节点前缀的公共前缀长度
		commonLen := commonPrefixLength(key, childPrefix)
		if commonLen > 0 {
			// 存在公共前缀
			if commonLen < len(childPrefix) {
				// 需要把 childPrefix 拆分
				// 原子节点前缀 => common 部分 + remainder
				remainder := childPrefix[commonLen:] // child 剩余部分
				// 创建新子节点，前缀为 remainder，继承childNode的子节点
				newChild := NewRadixNode(remainder)
				newChild.children = childNode.children
				newChild.isEnd = childNode.isEnd

				// 修改 childNode 的前缀为公共部分
				childNode.prefix = childPrefix[:commonLen]
				childNode.children = make(map[string]*RadixNode)
				childNode.isEnd = false

				// newChild 重新挂载到 childNode
				childNode.children[remainder] = newChild
			}

			// 对 key 的剩余部分继续插入
			keyRemainder := key[commonLen:]
			if len(keyRemainder) == 0 {
				// key 恰好匹配 childNode 的新 prefix
				childNode.isEnd = true
			} else {
				// 继续递归插入
				childNode.Insert(keyRemainder)
			}
			return
		}
	}

	// 如果没有找到公共前缀节点，则新建一个子节点
	newNode := NewRadixNode(key)
	newNode.isEnd = true
	node.children[key] = newNode
}

// Search 搜索 key 是否在树中
func (node *RadixNode) Search(key string) bool {
	// 若 key 为空，检查当前节点是否是单词结尾
	if len(key) == 0 {
		return node.isEnd
	}

	// 在子节点中寻找与 key 开头有前缀匹配的节点
	for childPrefix, childNode := range node.children {
		// 看 key 与 childPrefix 能匹配多少字符
		commonLen := commonPrefixLength(key, childPrefix)
		if commonLen == len(childPrefix) {
			// 匹配了整个 childPrefix
			// 那么在 childNode 中继续搜索剩余部分
			return childNode.Search(key[commonLen:])
		}
	}
	return false
}

// commonPrefixLength 计算两个字符串的公共前缀长度
func commonPrefixLength(a, b string) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	count := 0
	for i := 0; i < minLen; i++ {
		if a[i] == b[i] {
			count++
		} else {
			break
		}
	}
	return count
}

// PrintRadixTree 打印 Radix Trie 结构，方便观察
func (node *RadixNode) PrintRadixTree(level int) {
	indent := strings.Repeat("  ", level)
	// 标记是否是关键字终止
	endMark := ""
	if node.isEnd {
		endMark = " (End)"
	}
	fmt.Printf("%s- prefix: \"%s\"%s\n", indent, node.prefix, endMark)

	// 遍历子节点
	for _, child := range node.children {
		child.PrintRadixTree(level + 1)
	}
}

// --------------------- 测试示例 ---------------------
func main() {
	// 创建根节点，根节点的 prefix 可置空
	root := NewRadixNode("")

	// 一组测试字符串
	words := []string{"car", "cat", "dog", "catalog", "cattle", "do", "cart"}

	// 插入
	for _, w := range words {
		root.Insert(w)
	}

	// 打印 Radix Trie
	fmt.Println("Radix Trie Structure:")
	root.PrintRadixTree(0)
	fmt.Println()

	// 测试搜索
	tests := []string{"car", "cart", "cata", "catalog", "cattle", "dog", "do", "donut"}
	for _, t := range tests {
		fmt.Printf("Search \"%s\": %v\n", t, root.Search(t))
	}
}

package main

import (
	"fmt"
)

func main() {
	list := NewList()

	// 版本 0：空链表
	fmt.Println("在版本 0 读取第 0 个元素 =>", list.Get(0, 0)) // -1（表示空）

	// PushFront(10) => 版本 1
	v1 := list.PushFront(10)
	fmt.Println("新版本 v1 =", v1)
	fmt.Println("在版本 1 读取第 0 个元素 =>", list.Get(1, 0)) // 10

	// PushFront(20) => 版本 2
	v2 := list.PushFront(20)
	fmt.Println("新版本 v2 =", v2)
	fmt.Println("在版本 2 读取第 0 个元素 =>", list.Get(2, 0)) // 20
	fmt.Println("在版本 2 读取第 1 个元素 =>", list.Get(2, 1)) // 10
	fmt.Println("在版本 1 读取第 0 个元素 =>", list.Get(1, 0)) // 10 (不变)

	// PushFront(30) => 版本 3
	v3 := list.PushFront(30)
	fmt.Println("新版本 v3 =", v3)
	fmt.Println("在版本 3 读取第 0 个元素 =>", list.Get(3, 0)) // 30
	fmt.Println("在版本 3 读取第 1 个元素 =>", list.Get(3, 1)) // 20
	fmt.Println("在版本 3 读取第 2 个元素 =>", list.Get(3, 2)) // 10

	// 验证旧版本不受影响
	fmt.Println("在版本 2 读取第 0 个元素 =>", list.Get(2, 0)) // 20
	fmt.Println("在版本 1 读取第 0 个元素 =>", list.Get(1, 0)) // 10
}

// maxMods 表示每个节点可以记录多少条“修改”
const maxMods = 2

// modRecord 表示一次修改记录：哪个字段、在哪个版本、被赋予的新值
type modRecord struct {
	field   string      // "data" 或 "next"
	version int         // 该修改对应的版本号
	newVal  interface{} // 新值
}

// Node 是单链表中的节点
// - data: 节点的“当前（或初始化）”数据（read area）
// - next: 节点的“当前（或初始化）”下一指针
// - mods: 存储对本节点的历史修改（field, version, newVal）
// - parent: 存储唯一父指针(若有)，用于在节点分裂时修正父节点的引用
type Node struct {
	data   int
	next   *Node
	mods   []modRecord
	parent *Node
}

// List 维护了所有版本的“头指针”，以及当前最新版本号
type List struct {
	heads   []*Node // heads[v] 表示版本 v 对应的链表头指针
	version int     // 当前最新版本号（从 0 开始）
}

// NewList 创建一个空的部分持久化链表，初始版本号为 0
func NewList() *List {
	l := &List{
		heads:   make([]*Node, 1), // 先给一个空间存放 v=0 的头
		version: 0,
	}
	l.heads[0] = nil // 空链表
	return l
}

// readData 在给定 version 下读取节点 n 的 data
func readData(n *Node, version int) int {
	// 如果这个节点还没有任何 mods，直接返回 n.data
	if len(n.mods) == 0 {
		return n.data
	}

	// 找到“版本号不超过 version 的最新一次(field == "data")修改”
	latestVal := n.data
	latestVersion := -1
	for _, m := range n.mods {
		if m.field == "data" && m.version <= version && m.version > latestVersion {
			latestVal = m.newVal.(int)
			latestVersion = m.version
		}
	}
	return latestVal
}

// readNext 在给定 version 下读取节点 n 的 next 指针
func readNext(n *Node, version int) *Node {
	if len(n.mods) == 0 {
		return n.next
	}

	latestPtr := n.next
	latestVersion := -1
	for _, m := range n.mods {
		if m.field == "next" && m.version <= version && m.version > latestVersion {
			latestPtr = m.newVal.(*Node)
			latestVersion = m.version
		}
	}
	return latestPtr
}

// writeData 在 version 版本中，将节点 n 的 data 改为 newVal
// 返回实际生效的节点指针（可能是 n 自身，也可能是分裂后的新节点）
func writeData(n *Node, version int, newVal int) *Node {
	// 如果这个节点还有空位记录修改，则直接写日志
	if len(n.mods) < maxMods {
		n.mods = append(n.mods, modRecord{
			field:   "data",
			version: version,
			newVal:  newVal,
		})
		return n
	}

	// 否则就要进行“分裂”操作
	newNode := copyNode(n, version)
	// 修正 parent 的引用，让它在此版本下指向新Node
	fixParent(n, newNode, version)

	// 在 newNode 上做这次写操作
	newNode.mods = append(newNode.mods, modRecord{
		field:   "data",
		version: version,
		newVal:  newVal,
	})
	return newNode
}

// writeNext 在 version 版本中，将节点 n 的 next 指针改为 newNext
func writeNext(n *Node, version int, newNext *Node) *Node {
	if len(n.mods) < maxMods {
		n.mods = append(n.mods, modRecord{
			field:   "next",
			version: version,
			newVal:  newNext,
		})
		return n
	}

	newNode := copyNode(n, version)
	fixParent(n, newNode, version)

	newNode.mods = append(newNode.mods, modRecord{
		field:   "next",
		version: version,
		newVal:  newNext,
	})
	return newNode
}

// copyNode 根据当前最新状态（针对 version）复制节点 oldNode
func copyNode(oldNode *Node, version int) *Node {
	newNode := &Node{
		data:   readData(oldNode, version),
		next:   readNext(oldNode, version),
		mods:   make([]modRecord, 0, maxMods),
		parent: oldNode.parent, // 继承旧节点的 parent
	}
	return newNode
}

// fixParent 把 oldNode.parent 中指向 oldNode 的引用，改为指向 newNode
func fixParent(oldNode, newNode *Node, version int) {
	parent := oldNode.parent
	if parent == nil {
		// 说明 oldNode 可能是链表头，由 List.heads 来指向
		// 我们不在这里处理，而是由上层调用去更新 List.heads（见 PushFront 的处理）
		return
	}
	// 如果 parent 存在，则说明 parent.next 应该被改为 newNode
	// 我们需要一个 version 下写操作
	parentAfter := writeNext(parent, version, newNode)
	// 修正 newNode.parent
	newNode.parent = parentAfter
}

// PushFront 在链表头插入一个新节点，并产生新的版本
func (l *List) PushFront(value int) int {
	l.version++
	vNew := l.version

	// 如果 heads 还没扩容，则扩容
	if len(l.heads) <= vNew {
		newHeads := make([]*Node, vNew+1)
		copy(newHeads, l.heads)
		l.heads = newHeads
	}

	oldHead := l.heads[vNew-1] // 旧版本的头
	newNode := &Node{
		data:   value,
		next:   oldHead,
		mods:   make([]modRecord, 0, maxMods),
		parent: nil, // 它是新的头，没有父节点
	}

	// 如果旧头存在，让旧头的 parent = newNode，表示从新版本视角看，旧头被新头指向
	if oldHead != nil {
		oldHead.parent = newNode
	}

	l.heads[vNew] = newNode
	return vNew
}

// Get 返回在版本 version 下，链表第 index 个节点的 data；若越界则返回 -1
func (l *List) Get(version, index int) int {
	if version > l.version || version < 0 {
		fmt.Printf("Version %d 不存在\n", version)
		return -1
	}
	head := l.heads[version]
	if head == nil {
		return -1 // 空链表
	}

	curr := head
	for i := 0; i < index; i++ {
		if curr == nil {
			return -1 // 越界
		}
		curr = readNext(curr, version)
	}
	if curr == nil {
		return -1
	}
	return readData(curr, version)
}

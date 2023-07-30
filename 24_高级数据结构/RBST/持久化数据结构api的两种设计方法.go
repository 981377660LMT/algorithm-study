// !1. 所有修改操作返回新的根节点，不修改原树
// !2. 传入指向结点的指针 **node，修改结点的值.

package main

import "fmt"

func main() {
	node := &Node{foo: 0}
	node2 := Mutate1(node)
	fmt.Println(node.foo)  // 0
	fmt.Println(node2.foo) // 1

	node3 := &Node{foo: 0}
	node4 := node3
	Mutate2(&node4)
	fmt.Println(node3.foo, node4.foo) // 0 1
}

type Node struct {
	foo int
}

func Mutate1(node *Node) *Node {
	copy := *node
	copy.foo = 1
	return &copy
}

func Mutate2(node **Node) {
	*node = &Node{foo: 1}
}

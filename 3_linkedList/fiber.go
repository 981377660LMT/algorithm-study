// 单链表树遍历算法-fiber
// https://github.com/facebook/react/issues/7942

package main

type Fiber struct {
	child   *Fiber
	sibling *Fiber
	return_ *Fiber
}

func EnumerateFiber(root *Fiber, f func(f *Fiber)) {
	node := root
	for {
		f(node)
		if node.child != nil {
			node = node.child
			continue
		}
		if node == root {
			return
		}
		for node.sibling == nil {
			if node.return_ == nil || node.return_ == root {
				return
			}
			node = node.return_
		}
		node = node.sibling
	}
}

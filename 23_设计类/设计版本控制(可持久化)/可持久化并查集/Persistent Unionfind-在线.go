// https://judge.yosupo.jp/problem/persistent_unionfind
// 初始版本为-1
// 给定n个顶点的无向图,初始时都不连通
// 处理q个操作:
// 0 versioni u v => 在版本versioni上合并u和v
// 1 versioni u v => 在版本versioni上询问u和v是否连通 输出1/0
// n,q<=2e5
// 0<=version<i

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	git := make([]*Node, q+1)
	git[0] = Build(1, n)

	for i := 0; i < q; i++ {
		var op, version, u, v int
		fmt.Fscan(in, &op, &version, &u, &v)
		version++
		u++
		v++
		if op == 0 {
			newUnionFind := git[version].Union(u, v)
			git[i+1] = newUnionFind
		} else {
			if git[version].IsConnected(u, v) {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		}
	}
}

type Node struct {
	left, right           int
	parent, depth         int
	leftChild, rightChild *Node
}

// usage:
//  git := make([]*Node, maxVersion+1)  // restore all versions
//  git[0] = Build(1, n)  // init version 0
//  newUnionFind = git[0].Union(1, 2)  // version 1
//  git[1] = newUnionFind
//  newUnionFind = git[1].Union(2, 3)  // version 2
//  git[2] = newUnionFind
//  fmt.Println(git[2].IsConnected(1,3))  // true
// https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/union_find.go#L356
func Build(left, right int) *Node {
	o := &Node{left: left, right: right}
	if left == right {
		o.parent = left // !初始时i的父亲就是i自己(i>=1)
		return o
	}
	m := (left + right) >> 1
	o.leftChild = Build(left, m)
	o.rightChild = Build(m+1, right)
	return o
}

// !启发式合并：把深度小的合并到深度大的。若二者深度一样，则合并后的深度加一
//  1<=x,y<=n
func (o *Node) Union(x, y int) *Node {
	from, to := o.Find(x), o.Find(y)
	if from.parent == to.parent {
		return o
	}
	if from.depth > to.depth {
		from, to = to, from
	}
	p := o.changeParent(from.parent, to.parent)
	if from.depth == to.depth {
		p.addDepth(to.parent)
	}
	return p
}

//  1<=x<=n
func (o *Node) Find(x int) *Node {
	f := o.find(x)
	if f.parent == x {
		return f
	}
	return o.Find(f.parent)
}

//  1<=x,y<=n
func (o *Node) IsConnected(x, y int) bool {
	return o.Find(x).parent == o.Find(y).parent
}

func (o *Node) find(x int) *Node {
	if o.left == o.right {
		return o
	}
	if m := o.leftChild.right; x <= m {
		return o.leftChild.find(x)
	}
	return o.rightChild.find(x)
}

// !注意为了拷贝一份 pufNode，这里的接收器不是指针
func (o Node) changeParent(from, to int) *Node {
	if o.left == o.right {
		o.parent = to
		return &o
	}
	if m := o.leftChild.right; from <= m {
		o.leftChild = o.leftChild.changeParent(from, to)
	} else {
		o.rightChild = o.rightChild.changeParent(from, to)
	}
	return &o
}

func (o *Node) addDepth(x int) {
	if o.left == o.right {
		o.depth++
		return
	}
	if m := o.leftChild.right; x <= m {
		o.leftChild.addDepth(x)
	} else {
		o.rightChild.addDepth(x)
	}
}

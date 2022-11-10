// 笛卡尔树的单调栈建树方法
// https://oi-wiki.org/ds/cartesian-tree/

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

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	_, leftChild, rightChild := BuildCartesianTree2(nums)
	leftXor, rightXor := 0, 0
	for i := 0; i < n; i++ {
		leftXor ^= (leftChild[i] + 2) * (i + 1)
		rightXor ^= (rightChild[i] + 2) * (i + 1)
	}

	fmt.Fprintln(out, leftXor, rightXor)
}

type Node struct {
	weight, key int
	left, right *Node
}

// !指针版，返回根节点
func BuildCartesianTree1(insertNums []int) *Node {
	n := len(insertNums)
	if n == 0 {
		return nil
	}

	stack := []*Node{}
	for i, v := range insertNums {
		node := &Node{weight: i, key: v}
		var last *Node
		for len(stack) > 0 && stack[len(stack)-1].key > v {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			stack[len(stack)-1].right = node
		}
		if last != nil {
			node.left = last
		}

		stack = append(stack, node)
	}

	return stack[0]
}

// !非指针版，返回根节点索引以及每个节点的左右儿子的索引
//  !如果没有儿子，编号为-1
func BuildCartesianTree2(insertNums []int) (rootIndex int, leftChild, rightChild []int) {
	n := len(insertNums)
	stack := []int{}
	leftChild = make([]int, n)
	rightChild = make([]int, n)
	for i := 0; i < n; i++ {
		leftChild[i] = -1
		rightChild[i] = -1
	}

	for i, v := range insertNums {
		last := -1
		for len(stack) > 0 && insertNums[stack[len(stack)-1]] > v {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			rightChild[stack[len(stack)-1]] = i
		}
		if last != -1 {
			leftChild[i] = last
		}

		stack = append(stack, i)
	}

	rootIndex = stack[0]
	return
}

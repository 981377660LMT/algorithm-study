// https://rbtr.ee/fractional-cascading
// 给出k个长度为n的有序数组。
// 现在有q个查询: 给出数x，分别求出每个数组中大于等于x的最小的数，要求在线。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewReader(os.Stdin)
	var k, n, q int
	fmt.Fscanf(input, "%d %d %d\n", &k, &n, &q)

	arrs := make([][]int, k)
	for i := 0; i < k; i++ {
		arrs[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscanf(input, "%d", &arrs[i][j])
		}
		// 每个数组保证已排序，否则在这里 sort 一下
		// sort.Ints(arrs[i]) // 如果输入可能无序，这里需要 sort
	}

	// 构建 Fractional Cascading 结构
	fcList := buildFractionalCascading(arrs)

	// 在线处理 q 个查询
	for ; q > 0; q-- {
		var x int
		fmt.Fscanf(input, "%d\n", &x)
		res := fcQueryAll(fcList, x)

		// 输出结果
		// 若某数组中不存在 >= x, 就输出 -1；否则输出对应值
		for i, v := range res {
			if v == INF {
				fmt.Printf("-1")
			} else {
				fmt.Printf("%d", v)
			}
			if i+1 < len(res) {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

const INF = 1_000_000_000_000_000_000 // 一个很大的数

type FCNode struct {
	val     int
	origIdx int // 在原数组 A_i 中的位置(若本节点源自A_i)
	nextIdx int // 跳到 L_{i+1} 中对应位置
}

// 类似于一个跳表结构
type FCArray struct {
	nodes []FCNode // 增强数组 L_i
}

func buildFractionalCascading(arrs [][]int) []FCArray {
	k := len(arrs)
	result := make([]FCArray, k)

	// 先构建 L_{k-1} = 整个 A_{k-1}
	lastIndex := k - 1
	base := make([]FCNode, len(arrs[lastIndex]))
	for i, v := range arrs[lastIndex] {
		base[i] = FCNode{
			val:     v,
			origIdx: i,
			nextIdx: -1, // 没有下一个数组了
		}
	}
	result[lastIndex] = FCArray{nodes: base}

	// 再从后往前构建 L_{i}
	for i := k - 2; i >= 0; i-- {
		result[i] = mergeAndLinkFC(arrs[i], result[i+1])
	}

	return result
}
func mergeAndLinkFC(arr []int, nextFC FCArray) FCArray {
	n1 := len(arr)
	n2 := len(nextFC.nodes)

	// step1: 把 A_i 的所有元素转换为 FCNode
	arrNodes := make([]FCNode, n1)
	for i, v := range arr {
		arrNodes[i] = FCNode{
			val:     v,
			origIdx: i,
			nextIdx: -1,
		}
	}

	// step2: 从 nextFC 中抽取一半节点: 这里选取偶数下标
	pickNodes := make([]FCNode, 0)
	for idx := 0; idx < n2; idx += 2 {
		pickNodes = append(pickNodes, nextFC.nodes[idx])
	}

	// step3: 归并 arrNodes 和 pickNodes => merged
	merged := make([]FCNode, 0, len(arrNodes)+len(pickNodes))
	i1, i2 := 0, 0
	for i1 < len(arrNodes) && i2 < len(pickNodes) {
		if arrNodes[i1].val <= pickNodes[i2].val {
			merged = append(merged, arrNodes[i1])
			i1++
		} else {
			merged = append(merged, pickNodes[i2])
			i2++
		}
	}
	for i1 < len(arrNodes) {
		merged = append(merged, arrNodes[i1])
		i1++
	}
	for i2 < len(pickNodes) {
		merged = append(merged, pickNodes[i2])
		i2++
	}

	// step4: 构建 cross pointers (双向)
	// 我们需要在 merged 中和 nextFC.nodes 中找到同一个“值”的对应位置，互设 nextIdx
	// 但 merged 和 nextFC 都是有序的，我们可以用双指针扫描
	// （原则上只需要处理 pickNodes 对应的那部分）
	mp := merged
	np := nextFC.nodes

	im, in := 0, 0
	for im < len(mp) && in < len(np) {
		if mp[im].val == np[in].val {
			// 建立双向链接
			mp[im].nextIdx = in
			np[in].nextIdx = im
			im++
			in++
		} else if mp[im].val < np[in].val {
			im++
		} else {
			in++
		}
	}

	return FCArray{nodes: merged}
}

// 在 L_0（FCArray） 中找 "最小的 >= x" 的下标
func searchInFC(fc FCArray, x int) int {
	nodes := fc.nodes
	left, right := 0, len(nodes)
	for left < right {
		mid := (left + right) >> 1
		if nodes[mid].val >= x {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left // 可能等于 len(nodes)，表示都小于 x
}

func fcQueryAll(fcList []FCArray, x int) []int {
	k := len(fcList)
	ans := make([]int, k)

	// 1) 在 L_0 上做一次二分
	pos := searchInFC(fcList[0], x)
	var node FCNode

	// 2) 逐层级联
	for i := 0; i < k; i++ {
		L := fcList[i]
		if pos >= len(L.nodes) {
			// 越界 => 在 A_i 中无 >= x
			ans[i] = INF
			// 往下一层时，我们还是拿 pos = len(...) => nextIdx=-1
			// 这样下层要自己再做 fallback
			node = FCNode{val: INF, nextIdx: -1}
		} else {
			node = L.nodes[pos]
			ans[i] = node.val
		}

		// 准备跳到下一层
		if i+1 < k {
			if node.nextIdx >= 0 && node.nextIdx < len(fcList[i+1].nodes) {
				pos = node.nextIdx // 级联成功
			} else {
				// 没有有效nextIdx，做一次二分fallback
				pos = searchInFC(fcList[i+1], x)
			}
		}
	}

	return ans
}

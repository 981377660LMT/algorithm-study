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

	var q int
	fmt.Fscan(in, &q)
	operations := make([][2]int, q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		switch op {
		case "+":
			var x int
			fmt.Fscan(in, &x)
			operations[i] = [2]int{1, x}
		case "-":
			var k int
			fmt.Fscan(in, &k)
			operations[i] = [2]int{2, k}
		case "!":
			operations[i] = [2]int{3, 0}
		case "?":
			operations[i] = [2]int{4, 0}
		}
	}

	res := Rollbacks(operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/solution/CF1858E2?page=1
// 给定一个初始时为空的数组nums, 需要实现下面四种类型的操作：
// [1, x]: 将x添加到nums尾部
// [2, k]: 将尾部的k个数删除.保证存在k个数.
// [3, 0]: 撤销上一次操作1或2操作
// [4, 0]: 查询当前nums中有多少个不同的数
//
// 1<=q<=1e6,1<=x<1e6,询问次数不超过1e5
// 强制在线.
//
// https://zhuanlan.zhihu.com/p/650274675
// 因为要支持undo,这里将所有action(invert)全部保存到栈中。
//
// while(q--)
// 	if(kind == '+') {
// 		读取插入的元素 x

// 		如果 x 没有出现过 :
// 			更新 first[x]                          (1)
// 			更新 distinct[pos] = distinct[pos - 1] + 1 (2)
// 		否则
// 			更新 distinct[pos] = distinct[pos - 1]     (3)
// 		更新 nums[pos + 1] = x                          (4)
// 		更新 pos = pos + 1                             (5)

// 		将 (1) ~ (5) 的改动写入栈中, 方便还原。
// 	}

// 	if(kind == '-') {
// 		读取减少的数量 k

// 		更新 pos -= k;   (1)
// 		并将 (1) 写入栈中, 方便还原。
// 	}

//		if(kind == '!') {
//			根据栈中的数据，进行还原。
//		}
//		if(kind == '?') {
//			输出 distinct[pos]
//			continue;
//		}
//	}
func Rollbacks(operations [][2]int) []int {
	const INF int = 1e9
	const NONE int = INF
	const MAX int = 1e6

	q := len(operations)
	nums := make([]int, q+1) // 当前数组. nums[0]是虚拟节点(空).
	for i := range nums {
		nums[i] = NONE
	}
	first := [MAX + 5]int{} // first[x]表示x第一次出现的下标
	for i := range first {
		first[i] = INF
	}
	distinct := make([]int, q+1) // distinct[i]表示数组的前i个数中有多少个不同的数

	pos := 0 // 当前数组的长度

	type invert struct {
		ptr   *int
		value int
	}
	history := [][]invert{} // 记录action

	res := []int{}
	for _, op := range operations {
		kind, x := op[0], op[1]
		if kind == 1 {
			// !将x添加到nums尾部
			actions := []invert{}

			// 1. 处理懒删除造成的first不正确的情况
			//    如果nums[pos+1]不是最早出现的，那么就不用更新first
			//    否则将它更新为未出现状态
			if nums[pos+1] != NONE && first[nums[pos+1]] == pos+1 {
				actions = append(actions, invert{&first[nums[pos+1]], first[nums[pos+1]]})
				first[nums[pos+1]] = NONE
			}

			// 2. 更新first和distinct
			if first[x] <= pos {
				actions = append(actions, invert{&distinct[pos+1], distinct[pos+1]})
				distinct[pos+1] = distinct[pos]
			} else {
				actions = append(actions, invert{&first[x], first[x]})
				first[x] = pos + 1
				actions = append(actions, invert{&distinct[pos+1], distinct[pos+1]})
				distinct[pos+1] = distinct[pos] + 1
			}

			actions = append(actions, invert{&nums[pos+1], nums[pos+1]})
			nums[pos+1] = x
			actions = append(actions, invert{&pos, pos})
			pos++

			for i, j := 0, len(actions)-1; i < j; i, j = i+1, j-1 {
				actions[i], actions[j] = actions[j], actions[i]
			}
			history = append(history, actions)
		} else if kind == 2 {
			// !删除尾部的k个数，采用懒删除
			actions := []invert{}
			actions = append(actions, invert{&pos, pos})
			pos -= x
			history = append(history, actions)
		} else if kind == 3 {
			// !撤销上一次操作1或2操作
			actions := history[len(history)-1]
			history = history[:len(history)-1]
			for _, action := range actions {
				*action.ptr = action.value
			}
		} else {
			// !查询当前nums中有多少个不同的数
			res = append(res, distinct[pos])
		}
	}

	return res
}

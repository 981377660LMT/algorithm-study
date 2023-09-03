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
	R := NewRollbacksOnline(len(operations))
	res := []int{}
	for _, op := range operations {
		kind, x := op[0], op[1]
		switch kind {
		case 1:
			R.Append(x)
		case 2:
			R.Pop(x)
		case 3:
			R.Undo()
		case 4:
			res = append(res, R.Query())
		}
	}
	return res
}

const INF int = 1e9

type IStep interface {
	Apply()
	Invert() IStep // invert 产生一个 inverted 的 step
}

type SetStep struct {
	ptr   *int
	value int
}

func NewSetStep(ptr *int, value int) IStep {
	return &SetStep{ptr: ptr, value: value}
}

func (ss *SetStep) Apply() {
	*ss.ptr = ss.value
}

func (ss *SetStep) Invert() IStep {
	return &SetStep{ptr: ss.ptr, value: *ss.ptr}
}

type Transform struct {
	steps []IStep
}

func NewTransform(steps []IStep) *Transform {
	return &Transform{steps: steps}
}

func (bss *Transform) Apply() {
	for _, step := range bss.steps {
		step.Apply()
	}
}

func (bss *Transform) Invert() IStep {
	res := &Transform{}
	for i := len(bss.steps) - 1; i >= 0; i-- {
		res.steps = append(res.steps, bss.steps[i].Invert())
	}
	return res
}

type RollbacksOnline struct {
	nums      []int       // 保存nums中的数
	first     []int       // 保存nums中的数第一次出现的位置
	distinct  []int       // 每个版本中不同的数的个数
	hash      map[int]int // 将num映射到0~n-1
	undoStack []IStep
	redoStack []IStep
	pos       int // 当前数组的长度
}

func NewRollbacksOnline(maxOperation int) *RollbacksOnline {
	res := &RollbacksOnline{
		nums:     make([]int, maxOperation+1), // nums[0]是虚拟节点(空).
		first:    make([]int, maxOperation+1),
		distinct: make([]int, maxOperation+1),
		hash:     make(map[int]int),
	}
	for i := range res.first {
		res.first[i] = INF
	}
	return res
}

func (rb *RollbacksOnline) Append(num int) {
	num = rb._getHash(num)
	steps := []IStep{}
	v := rb.nums[rb.pos+1]

	// 1. 处理懒删除造成的first不正确的情况
	//    如果nums[pos+1]不是最早出现的，那么就不用更新first
	//    否则将它更新为未出现状态
	if v != INF && rb.first[v] == rb.pos+1 {
		steps = append(steps, NewSetStep(&rb.first[v], rb.first[v]))
		rb.first[v] = INF
	}

	// 2. 更新first和distinct
	if rb.first[num] <= rb.pos {
		steps = append(steps, NewSetStep(&rb.distinct[rb.pos+1], rb.distinct[rb.pos+1]))
		rb.distinct[rb.pos+1] = rb.distinct[rb.pos]
	} else {
		steps = append(steps, NewSetStep(&rb.first[num], rb.first[num]))
		rb.first[num] = rb.pos + 1
		steps = append(steps, NewSetStep(&rb.distinct[rb.pos+1], rb.distinct[rb.pos+1]))
		rb.distinct[rb.pos+1] = rb.distinct[rb.pos] + 1
	}

	steps = append(steps, NewSetStep(&rb.nums[rb.pos+1], rb.nums[rb.pos+1]))
	rb.nums[rb.pos+1] = num
	steps = append(steps, NewSetStep(&rb.pos, rb.pos))
	rb.pos++

	rb.undoStack = append(rb.undoStack, NewTransform(steps))

}

// 删除尾部len个元素.
func (rb *RollbacksOnline) Pop(len int) {
	step := NewSetStep(&rb.pos, rb.pos-len)
	step.Apply()
	rb.undoStack = append(rb.undoStack, step)
}

func (rb *RollbacksOnline) Undo() bool {
	if len(rb.undoStack) == 0 {
		return false
	}
	step := rb.undoStack[len(rb.undoStack)-1]
	rb.undoStack = rb.undoStack[:len(rb.undoStack)-1]
	step.Apply()
	rb.redoStack = append(rb.redoStack, step.Invert())
	return true
}

func (rb *RollbacksOnline) Redo() bool {
	if len(rb.redoStack) == 0 {
		return false
	}
	step := rb.redoStack[len(rb.redoStack)-1]
	rb.redoStack = rb.redoStack[:len(rb.redoStack)-1]
	step.Apply()
	rb.undoStack = append(rb.undoStack, step.Invert())
	return true
}

// 查询当前nums中有多少个不同的数.
func (rb *RollbacksOnline) Query() int {
	return rb.distinct[rb.pos]
}

func (rb *RollbacksOnline) Len() int {
	return rb.pos
}

func (rb *RollbacksOnline) _getHash(v int) int {
	res, ok := rb.hash[v]
	if ok {
		return res
	}
	res = len(rb.hash)
	rb.hash[v] = res
	return res
}

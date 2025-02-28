// 895. 最大频率栈
// https://leetcode.cn/problems/maximum-frequency-stack/solutions/1998430/mei-xiang-ming-bai-yi-ge-dong-hua-miao-d-oich/
// 设计一个类似堆栈的数据结构，将元素推入堆栈，并从堆栈中弹出出现频率最高的元素。
//
// 实现 FreqStack 类:
//
// FreqStack() 构造一个空的堆栈。
// void push(int val) 将一个整数 val 压入栈顶。
// int pop() 删除并返回堆栈中出现频率最高的元素。
// 如果出现频率最高的元素不只一个，则移除并返回最接近栈顶的元素。
//
// !把频率（出现次数）不同的元素，压入不同的栈中。每次出栈时，弹出含有频率最高元素的栈的栈顶。

package main

type FreqStack struct {
	stacks  [][]int
	counter map[int]int
}

func Constructor() FreqStack {
	return FreqStack{counter: make(map[int]int)}
}

func (f *FreqStack) Push(val int) {
	c := f.counter[val]
	if c == len(f.stacks) {
		f.stacks = append(f.stacks, []int{val})
	} else {
		f.stacks[c] = append(f.stacks[c], val)
	}
	f.counter[val]++
}

func (f *FreqStack) Pop() int {
	back2 := len(f.stacks) - 1
	stack := f.stacks[back2]
	back1 := len(stack) - 1
	val := stack[back1]
	if back1 == 0 {
		f.stacks = f.stacks[:back2]
	} else {
		f.stacks[back2] = stack[:back1]
	}
	f.counter[val]--
	return val
}

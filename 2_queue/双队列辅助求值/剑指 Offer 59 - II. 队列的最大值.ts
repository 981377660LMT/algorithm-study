class MaxQueue {
  private queue: number[]
  private maxQueue: number[] // 队头最大
  constructor() {
    this.queue = []
    this.maxQueue = []
  }

  max_value(): number {
    return this.maxQueue[0] ?? -1
  }

  push_back(value: number): void {
    this.queue.push(value)
    while (this.maxQueue.length && this.maxQueue[this.maxQueue.length - 1] < value) {
      this.maxQueue.pop()
    }
    this.maxQueue.push(value)
  }

  pop_front(): number {
    if (this.queue.length === 0) return -1
    const head = this.queue.shift()!
    head === this.maxQueue[0] && this.maxQueue.shift() // 最大值出队
    return head
  }
}

// 维护一个单调的双端队列
// 从队列尾部插入元素时，我们可以提前取出队列中所有比这个元素小的元素，
// 使得队列中只保留对结果有影响的数字

// 用一个队列保存正常元素，另一个双向队列保存单调递减的元素
// 入栈时，第一个队列正常入栈；第二个队列是递减队列，所以需要与之前的比较，从尾部把小于当前value的全部删除（因为用不到了）
// 出栈时，第一个队列正常出栈；第二个队列的头部与出栈的值作比较，如果相同，那么一起出栈
// https://leetcode-cn.com/problems/dui-lie-de-zui-da-zhi-lcof/solution/yi-miao-jiu-neng-du-dong-de-dong-hua-jie-b4de/

// 使最小元素位于栈顶
// 最多只能使用一个其他的临时栈存放数据

/**
 * 双栈的做法，核心就是用两个单调栈，一个栈顶大栈底小的栈装小数，一个栈顶小栈底大的栈装大数。
 * 其实这个和对顶堆求中位数的那个结构是很像的。
 */
class SortedStack {
  // 大顶栈increasedStack装小数 1,2,3 小顶栈decreasedStack装大数300,200,100
  constructor(private upStack: number[] = [], private downStack: number[] = []) {}

  push(val: number): void {
    while (this.upStack.length && this.upStack[this.upStack.length - 1] > val) {
      this.downStack.push(this.upStack.pop()!)
    }

    while (this.downStack.length && this.downStack[this.downStack.length - 1] < val) {
      this.upStack.push(this.downStack.pop()!)
    }

    this.downStack.push(val)
  }

  pop(): void {
    while (this.upStack.length) {
      this.downStack.push(this.upStack.pop()!)
    }
    this.downStack.pop()
  }

  peek(): number {
    while (this.upStack.length) {
      this.downStack.push(this.upStack.pop()!)
    }
    return this.downStack.length ? this.downStack[this.downStack.length - 1] : -1
  }

  isEmpty(): boolean {
    return this.upStack.length === 0 && this.downStack.length === 0
  }
}

/**
 * Your SortedStack object will be instantiated and called as such:
 * var obj = new SortedStack()
 * obj.push(val)
 * obj.pop()
 * var param_3 = obj.peek()
 * var param_4 = obj.isEmpty()
 */

const sortedStack = new SortedStack()

sortedStack.push(1)
sortedStack.push(3)
sortedStack.push(4)
sortedStack.push(2)
sortedStack.push(7)
console.log(sortedStack)
console.log(sortedStack.peek())
console.log(sortedStack)

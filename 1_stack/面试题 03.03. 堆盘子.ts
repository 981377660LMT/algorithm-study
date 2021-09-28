// 盘子堆到一定高度时，我们就会另外堆一堆盘子
// 1172. 餐盘栈
// 区别是这题当某个栈为空时，应当删除该栈；push是一直往后加而不是推入 从左往右第一个 没有满的栈。
class StackOfPlates {
  private cap: number
  private stack: number[][]

  constructor(cap: number) {
    this.cap = cap
    this.stack = []
  }

  push(val: number): void {
    if (this.cap <= 0) return
    if (!this.stack.length) {
      this.stack.push([val])
    } else {
      const last = this.stack[this.stack.length - 1]
      if (last.length === this.cap) {
        this.stack.push([val])
      } else {
        last.push(val)
      }
    }
  }

  pop(): number {
    return this.popAt(this.stack.length - 1)
  }

  popAt(index: number): number {
    if (this.stack.length - 1 < index || index < 0) return -1
    const res = this.stack[index].pop()!
    // 栈空则删除
    if (!this.stack[index].length) {
      this.stack.splice(index, 1)
    }

    return res
  }
}

// const SOP = new StackOfPlates(2)
// SOP.push(1)
// SOP.push(2)
// SOP.push(3)
// SOP.popAt(0)
// SOP.popAt(0)
// SOP.popAt(0)
// console.log(SOP)
const SOP = new StackOfPlates(1)
SOP.push(1)
SOP.push(2)
SOP.popAt(1)
SOP.pop()
SOP.pop()
console.log(SOP)

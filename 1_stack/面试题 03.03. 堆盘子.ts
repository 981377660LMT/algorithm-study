// 盘子堆到一定高度时，我们就会另外堆一堆盘子
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

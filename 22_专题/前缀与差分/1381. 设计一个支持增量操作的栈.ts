type Index = number
type Delta = number

class CustomStack {
  private stack: number[]
  private maxSize: number
  private diff: Map<Index, Delta>
  constructor(maxSize: number) {
    this.maxSize = maxSize
    this.stack = []
    this.diff = new Map()
  }

  push(x: number): void {
    if (this.stack.length < this.maxSize) {
      this.stack.push(x)
    }
  }

  pop(): number {
    if (this.stack.length === 0) return -1
    const index = this.stack.length - 1
    const delta = this.diff.get(index) ?? 0
    this.diff.delete(index)
    this.diff.set(index - 1, (this.diff.get(index - 1) ?? 0) + delta)
    return this.stack.pop()! + delta
  }

  /**
   *
   * @param k 栈底的 k 个元素的值都增加 val
   * @param val
   * 增量操作时只需要把增量存在 k 处那一个元素上
   * 我们只在出栈时才关心元素的值，所以在增量操作的时候，
   * 可以不用去更新栈内的元素，而是用一个 hashMap 来记录第几个元素需要增加多少
   * 出栈时，检查当前元素的下标是否在 hashMap 中有记录，有的话就加上增量再出栈
   */
  increment(k: number, val: number): void {
    const key = Math.min(this.stack.length, k) - 1
    this.diff.set(key, (this.diff.get(key) || 0) + val)
  }
}

const S = new CustomStack(3)
S.push(1)
S.push(2)
console.log(S.pop())
S.push(2)
S.push(3)
S.push(4)
console.log(S)
S.increment(5, 100)
S.increment(2, 100)

console.log(S)
console.log(S.pop())
console.log(S)
console.log(S.pop())
console.log(S)

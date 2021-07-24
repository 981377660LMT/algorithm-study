class Stack {
  private queue: unknown[]
  private tmp: unknown[]

  constructor() {
    this.queue = []
    this.tmp = []
  }

  push(x: number) {
    this.queue.push(x)
  }

  pop() {
    while (this.queue.length > 1) {
      this.tmp.push(this.queue.shift())
    }
    const ele = this.queue.shift()
    this.tmp.push(ele)
    this.queue = this.tmp
    this.tmp = []
    return ele
  }
}

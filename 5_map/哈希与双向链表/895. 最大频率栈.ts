// 参考LFU的freqToDoubleLinkedList与keyToNode
class FreqStack {
  private keyToFreq: Map<number, number>
  private freqToStack: Map<number, number[]>
  private maxFreq: number

  constructor() {
    this.keyToFreq = new Map()
    this.freqToStack = new Map()
    this.maxFreq = 0
  }

  push(val: number): void {
    this.keyToFreq.set(val, (this.keyToFreq.get(val) || 0) + 1)
    if (this.keyToFreq.get(val)! > this.maxFreq) {
      this.maxFreq = this.keyToFreq.get(val)!
    }
    const freq = this.keyToFreq.get(val)!
    !this.freqToStack.has(freq) && this.freqToStack.set(freq, [])
    this.freqToStack.get(freq)!.push(val)
  }

  /**
   *它移除并返回栈中出现最频繁的元素。
    如果最频繁的元素不只一个，则移除并返回最接近栈顶的元素。
   */
  pop(): number {
    const res = this.freqToStack.get(this.maxFreq)!.pop()!
    this.keyToFreq.set(res, this.keyToFreq.get(res)! - 1)
    if (this.freqToStack.get(this.maxFreq)!.length === 0) {
      this.maxFreq--
    }

    return res
  }

  static main() {
    const freqStack = new FreqStack()
    freqStack.push(5)
    freqStack.push(7)
    freqStack.push(5)
    freqStack.push(7)
    freqStack.push(4)
    console.log(freqStack)
    freqStack.push(5)
    console.log(freqStack.pop())
    console.log(freqStack.pop())
    console.log(freqStack.pop())
    console.log(freqStack.pop())
    console.log(freqStack.pop())
  }
}

FreqStack.main()

export {}

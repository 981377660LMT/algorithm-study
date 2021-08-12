class BSet {
  public state: number

  constructor() {
    this.state = 0
  }

  add(n: number) {
    this.state |= 1 << n
    return this
  }

  has(n: number) {
    return (this.state & (1 << n)) === 1 << n
  }

  delete(n: number) {
    if (this.has(n)) {
      this.state -= 1 << n
    }
    return this
  }

  // 位1的个数
  get size() {
    return this.state.toString(2).replace(/0/g, '').length
  }

  // 位1的个数
  get hammingWeight() {
    let sum = 0
    let n = this.state
    while (n) {
      sum += n & 1
      n = n >>> 1
    }
    return sum
  }
}

export { BSet }
console.log(2 ** 1000)
console.log(Number.MAX_VALUE)
console.log(Number.MAX_SAFE_INTEGER)

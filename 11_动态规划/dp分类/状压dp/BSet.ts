// console.log(0b10000000100 & 0b101)
console.log('a'.codePointAt(0)! - 97)
// 可用于判断字符串是否唯一
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

  get size() {
    return this.state.toString(2).replace(/0/g, '').length
  }
}

if (require.main === module) {
  const bset = new BSet()
  bset.add(4).add(2).add(3)
  console.log(bset.has(3))
  console.log(bset.size)
  bset.delete(4)
  console.log(bset.size, bset.has(4))
}

export { BSet }

class BIT {
  private size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array<number>(size + 1).fill(0)
  }

  add(x: number, k: number) {
    if (x <= 0) throw Error('add操作时树状数组索引应为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree[i] += k
    }
  }

  query(x: number) {
    let res = 0
    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += this.tree[i]
    }
    return res
  }

  private lowbit(x: number) {
    return x & -x
  }
}
// 每隔一段时间，你希望能找出数字 x 的秩(小于或等于 x 的值的个数)
// x <= 50000
class StreamRank {
  private bit: BIT
  constructor() {
    this.bit = new BIT(50001)
  }

  /**
   *
   * @param x 每读入一个数字都会调用该方法
   */
  track(x: number): void {
    this.bit.add(x + 1, 1)
  }

  /**
   *
   * @param x 返回小于或等于 x 的值的个数。
   */
  getRankOfNumber(x: number): number {
    return this.bit.query(x + 1)
  }
}

const streamRank = new StreamRank()

console.log(streamRank.getRankOfNumber(1))
console.log(streamRank.track(0))
console.log(streamRank.getRankOfNumber(0))

export {}

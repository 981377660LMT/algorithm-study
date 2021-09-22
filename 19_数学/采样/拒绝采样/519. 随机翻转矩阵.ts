import { randint } from '../randint'

// 尽量最少调用随机函数 Math.random()，并且优化时间和空间复杂度。
class Solution {
  private length: number
  private rows: number
  private cols: number
  private blackListMap: Map<number, number>

  constructor(m: number, n: number) {
    this.length = m * n
    this.rows = m
    this.cols = n
    this.blackListMap = new Map<number, number>()
  }

  /**
   * 均匀随机的将矩阵中的 0 变为 1，
   * 并返回该值的位置下标 [row_id,col_id]
   * @description
   * 一次调用random方法
   * 每次flip后都将flip的值记录到黑名单中，并将其映射到最后一个数上
   */
  flip(): number[] {
    if (this.length <= 0) return []
    this.length--
    const rand = randint(0, this.length)
    // rand 如果在黑名单 内 那么就把rand 对应的白名单里的那个数置为1 然后再映射到另一个白名单的数
    const res = this.blackListMap.get(rand) || rand
    this.blackListMap.set(rand, this.blackListMap.get(this.length) || this.length)

    return [~~(res / this.cols), res % this.cols]
  }

  /**
   * 将所有的值都重新置为 0
   */
  reset(): void {
    this.length = this.rows * this.cols
    this.blackListMap.clear()
  }
}

const s = new Solution(2, 3)
console.log(s.flip())
s.reset()
console.log(s.flip())
export {}

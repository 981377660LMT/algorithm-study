/**
 * 线性时间RMQ, 基于word-RAM model, 查询序列中任意区间的最小值的`索引`.
 * !和`SparseTable`差不多快,但是空间复杂度更低,为`O(n)`.
 * @see {@link  https://zhuanlan.zhihu.com/p/79423299}
 *      {@link  https://ei1333.github.io/library/structure/others/linear-rmq.hpp}
 *      {@link https://oi-wiki.org/topic/rmq/#%E5%9F%BA%E4%BA%8E%E7%8A%B6%E5%8E%8B%E7%9A%84%E7%BA%BF%E6%80%A7-rmq-%E7%AE%97%E6%B3%95}
 */
class LinearRMQ {
  private readonly _small: number[]
  private readonly _large: number[][]
  private readonly _less: (i: number, j: number) => boolean

  /**
   * @param n 序列长度.
   * @param less 索引i处的值是否小于索引j处的值.
   */
  constructor(n: number, less: (i: number, j: number) => boolean) {
    this._less = less
    let stack: number[] = []
    const small: number[] = []
    const large: number[][] = [[]]
    for (let i = 0; i < n; i++) {
      while (stack.length && !less(stack[stack.length - 1], i)) stack.pop()
      const tmp = stack.length ? small[stack[stack.length - 1]] : 0
      small.push(tmp | (1 << (i & 31)))
      stack.push(i)
      if (!((i + 1) & 31)) {
        large[0].push(stack[0])
        stack = []
      }
    }

    for (let i = 1; i << 1 <= n >> 5; i <<= 1) {
      const csz = (n >> 5) - (i << 1) + 1
      const v = Array(csz)
      for (let k = 0; k < csz; k++) {
        const back = large[large.length - 1]
        v[k] = this._getMin(back[k], back[k + i])
      }
      large.push(v)
    }

    this._small = small
    this._large = large
  }

  /**
   * 查询区间`[start, end)`中的最小值的`索引`.
   */
  query(start: number, end: number): number {
    if (start >= end) throw new Error(`start(${start}) should be less than end(${end})`)
    end--
    const l = (start >> 5) + 1
    const r = end >> 5
    if (l < r) {
      const msb = 31 - Math.clz32(r - l)
      const cache = this._large[msb]
      const i = ((l - 1) << 5) + trailingZeros32(this._small[(l << 5) - 1] & (~0 << (start & 31)))
      const cand1 = this._getMin(i, cache[l])
      const j = (r << 5) + trailingZeros32(this._small[end])
      const cand2 = this._getMin(cache[r - (1 << msb)], j)
      return this._getMin(cand1, cand2)
    }
    if (l === r) {
      const i = ((l - 1) << 5) + trailingZeros32(this._small[(l << 5) - 1] & (~0 << (start & 31)))
      const j = (l << 5) + trailingZeros32(this._small[end])
      return this._getMin(i, j)
    }
    return (r << 5) + trailingZeros32(this._small[end] & (~0 << (start & 31)))
  }

  private _getMin(i: number, j: number): number {
    return this._less(i, j) ? i : j
  }
}

function trailingZeros32(uint32: number): number {
  return 31 - Math.clz32(uint32 & -uint32)
}

export { LinearRMQ }

if (require.main === module) {
  console.log((1 << 31) >>> 0)
  console.log(1 << 31)
  console.log(1 & (1 << 31))

  const nums = Array(100)
    .fill(0)
    .map((_, i) => i)
  const rmq = new LinearRMQ(nums.length, (i, j) => nums[i] < nums[j])

  // check with brute force
  for (let i = 0; i < nums.length; i++) {
    for (let j = i + 1; j <= nums.length; j++) {
      const min = Math.min(...nums.slice(i, j))
      const minIndex = nums.indexOf(min, i)
      const rmqIndex = rmq.query(i, j)
      if (minIndex !== rmqIndex) {
        throw new Error('wrong')
      }
    }
  }
  console.log('pass')
}

// 1 <= val, inc, m <= 100
// 0 <= idx <= 1e5
// 总共最多会有 1e5 次对 append，addAll，multAll 和 getIndex 的调用。
// !TLE js的大数运算太慢了

import { AtcoderSegmentTree, Operation, useAtcoderLazySegmentTree } from '../AtcoderLazySegmentTree'

type Data = [sum: bigint, length: bigint]
type Lazy = [mul: bigint, add: bigint]

const MOD = BigInt(1e9 + 7)
const operation: Operation<Data, Lazy> = {
  e: () => [0n, 1n],
  id: () => [1n, 0n],
  op(data1, data2) {
    return [(data1[0] + data2[0]) % MOD, data1[1] + data2[1]]
  },
  // 区间和等于原来的区间和乘以mul加上区间的长度乘以add
  mapping(parentLazy, childData) {
    return [
      (childData[0] * parentLazy[0] + BigInt(childData[1]) * parentLazy[1]) % MOD,
      childData[1]
    ]
  },
  composition(parentLazy, childLazy) {
    return [
      (parentLazy[0] * childLazy[0]) % MOD,
      (parentLazy[0] * childLazy[1] + parentLazy[1]) % MOD
    ]
  }
}

// !注意单点查询 没有必要pushUp
class Fancy {
  private readonly tree: AtcoderSegmentTree<Data, Lazy>
  private length = 0

  constructor() {
    this.tree = useAtcoderLazySegmentTree(1e5 + 5, operation)
  }

  append(val: number): void {
    this.tree.update(this.length, this.length + 1, [1n, BigInt(val)])
    this.length++
  }

  addAll(inc: number): void {
    this.tree.update(0, this.length, [1n, BigInt(inc)])
  }

  multAll(m: number): void {
    this.tree.update(0, this.length, [BigInt(m), 0n])
  }

  getIndex(idx: number): number {
    if (idx >= this.length) return -1
    return Number(this.tree.query(idx, idx + 1)[0])
  }
}

/**
 * Your Fancy object will be instantiated and called as such:
 * var obj = new Fancy()
 * obj.append(val)
 * obj.addAll(inc)
 * obj.multAll(m)
 * var param_4 = obj.getIndex(idx)
 */
export {}

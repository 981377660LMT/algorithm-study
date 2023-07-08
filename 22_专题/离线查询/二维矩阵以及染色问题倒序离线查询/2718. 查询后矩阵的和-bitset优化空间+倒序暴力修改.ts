// 1e8很大,即使开全局数组+memset初始化也会超时
// 需要bitset优化空间
// !js暴力遍历,5e8=>8000ms(体现了js在力扣上的极限是5e8左右)

import { BitSet } from '../../../18_哈希/BitSet/BitSet'

function matrixSumQueries(n: number, queries: number[][]): number {
  const visited = new BitSet(n * n)
  let res = 0

  for (let i = queries.length - 1; ~i; i--) {
    const [type, index, val] = queries[i] // 比forEach快
    if (!type) {
      for (let j = 0; j < n; j++) {
        const pos = index * n + j
        if (!visited.has(pos)) {
          visited.add(pos)
          res += val
        }
      }
    } else {
      for (let j = 0; j < n; j++) {
        const pos = j * n + index
        if (!visited.has(pos)) {
          visited.add(pos)
          res += val
        }
      }
    }
  }

  return res
}

export {}

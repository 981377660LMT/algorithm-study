// 返回将 instructions 中所有元素依次插入 nums 后的 总最小代价
// 每一次插入操作的 代价 是以下两者的 较小值 ：
// nums 中 严格小于  instructions[i] 的数字数目。

import { BIT } from './BIT'

// nums 中 严格大于  instructions[i] 的数字数目。
// 1 <= instructions.length <= 10**5
// 1 <= instructions[i] <= 10**5
function createSortedArray(instructions: number[]): number {
  const MOD = 10 ** 9 + 7
  const bit = new BIT(10 ** 5)
  let res = 0
  for (const num of instructions) {
    const smaller = bit.sumRange(1, num - 1)
    const bigger = bit.sumRange(num + 1, bit.size)
    res += Math.min(smaller, bigger)
    res %= MOD
    bit.add(num, 1)
  }

  return res
}

console.log(createSortedArray([1, 5, 6, 2]))

export {}

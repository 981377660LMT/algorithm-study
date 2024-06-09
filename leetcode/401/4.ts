import { BitSet } from '../../18_哈希/BitSet/BitSet'

export {}

const INF = 2e9 // !超过int32使用2e15

// class Solution:
//     def maxTotalReward(self, rewardValues: List[int]) -> int:
//         rewardValues.sort()
//         dp = 1 << 0
//         for v in rewardValues:
//             low = dp & ((1 << v) - 1)
//             dp |= low << v
//         return dp.bit_length() - 1

function maxTotalReward(rewardValues: number[]): number {
  rewardValues = [...new Set(rewardValues)].sort((a, b) => a - b)
  const max = rewardValues[rewardValues.length - 1]
  const dp = new BitSet(max * 2)
  dp.add(0)
  rewardValues.forEach(v => {
    const low = dp.slice(0, v)
    dp.iorRange(v, 2 * v, low)
  })
  return dp.bitLength() - 1
}

if (require.main === module) {
  console.log(maxTotalReward([15, 20]))
}

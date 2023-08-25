// LIS的个数/LIS个数
// https://leetcode.com/problems/number-of-longest-increasing-subsequence/discuss/1643753/Python-O(nlogn)-solution-w-detailed-explanation-of-how-to-develop-a-binary-search-solution-from-300

import { DefaultDict } from '../../5_map/DefaultDict'
import { BIT1 } from '../../6_tree/树状数组/经典题/BIT'
import { bisectLeft, bisectRight } from '../../9_排序和搜索/二分/bisect'

function countLIS(nums: number[], isStrict = true): number {
  const n = nums.length
  if (n <= 1) return n

  const [rank, count] = sortedSet(nums)
  const bits = new DefaultDict(() => new BIT1(count + 10))
  const lis: number[] = []
  const bisect = isStrict ? bisectLeft : bisectRight
  for (let i = 0; i < n; i++) {
    const cur = nums[i]
    const pos = bisect(lis, cur)
    if (pos === lis.length) {
      lis.push(cur)
    } else {
      lis[pos] = cur
    }

    const preBit = bits.get(pos - 1)
    const curBit = bits.get(pos)
    const curRank = rank(cur)
    const smaller = preBit.query(curRank)
    curBit.add(curRank, smaller > 0 ? smaller : 1)
  }

  const lastPos = lis.length - 1
  return bits.get(lastPos).query(count)
}

function sortedSet(nums: number[]): [rank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const rank = (num: number) => {
    let left = 0
    let right = allNums.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (allNums[mid] >= num) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }
  return [rank, allNums.length]
}

if (require.main === module) {
  console.log(countLIS([1, 3, 5, 4, 7]))
  const arr1e5 = Array.from({ length: 1e5 }, (_, i) => Math.floor(Math.random() * 1e9))
  console.time('countLIS')
  console.log(countLIS(arr1e5))
  console.timeEnd('countLIS')
}

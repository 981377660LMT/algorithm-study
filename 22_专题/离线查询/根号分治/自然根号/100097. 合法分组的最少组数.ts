// 100097. 合法分组的最少组数
// https://leetcode.cn/problems/minimum-number-of-groups-to-create-a-valid-assignment/description/
//
// 给你一个长度为 n 下标从 0 开始的整数数组 nums 。
//
// 我们想将下标进行分组，使得 [0, n - 1] 内所有下标 i 都 恰好 被分到其中一组。
//
// 如果以下条件成立，我们说这个分组方案是合法的：
//
// 对于每个组 g ，同一组内所有下标在 nums 中对应的数值都相等。
// 对于任意两个组 g1 和 g2 ，两个组中 下标数量 的 差值不超过 1 。
// 请你返回一个整数，表示得到一个合法分组方案的 最少 组数。
//
// 最后每种频率需要拆成size和size+1两种
// !频率的种类数不超过根号n，因此可以直接枚举size
// !即 len(freqCounter) <= sqrt(n)

import { splitToKAndKPlusOne } from '../../../../19_数学/数论/扩展欧几里得/splitTo/splitToKAndKPlusOne'

function minGroupsForValidAssignment(nums: number[]): number {
  const n = nums.length
  const tmpCounter = new Map<number, number>()
  nums.forEach(v => tmpCounter.set(v, (tmpCounter.get(v) || 0) + 1))
  const freq = [...tmpCounter.values()]
  const freqCounter = new Map<number, number>()
  freq.forEach(v => freqCounter.set(v, (freqCounter.get(v) || 0) + 1))

  let res = n
  for (let size = 1; size < n; size++) {
    let ok = true
    let cand = 0
    for (const value of freqCounter.keys()) {
      const [count1, count2, ok_] = splitToKAndKPlusOne(value, size)
      if (!ok_) {
        ok = false
        break
      }
      cand += (count1 + count2) * freqCounter.get(value)!
    }
    if (ok) res = Math.min(res, cand)
  }

  return res
}

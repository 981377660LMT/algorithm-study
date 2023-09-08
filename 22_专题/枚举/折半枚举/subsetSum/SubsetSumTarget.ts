import { BitSet } from '../../../../18_哈希/BitSet/BitSet'
import { subsetSumSortedWithState } from './subsetSum'

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @param arr 非负整数数组.
 * @param target 非负整数.
 * @returns [res, ok] res为选出的数的下标，ok表示是否存在.
 */
function subsetSumTarget(arr: number[], target: number): [res: number[], ok: boolean] {
  if (target <= 0) return [[], false]
  const n = arr.length
  if (!n) return [[], false]
  let max = 0
  for (let i = 0; i < n; i++) {
    max = Math.max(max, arr[i])
  }

  const cost1 = n * max
  const cost2 = (n / 32) * target // 经验值
  let cost3 = 2e15
  if (n <= 50) {
    cost3 = 1 << ((n >>> 1) + 2)
  }

  const minCost = Math.min(cost1, cost2, cost3)
  if (minCost === cost1) return subsetSumTargetDp(arr, target)
  if (minCost === cost2) return subsetSumTargetBitset(arr, target)
  return subsetSumTargetMeetInMiddle(arr, target)
}

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @complexity O(n * max(arr)).
 */
function subsetSumTargetDp(arr: number[], target: number): [res: number[], ok: boolean] {
  if (target <= 0) {
    return [[], false]
  }

  const n = arr.length
  let max = 0
  for (let i = 0; i < n; i++) {
    max = Math.max(max, arr[i])
  }
  let right = 0
  let curSum = 0
  while (right < n && curSum + arr[right] <= target) {
    curSum += arr[right]
    right++
  }
  if (right === n && curSum !== target) {
    return [[], false]
  }

  const offset = target - max + 1
  let dp = Array<number>(2 * max).fill(-1)
  const pre = new Int32Array(n * (2 * max)).fill(-1)
  dp[curSum - offset] = right
  for (let i = right; i < n; i++) {
    const ndp = dp.slice()
    for (let j = 0; j < max; j++) {
      if (ndp[j + arr[i]] < dp[j]) {
        ndp[j + arr[i]] = dp[j]
        pre[i * (2 * max) + j + arr[i]] = -2
      }
    }
    for (let j = 2 * max - 1; j >= max; j--) {
      for (let k = ndp[j] - 1; k >= Math.max(dp[j], 0); k--) {
        if (ndp[j - arr[k]] < k) {
          ndp[j - arr[k]] = k
          pre[i * (2 * max) + j - arr[k]] = k
        }
      }
    }
    dp = ndp
  }

  if (dp[max - 1] === -1) {
    return [[], false]
  }

  const used = new Uint8Array(n)
  let i = n - 1
  let j = max - 1
  while (i >= right) {
    const p = pre[i * (2 * max) + j]
    if (p === -2) {
      used[i] ^= 1
      j -= arr[i]
      i--
    } else if (p === -1) {
      i--
    } else {
      used[p] ^= 1
      j += arr[p]
    }
  }

  while (i >= 0) {
    used[i] ^= 1
    i--
  }

  const res = []
  for (let i = 0; i < n; i++) {
    if (used[i]) res.push(i)
  }
  return [res, true]
}

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @complexity O(n * target / 32).
 */
function subsetSumTargetBitset(arr: number[], target: number): [res: number[], ok: boolean] {
  const n = arr.length
  const order = argSort(arr)
  let dp = new BitSet(1, 1)
  const last = new Int32Array(target + 1).fill(-1)
  for (let k = 0; k < n; k++) {
    const v = arr[order[k]]
    if (v > target) continue
    const newSize = Math.min(dp.size + v, target + 1)
    const ndp = dp.copy(newSize)
    dp.resize(newSize - v)
    ndp.iorRange(v, newSize, dp)
    for (let i = 0; i < ndp.bits.length; i++) {
      const updatedBits = i < dp.bits.length ? dp.bits[i] ^ ndp.bits[i] : ndp.bits[i]
      enumerateBits32(updatedBits, p => {
        last[(i << 5) | p] = order[k]
      })
    }
    dp = ndp
  }

  if (target >= dp.size || !dp.has(target)) {
    return [[], false]
  }

  const res: number[] = []
  while (target > 0) {
    const i = last[target]
    res.push(i)
    target -= arr[i]
  }
  return [res, true]

  function enumerateBits32(s: number, f: (bit: number) => void): void {
    while (s) {
      const i = 31 - Math.clz32(s & -s) // lowbit.bit_length() - 1
      f(i)
      s ^= 1 << i
    }
  }

  function argSort(nums: number[]): number[] {
    const order: number[] = Array(nums.length)
    for (let i = 0; i < order.length; i++) order[i] = i
    order.sort((i, j) => nums[i] - nums[j])
    return order
  }
}

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @complexity O(2^(n/2)).注意常数较大.
 */
function subsetSumTargetMeetInMiddle(arr: number[], target: number): [res: number[], ok: boolean] {
  const n = arr.length
  const mid = n >>> 1
  const dp1 = subsetSumSortedWithState(arr.slice(0, mid))
  const dp2 = subsetSumSortedWithState(arr.slice(mid))
  let left = 0
  let right = dp2.length - 1
  while (left < dp1.length && right >= 0) {
    const sum = dp1[left][0] + dp2[right][0]
    if (sum === target) {
      return [resolveState(mid, dp1[left][1], n - mid, dp2[right][1]), true]
    }
    if (sum < target) {
      left++
    } else {
      right--
    }
  }

  return [[], false]

  // eslint-disable-next-line max-len
  function resolveState(leftSize: number, leftState: number, rightSize: number, rightState: number): number[] {
    const res: number[] = []
    for (let i = 0; i < leftSize; i++) {
      if (leftState & (1 << i)) res.push(i)
    }
    for (let i = 0; i < rightSize; i++) {
      if (rightState & (1 << i)) res.push(i + leftSize)
    }
    return res
  }
}

export { subsetSumTarget, subsetSumTargetDp, subsetSumTargetBitset, subsetSumTargetMeetInMiddle }

if (require.main === module) {
  console.log(subsetSumTargetDp([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetBitset([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetMeetInMiddle([2, 3, 4, 5, 6, 7, 8, 9], 10))

  const max = 1e5
  const n = 200
  const nums = Array.from({ length: n }, () => max)
  const target = 3e6

  console.time('subsetSumTargetDp')
  subsetSumTargetDp(nums, target)
  console.timeEnd('subsetSumTargetDp')

  console.time('subsetSumTargetBitset')
  subsetSumTargetBitset(nums, target)
  console.timeEnd('subsetSumTargetBitset')

  const cost1 = n * max
  const cost2 = (n * target) / 32
  console.log(cost1, cost2)
}

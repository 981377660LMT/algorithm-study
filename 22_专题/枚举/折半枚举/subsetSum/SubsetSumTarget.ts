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
  let sum = 0
  for (let i = 0; i < n; i++) {
    max = Math.max(max, arr[i])
    sum += arr[i]
  }

  const cost1 = n * max
  const cost2 = (sum * Math.floor(Math.sqrt(sum))) / 100
  const cost3 = (n * target) / 32 // 经验值
  let cost5 = 2e15
  if (n <= 46) {
    cost5 = 1 << ((n >>> 1) + 1)
  }

  const minCost = Math.min(cost1, cost2, cost3, cost5)
  if (minCost === cost1) return subsetSumTargetDp1(arr, target)
  if (minCost === cost2) return subsetSumTargetDp2(arr, target)
  if (minCost === cost3) return subsetSumTargetDp3(arr, target)
  return subsetSumTargetDp5(arr, target)
}

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @complexity O(n * max(arr)).
 */
function subsetSumTargetDp1(arr: number[], target: number): [res: number[], ok: boolean] {
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
 * @complexity O(sum(arr) ^ 1.5).
 */
function subsetSumTargetDp2(arr: number[], target: number): [res: number[], ok: boolean] {
  let sum = 0
  for (let i = 0; i < arr.length; i++) {
    sum += arr[i]
  }
  if (target > sum) {
    return [[], false]
  }

  const counter = new Uint32Array(sum + 1)
  const dp = new Uint8Array(sum + 1)
  const last = new Uint32Array(sum + 1)
  const id: Map<number, number[]> = new Map()
  for (let i = 0; i < arr.length; i++) {
    const v = arr[i]
    if (v <= sum) {
      if (!id.has(v)) id.set(v, [])
      id.get(v)!.push(i)
      counter[v]++
    }
  }

  dp[0] = 1
  for (let i = 1; i <= sum; i++) {
    if (!counter[i]) continue
    for (let j = 0; j < i; j++) {
      let c = 0
      for (let k = j; k <= sum; k += i) {
        if (dp[k] === 1) {
          c = counter[i]
        } else if (c > 0) {
          dp[k] = 1
          c--
          last[k] = id.get(i)![c]
        }
      }
    }
  }

  if (!dp[target]) {
    return [[], false]
  }

  const res: number[] = []
  while (target > 0) {
    res.push(last[target])
    target -= arr[last[target]]
  }

  return [res, true]
}

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @complexity O(n * target / 32).
 */
function subsetSumTargetDp3(arr: number[], target: number): [res: number[], ok: boolean] {
  const n = arr.length
  const order = argSort(arr)
  let dp = new BitSet(1, 1)
  const last = new Int32Array(target + 1).fill(-1)
  for (let k = 0; k < n; k++) {
    const v = arr[order[k]]
    if (v > target) continue
    const newSize = dp.size + v
    const ndp = dp.copy(newSize)
    ndp.iorRange(v, newSize, dp)
    if (ndp.size > target + 1) ndp.resize(target + 1)
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
 * @complexity O(sum(arr) ^ 1.5 / 32).常数较大.
 * @deprecated
 */
function subsetSumTargetDp4(arr: number[], target: number): [res: number[], ok: boolean] {
  let sum = 0
  for (let i = 0; i < arr.length; i++) {
    sum += arr[i]
  }
  if (target > sum) {
    return [[], false]
  }
  const n = arr.length
  const ids: number[][] = Array(sum + 1)
  for (let i = 0; i < ids.length; i++) ids[i] = []
  for (let i = 0; i < n; i++) {
    ids[arr[i]].push(i)
  }
  const pre: { a: number; b: number }[] = Array(n)
  for (let i = 0; i < pre.length; i++) pre[i] = { a: -1, b: -1 }
  const grpVals: number[] = []
  const rawIdx: number[] = []
  for (let x = 1; x <= sum; x++) {
    const I = ids[x]
    while (I.length >= 3) {
      const a = I.pop()!
      const b = I.pop()!
      const c = pre.length
      pre.push({ a, b })
      ids[2 * x].push(c)
    }
    for (let i = 0; i < I.length; i++) {
      grpVals.push(x)
      rawIdx.push(I[i])
    }
  }
  const [I, tmp] = subsetSumTargetDp3(grpVals, target)
  if (!tmp) {
    return [[], false]
  }
  const res: number[] = []
  for (let i = 0; i < I.length; i++) {
    const st = [rawIdx[I[i]]]
    while (st.length > 0) {
      const c = st.pop()!
      if (c < n) {
        res.push(c)
        continue
      }
      const { a, b } = pre[c]
      st.push(a, b)
    }
  }
  return [res, true]
}

/**
 * 能否从arr中选出若干个数，使得它们的和为target.
 * @complexity O(2^(n/2)).注意常数较大.
 */
function subsetSumTargetDp5(arr: number[], target: number): [res: number[], ok: boolean] {
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

export {
  subsetSumTarget,
  subsetSumTargetDp1,
  subsetSumTargetDp2,
  subsetSumTargetDp3,
  subsetSumTargetDp3 as subsetSumTargetBitset,
  subsetSumTargetDp4,
  subsetSumTargetDp5,
  subsetSumTargetDp5 as subsetSumTargetMeetInMiddle
}

if (require.main === module) {
  console.log(subsetSumTargetDp1([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetDp2([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetDp3([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetDp4([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetDp5([2, 3, 4, 5, 6, 7, 8, 9], 10))
  console.log(subsetSumTargetDp1([3], 3))
  console.log(subsetSumTargetDp2([3], 3))
  console.log(subsetSumTargetDp3([3], 3))
  console.log(subsetSumTargetDp4([3], 3))
  console.log(subsetSumTargetDp5([3], 3))

  const max = 1e6
  const n = 10
  const nums = Array.from({ length: n }, () => max)
  const target = 3e6
  const sum = nums.reduce((a, b) => a + b, 0)

  console.time('subsetSumTargetDp1')
  subsetSumTargetDp1(nums, target)
  console.timeEnd('subsetSumTargetDp1')

  console.time('subsetSumTargetDp2')
  subsetSumTargetDp2(nums, target)
  console.timeEnd('subsetSumTargetDp2')

  console.time('subsetSumTargetDp3')
  subsetSumTargetDp3(nums, target)
  console.timeEnd('subsetSumTargetDp3')

  console.time('subsetSumTargetDp4')
  subsetSumTargetDp4(nums, target)
  console.timeEnd('subsetSumTargetDp4')

  // console.time('subsetSumTargetDp5')
  // subsetSumTargetDp5(nums, target)
  // console.timeEnd('subsetSumTargetDp5')

  const cost1 = n * max
  const cost2 = sum * Math.floor(Math.sqrt(sum))
  const cost3 = (n * target) / 32
  const cost4 = (sum * Math.floor(Math.sqrt(sum))) / 100
  const cost5 = 2 ** (n / 2)
  console.log(cost1, cost2, cost3, cost4, cost5)
}

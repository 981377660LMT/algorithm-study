import { enumeratePermutations } from '../../13_回溯算法/itertools/permutations'
import { SegmentTreeRangeUpdateRangeQuery } from '../../6_tree/线段树/template/atcoder_segtree/SegmentTreeRangeUpdateRangeQuery'

export {}

const INF = 2e15

function minimumTime(nums1: number[], nums2: number[], x: number): number {
  const initSum = nums1.reduce((pre, cur) => pre + cur, 0)
  if (initSum <= x) return 0
  const nums2Sum = nums2.reduce((pre, cur) => pre + cur, 0)
  let left = 0
  let right = nums1.length
  let ok = false
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) {
      right = mid - 1
      ok = true
    } else {
      left = mid + 1
    }
  }

  return ok ? left : -1
  // [
  //   0, 3, 6, 7,
  //   1, 2, 4, 5
  // ]
  // 每一轮删除 最大的 nums[i]+(mid-i)*nums2[i]
  function check(mid: number): boolean {
    if (mid === nums1.length) {
      // 按照 nums2 的顺序删除
      const perm = Array.from({ length: nums1.length }, (_, i) => i)
      perm.sort((a, b) => nums2[a] - nums2[b])
      const curNums1 = nums1.slice()
      let curSum = initSum
      for (let i = 0; i < mid; i++) {
        for (let j = 0; j < perm.length; j++) {
          curNums1[j] += nums2[j]
          curSum += nums2[j]
        }
        curSum -= curNums1[perm[i]]
        curNums1[perm[i]] = 0
        if (curSum <= x) return true
      }
      return curSum <= x
    } else {
      const curNums1 = nums1.map((num, i) => num + mid * nums2[i])
      let curSum = initSum + mid * nums2Sum
      if (curSum <= x) return true
      for (let i = 0; i < mid; i++) {
        let max = -INF
        let maxIndex = -1
        for (let j = 0; j < curNums1.length; j++) {
          if (curNums1[j] > max) {
            max = curNums1[j]
            maxIndex = j
          }
          curNums1[j] -= nums2[j]
        }
        curSum -= max
        curNums1[maxIndex] = 0
        if (curSum <= x) return true
      }
      return curSum <= x
    }
  }
}

if (require.main === module) {
  // nums1 = [1,2,3], nums2 = [1,2,3], x = 4
  // console.log(minimumTime([1, 2, 3], [1, 2, 3], 4))
  // [9,2,8,3,1,9,7,6]
  // [0,3,4,1,3,4,2,1]
  // 40

  console.log(minimumTime([9, 2, 8, 3, 1, 9, 7, 6], [0, 3, 4, 1, 3, 4, 2, 1], 40))
  const nums1 = [9, 2, 8, 3, 1, 9, 7, 6]
  const nums2 = [0, 3, 4, 1, 3, 4, 2, 1]
  const x = 40
  enumeratePermutations(
    Array.from({ length: nums1.length }, (_, i) => i),
    nums1.length,
    perm => {
      const curNums1 = nums1.slice()
      let curSum = nums1.reduce((pre, cur) => pre + cur, 0)
      for (let i = 0; i < perm.length; i++) {
        for (let j = 0; j < perm.length; j++) {
          curNums1[j] += nums2[j]
          curSum += nums2[j]
        }
        curSum -= curNums1[perm[i]]
        curNums1[perm[i]] = 0
        if (curSum <= x) {
          const cand = perm.slice(0, i + 1)
          const target = [5, 2, 1, 6, 4, 0, 7, 3].reverse()
          if (cand.length === target.length && cand.every((v, i) => v === target[i])) {
            console.log('cand', cand)
            console.log('curNums1', curNums1)
            console.log('curSum', curSum)
            console.log('i', i)
          }
          break
        }
      }
    }
  )
}

import { bisectRight } from '../../9_排序和搜索/二分api/7_二分搜索寻找最插右入位置'

/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number[]}
 * A 相对于 B 的优势可以用满足 A[i] > B[i] 的索引 i 的数目来描述。
 * 返回 A 的任意排列，使其相对于 B 的优势最大化。
 * 我们需要每次找一个比B[i]大的数，但是最接近B[i]的数，相当于A的上等马->B的中等马
   如果找不到，那么我们就给B[i]一个最小的数，相当于A的下等马->B的上等马
  
 */
const advantageCount1 = function (nums1: number[], nums2: number[]): number[] {
  nums1.sort((a, b) => a - b)

  const res: number[] = []

  for (const num of nums2) {
    const index = bisectRight(nums1, num)
    const choose = index < nums1.length ? nums1.splice(index, 1)[0] : nums1.shift()!
    res.push(choose)
  }

  return res
}
// 首先排序两个数组，然后一一比对，如果a比b大，那么a后面的所有数字都比b大
// ---> 说明a是大于b的数字中最小的一个数字
// 然后按照原顺序把答案映射回去
const advantageCount2 = function (nums1: number[], nums2: number[]): number[] {
  const assign = new Map<number, number[]>(nums2.map(v => [v, []]))
  const small: number[] = []

  let i = 0
  for (const num of nums1) {
    // a是大于b的数字中最小的一个数字
    if (num > nums2[i]) {
      assign.get(nums2[i])!.push(num)
      i++
    } else {
      small.push(num)
    }
  }

  const res: number[] = []
  console.log(assign)
  for (const num of nums2) {
    if (assign.has(num) && assign.get(num)!.length) res.push(assign.get(num)!.pop()!)
    else res.push(small.pop()!)
  }

  return res
}
console.log(advantageCount2([2, 7, 11, 15], [1, 10, 4, 11]))

/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number[]}
 * @description nums1、nums2有序时，可以采用双指针
 * 输出结果中每个元素出现的次数，应与元素在两个数组中出现次数的最小值一致。
 */
var intersect = function (nums1, nums2) {
  const m1 = new Map()
  const m2 = new Map()
  const res = []

  nums1.forEach(num => m1.set(num, m1.get(num) + 1 || 1))
  nums2.forEach(num => m2.set(num, m2.get(num) + 1 || 1))

  for (const key of m1.keys()) {
    if (m2.has(key)) {
      const count = Math.min(m1.get(key), m2.get(key))
      for (let i = 0; i < count; i++) {
        res.push(key)
      }
    }
  }

  return res
}

console.log(intersect([1, 2, 2, 1], [2, 2]))

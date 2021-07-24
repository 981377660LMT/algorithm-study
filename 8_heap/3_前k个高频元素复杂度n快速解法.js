/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 * @description 思路是将频率作为数组的index(key)，再倒序遍历数组
 * @description 时间复杂度O(n)
 */
var topKFrequent = function (nums, k) {
  const res = []
  // num映射到count的map
  const map = new Map()
  // bucket[i]存储count为i的数的集合
  const bucket = []

  nums.forEach(num => map.set(num, map.get(num) + 1 || 1))

  for (const [num, count] of map) {
    bucket[count] = (bucket[count] || new Set()).add(num)
  }

  for (let i = bucket.length - 1; i >= 0; i--) {
    if (bucket[i]) res.push(...bucket[i])
    if (res.length === k) break
  }

  return res
}

console.log(topKFrequent([1, 1, 1, 2, 2, 3], 2))

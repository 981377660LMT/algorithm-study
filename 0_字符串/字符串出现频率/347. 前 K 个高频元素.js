/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 */
var topKFrequent = function (nums, k) {
  const res = []
  const map = new Map()
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

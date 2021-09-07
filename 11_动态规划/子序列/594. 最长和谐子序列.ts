// 和谐数组是指一个数组里元素的最大值和最小值之间的差别 正好是 1 。
// 找到最长的和谐子序列的长度
type Value = number
type Count = number

function findLHS(nums: number[]): number {
  const map = new Map<Value, Count>()
  for (const num of nums) {
    map.set(num, (map.get(num) || 0) + 1)
  }

  let max = 0
  for (const key of map.keys()) {
    if (map.has(key + 1)) max = Math.max(max, map.get(key)! + map.get(key + 1)!)
  }

  return max
}

console.log(findLHS([1, 3, 2, 2, 5, 2, 3, 7]))
// 最长的和谐子序列是 [3,2,2,2,3]

export default 1

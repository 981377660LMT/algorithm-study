// 和谐数组是指一个数组里元素的最大值和最小值之间的差别 正好是 1 。
// 找到最长的和谐子序列的长度

// 1. 哈希表
function findLHS(nums: number[]): number {
  const counter = new Map<number, number>()
  for (const num of nums) {
    counter.set(num, (counter.get(num) ?? 0) + 1)
  }

  let max = 0
  for (const key of counter.keys()) {
    if (counter.has(key + 1)) max = Math.max(max, counter.get(key)! + counter.get(key + 1)!)
  }

  return max
}

// 2. 滑窗

console.log(findLHS([1, 3, 2, 2, 5, 2, 3, 7]))
// 最长的和谐子序列是 [3,2,2,2,3]

export default 1

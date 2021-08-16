/**
 * @param {number[]} nums
 * @return {number[][]}
 * 递增子序列中 至少有两个元素
 * 找出并返回所有该数组中不同的递增子序列，
 * 由于需要找到所有的递增子序列，因此动态规划就不行了，妥妥回溯就行了
 */
const findSubsequences = function (nums: number[]): number[][] {
  const res: number[][] = []

  const bt = (start: number, path: number[]) => {
    // 至少有两个元素
    if (path.length >= 2) {
      if (path[path.length - 1] >= path[path.length - 2]) {
        res.push(path.slice())
      } else return
    }

    const visited = new Set<number>()
    for (let i = start; i < nums.length; i++) {
      // if (i > start && nums[i] === nums[i - 1]) continue
      if (visited.has(nums[i])) continue
      visited.add(nums[i])
      path.push(nums[i])
      bt(i + 1, path)
      path.pop()
    }
  }
  bt(0, [])

  return res
}

console.log(findSubsequences([2, 3, 2, 2]))
// 输出：
// [[2,3],[2,2],[2,2,2]]

export default 1

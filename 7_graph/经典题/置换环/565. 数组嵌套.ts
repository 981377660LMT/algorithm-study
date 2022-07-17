/**
 *
 * 在置换群中找到最长的环
 *
 * @param {number[]} nums 0 - n-1 的 排列
 * @return {number}
 * 找到最大的集合S并返回其大小，其中 S[i] = {A[i], A[A[i]], A[A[A[i]]], ... }
 * @summary
 */
function arrayNesting(nums: number[]): number {
  const n = nums.length
  const visited = new Uint8Array(n)
  let res = 0

  for (const num of nums) {
    if (visited[num]) continue
    dfs(num, 0)
  }

  return res

  function dfs(cur: number, size: number): void {
    if (visited[cur]) {
      res = Math.max(res, size)
      return
    }

    visited[cur] = 1
    dfs(nums[cur], size + 1)
  }
}

console.log(arrayNesting([5, 4, 0, 3, 1, 6, 2]))
// 输出: 4
// 解释:
// A[0] = 5, A[1] = 4, A[2] = 0, A[3] = 3, A[4] = 1, A[5] = 6, A[6] = 2.
// 其中一种最长的 S[K]:
// S[0] = {A[0], A[5], A[6], A[2]} = {5, 6, 2, 0}

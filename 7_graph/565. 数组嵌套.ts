/**
 * @param {number[]} nums
 * @return {number}
 * 找到最大的集合S并返回其大小，其中 S[i] = {A[i], A[A[i]], A[A[A[i]]], ... }
 * @summary
 * 找到最长的环
 */
const arrayNesting = function (nums: number[]): number {
  let res = 0
  const n = nums.length
  const visited = Array<boolean>(n).fill(false)
  // const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  // for (let i = 0; i < n; i++) {
  //   adjList[i].push(nums[i])
  // }

  const dfs = (cur: number, steps: number) => {
    if (visited[cur]) return (res = Math.max(res, steps))
    visited[cur] = true
    dfs(nums[cur], steps + 1)
  }

  for (const num of nums) {
    if (visited[num]) continue
    dfs(num, 0)
  }

  return res
}

console.log(arrayNesting([5, 4, 0, 3, 1, 6, 2]))
// 输出: 4
// 解释:
// A[0] = 5, A[1] = 4, A[2] = 0, A[3] = 3, A[4] = 1, A[5] = 6, A[6] = 2.
// 其中一种最长的 S[K]:
// S[0] = {A[0], A[5], A[6], A[2]} = {5, 6, 2, 0}

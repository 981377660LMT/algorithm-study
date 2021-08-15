/**
 * @param {number[]} nums
 * @param {number} k
 * @return {boolean}
 * @summary 这题的visited数组可以用 位运算压缩 优化
 */
const canPartitionKSubsets = function (nums: number[], k: number): boolean {
  const total = nums.reduce((sum, num) => sum + num, 0)
  if (total % k !== 0) {
    return false
  }

  const target = total / k

  // 因为是要看每个位的数，所以不能用Set而要用数组记录每个数的状态
  const visited = Array<boolean>(nums.length).fill(false)
  function* bt(remain: number, start: number, curSum: number): Generator<boolean> {
    if (remain === 0) yield true
    // when a subset is found, we launch another search to find the
    // remaining subsets from the unvisited elements.
    if (curSum === target) yield* bt(remain - 1, 0, 0)
    for (let i = start; i < nums.length; i++) {
      const element = nums[i]
      if (visited[i]) continue
      visited[i] = true
      // launch a search to find other elements that will sum up to
      // the target with the current element.
      yield* bt(remain, i + 1, curSum + element)
      // reset to enable backtracking
      visited[i] = false
    }
  }

  return bt(k, 0, 0).next().value || false
}

console.log(canPartitionKSubsets([4, 3, 2, 3, 5, 2, 1], 4))

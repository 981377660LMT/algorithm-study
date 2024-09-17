/**
 * 遍历子集.
 * @complexity O(2^n), 2^27(1.3e8) => 1.1s.
 */
function powerset<T>(nums: ArrayLike<T>, f: (subset: readonly T[]) => void): void {
  const n = nums.length
  dfs(0, [])
  function dfs(index: number, path: T[]) {
    if (index === n) {
      f(path)
      return
    }
    dfs(index + 1, path)
    path.push(nums[index])
    dfs(index + 1, path)
    path.pop()
  }
}

if (require.main === module) {
  const nums = Array(20)
    .fill(0)
    .map((_, i) => i)

  console.time('enumerateSubset2')
  powerset(nums, subset => {})
  console.timeEnd('enumerateSubset2') // 1.1s

  powerset([1, 2, 3], console.log)
}

export { powerset }

import { useMinCostMaxFlow } from '../../../7_graph/网络流/4-费用流/useMinCostMaxFlow'

/**
 * 遍历子集.
 * @param copy 是否复制子集.
 * @complexity O(2^n), 2^27(1.3e8) => 1.1s.
 */
function enumerateSubset<T>(nums: ArrayLike<T>, callback: (subset: T[]) => void, copy = false): void {
  const n = nums.length
  dfs(0, [])
  function dfs(index: number, path: T[]) {
    if (index === n) {
      callback(copy ? path.slice() : path)
      return
    }
    dfs(index + 1, path)
    path.push(nums[index])
    dfs(index + 1, path)
    path.pop()
  }
}

if (require.main === module) {
  const nums = Array(30)
    .fill(0)
    .map((_, i) => i)

  console.time('enumerateSubset2')
  enumerateSubset(nums, subset => {})
  console.timeEnd('enumerateSubset2') // 1.1s
}

export { enumerateSubset }

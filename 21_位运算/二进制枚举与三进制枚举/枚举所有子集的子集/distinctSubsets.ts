/**
 * 不重复地枚举所有子集.
 */
function distinctSubsets(arr: number[], f: (subsetView: number[]) => void): void {
  arr = [...arr].sort((a, b) => a - b)

  const n = arr.length
  const subset: number[] = []

  function dfs(start: number): void {
    f(subset)
    for (let i = start; i < n; i++) {
      if (i > start && arr[i] === arr[i - 1]) {
        continue
      }
      subset.push(arr[i])
      dfs(i + 1)
      subset.pop()
    }
  }
  dfs(0)
}

export { distinctSubsets }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function subsetsWithDup(nums: number[]): number[][] {
    const res: number[][] = []
    distinctSubsets(nums, subsetView => {
      res.push(subsetView.slice())
    })
    return res
  }
}

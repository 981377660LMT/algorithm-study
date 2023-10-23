/**
 * 遍历所有的集合划分.
 */
function enumerateSetPartition<T>(arr: ArrayLike<T>, f: (groups: T[][]) => void): void {
  const n = arr.length
  const groups: T[][] = []
  const dfs = (pos: number): void => {
    if (pos === n) {
      f(groups)
      return
    }
    groups.push([arr[pos]])
    dfs(pos + 1)
    groups.pop()
    for (let i = 0; i < groups.length; ++i) {
      groups[i].push(arr[pos])
      dfs(pos + 1)
      groups[i].pop()
    }
  }
  dfs(0)
}

export { enumerateSetPartition }

if (require.main === module) {
  enumerateSetPartition([1, 2, 3], groups => console.log(groups))
  // [[1] [2] [3]]
  // [[1 3] [2]]
  // [[1] [2 3]]
  // [[1 2] [3]]
  // [[1 2 3]]
}

/**
 * 遍历所有的集合划分.
 * @param f 返回 true 时停止遍历.
 */
function enumerateSetPartition<T>(arr: ArrayLike<T>, f: (groups: T[][]) => boolean | void): void {
  const n = arr.length
  const groups: T[][] = []
  const dfs = (pos: number): boolean => {
    if (pos === n) {
      return !!f(groups)
    }
    groups.push([arr[pos]])
    dfs(pos + 1)
    groups.pop()
    for (let i = 0; i < groups.length; ++i) {
      groups[i].push(arr[pos])
      if (dfs(pos + 1)) return true
      groups[i].pop()
    }
    return false
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

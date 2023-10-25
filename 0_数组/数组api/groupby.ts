/**
 * 遍历连续相同元素的分组.
 * @alias groupBy
 * @example
 * ```ts
 * const list = [1, 1, 2, 3, 3, 4, 4, 5, 5, 5]
 * enumerateGroup(list, group => console.log(group)) // [1, 1], [2], [3, 3], [4, 4], [5, 5, 5]
 * ```
 */
function enumerateGroup<T>(arr: ArrayLike<T>, f: (group: T[], start: number, end: number) => boolean | void): void {
  const n = arr.length
  let ptr = 0
  while (ptr < n) {
    const leader = arr[ptr]
    const group = [leader]
    const start = ptr
    ptr++
    while (ptr < n && arr[ptr] === leader) {
      group.push(arr[ptr])
      ptr++
    }
    if (f(group, start, ptr)) return
  }
}

if (require.main === module) {
  enumerateGroup('abbcccdddd', (group, start, end) => {
    console.log(group, start, end)
  })
}

export { enumerateGroup }

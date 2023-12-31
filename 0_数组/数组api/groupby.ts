/* eslint-disable max-len */

/**
 * 遍历连续相同元素的分组(分组循环).
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

/**
 * 遍历连续key相同元素的分组.
 */
function enumerateGroupByKey<T>(
  arr: ArrayLike<T>,
  key: (index: number) => unknown,
  f: (group: T[], start: number, end: number) => boolean | void
): void {
  const n = arr.length
  let ptr = 0
  while (ptr < n) {
    const leader = key(ptr)
    const group = [arr[ptr]]
    const start = ptr
    ptr++
    while (ptr < n && key(ptr) === leader) {
      group.push(arr[ptr])
      ptr++
    }
    if (f(group, start, ptr)) return
  }
}

/**
 * 遍历分组(分组循环).
 * @param isDivider 判断当前元素是否为分组的分界点.如果返回true,则以当前元素为分界点,新建下一个分组.
 * @example
 * ```ts
 * // 每组最多3个元素
 * const list = [1, 1, 2, 3, 3, 4, 4]
 * enumerateGroupByDivider(list, (i, group) => group.length === 3, group => console.log(group)) // [1, 1, 2], [3, 3, 4], [4]
 * ```
 */
function enumerateGroupByDivider<T>(
  arr: ArrayLike<T>,
  isDivider: (index: number, curGroup: T[]) => boolean,
  f: (group: T[], start: number, end: number) => boolean | void
): void {
  const n = arr.length
  let ptr = 0
  while (ptr < n) {
    const leader = arr[ptr]
    const group = [leader]
    const start = ptr
    ptr++
    while (ptr < n && !isDivider(ptr, group)) {
      group.push(arr[ptr])
      ptr++
    }
    if (f(group, start, ptr)) return
  }
}

if (require.main === module) {
  console.log('enumerateGroup')
  enumerateGroup('abbcccdddd', (group, start, end) => {
    console.log(group, start, end)
  })

  console.log('enumerateGroupByKey')
  enumerateGroupByKey(
    'abbcccdddd',
    i => Math.floor(i / 2), // 按照下标某种规则分组
    (group, start, end) => {
      console.log(group, start, end)
    }
  )

  console.log('enumerateGroupByDivider')
  enumerateGroupByDivider(
    'abbcccdddd',
    (i, group) => group.length === 3, // 每组最多3个元素
    (group, start, end) => {
      console.log(group, start, end)
    }
  )
}

export { enumerateGroup, enumerateGroupByKey, enumerateGroupByDivider }

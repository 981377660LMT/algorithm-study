/* eslint-disable no-inner-declarations */
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
  isDivider: (elementIndex: number, curGroup: T[]) => boolean,
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

  // 2953. 统计完全子字符串
  // https://leetcode.cn/problems/count-complete-substrings/
  // 给你一个字符串 word 和一个整数 k 。
  // 如果 word 的一个子字符串 s 满足以下条件，我们称它是 完全字符串：
  // - s 中每个字符 恰好 出现 k 次。
  // - 相邻字符在字母表中的顺序 至多 相差 2 。也就是说，s 中两个相邻字符 c1 和 c2 ，它们在字母表中的位置相差 至多 为 2 。
  // 请你返回 word 中 完全 子字符串的数目。
  // !循环分组，在每个组内使用滑动窗口检查.
  function countCompleteSubstrings(word: string, k: number): number {
    const n = word.length
    const nums = new Uint8Array(n)
    for (let i = 0; i < n; i++) nums[i] = word.charCodeAt(i) - 97

    const groups: number[][] = []
    enumerateGroupByDivider(
      nums,
      index => Math.abs(nums[index] - nums[index - 1]) > 2,
      g => {
        groups.push(g)
      }
    )

    /**
     * 每个字符恰好出现k次的子字符串数目.
     */
    const solve = (group: number[]): number => {}

    let res = 0
    groups.forEach(g => {
      res += solve(g)
    })
    return res
  }
}

export { enumerateGroup, enumerateGroupByKey, enumerateGroupByDivider }

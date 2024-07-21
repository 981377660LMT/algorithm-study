/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

/**
 * 遍历连续相同元素的分组(分组循环).
 * @alias groupBy
 * @example
 * ```ts
 * const list = [1, 1, 2, 3, 3, 4, 4, 5, 5, 5]
 * enumerateGroup(list, (start, end) => console.log(list.slice(start, end))) // [1, 1], [2], [3, 3], [4, 4], [5, 5, 5]
 * ```
 */
function enumerateGroup(
  arr: ArrayLike<unknown>,
  f: (start: number, end: number) => boolean | void
): void {
  const n = arr.length
  let end = 0
  while (end < n) {
    const start = end
    const leader = arr[end]
    end++
    while (end < n && arr[end] === leader) end++
    if (f(start, end)) return
  }
}

/**
 * 遍历连续key相同元素的分组.
 *
 * @alias groupByKey
 */
function enumerateGroupByKey(
  arr: ArrayLike<unknown>,
  key: (index: number) => unknown,
  f: (start: number, end: number) => boolean | void
): void {
  const n = arr.length
  let end = 0
  while (end < n) {
    const start = end
    const leader = key(end)
    end++
    while (end < n && key(end) === leader) end++
    if (f(start, end)) return
  }
}

/**
 * 遍历分组(分组循环).
 *
 * @alias groupWhile
 * @param predicate 返回`true`表示`[left, curRight]`内的元素分为一组.
 * @param skipFalsySingleValueGroup 是否跳过 {@link predicate} 为`false`的单个元素的分组，默认为`false`.
 * @example
 * ```ts
 * // 每组最多3个元素
 * const list = [1, 1, 2, 3, 3, 4, 4]
 * enumerateGroupByDivider(list, (left, curRight) => curRight-left+1 <= 3, (start, end) => console.log(list.slice(start, end))) // [1, 1, 2], [3, 3, 4], [4]
 * ```
 */
function enumerateGroupByGroupWhile(
  n: number,
  predicate: (left: number, curRight: number) => boolean,
  consumer: (start: number, end: number) => void,
  skipFalsySingleValueGroup = false
): void {
  let end = 0
  while (end < n) {
    const start = end
    while (end < n && predicate(start, end)) end++
    const isFalsySingleValueGroup = start === end
    if (isFalsySingleValueGroup) {
      end++
      if (skipFalsySingleValueGroup) continue
    }
    consumer(start, end)
  }
}

if (require.main === module) {
  console.log('enumerateGroup')
  enumerateGroup('abbcccdddd', (start, end) => {
    console.log(start, end)
  })

  console.log('enumerateGroupByKey')
  enumerateGroupByKey(
    'abbcccdddd',
    i => Math.floor(i / 2), // 按照下标某种规则分组
    (start, end) => {
      console.log(start, end)
    }
  )

  console.log('enumerateGroupByDivider')
  const ss = 'abbcccdddd'
  enumerateGroupByGroupWhile(
    ss.length,
    (left, right) => right - left + 1 <= 3, // 每组最多3个元素
    (start, end) => {
      console.log(ss.slice(start, end))
    }
  )

  /**
   * 每组最多k个元素的分组.
   */
  function maxPartitionsAfterOperations(s: string, k: number): number {
    const n = s.length
    const ords = new Uint8Array(n)
    for (let i = 0; i < s.length; i++) ords[i] = s[i].codePointAt(0)! - 97
    let res = 0

    let ptr = 0
    while (ptr < n) {
      // !当前分组的第一个元素
      let visited = 1 << ords[ptr]
      let visitedCount = 1
      ptr++

      // !能否继续向后扩展
      while (ptr < n && visitedCount + (1 ^ ((visited >>> ords[ptr]) & 1)) <= k) {
        visitedCount += ((visited >>> ords[ptr]) & 1) ^ 1
        visited |= 1 << ords[ptr]
        ptr++
      }

      // !当前分组结束
      res++
    }

    return res
  }

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

    const groups: ArrayLike<number>[] = []
    enumerateGroupByGroupWhile(
      nums.length,
      (left, right) => left === right || Math.abs(nums[right] - nums[right - 1]) <= 2,
      (start, end) => {
        groups.push(nums.subarray(start, end))
      }
    )

    /**
     * 每个字符恰好出现k次的子字符串数目.
     */
    const solve = (group: ArrayLike<number>): number => {
      let res = 0

      for (let i = 1; i <= 26; i++) {
        const windowLength = i * k
        if (windowLength > group.length) break

        const counter = Array(26).fill(0)
        for (let right = 0; right < group.length; right++) {
          counter[group[right]]++
          if (right >= windowLength) counter[group[right - windowLength]]--
          if (right >= windowLength - 1) {
            let ok = true
            for (let j = 0; j < 26; j++) {
              if (counter[j] !== 0 && counter[j] !== k) {
                ok = false
                break
              }
            }
            if (ok) res++
          }
        }
      }

      return res
    }

    let res = 0
    groups.forEach(g => {
      res += solve(g)
    })
    return res
  }
}

export { enumerateGroup, enumerateGroupByKey, enumerateGroupByGroupWhile }

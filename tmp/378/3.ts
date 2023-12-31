export {}

const INF = 2e9 // !超过int32使用2e15
// 给你一个仅由小写英文字母组成的字符串 s 。

// 如果一个字符串仅由单一字符组成，那么它被称为 特殊 字符串。例如，字符串 "abc" 不是特殊字符串，而字符串 "ddd"、"zz" 和 "f" 是特殊字符串。

// 返回在 s 中出现 至少三次 的 最长特殊子字符串 的长度，如果不存在出现至少三次的特殊子字符串，则返回 -1 。

// 子字符串 是字符串中的一个连续 非空 字符序列。

/**
 * 遍历连续相同元素的分组.
 * @alias groupBy
 * @example
 * ```ts
 * const list = [1, 1, 2, 3, 3, 4, 4, 5, 5, 5]
 * enumerateGroup(list, group => console.log(group)) // [1, 1], [2], [3, 3], [4, 4], [5, 5, 5]
 * ```
 */
function enumerateGroup<T>(arr: ArrayLike<T>, f: (group: T[], start: number, end: number) => void): void {
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
    f(group, start, ptr)
  }
}

function maximumLength(s: string): number {
  const groups: [number, number, number][] = []
  enumerateGroup(s, (group, start, end) => {
    groups.push([group[0].codePointAt(0)! - 97, start, end])
  })

  let left = 1
  let right = s.length
  let ok = false
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) {
      left = mid + 1
      ok = true
    } else {
      right = mid - 1
    }
  }

  return ok ? right : -1

  function check(mid: number): boolean {
    const counter = Array(26).fill(0)
    for (let i = 0; i < groups.length; i++) {
      const { 0: ch, 1: start, 2: end } = groups[i]
      const len = end - start
      if (len >= mid) {
        counter[ch] += len - mid + 1
        if (counter[ch] >= 3) {
          return true
        }
      }
    }
    return false
  }
}

// "abcccccdddd"
console.log(maximumLength('abcccccdddd'))

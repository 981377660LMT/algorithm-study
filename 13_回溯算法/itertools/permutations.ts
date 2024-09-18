/* eslint-disable no-lone-blocks */
/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */

/**
 * 按照字典序遍历所有大小为`r`的排列.
 * @param n 排列总元素个数.
 * @param r 选取的元素个数.
 * @param f 回调函数, 用于处理每个排列的结果.返回`true`时停止遍历.
 * @example
 * ```ts
 * enumeratePermutations(3, 2, perm => {
 *  console.log(perm)
 * })
 * // [ 0, 1 ]
 * // [ 0, 2 ]
 * // [ 1, 0 ]
 * // [ 1, 2 ]
 * // [ 2, 0 ]
 * // [ 2, 1 ]
 * ```
 * @complexity 11!(4e7) => 2.2ms
 */
function enumeratePermutations(
  n: number,
  r: number,
  f: (indicesView: readonly number[]) => boolean | void
): void {
  const visited = new Uint8Array(n)
  const dfs = (path: number[]): boolean => {
    if (path.length === r) {
      return !!f(path)
    }
    for (let i = 0; i < n; i++) {
      if (visited[i]) continue
      visited[i] = 1
      path.push(i)
      if (dfs(path)) return true
      visited[i] = 0
      path.pop()
    }
    return false
  }

  dfs([])
}

/**
 * 无序遍历全排列.
 * @param n 全排列的长度.
 * @param f 回调函数, 用于处理每个排列的结果.返回`true`时停止遍历.
 * @example
 * ```ts
 * enumeratePermutationsAll(3, perm => {
 *  console.log(perm)
 * })
 * // [ 0, 1, 2 ]
 * // [ 0, 2, 1 ]
 * // [ 1, 0, 2 ]
 * // [ 1, 2, 0 ]
 * // [ 2, 1, 0 ]
 * // [ 2, 0, 1 ]
 * ```
 * @complexity 11!(4e7) => 570ms
 */
function enumeratePermutationsAll(
  n: number,
  f: (indicesView: readonly number[]) => boolean | void
): void {
  const dfs = (a: number[], i: number): boolean => {
    if (i === a.length) {
      return !!f(a)
    }
    dfs(a, i + 1)
    for (let j = i + 1; j < a.length; j++) {
      const tmp = a[i]
      a[i] = a[j]
      a[j] = tmp
      if (dfs(a, i + 1)) return true
      a[j] = a[i]
      a[i] = tmp
    }
    return false
  }
  const ids = Array.from({ length: n }, (_, i) => i)
  dfs(ids, 0)
}

if (require.main === module) {
  {
    console.time('enumeratePermutations')
    let count = 0
    enumeratePermutations(11, 11, indices => {
      count++
    })
    console.timeEnd('enumeratePermutations') // 2.2s
    console.log(count) // 39916800
  }

  {
    console.time('enumeratePermutationsAll')
    let count = 0
    enumeratePermutationsAll(11, indices => {
      count++
    })
    console.timeEnd('enumeratePermutationsAll') // 600ms
    console.log(count) // 39916800
  }
}

export { enumeratePermutations, enumeratePermutationsAll }

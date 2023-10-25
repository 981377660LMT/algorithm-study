/**
 * 遍历所有大小为`r`的排列.
 * @param arr 待遍历的数组.
 * @param r 排列大小.
 * @param cb 回调函数, 用于处理每个排列的结果.返回`true`时停止遍历.
 * @param copy 是否浅拷贝每个排列的结果, 默认为`false`.
 * @example
 * ```ts
 * enumeratePermutations([1, 2, 3], 2, perm => {
 *  console.log(perm)
 * })
 * // [ 1, 2 ]
 * // [ 1, 3 ]
 * // [ 2, 1 ]
 * // [ 2, 3 ]
 * // [ 3, 1 ]
 * // [ 3, 2 ]
 * ```
 * @complexity 11!(4e7) => 2.049s
 */
function enumeratePermutations<P>(arr: ArrayLike<P>, r: number, cb: (perm: readonly P[]) => boolean | void, copy = false): void {
  const n = arr.length
  bt([], new Uint8Array(n))

  function bt(path: P[], visited: Uint8Array): boolean {
    if (path.length === r) {
      return !!cb(copy ? path.slice() : path)
    }

    for (let i = 0; i < n; i++) {
      if (visited[i]) continue
      visited[i] = 1
      path.push(arr[i])
      if (bt(path, visited)) return true
      path.pop()
      visited[i] = 0
    }

    return false
  }
}

/**
 * 模拟python的`itertools.permutations`.
 * @complexity 11!(4e7) => 19.233s
 * @deprecated
 */
function* permutations<P>(arr: ArrayLike<P>, r = arr.length): Generator<P[]> {
  const n = arr.length
  yield* bt([], new Uint8Array(n))

  function* bt(path: P[], visited: Uint8Array): Generator<P[]> {
    if (path.length === r) {
      yield path.slice()
      return
    }

    for (let i = 0; i < n; i++) {
      if (visited[i]) continue
      visited[i] = 1
      path.push(arr[i])
      yield* bt(path, visited)
      path.pop()
      visited[i] = 0
    }
  }
}

if (require.main === module) {
  const n = 11
  const arr = Array.from({ length: n }, (_, i) => i)
  console.time('enumeratePermutations')
  let count = 0
  enumeratePermutations(arr, arr.length, () => {
    count++
  })
  console.timeEnd('enumeratePermutations') // 1.995s
  console.log(count) // 39916800
}

export { enumeratePermutations }

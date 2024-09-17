/**
 * 遍历子集.
 *
 * @complexity O(2^n), 2^27(1.3e8) => 1.3s.
 */
function powerset<T>(nums: ArrayLike<T>, f: (subset: readonly T[]) => boolean | void): void {
  const n = nums.length
  dfs(0, [])
  function dfs(index: number, path: T[]): boolean {
    if (index === n) return !!f(path)
    if (dfs(index + 1, path)) return true
    path.push(nums[index])
    if (dfs(index + 1, path)) return true
    path.pop()
    return false
  }
}

/**
 * 遍历数组所有的分割方案，按照分割点将数组分割成若干段.
 *
 * @complexity O(2^n), 2^27(1.3e8) => 630ms.
 * @example
 * ```ts
 * partitions(3, splits => {
 *   for (let i = 0; i < splits.length - 1; i++) {
 *     const start = splits[i]
 *     const end = splits[i + 1]
 *     console.log(arr.slice(start, end))
 *   }
 * })
 * ```
 */
function partitions(n: number, f: (splits: readonly number[]) => boolean | void): void {
  if (!n) return
  dfs(0, [0])
  function dfs(index: number, path: number[]): boolean {
    if (index === n - 1) {
      path.push(n)
      const stop = f(path)
      path.pop()
      return !!stop
    }
    if (dfs(index + 1, path)) return true
    path.push(index + 1)
    if (dfs(index + 1, path)) return true
    path.pop()
    return false
  }
}

/** @todo */
function distinctPermutations<T>(
  nums: T[],
  f: (permutation: readonly T[]) => boolean | void,
  r = nums.length
): void {
  const n = nums.length
  const used = Array(n).fill(false)
  const path: T[] = []
  dfs()
  function dfs() {
    if (path.length === n) {
      if (f(path)) return
    }
    for (let i = 0; i < n; i++) {
      if (used[i]) continue
      used[i] = true
      path.push(nums[i])
      dfs()
      path.pop()
      used[i] = false
    }
  }
}

/** @todo */
function distinctCombinations<T>(
  nums: T[],
  f: (combination: readonly T[]) => boolean | void,
  r: number
): void {
  const n = nums.length
  const path: T[] = []
  dfs(0, 0)
  function dfs(index: number, count: number) {
    if (count === r) {
      if (f(path)) return
    }
    for (let i = index; i < n; i++) {
      path.push(nums[i])
      dfs(i + 1, count + 1)
      path.pop()
    }
  }
}

export { powerset, partitions, distinctPermutations, distinctCombinations }

if (require.main === module) {
  console.time('powerset')
  powerset('a'.repeat(27), subset => {})
  console.timeEnd('powerset')

  console.time('partitions')
  partitions(27, splits => {})
  console.timeEnd('partitions')

  const arr = 'abc'
  partitions(arr.length, splits => {
    const cur = []
    for (let i = 0; i < splits.length - 1; i++) {
      const start = splits[i]
      const end = splits[i + 1]
      cur.push(arr.slice(start, end))
    }
    console.log(cur)
  })
}

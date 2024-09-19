/* eslint-disable no-lone-blocks */
/* eslint-disable @typescript-eslint/no-empty-function */
/* eslint-disable no-inner-declarations */

/**
 * 遍历子集.
 *
 * @complexity O(2^n), 2^27(1.3e8) => 1s.
 */
function powerset(n: number, f: (subset: readonly number[]) => boolean | void): void {
  const path: number[] = []
  dfs(0)
  function dfs(index: number): boolean {
    if (index === n) return !!f(path)
    if (dfs(index + 1)) return true
    path.push(index)
    if (dfs(index + 1)) return true
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
  const path: number[] = [0]
  dfs(0)
  function dfs(index: number): boolean {
    if (index === n - 1) {
      path.push(n)
      const stop = f(path)
      path.pop()
      return !!stop
    }
    if (dfs(index + 1)) return true
    path.push(index + 1)
    if (dfs(index + 1)) return true
    path.pop()
    return false
  }
}

/**
 * 将 {@link n} 个元素的集合分成 {@link k} 个部分，不允许为空.
 *
 * @complexity n=13,k=5 => 1e7、250ms.
 * @note 对每个元素，是放入之前的部分还是新建一个部分.
 */
function setPartitions(
  n: number,
  k: number,
  f: (parts: readonly number[][]) => boolean | void
): void {
  if (k < 1) throw new Error("Can't partition in a negative or zero number of groups")
  if (k > n) return

  const parts: number[][] = Array(k)
  for (let i = 0; i < k; i++) parts[i] = []
  dfs(0, 0)
  function dfs(index: number, count: number): boolean {
    if (index === n) {
      if (count === k) {
        return !!f(parts)
      }
      return false
    }
    for (let i = 0; i < count; i++) {
      parts[i].push(index)
      if (dfs(index + 1, count)) return true
      parts[i].pop()
    }
    if (count < k) {
      parts[count].push(index)
      if (dfs(index + 1, count + 1)) return true
      parts[count].pop()
    }
    return false
  }
}

/**
 * 将 {@link n} 个元素的集合分成任意个部分.
 *
 * @param f 分成 `k` 个部分，parts.slice(0, k) 为分割结果. 返回 `true` 则停止搜索.
 * @complexity n=13 => 2e7、250ms.
 */
function setPartitionsAll(
  n: number,
  f: (parts: readonly number[][], k: number) => boolean | void
): void {
  const parts: number[][] = Array(n)
  for (let i = 0; i < n; i++) parts[i] = []
  dfs(0, 0)
  function dfs(index: number, count: number): boolean {
    if (index === n) {
      return !!f(parts, count)
    }
    for (let i = 0; i < count; i++) {
      parts[i].push(index)
      if (dfs(index + 1, count)) return true
      parts[i].pop()
    }
    parts[count].push(index)
    if (dfs(index + 1, count + 1)) return true
    parts[count].pop()
    return false
  }
}

/**
 * 遍历无重复排列.
 *
 * @complexity 11!(4e7) => 200ms.
 */
function distinctPermutations<T extends number | string>(
  arr: ArrayLike<T>,
  r: number,
  f: (perm: readonly T[]) => boolean | void
): void {
  if (!arr.length || r < 1 || r > arr.length) return

  const copy = Array.from(arr)
  copy.sort((a, b) => {
    if (a < b) return -1
    if (a > b) return 1
    return 0
  })

  if (r === arr.length) {
    while (true) {
      if (f(copy)) break
      if (!nextPermutation(copy)) break
    }
    return
  }

  const head = copy.slice(0, r)
  const tail = copy.slice(r)
  while (tail.length < arr.length) tail.push(arr[0])
  const tailLen = arr.length - r
  while (true) {
    if (f(head)) return
    let pivot = tail[tailLen - 1]
    let i = r - 1
    let found = false
    for (; ~i; i--) {
      if (head[i] < pivot) {
        found = true
        break
      }
      pivot = head[i]
    }
    if (!found) return

    found = false
    for (let j = 0; j < tailLen; j++) {
      if (tail[j] > head[i]) {
        const tmp = head[i]
        head[i] = tail[j]
        tail[j] = tmp
        found = true
        break
      }
    }
    if (!found) {
      for (let j = r - 1; ~j; j--) {
        if (head[j] > head[i]) {
          const tmp = head[i]
          head[i] = head[j]
          head[j] = tmp
          break
        }
      }
    }

    // reversr head[i + 1:] and swap it with tail[:r - (i + 1)]
    for (let j = tailLen, k = r - 1; k > i; k--, j++) {
      tail[j] = head[k]
    }
    for (let j = i + 1, k = 0; j < r; j++, k++) {
      head[j] = tail[k]
    }
    for (let j = 0, k = r - i - 1; j < tailLen; j++, k++) {
      tail[j] = tail[k]
    }
  }

  /** full. */
  function nextPermutation(arr: T[]): boolean {
    if (!arr.length) return false
    let left = arr.length - 1
    while (left > 0 && arr[left - 1] >= arr[left]) left--
    if (!left) return false
    const last = left - 1
    let right = arr.length - 1
    while (arr[right] <= arr[last]) right--
    const tmp = arr[last]
    arr[last] = arr[right]
    arr[right] = tmp
    for (let i = last + 1, j = arr.length - 1; i < j; i++, j--) {
      const tmp = arr[i]
      arr[i] = arr[j]
      arr[j] = tmp
    }
    return true
  }
}

/**
 * 遍历无重复组合.
 *
 * @complexity C(30,10)(3e7) => 400ms.
 * @example
 * ```ts
 * distinctCombinations([0, 0, 1, 0], 2, console.log) // [0, 0], [0, 1], [1, 0]
 * ```
 */
function distinctCombinations<T extends number | string>(
  arr: ArrayLike<T>,
  r: number,
  f: (comb: readonly T[]) => boolean | void
): void {
  if (!arr.length || r < 1 || r > arr.length) return

  const uniqueEverSeen: number[][] = Array(arr.length)
  for (let i = 0; i < arr.length; i++) {
    const indexes: number[] = []
    const visited = new Set<T>()
    for (let j = i; j < arr.length; j++) {
      if (visited.has(arr[j])) continue
      visited.add(arr[j])
      indexes.push(j)
    }
    uniqueEverSeen[i] = indexes
  }

  const path: T[] = []
  function dfs(index: number, count: number): boolean {
    if (count === r) {
      return !!f(path)
    }
    if (index === arr.length) return false
    const indexes = uniqueEverSeen[index]
    for (let i = 0; i < indexes.length; i++) {
      const curIndex = indexes[i]
      path.push(arr[indexes[i]])
      if (dfs(curIndex + 1, count + 1)) return true
      path.pop()
    }
    return false
  }
  dfs(0, 0)
}

export {
  powerset,
  partitions,
  setPartitions,
  setPartitionsAll,
  distinctPermutations,
  distinctCombinations
}

if (require.main === module) {
  // 46. 全排列
  // https://leetcode.cn/problems/permutations/
  function permute(nums: number[]): number[][] {
    const res: number[][] = []
    distinctPermutations(nums, nums.length, perm => {
      res.push(perm.slice())
    })
    return res
  }

  // 78. 子集
  // https://leetcode.cn/problems/subsets/description/
  function subsets(nums: number[]): number[][] {
    const res: number[][] = []
    powerset(nums.length, subset => {
      const cur = Array<number>(subset.length)
      for (let i = 0; i < subset.length; i++) {
        cur[i] = nums[subset[i]]
      }
      res.push(cur)
    })
    return res
  }

  // 2698. 求一个整数的惩罚数
  // https://leetcode.cn/problems/find-the-punishment-number-of-an-integer/
  function punishmentNumber(n: number): number {
    /** v * v 的十进制表示的字符串可以分割成若干连续子字符串，且这些子字符串对应的整数值之和等于 v 。 */
    const toDigit = (nums: number[], start: number, end: number): number => {
      let res = 0
      for (let i = start; i < end; i++) res = res * 10 + nums[i]
      return res
    }

    const check = (v: number): boolean => {
      const nums = String(v * v)
        .split('')
        .map(Number)

      let ok = false
      partitions(nums.length, splits => {
        let sum = 0
        for (let i = 0; i < splits.length - 1; i++) {
          sum += toDigit(nums, splits[i], splits[i + 1])
        }

        if (sum === v) {
          ok = true
          return true
        }
      })

      return ok
    }

    let res = 0
    for (let i = 1; i <= n; i++) {
      if (check(i)) res += i * i
    }
    return res
  }

  function test1() {
    console.time('powerset')
    powerset(27, subset => {})
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

    let count = 0
    console.time('setPartitions')
    setPartitions(13, 6, parts => {
      count++
    })
    console.timeEnd('setPartitions')
    console.log(count)

    count = 0
    console.time('setPartitionsAll')
    setPartitionsAll(3, (parts, k) => {
      console.log(parts.slice(0, k))
    })
    console.log(count)
    console.timeEnd('setPartitionsAll')
  }

  console.time('distinctPermutations')
  const arr = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  let count = 0
  distinctPermutations(arr, arr.length - 1, perm => {
    count++
  })
  console.log(count)
  console.timeEnd('distinctPermutations')

  console.time('distinctCombinations')
  count = 0
  distinctCombinations(
    Array.from({ length: 30 }, (_, i) => i),
    10,
    comb => {
      count++
    }
  )
  console.log(count)
  console.timeEnd('distinctCombinations')
  // {
  //   let count = 0
  //   console.time('setPartitions')
  //   setPartitions(13, 6, parts => {
  //     count++
  //   })
  //   console.log(count)
  //   console.timeEnd('setPartitions')
  // }
  setPartitions(3, 2, parts => {
    console.log(parts)
  })
  console.log('done')
}

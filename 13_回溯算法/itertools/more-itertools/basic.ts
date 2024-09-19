/**
 * 将 {@link n} 个元素分成大小为 {@link size} 的批次，最后一个批次可能小于 {@link size}.
 * 适用于分片处理大量数据的场景.
 *
 * @returns 返回两个组的大小和组的个数.
 * @alias chunked
 * @example
 * ```ts
 * batched(10, 3, console.log) // [0, 3), [3, 6), [6, 9), [9, 10)
 * ```
 */
function batched(
  n: number,
  size: number,
  f?: (start: number, end: number) => void
): [size1: number, count1: number, size2: number, count2: number] {
  if (f) {
    for (let i = 0; i < n; i += size) {
      f(i, Math.min(i + size, n))
    }
  }
  const size1 = size
  const count1 = Math.floor(n / size)
  const size2 = n % size
  const count2 = +(size2 > 0)
  return [size1, count1, size2, count2]
}

/**
 * 将 {@link n} 个元素分成 {@link groupCount} 个组.
 * 每个组的大小尽可能均等分配，使得每个组的大小差距不超过1.
 * 适用于将任务分配给多个工作线程的场景.
 *
 * @returns 返回两个组的大小和组的个数.
 * @example
 * ```ts
 * distribute(10, 3, console.log) // [0, 4), [4, 7), [7, 10)
 * ```
 */
function distribute(
  n: number,
  groupCount: number,
  f?: (start: number, end: number) => void
): [size1: number, count1: number, size2: number, count2: number] {
  const q = Math.floor(n / groupCount)
  const r = n % groupCount
  if (f) {
    let start = 0
    for (let i = 0; i < groupCount; i++) {
      const end = start + q + +(i < r)
      f(start, end)
      start = end
    }
  }
  const size1 = q + 1
  const count1 = r
  const size2 = q
  const count2 = groupCount - r
  return [size1, count1, size2, count2]
}

export { batched, batched as chunked, distribute }

if (require.main === module) {
  const arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  batched(arr.length, 3, (start, end) => {
    console.log(arr.slice(start, end))
  })
  console.log('---')

  distribute(arr.length, 3, (start, end) => {
    console.log(arr.slice(start, end))
  })
}

/**
 * 在区间 `[start, end)` 中查找使 `f(i)` 为 `true` 的最小值 `i`.
 * @see https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/sort/search.go;l=58
 */
function sortSearch(
  start: number,
  end: number,
  f: (mid: number) => boolean
): [i: number, found: boolean] {
  let i = start
  let j = end
  while (i < j) {
    const h = (i + j) >>> 1
    if (!f(h)) {
      i = h + 1
    } else {
      j = h
    }
  }
  return [i, i < end && f(i)]
}

/**
 * @alias bisectLeft
 */
function sortSearchInts(arr: ArrayLike<number>, target: number): number {
  return sortSearch(0, arr.length, i => arr[i] >= target)[0]
}

export { sortSearch, sortSearchInts }

if (require.main === module) {
  console.log(sortSearchInts([1, 2, 3, 3], 2))
}

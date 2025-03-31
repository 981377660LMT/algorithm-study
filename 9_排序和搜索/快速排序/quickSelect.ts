/**
 * Rearranges items so that all items in the [left, k] are the smallest.
 * The k-th element will have the (k - left + 1)-th smallest value in [left, right].
 *
 * @param arr the array to partially sort (in place)
 * @param k middle index for partial sorting (as defined above)
 * @param left left index of the range to sort
 * @param right right index
 * @param compare compare function
 *
 * @see https://github.com/mourner/quickselect/blob/main/index.js
 */
function quickselect(
  arr: any[],
  k: number,
  left = 0,
  right = arr.length - 1,
  compare = defaultCompare
) {
  while (right > left) {
    if (right - left > 600) {
      const n = right - left + 1
      const m = k - left + 1
      const z = Math.log(n)
      const s = 0.5 * Math.exp((2 * z) / 3)
      const sd = 0.5 * Math.sqrt((z * s * (n - s)) / n) * (m - n / 2 < 0 ? -1 : 1)
      const newLeft = Math.max(left, Math.floor(k - (m * s) / n + sd))
      const newRight = Math.min(right, Math.floor(k + ((n - m) * s) / n + sd))
      quickselect(arr, k, newLeft, newRight, compare)
    }

    const t = arr[k]
    let i = left
    let j = right

    swap(arr, left, k)
    if (compare(arr[right], t) > 0) swap(arr, left, right)

    while (i < j) {
      swap(arr, i, j)
      i++
      j--
      while (compare(arr[i], t) < 0) i++
      while (compare(arr[j], t) > 0) j--
    }

    if (compare(arr[left], t) === 0) swap(arr, left, j)
    else {
      j++
      swap(arr, j, right)
    }

    if (j <= k) left = j + 1
    if (k <= j) right = j - 1
  }
}

function swap(arr: any[], i: number, j: number) {
  const tmp = arr[i]
  arr[i] = arr[j]
  arr[j] = tmp
}

function defaultCompare(a: number, b: number) {
  // eslint-disable-next-line no-nested-ternary
  return a < b ? -1 : a > b ? 1 : 0
}

export { quickselect, quickselect as nthElement }

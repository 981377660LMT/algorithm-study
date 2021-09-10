/**
 * @param {number[]} arr
 * @return {number}
 */
var lenLongestFibSubseq = function (arr: number[]): number {
  const len = arr.length
  if (len < 3) return 0
  const set = new Set(arr)
  const map = new Map<string, number>()

  for (let r = 2; r < len; r++) {
    for (let m = 1; m < r; m++) {
      const lVal = arr[r] - arr[m]
      if (set.has(lVal) && lVal < arr[m]) {
        const key = `${arr[m]}#${arr[r]}`
        const preKey = `${lVal}#${arr[m]}`
        map.set(key, (map.get(preKey) || 2) + 1)
      }
    }
  }

  return Math.max.apply(null, [...map.values(), 0])
}

console.log(lenLongestFibSubseq([1, 2, 3, 4, 5, 6, 7, 8]))

export {}

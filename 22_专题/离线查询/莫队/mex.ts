/**
 * 获取数组的 mex 值.
 */
function mex(arr: ArrayLike<number>, mexStart = 0): number {
  const n = arr.length
  const counter = new Uint32Array(mexStart + n + 1)
  for (let i = 0; i < n; i++) {
    const num = arr[i]
    if (num < mexStart || num > mexStart + n) {
      continue
    }
    counter[num - mexStart]++
  }
  let mex = mexStart
  while (counter[mex - mexStart]) {
    mex++
  }
  return mex
}

export { mex }

if (require.main === module) {
  console.log(mex([0, 1, 2, 3, 4, 5, 6, 7, 8, 9]))
}

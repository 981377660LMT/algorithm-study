/**
 * 前后缀分解模板.
 */

function solve(nums: number[]): number {
  const makeDp = (arr: number[]): number[] => {
    const m = arr.length
    const dp: number[] = Array(m + 1).fill(0)
    for (let i = 1; i <= m; i++) {
      const cur = arr[i - 1]
      // ...
    }
    return dp
  }

  const preDp = makeDp(nums)
  const sufDp = makeDp(nums.slice().reverse()).reverse()
  let res = 0
  for (let i = 0; i < nums.length + 1; i++) {
    res += preDp[i] * sufDp[i] // [0,i) x [i,n)
  }
  return res
}

export {}

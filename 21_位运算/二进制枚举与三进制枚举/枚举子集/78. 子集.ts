/**
 *
 * @param nums
 * @returns 2^n 个子集
 * @description 2^n*n时间复杂度
 */
function subsets<T>(nums: T[]): T[][] {
  const n = nums.length
  const res: T[][] = []

  for (let state = 0; state < 1 << n; state++) {
    const cands: T[] = []
    for (let j = 0; j < nums.length; j++) {
      if (state & (1 << j)) cands.push(nums[j])
    }

    res.push(cands)
  }

  return res
}

// 2^n 时间复杂度
function subsets2<T>(nums: T[]): T[][] {
  const n = nums.length
  const res: T[][] = []
  dfs(0, [])
  return res

  function dfs(index: number, path: T[]): void {
    if (index === n) {
      res.push(path.slice())
      return
    }

    dfs(index + 1, path)

    path.push(nums[index])
    dfs(index + 1, path)
    path.pop()
  }
}

function GosperHack(n: number, k: number) {
  let x = (1 << k) - 1
  const limit = 1 << n
  while (x < limit) {
    console.log(x.toString(2).padStart(5, '0'))
    const lowbit = x & -x
    const r = x + lowbit
    // xor
    x = r | (((x ^ r) >> 2) / lowbit)
  }
}

export {}

if (require.main === module) {
  console.log(subsets([1, 2, 3]))
  console.log(subsets2([1, 2, 3]))
}

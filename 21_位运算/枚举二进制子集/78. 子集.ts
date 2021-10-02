function subsets(nums: number[]): number[][] {
  const n = 1 << nums.length
  const res: number[][] = []

  for (let state = 0; state < n; state++) {
    const tmp: number[] = []
    for (let j = 0; j < nums.length; j++) {
      if (state & (1 << j)) tmp.push(nums[j])
    }
    res.push(tmp)
  }

  return res
}

console.log(subsets([1, 2, 3]))

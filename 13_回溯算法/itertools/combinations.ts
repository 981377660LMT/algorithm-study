function combinations(nums: number[], k: number) {
  const res: number[][] = []

  const bt = (cur: number, path: number[]) => {
    if (path.length === k) return res.push(path.slice())

    for (let i = cur; i < nums.length; i++) {
      path.push(nums[i])
      bt(i + 1, path)
      path.pop()
    }
  }

  bt(0, [])

  return res
}

if (require.main === module) {
  console.log(combinations([1, 1, 3, 4], 2))
}

export { combinations }

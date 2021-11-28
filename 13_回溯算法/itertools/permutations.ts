function combinations(nums: number[]) {
  const res: number[][] = []

  const bt = (path: number[], visited: boolean[]) => {
    if (path.length === nums.length) return res.push(path.slice())

    for (let i = 0; i < nums.length; i++) {
      if (visited[i]) continue
      visited[i] = true
      path.push(nums[i])
      bt(path, visited)
      path.pop()
      visited[i] = false
    }
  }
  bt([], [])

  return res
}

if (require.main === module) {
  console.log(combinations([1, 1, 0]))
}

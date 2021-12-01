function combinationsWithReplacement(nums: number[], k: number) {
  const res: number[][] = []

  const bt = (cur: number, path: number[]) => {
    if (path.length === k) return res.push(path.slice())

    for (let i = cur; i < nums.length; i++) {
      path.push(nums[i])
      bt(i, path) // 唯一的区别在此：可取重复的元素
      path.pop()
    }
  }

  bt(0, [])

  return res
}

if (require.main === module) {
  console.log(combinationsWithReplacement([1, 1, 3, 4], 2))
}

export {}

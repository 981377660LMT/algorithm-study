// 从给定的无需、不重复的数组A中，取出N个数，使其相加和为M
const nSum = (nums: number[], n: number, m: number) => {
  const len = nums.length
  const res: number[][] = []

  const bt = (index: number, path: number[], sum: number) => {
    if (path.length === n) {
      if (sum === m) res.push(path.slice())
      return
    }

    for (let i = index; i < len; i++) {
      path.push(nums[i])
      bt(i + 1, path, sum + nums[i])
      path.pop()
    }
  }

  bt(0, [], 0)
  return res
}

console.log(nSum([1, 4, 7, 11, 9, 8, 10, 6], 3, 27))

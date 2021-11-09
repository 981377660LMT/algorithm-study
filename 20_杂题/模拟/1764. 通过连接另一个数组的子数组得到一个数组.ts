function canChoose(groups: number[][], nums: number[]): boolean {
  const targetString = nums.map(num => `#${num}#`).join('')
  let position = 0

  for (const group of groups) {
    const searchString = group.map(num => `#${num}#`).join('')
    const hitIndex = targetString.indexOf(searchString, position)
    if (hitIndex === -1) return false
    else position = hitIndex + searchString.length
  }

  return true
}

console.log(
  canChoose(
    [
      [1, -1, -1],
      [3, -2, 0],
    ],
    [1, -1, 0, 1, -1, -1, 3, -2, 0]
  )
)
// 输入：groups = [[1,-1,-1],[3,-2,0]], nums = [1,-1,0,1,-1,-1,3,-2,0]
// 输出：true
// 解释：你可以分别在 nums 中选出第 0 个子数组 [1,-1,0,`1,-1,-1`,3,-2,0] 和第 1 个子数组 [1,-1,0,1,-1,-1,`3,-2,0`] 。
// 这两个子数组是不相交的，因为它们没有任何共同的元素。
console.log(canChoose([[2, 1]], [12, 1]))

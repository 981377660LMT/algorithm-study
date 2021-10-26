// 给定一个可包含重复数字的序列 nums ，按任意顺序 返回所有不重复的全排列。
// 记录每一个被选取的索引，而不是值，这就保证了同一个数字不会被选取多次，并且可以选取所有数字了
const permute = (nums: number[]) => {
  const res: number[][] = []
  nums.sort((a, b) => a - b)

  const bt = (path: number[], visited: boolean[]) => {
    // 1.递归终点
    if (path.length === nums.length) return res.push(path.slice())

    for (let i = 0; i < nums.length; i++) {
      // 同一个数字不能用两次
      if (visited[i]) continue
      // 规定每个重复的元素只能在开头第一个(i===index)被使用
      // 同样值的数字不能用两次 之前已经看过了,那就要防止[前一个1，后一个1]和[后一个1，前一个1]这种情况的出现 这里约定只允许[前一个1，后一个1]
      if (i > 0 && nums[i] === nums[i - 1] && visited[i - 1]) continue
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

console.log(permute([1, 1, 2]))
// nums = [1,1,2]
// 输出：
// [[1,1,2],
//  [1,2,1],
//  [2,1,1]]
export {}

/**
 *
 * @description 回溯算法,递归输出全排列，遇到包含重读元素的的情况就回溯；返回到达递归终点的情况
 * @description 时间复杂度O(n!)
 * @description 空间复杂度O(n) 递归堆栈层数
 */
const permute = (nums: number[]) => {
  const res: number[][] = []

  const bt = (path: number[]) => {
    // 1.递归终点
    if (path.length === nums.length) return res.push(path.slice())

    nums.forEach(num => {
      // 2.排除死路
      if (path.includes(num)) return
      // 3. 递归
      path.push(num)
      bt(path)
      path.pop()
    })
  }
  bt([])

  return res
}

console.log(permute([1, 2, 3]))

export {}

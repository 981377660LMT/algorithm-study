// 给你一个按升序排序的整数数组 num（可能包含重复数字），
// 请你将它们分割成一个或多个长度至少为 3 的子序列
/**
 * @summary
 * 每遍历一个数，将该数加入能加入的长度最短的序列中，不能加入序列则新建一个序列；然后更新字典
 * @link https://leetcode-cn.com/problems/split-array-into-consecutive-subsequences/solution/zui-jian-dan-de-pythonban-ben-by-semirondo/
 */
function isPossible(nums: number[]): boolean {
  const n = nums.length
  if (n < 3) return false

  const res: number[][] = [[nums[0]]]
  for (let i = 1; i < n; i++) {
    const cur = nums[i]
    let shouldNewArray = true
    // 注意新插入的数字优先少的 所以是从后往前遍历数组
    for (let j = res.length - 1; j >= 0; j--) {
      const arr = res[j]
      if (arr[arr.length - 1] + 1 === cur) {
        arr.push(cur)
        shouldNewArray = false
        break
      }
    }
    shouldNewArray && res.push([cur])
  }
  console.log(res)
  return res.every(arr => arr.length >= 3)
}

// console.log(isPossible([1, 2, 3, 3, 4, 4, 5, 5]))
console.log(isPossible([1, 2, 3, 4, 4, 5]))

export default 1

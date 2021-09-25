// 以数组形式返回能使 arr 有序的煎饼翻转操作所对应的 k 值序列

// 每次将最大的煎饼翻转到最后，需要两步翻转：
// 1.找到最大值的索引idx，并将idx之前的煎饼翻转；
// 2.将arr整体翻转；（除了最后已经翻转完的煎饼）
function pancakeSort(arr: number[]): number[] {
  let len = arr.length
  if (len <= 1) return []

  const res: number[] = []
  while (len > 0) {
    const index = arr.indexOf(max(arr.slice(0, len)))
    reverse(arr, 0, index)
    reverse(arr, 0, len - 1)
    res.push(index + 1) // 题中的k从1开始
    res.push(len)
    len--
  }

  return res
  function reverse(nums: number[], l: number, r: number) {
    if (l >= r) return
    while (l < r) {
      ;[nums[l], nums[r]] = [nums[r], nums[l]]
      l++
      r--
    }
    return nums
  }

  function max(nums: number[]) {
    if (!nums.length) throw new Error('max() arg is an empty sequence')
    return Math.max.apply(null, nums)
  }
}

console.log(pancakeSort([3, 2, 4, 1]))
// 输出：[4,2,4,3]
// 解释：
// 我们执行 4 次煎饼翻转，k 值分别为 4，2，4，和 3。
// 初始状态 arr = [3, 2, 4, 1]
// 第一次翻转后（k = 4）：arr = [1, 4, 2, 3]
// 第二次翻转后（k = 2）：arr = [4, 1, 2, 3]
// 第三次翻转后（k = 4）：arr = [3, 2, 1, 4]
// 第四次翻转后（k = 3）：arr = [1, 2, 3, 4]，此时已完成排序。

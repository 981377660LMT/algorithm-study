/**
 * @description 计算nums全部子序列和
 * @summary 时间复杂度O(2^n) 小于取所有子集的复杂度O(2^n*n)
 * 实在看不懂，也可用dfs来做
 */
function getSubArraySum(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(1 << n).fill(0)

  // 外层遍历数组每个元素，遍历到时，表示取该元素
  for (let i = 0; i < n; i++) {
    // 内层遍历从0到外层元素之间到每一个元素，表示能取到的元素，由于前面的结果已经计算过，因此可以直接累加
    for (let pre = 0; pre < 1 << i; pre++) {
      res[pre + (1 << i)] = res[pre] + nums[i]
    }
  }

  return res
}

function getSubArraySum2(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(1 << n).fill(0)
  dfs(0, 0, 0)
  return res

  function dfs(index: number, curSum: number, state: number): void {
    if (index === n) {
      res[state] = curSum
      return
    }

    dfs(index + 1, curSum, state)
    dfs(index + 1, curSum + nums[index], state | (1 << index))
  }
}

if (require.main === module) {
  console.log(getSubArraySum([1, 2]))
  console.log(getSubArraySum2([1, 2]))
}

export { getSubArraySum }

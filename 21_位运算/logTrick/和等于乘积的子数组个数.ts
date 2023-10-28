/**
 * 给一数组 a，元素均为正整数，求子数组和等于子数组积的子数组个数.
 */
function countSubarrayWithSumEqualToProduct(arr: ArrayLike<number>): number {
  type Interval = { leftStart: number; leftEnd: number; value: number }
  let allSum = 0
  for (let i = 0; i < arr.length; i++) allSum += arr[i]

  const preSum = new Map<number, number>([[0, 0]])
  let curSum = 0
  let res = 0
  const dp: Interval[] = []
  for (let pos = 0; pos < arr.length; pos++) {
    const cur = arr[pos]
    curSum += cur
    for (let i = 0; i < dp.length; i++) dp[i].value *= cur
    dp.push({ leftStart: pos, leftEnd: pos + 1, value: cur })

    // 去重
    let ptr = 0
    for (let i = 1; i < dp.length; i++) {
      if (dp[i].value !== dp[ptr].value) {
        ptr++
        dp[ptr] = dp[i]
      } else {
        dp[ptr].leftEnd = dp[i].leftEnd
      }
    }
    dp.length = ptr + 1

    // 去掉超过 allSum 的，从而保证 dp 中至多有 O(log(allSum)) 个元素
    while (dp.length && dp[0].value > allSum) {
      dp.shift()
    }

    // 将区间[0,pos]分成了dp.length个左闭右开区间.
    // 每一段区间的左端点left范围 在 [dp[i].leftStart,dp[i].leftEnd) 中。
    // 对应子数组 arr[left:pos+1] 的 op 值为 dp[i].value.
    for (let i = 0; i < dp.length; i++) {
      const { leftStart, leftEnd, value } = dp[i]
      const target = curSum - value
      const pos = preSum.get(target)!
      if (pos != undefined && leftStart <= pos && pos < leftEnd) {
        res++
      }
    }

    preSum.set(curSum, pos + 1)
  }

  return res
}

export {}

if (require.main === module) {
  console.log(countSubarrayWithSumEqualToProduct([1, 3, 2, 2])) // 6
}

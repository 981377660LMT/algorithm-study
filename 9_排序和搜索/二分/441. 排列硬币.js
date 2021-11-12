// 你总共有 n 枚硬币，并计划将它们按阶梯状排列。
// 对于一个由 k 行组成的阶梯，其第 i 行必须正好有 i 枚硬币。阶梯的最后一行 可能 是不完整的。
// 给你一个数字 n ，计算并返回可形成 完整阶梯行 的总行数。

function arrangeCoins(n) {
  let l = 1
  let r = n
  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    const sum = (mid * (mid + 1)) / 2
    if (sum === n) return mid
    else if (sum < n) l = mid + 1
    else r = mid - 1
  }
  return l - 1
}

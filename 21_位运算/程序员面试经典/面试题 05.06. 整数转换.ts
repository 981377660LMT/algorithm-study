// 确定需要改变几个位才能将整数A转成整数B。
function convertInteger(A: number, B: number): number {
  // 汉明距离
  let xor = A ^ B
  let res = 0
  while (xor) {
    xor &= xor - 1
    res++
  }

  return res
}

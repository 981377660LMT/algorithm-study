function largestVariance(s: string): number {
  const allChars = [...new Set(s)]

  let res = 0
  for (const s1 of allChars) {
    for (const s2 of allChars) {
      if (s1 === s2) continue
      res = Math.max(res, cal(s1, s2))
    }
  }

  return res

  // 最多的s1次数与最少的s2次数之差。
  function cal(s1: string, s2: string): number {
    let res = 0
    let [maxSum1, maxSum2] = [0, -Infinity] // 当前的和，包含s2的和

    for (const char of s) {
      if (char === s1) {
        maxSum1++
        maxSum2++
      } else if (char === s2) {
        maxSum1--
        maxSum2 = maxSum1 // 更新当前包含s2的最大值
        if (maxSum1 < 0) maxSum1 = 0 // 前面一截全不要了
      }

      res = Math.max(res, maxSum2)
    }

    return res
  }
}

export {}

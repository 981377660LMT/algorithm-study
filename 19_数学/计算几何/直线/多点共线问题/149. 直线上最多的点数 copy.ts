function maxPoints(points: number[][]): number {
  const n = points.length

  let res = 0

  // 考虑通过每个点的所有直线（和这个点之后的所有点形成的线段斜率），使用最大公约数归一化斜率为分数：
  for (let i = 0; i < n; i++) {
    const counter = new Map<string, number>()
    const [x1, y1] = points[i]
    let curMax = 0
    for (let j = i + 1; j < n; j++) {
      const [x2, y2] = points[j]

      let A = y2 - y1
      let B = x2 - x1
      const divide = gcd(A, B)
      A /= divide // 如果divide为0，则为Infinity
      B /= divide
      const key = `${A}#${B}`

      counter.set(key, (counter.get(key) || 0) + 1)
      curMax = Math.max(curMax, counter.get(key)!)
    }

    res = Math.max(res, curMax + 1)
  }

  return res
}

function gcd(...nums: number[]) {
  const twoNumGcd = (a: number, b: number): number => {
    return b === 0 ? a : twoNumGcd(b, a % b)
  }
  return nums.reduce(twoNumGcd)
}

export {}

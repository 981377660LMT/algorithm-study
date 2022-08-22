/* eslint-disable no-param-reassign */
/* eslint-disable no-multi-assign */
// 41312432 或 23421314
// 有8个数，11223344
// 将其排列，要求结果满足：两个1之间有一个数，两个2之间有两个数，两个3之间有三个数，两个4之间有四个数。问这个结果是多少？
// 可以证明 n%4 = 1和n%4 =2一定没有解

// n<=20
function generateLangFordSequence(n: number): number[][] {
  if (n % 4 === 1 || n % 4 === 2) return []

  const res: number[][] = []
  bt(0, Array(2 * n).fill(-1))
  return res

  function bt(index: number, path: number[]): void {
    if (index === 2 * n) {
      res.push(path.slice())
      return
    }

    if (path[index] !== -1) {
      bt(index + 1, path)
      return
    }

    for (let num = n; num > 0; num--) {
      // 已经选过这个数了
      if (path.includes(num)) continue

      const gap = num
      const anotherIndex = index + gap
      if (anotherIndex < path.length && path[anotherIndex] === -1) {
        path[index] = path[anotherIndex] = num
        bt(index + 1, path)
        path[index] = path[anotherIndex] = -1
      }
    }
  }
}

if (require.main === module) {
  console.log(generateLangFordSequence(4))
}

export { generateLangFordSequence }

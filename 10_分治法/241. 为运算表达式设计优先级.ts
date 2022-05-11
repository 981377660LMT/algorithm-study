// 为表达式添加括号，改变其运算优先级以求出不同的结果。
// 你需要给出所有可能的组合的结果。
// 有效的运算符号包含 +, - 以及 * 。

const opts = new Set(['+', '-', '*'])
const isDigit = (x: any) => !isNaN(parseFloat(x)) && isFinite(x)

function diffWaysToCompute(expression: string): number[] {
  // 注意返回如果是数字直接返回数字结果
  if (isDigit(expression)) return [Number(expression)]
  // 采用后序遍历(归并排序思路)
  const res: number[] = []

  for (let i = 0; i < expression.length; i++) {
    if (opts.has(expression[i])) {
      const left = diffWaysToCompute(expression.slice(0, i))
      const right = diffWaysToCompute(expression.slice(i + 1))
      const opt = expression[i]

      for (let l of left) {
        for (let r of right) {
          l = Number(l)
          r = Number(r)

          if (opt === '+') {
            res.push(l + r)
          } else if (opt === '-') {
            res.push(l - r)
          } else {
            res.push(l * r)
          }
        }
      }
    }
  }

  return res
}

console.log(diffWaysToCompute('2-1-1'))
// 输出: [0, 2]
// 解释:
// ((2-1)-1) = 0
// (2-(1-1)) = 2
export {}

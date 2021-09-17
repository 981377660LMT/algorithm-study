/**
 * @param {string} s
 * @param {number} result
 * @return {number}
 * 布尔表达式由 0 (0)、1 (1)、& (AND)、 | (OR) 和 ^ (XOR) 符号组成。
 * 实现一个函数，算出有几种可使该表达式得出 result 值的括号方法。

 */
const countEval = function (s: string, result: number): number {
  const memo = new Map<string, number>()
  const opts = {
    '&': {
      1: [[1, 1]],
      0: [
        [1, 0],
        [0, 1],
        [0, 0],
      ],
    },
    '|': {
      1: [
        [1, 0],
        [0, 1],
        [1, 1],
      ],
      0: [[0, 0]],
    },
    '^': {
      1: [
        [1, 0],
        [0, 1],
      ],
      0: [
        [1, 1],
        [0, 0],
      ],
    },
  } as Record<string, Record<string, [number, number][]>>

  const inner = (s: string, result: number, memo: Map<string, number>): number => {
    const key = `${s}#${result}`
    if (memo.has(key)) return memo.get(key)!
    if (!s.length) return 0
    if (s.length === 1) return Number(s === result.toString())
    // if (s.length === 3) return eval(s) === result ? 1 : 0

    let res = 0
    for (let i = 1; i < s.length - 1; i += 2) {
      const opt = s[i]
      const left = s.slice(0, i)
      const right = s.slice(i + 1)
      const needList = opts[opt][result]
      for (const [leftVal, rightVal] of needList) {
        res += inner(left, leftVal, memo) * inner(right, rightVal, memo)
      }
    }

    memo.set(key, res)
    return res
  }

  return inner(s, result, memo)
}

console.log(countEval('0^0|0|0&1|0|0^1&1&1^0', 0))

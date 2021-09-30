/**
 *
 * @param num  1 <= num.length <= 10
 * @param target
 * 在 num 的数字之间添加 二元 运算符（不是一元）+、- 或 * ，返回所有能够得到目标值的表达式。
 */
function addOperators(num: string, target: number): string[] {
  const res: string[] = []
  dfs(0, 0, 0, '')
  return res

  function dfs(sum: number, pre: number, index: number, str: string) {
    if (index === num.length) {
      if (sum === target) res.push(str)
      return
    }

    // 插入位置
    for (let len = 1; len + index <= num.length; len++) {
      const slice = num.slice(index, index + len)
      if (slice[0] === '0' && slice.length > 1) continue // prevent "00*" as a number
      const cur = Number(slice)
      if (index === 0) {
        dfs(cur, cur, index + len, slice)
      } else {
        dfs(sum + cur, cur, index + len, str + '+' + slice)
        dfs(sum - cur, -cur, index + len, str + '-' + slice)
        dfs(sum + pre * cur - pre, pre * cur, index + len, str + '*' + slice)
      }
    }
  }
}

console.log(addOperators('123', 6))
// 输入: num = "123", target = 6
// 输出: ["1+2+3", "1*2*3"]

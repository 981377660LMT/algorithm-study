/**
 * @param {number} n
 * @return {string}
 * 从数字 1 开始，序列中的每一项都是对前一项的描述。
 */
const countAndSay = function (n: number): string {
  let res = '1'
  for (let i = 1; i < n; i++) {
    res = res.replace(/(\d)\1*/g, (match, g1) => {
      // 几个匹配项
      return `${match.length}${g1}`
    })
  }
  return res
}

console.log(countAndSay(3))
// 输出："1211"
// 解释：
// countAndSay(1) = "1"
// countAndSay(2) = 读 "1" = 一 个 1 = "11"
// countAndSay(3) = 读 "11" = 二 个 1 = "21"
// countAndSay(4) = 读 "21" = 一 个 2 + 一 个 1 = "12" + "11" = "1211"

export default 1

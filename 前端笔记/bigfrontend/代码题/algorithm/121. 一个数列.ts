/**
 * @param {number} n - integer
 * @returns {string}
 */
function getNthNum(n: number): string {
  // your code here
  let res = '1'
  for (let i = 1; i < n; i++) {
    res = res.replace(/(\d)\1*/g, (match, g1) => {
      // 几个匹配项
      return `${match.length}${g1}`
    })
  }
  return res
}

console.log(getNthNum(5))

// 按照以下规则可以生成一个数列。

// '1'，第一个是1
// '11'，前一个数包含1个1
// '21'，前一个数包含2个1
// '1211'，前一个数包含1个2，1个1
// '111221'，前一个数包含1个1，1个2和2个1
// '312211'，前一个数包含3个1，2个2和1个1
// ....
// 也就是说通过计数前面的数字可以得到下一个数。

// 请实现getNthNum(n)来返回该数列中的第n个数，n从1开始。

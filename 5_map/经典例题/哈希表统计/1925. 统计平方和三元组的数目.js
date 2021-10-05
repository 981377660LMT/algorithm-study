/**
 * @param {number} n  1 <= n <= 250
 * @return {number}
 */
var countTriples = function (n) {
  let res = 0
  let set = new Set()
  for (let i = 1; i <= n; i++) {
    set.add(i ** 2)
  }

  for (let i = 1; i < n; i++) {
    for (let j = 1; j < n; j++) {
      if (set.has(i ** 2 + j ** 2)) res++
    }
  }
  return res
}

console.log(countTriples(10))
// 输入：n = 10
// 输出：4
// 解释：平方和三元组为 (3,4,5)，(4,3,5)，(6,8,10) 和 (8,6,10) 。

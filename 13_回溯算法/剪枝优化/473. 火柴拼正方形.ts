/**
 * @param {number[]} matchsticks  火柴数组的长度不超过15。
 * @return {boolean}
 * 不能折断火柴，可以把火柴连接起来，并且每根火柴都要用到。
 * @summary
 * 如果先放一个权值大的，那么选择就会少很多，因此递归树的规模就会小很多。
 * 降序排序，优先选择权值大的可以减少搜索树的规模。
 */
var makesquare = function (matchsticks: number[]): boolean {
  const sum = matchsticks.reduce((pre, cur) => pre + cur)
  const div = sum / 4
  const mod = sum % 4
  if (mod) return false

  const res = [0, 0, 0, 0]
  matchsticks.reverse() // 降序排序，优先选择权值大的可以减少搜索树的规模。

  const bt = (index: number): boolean => {
    if (index === matchsticks.length) return res.every(v => v === div)
    for (let j = 0; j < 4; j++) {
      if (matchsticks[index] + res[j] <= div) {
        res[j] += matchsticks[index]
        if (bt(index + 1)) return true
        res[j] -= matchsticks[index]
      }
    }
    return false
  }
  return bt(0)
}

console.log(makesquare([1, 1, 2, 2, 2]))

export {}

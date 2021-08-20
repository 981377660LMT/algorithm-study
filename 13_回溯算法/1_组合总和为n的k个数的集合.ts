/**
 * @param {number} k
 * @param {number} n
 * @return {number[][]}
 * @description 找出所有相加之和为 n 的 k 个数的组合。组合中只允许含有 1 - 9 的正整数，并且每种组合中不存在重复的数字。
 * @description 使用每次index+1 保证不取到重读数字
 */
const combinationSum = (k: number, n: number): number[][] => {
  const res: number[][] = []

  const bt = (path: number[], sum: number, index: number) => {
    if (path.length === k) {
      if (sum === n) res.push(path.slice())
      return
    }

    for (let i = index; i < 10; i++) {
      path.push(i)
      bt(path, sum + i, i + 1)
      path.pop()
    }
  }
  bt([], 0, 1)

  return res
}

console.log(combinationSum(3, 9))
// [[1,2,6], [1,3,5], [2,3,4]]
export {}

/**
 * @param {number[]} candidates
 * @param {number} target
 * @return {number[][]}
 * @description candidates 中可以有重复数字，但是每个数字在每个组合中只能使用一次。
 * @description 如何去除重复的元素?
 * @description 1.排好序 2.使用index每次加1/visited剪枝
 */
const combinationSum = (candidates: number[], target: number): number[][] => {
  const res: number[][] = []
  const len = candidates.length
  candidates.sort((a, b) => a - b)

  const bt = (path: number[], sum: number, index: number) => {
    if (sum > target) {
      return
    } else if (sum === target) {
      return res.push(path.slice())
    }

    for (let i = index; i < len; i++) {
      // 规定每个重复的元素只能在开头第一个(i===index)被使用
      if (i > index && candidates[i] === candidates[i - 1]) continue

      const next = candidates[i]
      path.push(next)
      // 注意这里i+1限制不能取到重复的元素
      bt(path, sum + next, i + 1)
      path.pop()
    }
  }
  bt([], 0, 0)

  return res
}

console.log(combinationSum([10, 1, 1, 1, 1, 2, 7, 6, 1, 5], 8))
// [
//   [1,1,6],
//   [1,2,5],
//   [1,7],
//   [2,6]
//   ]
export {}

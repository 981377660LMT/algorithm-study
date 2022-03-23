/**
 * @param {number[]} candidates
 * @param {number} target
 * @return {number[][]}
 * @description 无重复元素的数组candidate 中的每个元素都是独一无二的 每个元素可使用多次。
 * @description 如何去除重复的元素? 例如 [2,2,3] [2,3,2] [3,2,2] 使用index保持索引不减
 */
const combinationSum = (candidates: number[], target: number): number[][] => {
  const res: number[][] = []
  const len = candidates.length
  // candidates.sort((a, b) => a - b)

  const bt = (path: number[], sum: number, index: number) => {
    if (sum > target) {
      return
    } else if (sum === target) {
      return res.push(path.slice())
    }

    for (let i = index; i < len; i++) {
      const next = candidates[i]
      path.push(next)
      // i 数字可以重复使用
      bt(path, sum + next, i + 1)
      path.pop()
    }
  }
  bt([], 0, 0)

  return res
}

console.log(combinationSum([2, 3, 6, 7], 7))

export {}

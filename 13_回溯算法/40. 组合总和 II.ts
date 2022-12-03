// 给定一个候选人编号的集合 candidates 和一个目标数 target ，
// 找出 candidates 中所有可以使数字和为 target 的组合。
// candidates 中的每个数字在每个组合中只能使用 一次 。

// 1 <= candidates.length <= 100
// 1 <= candidates[i] <= 50
// 1 <= target <= 30

// 注意：解集不能包含重复的组合。
function combinationSum2(candidates: number[], target: number): number[][] {
  const n = candidates.length
  const res: number[][] = []
  candidates.sort((a, b) => a - b)

  const bt = (path: number[], sum: number, index: number): void => {
    if (sum > target) return
    if (sum === target) {
      res.push(path.slice())
      return
    }

    for (let i = index; i < n; i++) {
      // !规定一组重复的元素中，只有开头第一个(i===index)被使用 避免[1,2] [1,2] 这种情况
      if (i > index && candidates[i] === candidates[i - 1]) continue

      const next = candidates[i]
      path.push(next)
      bt(path, sum + next, i + 1) // i+1 限制不能取到重复的元素
      path.pop()
    }
  }

  bt([], 0, 0)
  return res
}

console.log(combinationSum2([10, 1, 1, 1, 1, 2, 7, 6, 1, 5], 8))
// [
//   [1,1,6],
//   [1,2,5],
//   [1,7],
//   [2,6]
//   ]
export {}

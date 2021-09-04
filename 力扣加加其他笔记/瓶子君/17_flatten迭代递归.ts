interface NestedArray extends Array<NestedArray | number> {}
// 递归
const flatten1 = (nums: NestedArray): number[] =>
  nums.reduce<number[]>(
    (pre, cur) => (Array.isArray(cur) ? [...pre, ...flatten(cur)] : [...pre, cur]),
    []
  )

// 迭代
const flatten2 = (nums: NestedArray): number[] => {
  const res: number[] = []
  const queue: NestedArray = []

  for (let i = 0; i < nums.length; i++) {
    queue.push(nums[i])
    while (queue.length) {
      const head = queue.shift()!
      if (!Array.isArray(head)) res.push(head)
      else queue.push(...head)
    }
  }

  return res
}

console.log(flatten2([1, 2, [3, 4, 5, 6, [5]]]))

export default 1

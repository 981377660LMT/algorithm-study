interface NestedArray extends Array<NestedArray | number> {}
// 递归
// const flatten1 = (nums: NestedArray): number[] =>
//   nums.reduce<number[]>(
//     (pre, cur) => (Array.isArray(cur) ? [...pre, ...flatten(cur)] : [...pre, cur]),
//     []
//   )

// bfs迭代，可带层数限制
const flatten2 = (nums: NestedArray, depth: number = 1): (NestedArray | number)[] => {
  const res: (NestedArray | number)[] = []
  const queue: NestedArray = nums.slice()

  let step = 0
  while (queue.length && step++ < depth) {
    const length = queue.length
    for (let _ = 0; _ < length; _++) {
      const head = queue.shift()!
      if (!Array.isArray(head)) res.push(head)
      else queue.push(...head)
    }
  }

  return [...res, ...queue]
}

console.log(flatten2([1, 2, [3, 4, 5, 6, [5]]]))

export default 1

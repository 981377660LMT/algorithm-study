// 验证拓扑排序唯一性
// 请判断原始的序列 org 是否可以从序列集 seqs 中唯一地 重建 。
// 总结：路径唯一 ==== queue长度始终为1
function sequenceReconstruction(org: number[], seqs: number[][]): boolean {
  const adjMap = new Map<number, number[]>()
  const indegrees = new Map<number, number>()
  const vertex = new Set<number>()

  for (const seq of seqs) {
    vertex.add(seq[0])
    for (let i = 0; i < seq.length - 1; i++) {
      const [cur, next] = [seq[i], seq[i + 1]]
      !adjMap.has(cur) && adjMap.set(cur, [])
      adjMap.get(cur)!.push(next)
      indegrees.set(next, (indegrees.get(next) || 0) + 1)
      vertex.add(next)
    }
  }

  // 顶点不同 直接返回false
  if (!isSameSet(new Set(org), vertex)) return false

  // 拓扑排序
  const queue: number[] = []
  const res: number[] = []
  for (const v of vertex) {
    if (!indegrees.has(v)) queue.push(v)
  }

  while (queue.length) {
    if (queue.length !== 1) return false
    const cur = queue.shift()!
    res.push(cur)
    for (const next of adjMap.get(cur) || []) {
      indegrees.set(next, (indegrees.get(next) || 0) - 1)
      if (indegrees.get(next) === 0) queue.push(next)
    }
  }

  return isSameArray(res, org)

  function isSameSet<T>(set1: Set<T>, set2: Set<T>): boolean {
    if (set1.size !== set2.size) return false
    for (const item of set1) {
      if (!set2.has(item)) return false
    }
    return true
  }

  function isSameArray<T>(arr1: Array<T>, arr2: Array<T>): boolean {
    if (arr1.length !== arr2.length) return false
    return arr1.every((value, index) => value === arr2[index])
  }
}

console.log(sequenceReconstruction([1], [[1, 1]]))

// console.log(
//   sequenceReconstruction(
//     [1, 2, 3],
//     [
//       [1, 2],
//       [1, 3],
//     ]
//   )
// )
// // false

// console.log(
//   sequenceReconstruction(
//     [1, 2, 3],
//     [
//       [1, 2],
//       [1, 3],
//       [2, 3],
//     ]
//   )
// )
// // true

// console.log(sequenceReconstruction([1, 2, 3], [[1, 2]])) // false

// console.log(
//   sequenceReconstruction(
//     [4, 1, 5, 2, 6, 3],
//     [
//       [5, 2, 6, 3],
//       [4, 1, 5, 2],
//     ]
//   )
// )
// // 输出: true

import { GCD } from '../../19_数学/最大公约数/gcd'

type NodeValue = number
type Level = number
type NodeIndex = number

/**
 * @param {number[]} nums  nums[i] 表示第 i 个点的值  1 <= nums[i] <= 50
 * @param {number[][]} edges  edges[j] = [uj, vj] 表示节点 uj 和节点 vj 在树中有一条边。
 * @return {number[]}
 * @summary
 * 枚举满足 1 <= y <= 50  且 gcd(x,y) = 1 的 y，
 * 并对每个 y 找出离着节点 i 最近的点，最后再在这些点中求出离着当前点最近的点即可。
 * 这样只需检查 50 次即可。而不用对每个节点向上dfs
 */
const getCoprimes = function (nums: number[], edges: number[][]): number[] {
  const res: number[] = []
  const adjList = Array.from<unknown, number[]>({ length: nums.length }, () => [])
  edges.forEach(([u, v]) => {
    adjList[u].push(v)
    adjList[v].push(u)
  })

  // 对每个数字 1∼50 维护一个栈
  // 向下dfs过程中栈顶就是数字 x 的层数最深的节点
  // dfs 完成后要将之前 push 进去的元素 pop 出来
  const valueToNodeDetail = Array.from<NodeValue, [Level, NodeIndex][]>({ length: 55 }, () => [])

  const dfs = (cur: number, level: number, nodeValues: number[], visited: Set<number>) => {
    if (visited.has(cur)) return
    visited.add(cur)

    let preNodeLevel = -1
    let preNodeIndex = -1
    for (let value = 1; value <= 50; value++) {
      if (!valueToNodeDetail[value].length || GCD(nodeValues[cur], value) !== 1) continue
      const preNodes = valueToNodeDetail[value]
      const [level, index] = preNodes[preNodes.length - 1]
      if (level > preNodeLevel) {
        preNodeLevel = level
        preNodeIndex = index
      }
    }

    res[cur] = preNodeIndex

    for (const next of adjList[cur]) {
      valueToNodeDetail[nums[cur]].push([level, cur])
      dfs(next, level + 1, nums, visited)
      valueToNodeDetail[nums[cur]].pop()
    }
  }

  dfs(0, 0, nums, new Set())
  return res
}

console.log(
  getCoprimes(
    [5, 6, 10, 2, 3, 6, 15],
    [
      [0, 1],
      [0, 2],
      [1, 3],
      [1, 4],
      [2, 5],
      [2, 6],
    ]
  )
)

// 总结：
// 1. dfs的过程 使用栈记录节点的信息 栈顶总是层数最深的
// 2. 换个角度遍历 小集合遍历

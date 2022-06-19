type NodeValue = number
type NodeLevel = number
type NodeId = number

/**
 * @param {number[]} nums  nums[i] 表示第 i 个点的值  1 <= nums[i] <= 50
 * @param {number[][]} edges  edges[j] = [uj, vj] 表示节点 uj 和节点 vj 在树中有一条边。
 * @return {number[]}
 * @summary
 * 枚举满足 1 <= y <= 50  且 gcd(x,y) = 1 的 y，
 * 并对每个 y 找出离着节点 i 最近的点，最后再在这些点中求出离着当前点最近的点即可。
 * 这样只需检查 50 次即可。而不用对每个节点向上dfs
 */
function getCoprimes(nums: number[], edges: number[][]): number[] {
  // 建图
  const adjList = Array.from<unknown, number[]>({ length: nums.length }, () => [])
  edges.forEach(([u, v]) => {
    adjList[u].push(v)
    adjList[v].push(u)
  })

  // 对每个数字 1∼50 维护一个栈
  // !向下dfs过程中栈顶就是数字 x 的层数最深的节点
  // dfs 完成后要将之前 push 进去的元素 pop 出来
  const stacks = Array.from<NodeValue, [NodeLevel, NodeId][]>({ length: 51 }, () => [])

  const res: number[] = []
  dfs(0, 0, new Set())
  return res

  function dfs(cur: number, level: number, visited: Set<number>): void {
    if (visited.has(cur)) return
    visited.add(cur)

    let parentLevel = -1
    let parentId = -1
    for (let n = 1; n <= 50; n++) {
      if (stacks[n].length === 0 || gcd(nums[cur], n) !== 1) continue
      const parents = stacks[n]
      const [level, index] = parents[parents.length - 1]
      if (level > parentLevel) {
        parentLevel = level
        parentId = index
      }
    }

    res[cur] = parentId

    for (const next of adjList[cur]) {
      stacks[nums[cur]].push([level, cur])
      dfs(next, level + 1, visited)
      stacks[nums[cur]].pop()
    }
  }
}

function gcd(...nums: number[]): number {
  const _gcd = (a: number, b: number): number => (b === 0 ? a : gcd(b, a % b))
  return nums.reduce(_gcd)
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

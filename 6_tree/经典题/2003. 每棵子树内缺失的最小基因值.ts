// 2 <= n <= 105
// parents[0] == -1
// 1 <= nums[i] <= 105
// nums[i] 互不相同。
function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = parents.length
  const res = Array<number>(n).fill(1)
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  const visitedGene = new Set<number>()

  for (let i = 1; i < n; i++) {
    adjList[parents[i]].push(i)
  }

  // 先找基因值=1的结点
  // 如果没有基因值为1的结点。则所有点的缺失gene是1
  // 1.root的子树肯定缺失1
  // 2.root不断向上走 minLost逐渐递增 不断dfs 看visitedGene 有哪些没有
  const gene1Node = nums.indexOf(1)
  let root = gene1Node
  let minLost = 1
  while (root !== -1) {
    visitSubtree(root)
    while (visitedGene.has(minLost)) {
      minLost++
    }

    res[root] = minLost
    root = parents[root]
  }

  return res

  function visitSubtree(cur: number): void {
    visitedGene.add(nums[cur])
    for (const next of adjList[cur]) {
      if (visitedGene.has(nums[next])) continue
      visitSubtree(next)
    }
  }
}

console.log(smallestMissingValueSubtree([-1, 0, 1, 0, 3, 3], [5, 4, 6, 2, 1, 3]))
// 输出：[7,1,1,4,2,1]
// 解释：每个子树答案计算结果如下：
// - 0：子树内包含节点 [0,1,2,3,4,5] ，基因值分别为 [5,4,6,2,1,3] 。7 是缺失的最小基因值。
// - 1：子树内包含节点 [1,2] ，基因值分别为 [4,6] 。 1 是缺失的最小基因值。
// - 2：子树内只包含节点 2 ，基因值为 6 。1 是缺失的最小基因值。
// - 3：子树内包含节点 [3,4,5] ，基因值分别为 [2,1,3] 。4 是缺失的最小基因值。
// - 4：子树内只包含节点 4 ，基因值为 1 。2 是缺失的最小基因值。
// - 5：子树内只包含节点 5 ，基因值为 3 。1 是缺失的最小基因值。

// 总结：
// 0.初始化res所有节点默认最小缺失值为1
// 1.先找基因值为1的结点，找不到1直接返回全部缺失1。
// 2.如果找到gene1Node，从gene1Node不断向上走，这个过程中不断对新的parent做dfs，把子树的基因加入visited，第一个没有出现在visited里的就是缺失的最小基因。
// 注意：
// 我们可以把所有子树分成三种:

// - 基因值为1的结点的子树最小缺失值都为1；
// - 子树中没有1的树，最小缺失值为1；
// - 基因值为1的父辈节点的最小缺失值不为1，需要dfs看他的子树里有哪些基因。

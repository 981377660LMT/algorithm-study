/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

/** DAG最长路. */
function longestPathInDag(n: number, adjList: ArrayLike<ArrayLike<number>>, getWeight: (from: number, to: number) => number = () => 1): number[] {
  const indeg = new Uint32Array(n)
  for (let i = 0; i < n; i++) {
    const nexts = adjList[i]
    for (let j = 0; j < nexts.length; j++) {
      indeg[nexts[j]]++
    }
  }

  const queue = new Uint32Array(n)
  let head = 0
  let tail = 0
  for (let i = 0; i < n; i++) {
    if (indeg[i] === 0) {
      queue[tail++] = i
    }
  }

  const dp = Array(n).fill(0)
  while (head < tail) {
    const cur = queue[head++]
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      dp[next] = Math.max(dp[next], dp[cur] + getWeight(cur, next))
      if (--indeg[next] === 0) {
        queue[tail++] = next
      }
    }
  }

  return dp
}

const INF = 2e15

/**
 * 从起点出发到终点的最长路.如果无法到达,则返回-1.
 * @param start 起点
 * @param end 终点
 */
function longestPathInDagWithStart(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  start: number,
  end: number,
  getWeight: (from: number, to: number) => number = () => 1
): number {
  const dp = new Float64Array(n)
  const visited = new Uint8Array(n)
  const dfs = (cur: number): number => {
    if (visited[cur]) return dp[cur]
    visited[cur] = 1
    if (cur === end) return 0

    let res = -INF
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      res = Math.max(res, dfs(next) + getWeight(cur, next))
    }

    dp[cur] = res
    return res
  }

  const res = dfs(start)
  return visited[end] ? res : -1
}

export { longestPathInDag, longestPathInDagWithStart }

if (require.main === module) {
  // 2050. 并行课程 III
  // https://leetcode.cn/problems/parallel-courses-iii/description/
  function minimumTime(n: number, relations: number[][], time: number[]): number {
    const dummy = n
    const adjList: number[][] = Array(n + 1)
    for (let i = 0; i < n + 1; i++) adjList[i] = []
    relations.forEach(v => {
      adjList[v[0] - 1].push(v[1] - 1)
    })
    for (let i = 0; i < n; i++) {
      adjList[dummy].push(i)
    }
    const dp = longestPathInDag(n + 1, adjList, (_, to) => time[to])
    return Math.max(...dp)
  }

  // 2770. 达到末尾下标所需的最大跳跃次数
  // https://leetcode.cn/problems/maximum-number-of-jumps-to-reach-the-last-index/
  function maximumJumps(nums: number[], target: number): number {
    const n = nums.length
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    for (let i = 0; i < n; i++) {
      for (let j = i + 1; j < n; j++) {
        const diff = nums[j] - nums[i]
        if (diff >= -target && diff <= target) {
          adjList[i].push(j)
        }
      }
    }

    return longestPathInDagWithStart(n, adjList, 0, n - 1)
  }

  // [1,3,6,4,1,2]
  console.log(maximumJumps([1, 3, 6, 4, 1, 2], 2))
}

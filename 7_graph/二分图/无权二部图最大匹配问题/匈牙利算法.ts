/* eslint-disable no-shadow */
// 二分图的最大匹配问题
// 对左侧每一个尚未匹配的点，不断地寻找可以增广的交替路 (有点像点引线，连环触发)
// 如果是完美匹配 则匹配对数*2等于顶点数
// https://blog.csdn.net/kaisa158/article/details/48718403?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522161564099216780266286846%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=161564099216780266286846&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~baidu_landing_v2~default-1-48718403.first_rank_v2_pc_rank_v29&utm_term=%E5%8C%88%E7%89%99%E5%88%A9dfs

/**
 * 匈牙利算法求无权二分图最大匹配
 *
 * @param row 男孩的个数
 * @param col 女孩的个数
 * @complexity O(V*E)
 */
function useHungarian(row: number, col: number) {
  const adjList = Array.from<unknown, number[]>({ length: row }, () => [])
  const rowMatching = new Int32Array(row).fill(-1)
  const colMathching = new Int32Array(col).fill(-1)
  let matchingEdges: [boy: number, girl: number][]

  function addEdge(boy: number, girl: number): void {
    if (!(boy >= 0 && boy < row && girl >= 0 && girl < col)) {
      throw new RangeError(`${boy},${girl} out of range ${row},${col}`)
    }
    adjList[boy].push(girl)
  }

  function work(): number {
    // 寻找增广路
    const dfs = (cur: number): boolean => {
      if (visited[cur]) return false

      visited[cur] = 1
      for (let i = 0; i < adjList[cur].length; i++) {
        const next = adjList[cur][i]
        if (colMathching[next] === -1 || dfs(colMathching[next])) {
          colMathching[next] = cur
          rowMatching[cur] = next
          return true
        }
      }

      return false
    }

    let res = 0
    let hasUpdated = true
    const visited = new Uint8Array(row).fill(0)
    while (hasUpdated) {
      hasUpdated = false
      for (let cur = 0; cur < row; cur++) {
        if (rowMatching[cur] === -1 && dfs(cur)) {
          hasUpdated = true
          res++
        }
      }

      if (hasUpdated) visited.fill(0)
    }

    return res
  }

  function getMathingEdges(): [boy: number, girl: number][] {
    if (matchingEdges) return matchingEdges

    matchingEdges = []
    for (let cur = 0; cur < row; cur++) {
      if (rowMatching[cur] !== -1) {
        matchingEdges.push([cur, rowMatching[cur]])
      }
    }

    return matchingEdges
  }

  return {
    /**
     * 男孩向女孩连边
     */
    addEdge,
    /**
     * 求二分图最大匹配数
     */
    work,
    /**
     * 取得最大匹配时的匹配对
     */
    getMathingEdges
  }
}

export { useHungarian }

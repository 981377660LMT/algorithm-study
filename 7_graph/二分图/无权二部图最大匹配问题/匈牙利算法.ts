// 二分图的最大匹配问题
// 对左侧每一个尚未匹配的点，不断地寻找可以增广的交替路 (有点像点引线，连环触发)
// 如果是完美匹配 则匹配对数*2等于顶点数
// https://blog.csdn.net/kaisa158/article/details/48718403?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522161564099216780266286846%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=161564099216780266286846&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~baidu_landing_v2~default-1-48718403.first_rank_v2_pc_rank_v29&utm_term=%E5%8C%88%E7%89%99%E5%88%A9dfs

function hungarian(adjList: number[][]): number {
  let maxMatching = 0
  let visited = Array<boolean>(adjList.length).fill(false)
  const matching = Array<number>(adjList.length).fill(-1)

  const colors = bisect(adjList)
  for (let i = 0; i < adjList.length; i++) {
    // 从左侧还没有匹配到的男生出发，并重置visited
    if (colors[i] === 0 && matching[i] === -1) {
      visited = Array<boolean>(adjList.length).fill(false)
      if (dfs(i)) maxMatching++
    }
  }

  return maxMatching

  // 匈牙利算法核心:寻找增广路径 找到的话最大匹配加一
  // dfs(cur) 表示给cur找匹配
  function dfs(cur: number): boolean {
    if (visited[cur]) return false
    visited[cur] = true

    for (const next of adjList[cur]) {
      // 是增广路径或者dfs找到增广路径
      if (matching[next] === -1 || dfs(matching[next])) {
        matching[cur] = next
        matching[next] = cur
        return true
      }
    }

    return false
  }

  // 二分图检测、获取colors
  function bisect(adjList: number[][]) {
    const colors = Array<number>(adjList.length).fill(-1)

    const dfs = (cur: number, color: number): void => {
      colors[cur] = color
      for (const next of adjList[cur]) {
        if (colors[next] === -1) {
          dfs(next, color ^ 1)
        } else {
          if (colors[next] === colors[cur]) {
            throw new Error('不是二分图')
          }
        }
      }
    }

    for (let i = 0; i < adjList.length; i++) {
      if (colors[i] === -1) dfs(i, 0)
    }

    return colors
  }
}

function isPerfectMatching(adjList: number[][]): boolean {
  const maxMatching = hungarian(adjList)
  return maxMatching * 2 === adjList.length
}

if (require.main === module) {
  console.log(
    hungarian([
      [1, 3],
      [0, 2],
      [1, 3],
      [0, 2],
    ])
  ) // 2
  console.log(hungarian([[1, 3], [0, 2], [1], [0], [], []])) // 2
}

export { hungarian }

function dfsNonrecursiveOnGraph(graph: number[][], start: number, down: (v: number) => void, up: (child: number, parent: number) => void): void {
  const n = graph.length
  const parent = new Int32Array(n).fill(-1)
  const iter = new Uint32Array(n)
  let cur = start
  parent[cur] = -2
  while (cur >= 0) {
    if (iter[cur] === 0) {
      down(cur) // 首次访问
    }

    if (iter[cur] === graph[cur].length) {
      const p = parent[cur]
      if (p >= 0) {
        up(cur, p) // 回溯
      }
      cur = p
      continue
    }

    const next = graph[cur][iter[cur]++]
    if (parent[next] !== -1) {
      // 返回边
      continue
    }

    // DFS树边
    parent[next] = cur
    cur = next
  }
}

export {}

if (require.main === module) {
  {
    // dfs 求子树大小
    //   0
    //  / \
    // 1   2
    //      \
    //       3
    const graph = [[1, 2], [], [3], []]
    const subsize = new Uint32Array(graph.length)
    dfsNonrecursiveOnGraph(
      graph,
      0,
      v => {
        console.log('down', v)
        subsize[v] = 1
      },
      (child, parent) => {
        console.log('up', child, parent)
        subsize[parent] += subsize[child]
      }
    )

    console.log(subsize)
  }
}

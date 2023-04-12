function useHungarian(graph?: number[][]) {
  const to: number[][] = []
  let row = 0
  let col = 0

  if (graph) {
    const n = graph.length
    const [colors, ok] = isBipartite(n, graph)
    if (!ok) {
      throw new Error('not bipartite')
    }
    for (let i = 0; i < n; i++) {
      if (colors[i] === 0) {
        graph[i].forEach(j => {
          if (colors[j] === 1) {
            addEdge(i, j)
          }
        })
      }
    }
  }

  function addEdge(boy: number, girl: number): void {
    row = Math.max(row, boy + 1)
    col = Math.max(col, girl + 1)
    while (to.length <= boy) {
      to.push([])
    }
    to[boy].push(girl)
  }

  function work(): [boy: number, girl: number][] {
    const pre = new Int32Array(row).fill(-1)
    const root = new Int32Array(row).fill(-1)
    const p = new Int32Array(row).fill(-1)
    const queue = new Int32Array(col).fill(-1)
    let updated = true

    while (updated) {
      updated = false
      const s: number[] = []
      let sFront = 0
      for (let i = 0; i < row; i++) {
        if (p[i] === -1) {
          root[i] = i
          s.push(i)
        }
      }

      while (sFront < s.length) {
        let v = s[sFront]
        sFront++
        if (p[root[v]] !== -1) {
          continue
        }

        for (let u of to[v]) {
          if (queue[u] === -1) {
            while (u !== -1) {
              queue[u] = v
              const tmp = p[v]
              p[v] = u
              u = tmp
              v = pre[v]
            }
            updated = true
            break
          }

          u = queue[u]
          if (pre[u] !== -1) {
            continue
          }

          pre[u] = v
          root[u] = root[v]
          s.push(u)
        }
      }

      if (updated) {
        for (let i = 0; i < row; i++) {
          pre[i] = -1
          root[i] = -1
        }
      }
    }

    const res: [boy: number, girl: number][] = []
    for (let i = 0; i < row; i++) {
      if (p[i] !== -1) {
        res.push([i, p[i]])
      }
    }
    return res
  }

  return {
    addEdge,
    work
  }
}

function isBipartite(n: number, adjList: number[][]): [colors: Int8Array, ok: boolean] {
  const colors = new Int8Array(n).fill(-1)
  for (let i = 0; i < n; i++) {
    if (colors[i] === -1 && !dfs(i, 0)) {
      return [colors, false]
    }
  }
  return [colors, true]

  function dfs(cur: number, color: number): boolean {
    colors[cur] = color
    for (let i = 0; i < adjList[cur].length; i++) {
      const next = adjList[cur][i]
      if (colors[next] === -1) {
        if (!dfs(next, color ^ 1)) {
          return false
        }
      } else if (colors[next] === color) {
        return false
      }
    }
    return true
  }
}

export { useHungarian }

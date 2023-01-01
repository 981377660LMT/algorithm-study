function useHungarian(row: number, col: number) {
  const to: number[][] = Array.from({ length: row }, () => [])

  function addEdge(boy: number, girl: number): void {
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
              ;[p[v], u] = [u, p[v]]
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

export { useHungarian }

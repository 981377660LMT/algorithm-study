function maxCollectedFruits(fruits: number[][]): number {
  const n = fruits.length
  if (n === 0) return 0

  const table1: { [key: number]: [number, number][] } = {
    1: [
      [1, 0],
      [0, 1],
      [1, 1]
    ],
    2: [
      [1, -1],
      [1, 0],
      [1, 1]
    ],
    3: [
      [-1, 1],
      [0, 1],
      [1, 1]
    ]
  }

  const table2: { [key: number]: [number, number] } = {
    1: [0, 0],
    2: [0, n - 1],
    3: [n - 1, 0]
  }

  function generatePermutations(arr: number[]): number[][] {
    const results: number[][] = []

    function bt(path: number[], used: boolean[]) {
      if (path.length === arr.length) {
        results.push([...path])
        return
      }
      for (let i = 0; i < arr.length; i++) {
        if (used[i]) continue
        used[i] = true
        path.push(arr[i])
        bt(path, used)
        path.pop()
        used[i] = false
      }
    }

    bt([], Array(arr.length).fill(false))
    return results
  }

  function solve(num: number, grid: number[][]): { total: number; path: Set<string> } {
    const [startI, startJ] = table2[num]
    const dp: Array<Map<string, number>> = Array.from({ length: n }, () => new Map())
    dp[0].set(`${startI},${startJ}`, grid[startI][startJ])

    for (let step = 0; step < n - 1; step++) {
      const dpNext = new Map<string, number>()
      for (const [position, total] of dp[step]) {
        const [i, j] = position.split(',').map(Number)
        for (const [di, dj] of table1[num]) {
          const ni = i + di
          const nj = j + dj
          if (ni >= 0 && ni < n && nj >= 0 && nj < n) {
            const newTotal = total + grid[ni][nj]
            const newPosition = `${ni},${nj}`
            if (!dpNext.has(newPosition) || dpNext.get(newPosition)! < newTotal) {
              dpNext.set(newPosition, newTotal)
            }
          }
        }
      }
      dp[step + 1] = dpNext
    }

    let maxTotal = 0
    const endPosition = `${n - 1},${n - 1}`
    if (dp[n - 1].has(endPosition)) {
      maxTotal = dp[n - 1].get(endPosition)!

      const dpPrev: Array<Map<string, string | null>> = Array.from({ length: n }, () => new Map())
      dpPrev[0].set(`${startI},${startJ}`, null)

      for (let step = 0; step < n - 1; step++) {
        const dpNext = new Map<string, number>()
        const dpPrevNext = new Map<string, string>()
        for (const [position, total] of dp[step]) {
          const [i, j] = position.split(',').map(Number)
          for (const [di, dj] of table1[num]) {
            const ni = i + di
            const nj = j + dj
            if (ni >= 0 && ni < n && nj >= 0 && nj < n) {
              const newTotal = total + grid[ni][nj]
              const newPosition = `${ni},${nj}`
              if (!dpNext.has(newPosition) || dpNext.get(newPosition)! < newTotal) {
                dpNext.set(newPosition, newTotal)
                dpPrevNext.set(newPosition, position)
              }
            }
          }
        }
        dp[step + 1] = dpNext
        dpPrev[step + 1] = dpPrevNext
      }

      const pathSet = new Set<string>()
      let currentPos = endPosition
      pathSet.add(currentPos)
      for (let step = n - 1; step > 0; step--) {
        const prevPos = dpPrev[step].get(currentPos)!
        if (!prevPos) break
        pathSet.add(prevPos)
        currentPos = prevPos
      }

      return { total: maxTotal, path: pathSet }
    }
    return { total: 0, path: new Set() }
  }

  let res = 0
  const permutations = generatePermutations([1, 2, 3])
  for (const perm of permutations) {
    const gridCopy = fruits.map(row => row.slice())
    let totalFruits = 0
    for (const kidNum of perm) {
      const { total, path } = solve(kidNum, gridCopy)
      totalFruits += total
      for (const coord of path) {
        const [i, j] = coord.split(',').map(Number)
        gridCopy[i][j] = 0
      }
    }
    if (totalFruits > res) {
      res = totalFruits
    }
  }

  return res
}

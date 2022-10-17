// 1 <= m, n <= 200
// 1 <= hits.length <= 4 * 104

enum State {
  Empty = 0,
  BrokenOrLastAlive = 1,
  Stable = 2
}

const DIR4 = [
  [0, 1],
  [0, -1],
  [1, 0],
  [-1, 0]
]

// 逆推:
// 整体思路是倒推，看每次打砖块造成了多少砖块本应该是稳定的，变为了不稳定
// 1.确定最后有哪些砖块剩余(从顶部dfs即可)，做好标记
// 2. 倒推打砖块，如果击打位置没有砖块或者不与稳定砖块相连，无事发生
// 3.如果击打位置的砖块与稳定砖块相连，说明这个砖块是被打下来的,dfs+1并重置那些砖块为稳定砖块
function hitBricks(grid: number[][], hits: number[][]): number[] {
  const [ROW, COL] = [grid.length, grid[0].length]
  const res = Array<number>(hits.length).fill(0)

  // 1. 标记，确定最后哪些存活
  for (const [hitR, hitC] of hits) grid[hitR][hitC]--
  for (let c = 0; c < COL; c++) {
    dfs(0, c)
  }

  // 2.倒推打砖块
  for (let i = hits.length - 1; ~i; i--) {
    const [hitR, hitC] = hits[i]
    // 标记还原
    grid[hitR][hitC]++
    // 3.如果击打位置没有砖块或者不与稳定砖块相连，无事发生；否则dfs重置为稳定
    if (grid[hitR][hitC] === State.Empty || !isConnectToStable(hitR, hitC)) continue
    res[i] = dfs(hitR, hitC) - 1
  }

  return res

  // 将砖块还原成稳定
  function dfs(r: number, c: number): number {
    let res = 0

    if (r >= 0 && r < ROW && c >= 0 && c < COL && grid[r][c] === State.BrokenOrLastAlive) {
      res++
      grid[r][c] = State.Stable
      for (const [dr, dc] of DIR4) {
        const nr = r + dr
        const nc = c + dc
        res += dfs(nr, nc)
      }
    }

    return res
  }

  function isConnectToStable(r: number, c: number): boolean {
    if (r === 0) return true
    for (const [dr, dc] of DIR4) {
      const [nr, nc] = [r + dr, c + dc]
      if (grid[nr]?.[nc] === State.Stable) return true
    }
    return false
  }
}

console.log(
  hitBricks(
    [
      [1, 0, 0, 0],
      [1, 1, 1, 0]
    ],
    [[1, 0]]
  )
)

// 输出：[2]
// 解释：
// 网格开始为：
// [[1,0,0,0]，
//  [1,1,1,0]]
// 消除 (1,0) 处加粗的砖块，得到网格：
// [[1,0,0,0]
//  [0,1,1,0]]
// 两个加粗的砖不再稳定，因为它们不再与顶部相连，也不再与另一个稳定的砖相邻，因此它们将掉落。得到网格：
// [[1,0,0,0],
//  [0,0,0,0]]
// 因此，结果为 [2] 。
export {}

// TODO 打砖块

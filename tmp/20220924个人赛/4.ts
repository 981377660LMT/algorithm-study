// DIR4 = {
//   0: (0, 1),
//   1: (1, 0),
//   2: (0, -1),
//   3: (-1, 0),
// }  # 顺时针

const DIR4 = new Map()
DIR4.set(0, [0, 1])
DIR4.set(1, [1, 0])
DIR4.set(2, [0, -1])
DIR4.set(3, [-1, 0])

function ballGame(num: number, plate: string[]): number[][] {
  const ROW = plate.length
  const COL = plate[0].length
  let queue: [number, number, number, number][] = []
  const visited = Array.from({ length: ROW * COL }, () => new Uint8Array(4))

  for (let i = 0; i < ROW; i++) {
    for (let j = 0; j < COL; j++) {
      if (plate[i][j] === 'O') {
        for (let dir = 0; dir < 4; dir++) {
          queue.push([i, j, dir, 0])
          visited[i * COL + j][dir] = 1
        }
      }
    }
  }

  const res: number[][] = []
  const BAD = new Set([0, COL - 1, (ROW - 1) * COL, (ROW - 1) * COL + COL - 1])

  while (queue.length) {
    const nextQueue: [number, number, number, number][] = []
    const { length } = queue
    for (let i = 0; i < length; i++) {
      const [curRow, curCol, curDir, curStep] = queue.pop()!

      if (curStep > num) continue

      let nextDir: number
      if (plate[curRow][curCol] === 'W') {
        nextDir = (curDir + 1) % 4 // 顺时针
      } else if (plate[curRow][curCol] === 'E') {
        nextDir = (curDir - 1 + 4) % 4 // 逆时针
      } else {
        nextDir = curDir
      }

      const [dr, dc] = DIR4.get(nextDir)!
      const [nextRow, nextCol] = [curRow + dr, curCol + dc]
      if (nextRow < 0 || nextRow >= ROW || nextCol < 0 || nextCol >= COL) {
        // !四个角除外
        const hash = curRow * COL + curCol
        if (!BAD.has(hash) && plate[curRow][curCol] === '.') {
          res.push([curRow, curCol])
        }
      } else {
        const hash = nextRow * COL + nextCol
        if (!visited[hash][nextDir] && curStep + 1 <= num) {
          nextQueue.push([nextRow, nextCol, nextDir, curStep + 1])
          visited[hash][nextDir] = 1
        }
      }
    }

    queue = nextQueue
  }

  return res
}

console.log(ballGame(4, ['..E.', '.EOW', '..W.']))
export {}

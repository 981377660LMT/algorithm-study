/* eslint-disable eqeqeq */

const INF = 2e15

/**
 * bfs求最无权图短路.
 * @param n 节点数.
 * @param adjList 邻接表.
 * @param start 起点.
 * @param end 终点.
 */
function bfs(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  start: number | ArrayLike<number>
): number[]
function bfs(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  start: number | ArrayLike<number>,
  end: number
): number
function bfs(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  start: number | ArrayLike<number>,
  end?: number
): number | number[] {
  const dist = Array(n).fill(INF)
  const queue = new Uint32Array(n)
  let left = 0
  let right = 0
  start = Array.isArray(start) ? start : [start]
  for (let i = 0; i < start.length; ++i) {
    dist[start[i]] = 0
    queue[right++] = start[i]
  }

  if (end == undefined) {
    while (left < right) {
      const cur = queue[left++]
      const nexts = adjList[cur]
      for (let i = 0; i < nexts.length; ++i) {
        const next = nexts[i]
        const cand = dist[cur] + 1
        if (cand < dist[next]) {
          dist[next] = cand
          queue[right++] = next
        }
      }
    }
    return dist
  }

  while (left < right) {
    const cur = queue[left++]
    if (cur === end) return dist[cur]
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; ++i) {
      const next = nexts[i]
      const cand = dist[cur] + 1
      if (cand < dist[next]) {
        dist[next] = cand
        queue[right++] = next
      }
    }
  }
  return INF
}

/**
 * bfs求起点到终点的最短距离和路径.
 */
function bfsPath(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  start: number,
  end: number
): [dist: number, path: number[]] {
  const dist = Array(n).fill(INF)
  dist[start] = 0
  const pre = new Int32Array(n).fill(-1)
  const queue = new Uint32Array(n)
  let left = 0
  let right = 0
  queue[right++] = start
  while (left < right) {
    const cur = queue[left++]
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; ++i) {
      const next = nexts[i]
      const cand = dist[cur] + 1
      if (cand < dist[next]) {
        dist[next] = cand
        pre[next] = cur
        queue[right++] = next
      }
    }
  }

  if (dist[end] === INF) return [INF, []]
  const path: number[] = []
  let cur = end
  while (pre[cur] !== -1) {
    path.push(cur)
    cur = pre[cur]
  }
  path.push(start)
  path.reverse()
  return [dist[end], path]
}

/**
 * 返回距离start为dist的结点.
 */
function bfsDepth(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  start: number,
  dist: number
): number[] {
  if (dist < 0) return []
  if (!dist) return [start]
  let queue: number[] = [start]
  const visited = new Uint8Array(n)
  let todo = dist
  while (queue.length && todo) {
    const len_ = queue.length
    const nextQueue: number[] = []
    for (let i = 0; i < len_; i++) {
      const cur = queue[i]
      const nexts = adjList[cur]
      for (let j = 0; j < nexts.length; j++) {
        const next = nexts[j]
        if (!visited[next]) {
          visited[next] = 1
          nextQueue.push(next)
        }
      }
    }
    todo--
    queue = nextQueue
  }
  return queue
}

/**
 * 网格图bfs, 返回每个格子到起点的最短距离.
 */
function bfsGrid(row: number, col: number, start: ArrayLike<[r: number, c: number]>): number[][] {
  const DIR4 = [
    [0, 1],
    [0, -1],
    [1, 0],
    [-1, 0]
  ]

  const dist: number[][] = Array(row)
  for (let i = 0; i < row; i++) dist[i] = Array(col).fill(INF)
  let queue: [r: number, c: number][] = Array(start.length)
  for (let i = 0; i < start.length; i++) {
    const { 0: r, 1: c } = start[i]
    dist[r][c] = 0
    queue[i] = [r, c]
  }

  while (queue.length) {
    const len = queue.length
    const nextQueue: [r: number, c: number][] = []
    for (let i = 0; i < len; i++) {
      const { 0: curR, 1: curC } = queue[i]
      for (let j = 0; j < 4; j++) {
        const nextR = curR + DIR4[j][0]
        const nextC = curC + DIR4[j][1]
        const cand = dist[curR][curC] + 1
        if (nextR >= 0 && nextR < row && nextC >= 0 && nextC < col && cand < dist[nextR][nextC]) {
          dist[nextR][nextC] = cand
          nextQueue.push([nextR, nextC])
        }
      }
    }
    queue = nextQueue
  }

  return dist
}

export { bfs, bfsPath, bfsDepth, bfsGrid }

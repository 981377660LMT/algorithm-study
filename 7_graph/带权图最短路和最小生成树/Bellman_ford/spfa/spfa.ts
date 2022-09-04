/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

type AdjList = [next: number, weight: number][][]

const INF = 2e15

/**
 * spfa求以`虚拟节点`为起点的单源`最长路` 并检测正环
 */
function spfa1(n: number, adjList: AdjList): [ok: boolean, dist: number[]] {
  const dist = Array<number>(n).fill(0)
  const inQueue = new Uint8Array(n).fill(1)
  const count = new Uint32Array(n)
  let queue = Array<number>(n)
    .fill(0)
    .map((_, i) => i)

  while (queue.length) {
    const nextQueue: number[] = []
    const step = queue.length

    for (let _ = 0; _ < step; _++) {
      const cur = queue.pop()!
      inQueue[cur] = 0
      for (let i = 0; i < adjList[cur].length; i++) {
        const [next, weight] = adjList[cur][i]
        const cand = dist[cur] + weight
        if (cand > dist[next]) {
          dist[next] = cand
          count[next] = count[cur] + 1
          if (count[next] >= n) return [false, []]
          if (!inQueue[next]) {
            inQueue[next] = 1
            nextQueue.push(next)
          }
        }
      }
    }

    queue = nextQueue
  }

  return [true, dist]
}

/**
 * spfa求单源`最短路`(图中有负边无负环)
 *
 * 适用于解决带有负权重的图,是Bellman-ford的常数优化版
 */
function spfa2(n: number, adjList: AdjList, start: number): number[] {
  const dist = Array<number>(n).fill(INF)
  dist[start] = 0
  let queue = [start]
  const inQueue = new Uint8Array(n)
  inQueue[start] = 1

  while (queue.length) {
    const nextQueue: number[] = []
    const step = queue.length

    for (let _ = 0; _ < step; _++) {
      const cur = queue.pop()!
      inQueue[cur] = 0
      for (let i = 0; i < adjList[cur].length; i++) {
        const next = adjList[cur][i][0]
        const weight = adjList[cur][i][1]
        const cand = dist[cur] + weight
        if (cand < dist[next]) {
          dist[next] = cand
          if (!inQueue[next]) {
            inQueue[next] = 1
            nextQueue.push(next)
          }
        }
      }
    }

    queue = nextQueue
  }

  return dist
}

if (require.main === module) {
  // LCP 32. 批量处理任务
  // https://leetcode.cn/problems/t3fKg1/
  // !TLE 因为需要去除重边
  function processTasks(tasks: number[][]): number {
    const allNums = new Set<number>()
    tasks.forEach(([start, end]) => {
      allNums.add(start - 1)
      allNums.add(end)
    })

    const nums = [...allNums].sort((a, b) => a - b)
    const n = nums.length
    const mp = new Map<number, number>()
    for (let i = 0; i < n; i++) {
      mp.set(nums[i], i)
    }

    const adjList: AdjList = Array.from({ length: n }, () => [])
    tasks.forEach(([start, end, period]) => {
      const u = mp.get(start - 1)!
      const v = mp.get(end)!
      adjList[u].push([v, period]) // v - u >= period
    })

    for (let i = 1; i < n; i++) {
      adjList[i - 1].push([i, 0]) // i - (i-1) >= 0
      adjList[i].push([i - 1, nums[i - 1] - nums[i]]) // i - (i-1) <= nums[i] - nums[i-1]
    }

    const [ok, dist] = spfa1(n, adjList)
    if (!ok) return -1
    return dist[n - 1]
  }
}

export { spfa1, spfa2 }

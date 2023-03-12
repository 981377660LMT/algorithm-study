/* eslint-disable consistent-return */
/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

const INF = 2e15

/**
 * spfa求单源`最短路`
 * 如果有负环, 返回`null
 *
 * 适用于解决带有负权边的图,是Bellman-ford的常数优化版
 */
function spfa(
  n: number,
  adjList: [next: number, weight: number][][],
  start: number
): [dist: number[], ok: boolean] {
  const dist = Array<number>(n).fill(INF)
  const inQueue = new Uint8Array(n)
  const relaxedConut = new Uint32Array(n)
  let queue = [start] // !TODO:queue用长为n的数组+left,right指针实现
  let ql = 0
  let qr = 1
  dist[start] = 0
  inQueue[start] = 1
  relaxedConut[start] = 1

  while (queue.length) {
    const nextQueue: number[] = []
    const step = queue.length
    for (let _ = 0; _ < step; _++) {
      const cur = queue.pop()!
      inQueue[cur] = 0
      adjList[cur].forEach(([next, weight]) => {
        const cand = dist[cur] + weight
        // !如果要最长路这里需要改成 >
        if (cand < dist[next]) {
          dist[next] = cand
          if (!inQueue[next]) {
            relaxedConut[next]++
            if (relaxedConut[next] >= n) {
              // 找到从起点出发可达的负环
              return [[], false]
            }
            inQueue[next] = 1
            nextQueue.push(next)
          }
        }
      })
    }

    queue = nextQueue
  }

  return [dist, true]
}

export { spfa }

/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/** 拓扑排序环检测. */
function hasCycle(n: number, adjList: ArrayLike<ArrayLike<number>>, directed = true): boolean {
  const startDeg = directed ? 0 : 1
  const deg = new Uint32Array(n)
  if (directed) {
    for (let i = 0; i < n; i++) {
      const nexts = adjList[i]
      for (let j = 0; j < nexts.length; j++) {
        deg[nexts[j]]++
      }
    }
  } else {
    for (let i = 0; i < n; i++) {
      deg[i] = adjList[i].length
    }
  }

  let count = 0
  const queue = new Uint32Array(n)
  let head = 0
  let tail = 0
  for (let i = 0; i < n; i++) {
    if (deg[i] === startDeg) {
      queue[tail++] = i
    }
  }

  while (head < tail) {
    const cur = queue[head++]
    count++
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (--deg[next] === startDeg) {
        queue[tail++] = next
      }
    }
  }

  return count < n
}

/** 拓扑排序求方案. */
function topoSort(n: number, adjList: ArrayLike<ArrayLike<number>>, directed = true): [order: number[], hasCycle: boolean] {
  const startDeg = directed ? 0 : 1
  const deg = new Uint32Array(n)
  if (directed) {
    for (let i = 0; i < n; i++) {
      const nexts = adjList[i]
      for (let j = 0; j < nexts.length; j++) {
        deg[nexts[j]]++
      }
    }
  } else {
    for (let i = 0; i < n; i++) {
      deg[i] = adjList[i].length
    }
  }

  const queue = new Uint32Array(n)
  let head = 0
  let tail = 0
  for (let i = 0; i < n; i++) {
    if (deg[i] === startDeg) {
      queue[tail++] = i
    }
  }

  const order: number[] = []
  while (head < tail) {
    const cur = queue[head++]
    order.push(cur)
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (--deg[next] === startDeg) {
        queue[tail++] = next
      }
    }
  }

  return order.length < n ? [[], false] : [order, true]
}

export { hasCycle, topoSort }

if (require.main === module) {
  // https://leetcode.cn/problems/course-schedule/description/
  // 207. 课程表
  function canFinish(numCourses: number, prerequisites: number[][]): boolean {
    const n = numCourses
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    for (let i = 0; i < prerequisites.length; i++) {
      const [a, b] = prerequisites[i]
      adjList[b].push(a)
    }

    return !hasCycle(n, adjList)
  }
}

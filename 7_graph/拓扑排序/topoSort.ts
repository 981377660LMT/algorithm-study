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
function topoSort(
  n: number,
  adjList: ArrayLike<ArrayLike<number>>,
  directed = true
): [order: number[], hasCycle: boolean] {
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

/**
 * 拓扑排序求方案.
 */
function topoSortMap<T extends PropertyKey>(
  allVertices: Set<T>,
  edges: [from: T, to: T][],
  directed = true
): T[] | undefined {
  edges.forEach(([from, to]) => {
    if (!allVertices.has(from) || !allVertices.has(to)) {
      throw new Error('Invalid vertex')
    }
  })

  const deg = new Map<T, number>()
  const graph = new Map<T, T[]>()
  allVertices.forEach(v => {
    deg.set(v, 0)
    graph.set(v, [])
  })

  const addDirectedEdge = (from: T, to: T): void => {
    deg.set(to, deg.get(to)! + 1)
    graph.get(from)!.push(to)
  }

  if (directed) {
    edges.forEach(([from, to]) => {
      addDirectedEdge(from, to)
    })
  } else {
    edges.forEach(([from, to]) => {
      addDirectedEdge(from, to)
      addDirectedEdge(to, from)
    })
  }

  const startDeg = directed ? 0 : 1
  let queue: T[] = []
  allVertices.forEach(v => {
    if (deg.get(v) === startDeg) {
      queue.push(v)
    }
  })

  const order: T[] = []
  while (queue.length) {
    const nextQueue: T[] = []
    queue.forEach(v => {
      order.push(v)
      const nexts = graph.get(v)!
      nexts.forEach(next => {
        deg.set(next, deg.get(next)! - 1)
        if (deg.get(next) === startDeg) {
          nextQueue.push(next)
        }
      })
    })
    queue = nextQueue
  }

  return order.length < allVertices.size ? undefined : order
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

  function canFinish2(numCourses: number, prerequisites: number[][]): boolean {
    const allVertices = new Set<number>()
    for (let i = 0; i < numCourses; i++) {
      allVertices.add(i)
    }
    const edges: [from: number, to: number][] = []
    for (let i = 0; i < prerequisites.length; i++) {
      const [a, b] = prerequisites[i]
      edges.push([a, b])
    }
    return !!topoSortMap(allVertices, edges)
  }
}

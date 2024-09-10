/* eslint-disable no-constant-condition */

/**
 * 双向bfs.
 */
function biBfs<S extends string | number>(
  start: S,
  target: S,
  getNextStates: (curState: S) => S[]
): number {
  let queue1 = new Set<S>([start])
  let queue2 = new Set<S>([target])
  const visited = new Set<S>()
  let steps = 0

  while (queue1.size && queue2.size) {
    if (queue1.size > queue2.size) {
      const tmp = queue1
      queue1 = queue2
      queue2 = tmp
    }

    const nextQueue = new Set<S>()
    for (const cur of queue1) {
      if (queue2.has(cur)) return steps
      if (visited.has(cur)) continue
      visited.add(cur)
      getNextStates(cur).forEach(next => {
        nextQueue.add(next)
      })
    }

    steps++
    queue1 = queue2
    queue2 = nextQueue
  }

  return -1
}

/**
 * 双向bfs，返回路径.
 */
function biBfsPath<S extends number | string>(
  start: S,
  target: S,
  getNextStates: (cur: S) => S[]
): S[] {
  const queue = [new Set<S>([start]), new Set<S>([target])]
  const pre = [new Map<S, S>(), new Map<S, S>()]
  const visited = [new Set<S>([start]), new Set<S>([target])]

  let curQueue: Set<S>
  let curVisited: Set<S>
  let curPre: Map<S, S>
  let otherQueue: Set<S>
  while (queue[0].size && queue[1].size) {
    const qi = +(queue[0].size > queue[1].size)

    const nextQueue = new Set<S>()
    curQueue = queue[qi]
    curVisited = visited[qi]
    curPre = pre[qi]
    otherQueue = queue[qi ^ 1]

    for (const cur of curQueue) {
      if (otherQueue.has(cur)) return restorePath(cur)
      for (const next of getNextStates(cur)) {
        if (curVisited.has(next)) continue
        curVisited.add(next)
        nextQueue.add(next)
        curPre.set(next, cur)
      }
    }

    queue[qi] = nextQueue
  }

  return []

  function restorePath(mid: S): S[] {
    const pre1 = pre[0]
    const path1 = [mid]
    let cur1 = mid
    while (true) {
      const p = pre1.get(cur1)
      if (p === undefined) break
      cur1 = p
      path1.push(cur1)
    }

    const pre2 = pre[1]
    const path2: S[] = []
    let cur2 = mid
    while (true) {
      const p = pre2.get(cur2)
      if (p === undefined) break
      cur2 = p
      path2.push(cur2)
    }

    path1.reverse()
    path1.push(...path2)
    return path1
  }
}

export { biBfs, biBfsPath }

interface GetNextState<State> {
  (curState: State): State[]
}

/**
 *
 * @param start
 * @param target
 * @param getNextState
 * @returns start 到 target 的最短距离
 *
 */
function bibfs<State = string>(
  start: State,
  target: State,
  getNextState: GetNextState<State>
): number {
  // 用集合不用队列，可以快速判断元素是否存在
  let queue1 = new Set<State>([start])
  let queue2 = new Set<State>([target])
  const visited = new Set()
  let steps = 0

  while (queue1.size && queue2.size) {
    if (queue1.size > queue2.size) {
      ;[queue1, queue2] = [queue2, queue1]
    }

    // 本层搜出来的结果
    const nextQueue = new Set<State>()

    for (const cur of queue1) {
      if (queue2.has(cur)) return steps

      if (visited.has(cur)) continue
      visited.add(cur)

      for (const next of getNextState(cur)) {
        nextQueue.add(next)
      }
    }

    steps++
    ;[queue1, queue2] = [queue2, nextQueue]
  }

  return -1
}

export { bibfs }
// 双向 BFS 在无解的情况下不如单向 BFS
// 们可以先使用「并查集」进行预处理，判断「起点」和「终点」是否连通，如果不联通，直接返回 -1−1，有解才调用双向 BFS
// 使用「并查集」预处理的复杂度与建图是近似的，增加这样的预处理并不会越过我们时空复杂度的上限，因此这样的预处理是有益的

interface WindowManager<T, Q> {
  // 使用 this:unknown 禁止在外部调用this
  add(this: unknown, value: T, index: number, qLeft: number, qRight: number): void
  remove(this: unknown, value: T, index: number, qLeft: number, qRight: number): void
  query(this: unknown, index: number, qLeft: number, qRight: number): Q
}

/**
 * 静态查询区间的莫队算法
 *
 * @param arrayLike 查询的原数据
 * @param windowManager 区间的维护方式
 * @complexity `O(n*sqrt(q))`
 */
function useMoAlgo<T, Q>(arrayLike: Readonly<ArrayLike<T>>, windowManager: WindowManager<T, Q>) {
  const queries: [qIndex: number, qLeft: number, qRight: number][] = []

  /**
   * 添加查询区间
   *
   * 0 <= left <= right < {@link arrayLike}.length
   */
  function addQuery(left: number, right: number): void {
    queries.push([queries.length, left, right + 1]) // 注意这里的 right + 1
  }

  /**
   * 返回每个查询的结果
   */
  function work(): Q[] {
    sortQueries()

    const res: Q[] = Array(queries.length).fill(null)
    let left = 0
    let right = 0

    // const { add, remove, query } = windowManager // !注意解构会使this指向不正确
    const add = windowManager.add.bind(windowManager)
    const remove = windowManager.remove.bind(windowManager)
    const query = windowManager.query.bind(windowManager)

    for (let i = 0; i < queries.length; i++) {
      // 不使用解构来加速
      const qIndex = queries[i][0]
      const qLeft = queries[i][1]
      const qRight = queries[i][2]

      // !窗口收缩
      while (right > qRight) {
        right--
        remove(arrayLike[right], right, qLeft, qRight - 1)
      }

      while (left < qLeft) {
        remove(arrayLike[left], left, qLeft, qRight - 1)
        left++
      }

      // !窗口扩张
      while (right < qRight) {
        add(arrayLike[right], right, qLeft, qRight - 1)
        right++
      }

      while (left > qLeft) {
        left--
        add(arrayLike[left], left, qLeft, qRight - 1)
      }

      res[qIndex] = query(left, qLeft, qRight - 1)
    }

    return res
  }

  function sortQueries(): void {
    const chunkSize = Math.max(1, Math.floor(arrayLike.length / Math.sqrt(queries.length)))
    queries.sort(
      (q1, q2) => Math.floor(q1[1] / chunkSize) - Math.floor(q2[1] / chunkSize) || q1[2] - q2[2]
    )
  }

  return {
    addQuery,
    work
  }
}

export { useMoAlgo, WindowManager }

interface WindowManager<T, Q> {
  // 使用 this:unknown 禁止在外部调用this
  add(this: unknown, value: T, index: number, qLeft: number, qRight: number): void
  remove(this: unknown, value: T, index: number, qLeft: number, qRight: number): void
  query(this: unknown, index: number, qLeft: number, qRight: number): Q
}

type Query = [qIndex: number, qLeft: number, qRight: number]

/**
 * 静态查询区间的莫队算法
 *
 * @param arrayLike 查询的原数据
 * @param windowManager 区间的维护方式
 * @complexity `O(n*sqrt(q))`
 *
 * 左端点分桶，右端点排序
 */
function useMoAlgo<T, Q>(arrayLike: Readonly<ArrayLike<T>>, windowManager: WindowManager<T, Q>) {
  const n = arrayLike.length
  const chunkSize = Math.ceil(Math.sqrt(n)) // const chunkSize = Math.ceil(n / Math.sqrt(2 * q))
  const buckets = Array.from<unknown, Query[]>({ length: Math.floor(n / chunkSize) + 1 }, () => [])
  let queryOrder = 0

  /**
   * 添加查询区间
   *
   * 0 <= left <= right < {@link arrayLike}.length
   */
  function addQuery(left: number, right: number): void {
    const index = Math.floor(left / chunkSize)
    buckets[index].push([queryOrder++, left, right + 1]) // 注意这里的 right + 1
  }

  /**
   * 返回每个查询的结果
   */
  function work(): Q[] {
    for (let i = 0; i < buckets.length; i++) {
      buckets[i].sort((a, b) => (i & 1 ? -(a[2] - b[2]) : a[2] - b[2])) // 块内按区间右端点排序
    }

    const res: Q[] = Array(queryOrder).fill(null)
    let left = 0
    let right = 0

    // const { add, remove, query } = windowManager // !注意解构会使this指向不正确
    const add = windowManager.add.bind(windowManager)
    const remove = windowManager.remove.bind(windowManager)
    const query = windowManager.query.bind(windowManager)

    for (let i = 0; i < buckets.length; i++) {
      const bucket = buckets[i]
      for (let j = 0; j < bucket.length; j++) {
        // 不使用解构来加速
        const qIndex = bucket[j][0]
        const qLeft = bucket[j][1]
        const qRight = bucket[j][2]

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
    }

    return res
  }

  return {
    addQuery,
    work
  }
}

export { useMoAlgo, WindowManager }

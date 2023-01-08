interface WindowManager<Q> {
  // 使用 `this:void` 禁止在外部调用this
  add(this: void, index: number, delta: -1 | 1): void
  remove(this: void, index: number, delta: -1 | 1): void
  query(this: void, qLeft: number, qRight: number): Q
}

type Query = [index: number, left: number, right: number]

/**
 * 静态查询区间的莫队算法
 * 左端点分桶，右端点排序
 *
 * @param n 数组长度
 * @param q 查询区间个数
 * @param windowManager 窗口管理器
 * @complexity `O(n*sqrt(q))`
 */
function useMoAlgo<Q>(n: number, q: number, windowManager: WindowManager<Q>) {
  const isqrt = Math.floor(Math.sqrt(q))
  const chunkSize = Math.max(1, Math.floor(n / isqrt))
  const buckets = Array.from<unknown, Query[]>({ length: Math.floor(n / chunkSize) + 1 }, () => [])
  let queryOrder = 0

  /**
   * 添加查询区间
   *
   * 0 <= left <= right < {@link n}
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

    const { add, remove, query } = windowManager // !不使用bind,减小开销

    for (let i = 0; i < buckets.length; i++) {
      const bucket = buckets[i]
      for (let j = 0; j < bucket.length; j++) {
        // 不使用解构来加速
        const qIndex = bucket[j][0]
        const qLeft = bucket[j][1]
        const qRight = bucket[j][2]

        // !窗口扩张
        while (left > qLeft) {
          left--
          add(left, -1)
        }
        while (right < qRight) {
          add(right, 1)
          right++
        }

        // !窗口收缩
        while (left < qLeft) {
          remove(left, 1)
          left++
        }
        while (right > qRight) {
          right--
          remove(right, -1)
        }

        res[qIndex] = query(qLeft, qRight - 1)
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

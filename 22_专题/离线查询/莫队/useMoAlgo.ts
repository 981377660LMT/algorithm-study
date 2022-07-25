<<<<<<< HEAD
interface WindowManager<T, Q> {
  add(value: T, index: number, qLeft: number, qRight: number): void
  remove(value: T, index: number, qLeft: number, qRight: number): void
  query(): Q
}

=======
>>>>>>> 2d83813011271daf312944be57f849d59f1ccf70
/**
 * 静态查询区间的莫队算法
 *
 * @param data 原始数据
<<<<<<< HEAD
 * @param windowManager 区间的维护方式
 */
function useMoAlgo<T = number, Q = number>(data: T[], windowManager: WindowManager<T, Q>) {
  const queries: [qIndex: number, qLeft: number, qRight: number][] = []

  /**
   * 添加查询区间
   *
   * 0 <= left <= right < {@link data}.length
   */
  function addQuery(left: number, right: number) {
    queries.push([queries.length, left, right + 1]) // 注意这里的 right + 1
  }

  /**
   * 返回每个查询的结果
   */
  function work(): Q[] {
    sortQueries()

    const res: Q[] = Array(queries.length).fill(undefined)
    let [left, right] = [0, 0]
    const { add, remove, query } = windowManager

    for (let i = 0; i < queries.length; i++) {
      const [qIndex, qLeft, qRight] = queries[i]

      // !窗口收缩
      while (right > qRight) {
        right--
        remove(data[right], right, qLeft, qRight - 1)
      }

      while (left < qLeft) {
        remove(data[left], left, qLeft, qRight - 1)
        left++
      }

      // !窗口扩张
      while (right < qRight) {
        add(data[right], right, qLeft, qRight - 1)
        right++
      }

      while (left > qLeft) {
        left--
        add(data[left], left, qLeft, qRight - 1)
      }

      res[qIndex] = query()
    }

    return res
  }

  function sortQueries(): void {
    const chunkSize = Math.max(1, Math.floor(data.length / Math.sqrt(queries.length)))
    queries.sort(
      ([_q1Index, q1Left, q1Right], [_q2Index, q2Left, q2Right]) =>
        Math.floor(q1Left / chunkSize) - Math.floor(q2Left / chunkSize) || q1Right - q2Right
    )
  }

  return {
    addQuery,
    work,
  }
}

export { useMoAlgo, WindowManager }
=======
 */
function useMoAlgo<T = number, Q = number>(data: T[]) {
  function addQuery(params: type) {}
  return {
    addQuery,
  }
}

export { useMoAlgo }
>>>>>>> 2d83813011271daf312944be57f849d59f1ccf70

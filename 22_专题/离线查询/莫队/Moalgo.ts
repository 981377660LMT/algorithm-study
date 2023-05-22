class MoAlgo {
  private readonly _chunkSize: number
  private readonly _buckets: [qi: number, left: number, right: number][][]
  private _queryOrder = 0

  constructor(n: number, q: number) {
    const sqrt = Math.sqrt((q * 2) / 3) | 0
    const chunkSize = Math.max(1, (n / Math.max(1, sqrt)) | 0)
    const buckets = Array(((n / chunkSize) | 0) + 1)
    for (let i = 0; i < buckets.length; i++) {
      buckets[i] = []
    }
    this._chunkSize = chunkSize
    this._buckets = buckets
  }

  /**
   * 添加一个查询，查询范围为`左闭右开区间` [left, right).
   * 0 <= left <= right <= n
   */
  addQuery(left: number, right: number): void {
    const index = (left / this._chunkSize) | 0
    this._buckets[index].push([this._queryOrder, left, right])
    this._queryOrder++
  }

  /**
   * 返回每个查询的结果.
   * @param add 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
   * @param remove 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
   * @param query 查询窗口内的数据.
   */
  run(
    add: (index: number, delta: -1 | 1) => void,
    remove: (index: number, delta: -1 | 1) => void,
    query: (qid: number) => void
  ): void {
    let left = 0
    let right = 0

    this._buckets.forEach((bucket, i) => {
      if (i & 1) {
        bucket.sort((a, b) => a[2] - b[2])
      } else {
        bucket.sort((a, b) => b[2] - a[2])
      }

      bucket.forEach(([qi, ql, qr]) => {
        // !窗口扩张
        while (left > ql) {
          add(--left, -1)
        }
        while (right < qr) {
          add(right++, 1)
        }

        // !窗口收缩
        while (left < ql) {
          remove(left++, 1)
        }
        while (right > qr) {
          remove(--right, -1)
        }

        query(qi)
      })
    })
  }
}

export { MoAlgo }

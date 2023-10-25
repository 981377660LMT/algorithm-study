class MoAlgo {
  private readonly _chunkSize: number
  private readonly _buckets: { qi: number; left: number; right: number }[][]
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
   * 添加一个查询，查询范围为`左闭右开区间` [start, end).
   * 0 <= start <= end <= n
   */
  addQuery(start: number, end: number): void {
    const index = (start / this._chunkSize) | 0
    this._buckets[index].push({ qi: this._queryOrder, left: start, right: end })
    this._queryOrder++
  }

  run(
    addLeft: (index: number) => void,
    addRight: (index: number) => void,
    removeLeft: (index: number) => void,
    removeRight: (index: number) => void,
    query: (qid: number) => void
  ): void {
    let left = 0
    let right = 0

    this._buckets.forEach((bucket, i) => {
      if (i & 1) {
        bucket.sort((a, b) => a.right - b.right)
      } else {
        bucket.sort((a, b) => b.right - a.right)
      }

      bucket.forEach(({ qi, left: ql, right: qr }) => {
        // !窗口扩张
        while (left > ql) {
          addLeft(--left)
        }
        while (right < qr) {
          addRight(right++)
        }

        // !窗口收缩
        while (left < ql) {
          removeLeft(left++)
        }
        while (right > qr) {
          removeRight(--right)
        }

        query(qi)
      })
    })
  }
}

export { MoAlgo }

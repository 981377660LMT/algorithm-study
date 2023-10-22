/**
 * 分治法处理离线RMQ问题.
 */
function offlineRMQDivideAndConquer<E>(
  arr: ArrayLike<E>,
  queries: ArrayLike<[start: number, end: number]>,
  e: () => E,
  op: (a: E, b: E) => E
): E[] {
  const n = arr.length
  const q = queries.length
  const res: E[] = Array(q).fill(e())
  const qid: number[] = Array(q)
  for (let i = 0; i < q; ++i) qid[i] = i
  const leftRes: E[] = Array(n).fill(e())
  const rightRes: E[] = Array(n).fill(e())

  const dfs = (left: number, right: number, qid: number[]): void => {
    if (!qid.length) return
    if (left === right) {
      for (let i = 0; i < qid.length; ++i) {
        res[qid[i]] = arr[left]
      }
      return
    }

    const mid = (left + right) >>> 1
    leftRes[mid] = arr[mid]
    for (let i = mid - 1; i >= left; --i) leftRes[i] = op(arr[i], leftRes[i + 1])
    rightRes[mid + 1] = arr[mid + 1]
    for (let i = mid + 2; i <= right; ++i) rightRes[i] = op(rightRes[i - 1], arr[i])
    const todo: number[][] = [[], []]
    for (let i = 0; i < qid.length; ++i) {
      const id = qid[i]
      const { 0: a, 1: b } = queries[id]
      if (a <= mid && mid < b - 1) {
        res[id] = op(op(leftRes[a], arr[mid]), rightRes[b - 1])
        continue
      }
      todo[+(a > mid)].push(id)
    }
    dfs(left, mid, todo[0])
    dfs(mid + 1, right, todo[1])
  }

  dfs(0, n - 1, qid)
  return res
}

export {}

if (require.main === module) {
  const res = offlineRMQDivideAndConquer(
    [1, 2, 3, 4, 5],
    [
      [0, 5],
      [1, 3],
      [2, 3]
    ],
    () => 0,
    Math.max
  )
  console.log(res)
}

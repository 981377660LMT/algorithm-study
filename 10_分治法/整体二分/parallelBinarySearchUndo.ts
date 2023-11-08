/**
 * !如果undo操作(nlong次)比reset(logn次)更优，则可以使用这个版本.
 * 整体二分解决这样一类问题:
 *  - 给定一个长度为n的操作序列, 按顺序执行这些操作;
 *  - 给定q个查询,每个查询形如:"条件qi为真(满足条件)是在第几次操作之后?".
 *  !要求对条件为真的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真.
 *
 * @param n 操作序列的长度.
 * @param q 查询的个数.
 * @param mutate 执行第`mutationId`次操作.一共调用`O(nlogn)`次.
 * @param undo 撤销上一次`mutate`操作.一共调用`O(nlogn)`次.
 * @param predicate 判断第`queryId`次查询是否满足条件.一共调用`O(qlogn)`次.
 * @returns
 *  - `-1` => 不需要操作就满足条件的查询.
 *  - `[0, n)` => 满足条件的最早的操作的编号(0-based).
 *  - `n` => 执行完所有操作后都不满足条件的查询.
 *
 * @see https://betrue12.hateblo.jp/entry/2019/08/14/152227
 */
function parallelBinarySearchUndo(
  n: number,
  q: number,
  options: {
    mutate: (mutationId: number) => void
    undo: () => void
    predicate: (queryId: number) => boolean
  } & ThisType<void>
): Int32Array {
  const { mutate, undo, predicate } = options
  const left = new Int32Array(q)
  const right = new Int32Array(q).fill(n)

  // 不需要操作就满足条件的查询
  for (let i = 0; i < q; i++) {
    if (predicate(i)) {
      right[i] = -1
    }
  }

  while (true) {
    const mids = new Int32Array(q).fill(-1)
    for (let i = 0; i < q; i++) {
      if (left[i] <= right[i]) {
        mids[i] = (left[i] + right[i]) >>> 1
      }
    }

    // csr 数组保存二元对 (qi,mid).
    const indeg = new Uint32Array(n + 2)
    for (let i = 0; i < q; i++) {
      const mid = mids[i]
      if (mid !== -1) {
        indeg[mid + 1]++
      }
    }
    for (let i = 0; i < indeg.length - 1; i++) {
      indeg[i + 1] += indeg[i]
    }
    const total = indeg[indeg.length - 1]
    if (total === 0) {
      break
    }

    const counter = indeg.slice()
    const csr = new Uint32Array(total)
    for (let i = 0; i < q; i++) {
      const mid = mids[i]
      if (mid !== -1) {
        csr[counter[mid]++] = i
      }
    }

    let times = 0
    for (let i = 0; i < csr.length; i++) {
      const pos = csr[i]
      while (times < mids[pos]) {
        mutate(times)
        times++
      }
      if (predicate(pos)) {
        right[pos] = times - 1
      } else {
        left[pos] = times + 1
      }
    }

    for (let _ = 0; _ < times; _++) {
      undo()
    }
  }

  return right
}

export { parallelBinarySearchUndo }

if (require.main === module) {
  let curSum = 0
  const history = Array<number>(10).fill(0)
  let head = 0
  const res = parallelBinarySearchUndo(10, 10, {
    mutate(mutationId) {
      history[head++] = curSum
      curSum += mutationId + 1
    },
    undo() {
      curSum = history[--head]
    },
    predicate(queryId) {
      return curSum >= 55
    }
  })

  console.log(res)
  console.log(1)
}

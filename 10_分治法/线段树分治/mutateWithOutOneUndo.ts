/* eslint-disable no-inner-declarations */

/**
 * 分治删点.类似`除自身以外数组的乘积`.
 * 调用 `query` 时，`state` 为对除了 `index` 以外所有点均调用过了 `mutate` 的状态。但不保证调用 `mutate` 的顺序。
 * 总计会调用 $O(NlgN)$ 次的 `mutate` 和 `undo`, 以及 $O(N)$ 次的 `query`.
 */
function mutateWithoutOneUndo(
  start: number,
  end: number,
  options: {
    /** 这里的 index 也就是 time. */
    mutate: (index: number) => void
    undo: () => void
    query: (index: number) => void
  } & ThisType<void>
): void {
  const { mutate, undo, query } = options
  const dfs = (curStart: number, curEnd: number): void => {
    if (curEnd === curStart + 1) {
      query(curStart)
      return
    }

    const mid = Math.floor((curStart + curEnd) / 2)

    for (let i = curStart; i < mid; i++) mutate(i)
    dfs(mid, curEnd)
    for (let i = curStart; i < mid; i++) undo()

    for (let i = mid; i < curEnd; i++) mutate(i)
    dfs(curStart, mid)
    for (let i = mid; i < curEnd; i++) undo()
  }

  dfs(start, end)
}

export { mutateWithoutOneUndo }

if (require.main === module) {
  // https://leetcode.cn/problems/product-of-array-except-self/description/
  function productExceptSelf(nums: number[]): number[] {
    const res: number[] = Array(nums.length).fill(1)
    const history: number[] = []
    let mul = 1
    mutateWithoutOneUndo(0, nums.length, {
      mutate(index) {
        history.push(mul)
        mul *= nums[index]
      },
      undo() {
        mul = history.pop()!
      },
      query(index) {
        res[index] = mul
      }
    })
    return res
  }
}

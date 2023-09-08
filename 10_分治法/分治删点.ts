/* eslint-disable no-inner-declarations */

/**
 * 分治删点.
 * 调用 `query` 时，`state` 为对除了 `index` 以外所有点均调用过了 `mutate` 的状态。但不保证调用 `mutate` 的顺序。
 * 总计会调用 $O(NlgN)$ 次的 `mutate` 和 `query`.
 * @link https://github.com/tdzl2003/leetcode_live/blob/master/templates/%E5%9F%BA%E7%A1%80/%E5%88%86%E6%B2%BB%E5%88%A0%E7%82%B9.cpp
 */
function divideConquer<S>(
  state: S,
  start: number,
  end: number,
  options: {
    copy: (state: S) => S
    mutate: (state: S, index: number) => void
    query: (state: S, index: number) => void
  } & ThisType<void>
): void {
  const { copy, mutate, query } = options
  const dfs = (state: S, curStart: number, curEnd: number): void => {
    if (curEnd === curStart + 1) {
      query(state, curStart)
      return
    }
    const mid = Math.floor((curStart + curEnd) / 2)

    const copy1 = copy(state)
    for (let i = curStart; i < mid; i++) {
      mutate(copy1, i)
    }
    dfs(copy1, mid, curEnd)

    const copy2 = copy(state)
    for (let i = mid; i < curEnd; i++) {
      mutate(copy2, i)
    }
    dfs(copy2, curStart, mid)
  }

  dfs(state, start, end)
}

export { divideConquer }

if (require.main === module) {
  // 238. 除自身以外数组的乘积
  // https://leetcode.cn/problems/product-of-array-except-self/
  function productExceptSelf(nums: number[]): number[] {
    const n = nums.length
    const res = Array(n).fill(1)
    divideConquer({ mul: 1 }, 0, n, {
      copy: state => ({ mul: state.mul }),
      mutate: (state, index) => {
        state.mul *= nums[index]
      },
      query: (state, index) => {
        res[index] = state.mul
      }
    })
    return res
  }

  console.log(productExceptSelf([1, 2, 3, 4]))
}

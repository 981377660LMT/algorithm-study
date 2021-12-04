/**
 *
 * @param nums
 * @param size  可操作区间长度上限
 * @param k  至多k次操作
 * @returns
 * 题目让你求的是经过这样的 k 次+1 操作，数组 nums 的最小值最大可以达到多少。
 * 最右能力二分
 */
function minNumberOperations(nums: number[], size: number, k: number): number {
  let l = Math.min(...nums)
  let r = Math.max(...nums) + k

  while (l <= r) {
    const mid = (l + r) >> 1
    if (possible(mid)) l = mid + 1
    else r = mid - 1
  }

  return r

  function possible(target: number): boolean {
    let add = 0
    let moves = 0
    const diff = Array(nums.length).fill(0)
    for (let i = 0; i < nums.length; i++) {
      add += diff[i]
      const delta = target - (add + nums[i])
      if (delta > 0) {
        moves += delta
        add += delta // 到目前为止窗口内已经增加的高度
        if (i + size < nums.length) diff[i + size] -= delta // 之后超出k窗口范围时需要还原
      }
    }

    return moves <= k
  }
}

console.log(minNumberOperations([1, 4, 1, 1, 6], 3, 2))

// 输出：4
// 这道题我们可以逆向思考返回从 target 得到 initial 的最少操作次数。

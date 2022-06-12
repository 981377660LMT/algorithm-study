// !数组元素分成k份 最小化求每份和的最大值
// 1 <= cookies.length <= 8
function distributeCookies(nums: number[], k: number): number {
  const n = nums.length
  let res = Infinity
  nums.sort((a, b) => -(a - b))
  bt(0, 0, new Uint32Array(k))
  return res

  function bt(index: number, curMax: number, groups: Uint32Array): void {
    if (curMax > res) return
    if (index === n) {
      res = curMax
      return
    }

    for (let i = 0; i < k; i++) {
      if (i > 0 && groups[i] === groups[i - 1]) continue
      groups[i] += nums[index]
      bt(index + 1, Math.max(curMax, groups[i]), groups)
      groups[i] -= nums[index]
    }
  }
}

export {}

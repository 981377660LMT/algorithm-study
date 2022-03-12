// 乘积为正数的子数组  子数组：取还是只取
function foo(nums: number[]): number {
  const n = nums.length
  const pos = Array<number>(n).fill(0)
  const neg = Array<number>(n).fill(0)

  if (nums[0] < 0) {
    neg[0] = 1
  } else if (nums[0] > 0) {
    pos[0] = 1
  }

  for (let i = 1; i < n; i++) {
    const cur = nums[i]
    if (cur > 0) {
      pos[i] = pos[i - 1] + 1
      if (neg[i - 1] !== 0) neg[i] = neg[i - 1]
    } else if (cur < 0) {
      neg[i] = pos[i - 1] + 1
      if (neg[i - 1] !== 0) pos[i] = neg[i - 1]
    }
  }

  return pos.reduce((pre, cur) => pre + cur, 0)
}

console.log(foo([1, 1, -1, -1]))

export {}

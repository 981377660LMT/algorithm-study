// 如果 恰好 删除 一个 元素后，数组 严格递增 ，那么请你返回 true

// 找到第一对不满足单调递增的，必然要删除一个。分别判断删除后是否单调递增即可。
function canBeIncreasing(nums: number[]): boolean {
  let breakPoint = Infinity

  for (let i = 0; i < nums.length - 1; i++) {
    if (nums[i] >= nums[i + 1]) {
      breakPoint = i
      break
    }
  }

  if (breakPoint === Infinity) return true

  return (
    check(nums.filter((_, i) => i !== breakPoint)) ||
    check(nums.filter((_, i) => i !== breakPoint + 1))
  )

  function check(nums: number[]): boolean {
    for (let i = 0; i < nums.length - 1; i++) {
      if (nums[i] >= nums[i + 1]) return false
    }
    return true
  }
}

export {}
console.log(canBeIncreasing([2, 3, 1, 2]))

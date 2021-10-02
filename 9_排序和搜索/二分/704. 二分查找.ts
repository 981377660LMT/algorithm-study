function search(nums: number[], target: number): number {
  let l = 0
  let r = nums.length - 1

  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = nums[mid]
    if (midElement === target) return mid
    else if (midElement < target) l = mid + 1
    else if (midElement > target) r = mid - 1
  }

  return -1
}

export {}

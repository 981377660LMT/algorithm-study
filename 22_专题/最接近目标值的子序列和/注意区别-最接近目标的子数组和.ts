function getMinDiff(nums: number[], target: number) {
  let cand: number[] = []
  let minDiff = Infinity
  let curSum = 0
  let left = 0

  for (let right = 0; right < nums.length; right++) {
    curSum += nums[right]
    let diff = Math.abs(curSum - target)

    while (curSum > target) {
      curSum -= nums[left]
      left++
    }

    diff = Math.abs(curSum - target)
    if (diff < minDiff) {
      cand = nums.slice(left, right + 1)
      minDiff = diff
    }
  }

  return cand
}

console.log(getMinDiff([1, 2, 3, 4, 5], 8))

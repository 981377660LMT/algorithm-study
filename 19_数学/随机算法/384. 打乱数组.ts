class Solution {
  private nums: number[]
  constructor(nums: number[]) {
    this.nums = nums
  }

  reset(): number[] {
    return this.nums
  }

  shuffle(): number[] {
    const clone = this.nums.slice()
    let len = clone.length
    while (len) {
      const rand = Math.floor(Math.random() * len)
      ;[clone[len - 1], clone[rand]] = [clone[rand], clone[len - 1]]
      len--
    }
    return clone
  }
}

export {}

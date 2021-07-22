/**
 * @param {number[]} nums
 * @return {number[][]}
 */
const subsets = (nums: number[]): number[][] => {
  const res: number[][] = []
  nums.sort()

  const bt = (volume: number, path: number[] = [], index: number = 0) => {
    if (path.length === volume) {
      return res.push(path.slice())
    }

    for (let i = index; i < nums.length; i++) {
      if (i !== index && nums[i] === nums[i - 1]) continue
      path.push(nums[i])
      bt(volume, path, i + 1)
      path.pop()
    }
  }
  for (let i = 0; i <= nums.length; i++) {
    bt(i)
  }

  return res
}

console.log(subsets([1, 2, 2]))

export {}

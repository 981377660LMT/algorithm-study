/**
 * @param {number[]} nums
 * @return {number[][]}
 * 排序+规定每个重复的元素只能在开头第一个(i===index)被使用
 */
const subsets = (nums: number[]): number[][] => {
  const res: number[][] = []
  nums.sort((a, b) => a - b)

  const bt = (volume: number, path: number[] = [], index: number = 0) => {
    if (path.length === volume) {
      return res.push(path.slice())
    }

    for (let i = index; i < nums.length; i++) {
      // 规定每个重复的元素只能在开头第一个(i===index)被使用
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

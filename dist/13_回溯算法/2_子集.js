'use strict'
Object.defineProperty(exports, '__esModule', { value: true })
const subsets = nums => {
  const res = []
  const bt = (path, volumn, start) => {
    if (path.length === volumn) {
      res.push(path)
      return
    }
    for (let index = start; index < nums.length; index++) {
      bt(path.concat(nums[index]), volumn, start + 1)
    }
  }
  for (let index = 0; index <= nums.length; index++) {
    bt([], index, 0)
  }
  return res
}
console.log(subsets([1, 2, 3]))

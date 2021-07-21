// 不推荐这样做
// javascript heap out of memory
// 参见双指针指针
const threeSum = (nums: number[]) => {
  if (nums.length < 3) return []
  const res: number[][] = []
  const map = new Map<number, { val: number; index: number }[][]>()

  for (let i = 0; i < nums.length; i++) {
    for (let j = i + 1; j < nums.length; j++) {
      const need = 0 - nums[i] - nums[j]
      if (map.has(need)) {
        map.set(
          need,
          map.get(need)!.concat([
            [
              { val: nums[i], index: i },
              { val: nums[j], index: j },
            ],
          ])
        )
      } else {
        map.set(need, [
          [
            { val: nums[i], index: i },
            { val: nums[j], index: j },
          ],
        ])
      }
    }
  }

  // console.dir(map, { depth: null })
  // 收集
  for (let i = 0; i < nums.length; i++) {
    const element = nums[i]
    if (map.has(element)) {
      const towNumArr = map.get(element)!
      towNumArr.forEach(arr => {
        if (arr[0].index !== i && arr[1].index !== i) {
          res.push([arr[0].val, arr[1].val, element])
        }
      })
    }
  }

  // 二维数组去重
  return [...new Set(res.map(arr => arr.sort()).map(arr => JSON.stringify(arr)))].map(str =>
    JSON.parse(str)
  )
}
console.log(threeSum([-1, 0, 1, 2, -1, -4]))
// 输出：[[-1,-1,2],[-1,0,1]]
export {}

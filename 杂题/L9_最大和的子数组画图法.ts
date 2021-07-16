// 将点的和画成折线图，找出图上差最大的的那两点即可
const maxSubArray = (nums: number[]) => {
  if (nums.length <= 1) return nums

  let curSum = 0
  const map = new Map<number, number>()
  for (let index = 0; index < nums.length; index++) {
    curSum += nums[index]
    map.set(index, curSum)
  }

  const sortedArray = [...map].sort((a, b) => a[1] - b[1])

  return nums.slice(sortedArray[0][0] + 1, sortedArray[sortedArray.length - 1][0] + 1)
}

console.log(maxSubArray([-2, 2, -3, 4, -1, 2, 1, -5, 3]))
console.log(maxSubArray([1, -2, 3, 10, -4, 7, 2, -5]))

export {}

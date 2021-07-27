const sortColors = (nums: number[]) => {
  const count: number[] = Array(nums.length).fill(0)
  const indexArr: number[] = Array(nums.length + 1).fill(0)
  const res: number[] = []

  for (let index = 0; index < nums.length; index++) {
    const element = nums[index]
    count[element]++
  }

  for (let i = 0; i < nums.length; i++) {
    indexArr[i + 1] = indexArr[i] + count[i]
  }

  for (let i = 0; i < count[0]; i++) {
    res[i] = 0
  }
  for (let i = count[0]; i < count[0] + count[1]; i++) {
    res[i] = 1
  }
  for (let i = count[0] + count[1]; i < nums.length; i++) {
    res[i] = 2
  }
  console.log(indexArr)
  return res
}

console.log(sortColors([0, 0, 1, 2, 0, 0, 1, 0, 2, 1, 1]))

export {}

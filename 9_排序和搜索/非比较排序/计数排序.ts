const sortColors = (nums: number[]) => {
  const count: number[] = [0, 0, 0]
  const res: number[] = []

  for (let index = 0; index < nums.length; index++) {
    const element = nums[index]
    count[element]++
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

  return res
}

console.log(sortColors([0, 0, 1, 2, 0, 0, 1, 0, 2, 1, 1]))

export {}

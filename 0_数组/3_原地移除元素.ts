const removeElement = (nums: number[], val: number) => {
  for (let i = 0; i < nums.length; i++) {
    const element = nums[i]
    if (element === val) {
      nums.splice(i, 1)
      i--
    }
  }

  return nums
}

console.log(removeElement([0, 1, 2, 2, 3, 0, 4, 2], 2))

export {}

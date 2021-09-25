// 原地 移除所有数值等于 val 的元素，并返回移除后数组的新长度。

const removeElement0 = (nums: number[], val: number) => {
  for (let j = 0; j < nums.length; j++) {
    const element = nums[j]
    if (element === val) {
      nums.splice(j, 1)
      j--
    }
  }

  return nums
}

// 双指针 没见过的就搬过来
const removeElement = (nums: number[], val: number) => {
  const n = nums.length
  let i = 0
  for (let j = 0; j < n; j++) {
    if (nums[j] !== val) {
      nums[i] = nums[j]
      i++
    }
  }

  return i
}
console.log(removeElement([0, 1, 2, 2, 3, 0, 4, 2], 2))
console.log(removeElement([3, 2, 2, 3], 3))

export {}

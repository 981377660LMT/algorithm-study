// 给定一个大小为 n 的整数数组，找出其中所有出现超过 ⌊ n/3 ⌋ 次的元素。
function majorityElement(nums: number[]): number[] {
  const count = (nums: number[], target: number) => {
    let res = 0
    const n = nums.length
    for (let i = 0; i < n; i++) {
      nums[i] === target && res++
    }
    return res
  }
  const n = nums.length
  if (n < 1) return []
  if (n < 2) return nums

  let count1 = 0
  let count2 = 0
  let candidate1 = 0
  let candidate2 = 1
  for (const num of nums) {
    if (num === candidate1) count1++
    else if (num === candidate2) count2++
    else if (count1 === 0) (candidate1 = num), (count1 = 1)
    else if (count2 === 0) (candidate2 = num), (count2 = 1)
    else count1--, count2--
  }

  const res: number[] = []
  count(nums, candidate1) > n / 3 && res.push(candidate1)
  count(nums, candidate2) > n / 3 && res.push(candidate2)
  return res
}

console.log(majorityElement([1, 1, 1, 3, 3, 2, 2, 2]))
export {}

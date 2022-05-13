/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 */
const numberOfSubarrays = function (nums: number[], k: number): number {
  // 此时的奇数个数,出现了几个
  const map = new Map<number, number>([[0, 1]])
  let sum = 0
  let res = 0

  for (let i = 0; i < nums.length; i++) {
    nums[i] % 2 === 1 && sum++
    const pre = sum - k
    if (map.has(pre)) res += map.get(pre)!
    map.set(sum, map.get(sum)! + 1 || 1)
  }

  return res
}

console.log(numberOfSubarrays([1, 1, 2, 1, 1], 3))

/**
 * @param {number[]} nums
 * @param {number[][]} queries
 * @return {number[]}
 */
function maximizeXor(nums: number[], queries: number[][]) {
  nums.sort((a, b) => a - b)
  const res: number[] = []

  for (const query of queries) {
    let max = -1
    for (let i = 0; i < nums.length; i++) {
      if (nums[i] > query[1]) break
      max = Math.max(max, query[0] ^ nums[i])
    }

    res.push(max)
  }

  return res
}

export {}

function minOperations(nums: number[], k: number): number {
  for (let num of nums) {
    if (num < k) return -1
  }

  const unique = Array.from(new Set(nums)).sort((a, b) => b - a)
  if (unique.length === 1 && unique[0] === k) return 0
  let res = 0
  for (let d of unique) {
    if (d > k) res++
  }
  return res
}

/**
 *
 * @param nums 给定一个整数数组和一个整数 k，你需要找到该数组中和为 k 的连续的子数组的个数。
 * @param k
 * @description 连续子数组的和:前缀和
 */
const subarraySum = (nums: number[], k: number): number => {
  let res = 0
  const pre = nums.slice()
  pre.reduce((pre, _, index, array) => (array[index] += pre))

  // 两数之差问题:
  // map表示前缀和为某个数的个数
  const map = new Map<number, number>([[0, 1]])
  for (let i = 0; i < nums.length; i++) {
    const sum = pre[i]
    res += map.get(sum - k) || 0
    map.set(sum, (map.get(sum) || 0) + 1)
  }

  console.log(pre, map)
  return res
}

console.log(subarraySum([1, 2, 1, 2, 1], 3))

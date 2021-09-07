// 请你判断是否存在 两个不同下标 i 和 j，使得 abs(nums[i] - nums[j]) <= t ，
// 同时又满足 abs(i - j) <= k 。
// 考虑到元素差限定，分桶
/**
 * @summary
 * 1.我们可使用 (t + 1) 个桶。将所有数除以 (t+1) 的结果作为编号存到一个哈希表中
 * 如果两个数字的编号相同，那么意味着其绝对值差小于等于 t。(eg:0和t)
 * 如果两个数字的编号不同,相邻编号也可能是绝对值差小于等于 t,需要检查当前和左右相邻
 *
 * 2.由于题目限定是索引差小于等于 k，因此我们可以固定一个窗口大小为 k 的滑动窗口，每次都仅处理窗口内的元素
 */
const containsNearbyAlmostDuplicate = (nums: number[], k: number, t: number): boolean => {
  if (k === 0) return false
  if (nums.length <= 1) return false

  // 关键点：我们的一个桶内同一时刻只会有一个元素 多了我们会直接返回结果
  const record = new Map<number, number>() // hash值 实际值

  for (let i = 0; i < nums.length; i++) {
    const num = nums[i]
    const hash = Math.floor(num / (t + 1)) // 这里不能用~~ 因为可能是负数
    if (record.has(hash)) return true
    if (record.has(hash - 1) && Math.abs(num - record.get(hash - 1)!) <= t) return true
    if (record.has(hash + 1) && Math.abs(num - record.get(hash + 1)!) <= t) return true
    record.set(hash, num)

    if (i >= k) record.delete(~~(nums[i - k] / (t + 1)))
  }
  return false
}

// console.log(containsNearbyAlmostDuplicate([1, 5, 9, 1, 5, 9], 2, 3))
console.log(containsNearbyAlmostDuplicate([2147483640, 2147483641], 1, 100))

export {}

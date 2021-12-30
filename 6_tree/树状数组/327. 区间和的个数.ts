class BIT {
  private size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array(size + 1).fill(0)
  }

  add(x: number, k: number) {
    if (x <= 0) throw Error('树状数组索引必须为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree[i] += k
    }
  }

  query(x: number) {
    let res = 0
    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += this.tree[i]
    }
    return res
  }

  sumRange(left: number, right: number) {
    return this.query(right) - this.query(left - 1)
  }

  private lowbit(x: number) {
    return x & -x
  }
}

/**
 *
 * @param nums 1 <= nums.length <= 10**5
 * @param lower
 * @param upper
 * @description
 * 数组 A 有多少个连续的子数组，其元素只和在 [lower, upper]的范围内。
 * 即：前缀和之差不超过[lower,upper]
 * @summary
 * 0.准备好需要用于树状数组查询/更新的所有值(即前缀和和前缀和减去upper/lower)
   1.set+map离散化所有值
   2.对每一个前缀和pre[i]，查询[pre[i]-upper,pre[i]-lower]间的个数，查询完后更新树状数组
   @link
   https://leetcode-cn.com/problems/count-of-range-sum/solution/jstsshu-zhuang-shu-zu-jie-fa-by-cao-mei-0icur/
 */
function countRangeSum(nums: number[], lower: number, upper: number): number {
  const pre = [0]
  for (const num of nums) {
    pre.push(pre[pre.length - 1] + num)
  }

  // 离散化
  const allNums = new Set<number>()
  for (const val of pre) {
    allNums
      .add(val)
      .add(val - upper)
      .add(val - lower)
  }

  const map = new Map()
  // 把要用到的值利用哈希表进行离散化
  for (const [key, realValue] of [...allNums].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)
  }
  console.log(pre, allNums, map)
  // Map(8) {
  //   -4 => 1,
  //   -2 => 2,
  //   0 => 3,
  //   1 => 4,
  //   2 => 5,
  //   3 => 6,
  //   4 => 7,
  //   5 => 8
  // }

  let res = 0
  const bit = new BIT(map.size)
  for (let i = 0; i < pre.length; i++) {
    const realValue = pre[i]
    const left = map.get(realValue - upper)!
    const right = map.get(realValue - lower)!
    res += bit.sumRange(left, right)
    bit.add(map.get(realValue)!, 1)
  }

  return res
}

console.log(countRangeSum([-2, 5, -1], -2, 2))
// 输出：3
// 解释：存在三个区间：[0,0]、[2,2] 和 [0,2] ，对应的区间和分别是：-2 、-1 、2 。

export {}

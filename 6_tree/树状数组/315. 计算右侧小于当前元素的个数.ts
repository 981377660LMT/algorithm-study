class BIT {
  private size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array(size + 1).fill(0)
  }

  // 最好x都离散化正数
  add(x: number, k: number) {
    if (x <= 0) throw Error('查询索引应为正整数')
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

  // sumRange(left: number, right: number) {
  //   return this.query(right + 1) - this.query(left)
  // }

  private lowbit(x: number) {
    return x & -x
  }
}

/**
 * @param {number[]} nums
 * @return {number[]}
 * 按要求返回一个新数组 counts 。数组 counts 有该性质：
 *  counts[i] 的值是  nums[i] 右侧小于 nums[i] 的元素的数量。
 * @link
 * https://leetcode-cn.com/problems/count-of-smaller-numbers-after-self/solution/shu-zhuang-shu-zu-de-xiang-xi-fen-xi-by-yangbingji/
 */
var countSmaller = function (nums: number[]): number[] {
  // 离散化
  const set = new Set(nums)
  const map = new Map<number, number>()
  for (const [key, realValue] of [...set].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)
  }
  console.log(map)

  const res = Array<number>(nums.length).fill(0)
  const bit = new BIT(map.size)
  // 倒序遍历 动态更新 查询小于当前数的有多少个
  for (let i = nums.length - 1; i >= 0; i--) {
    const realValue = nums[i]
    const discretizedValue = map.get(realValue)!
    res[i] = bit.query(discretizedValue - 1)
    bit.add(discretizedValue, 1)
  }

  return res
}

console.log(countSmaller([5, 2, 6, 1]))
console.log(countSmaller([-1, -1]))
console.log(countSmaller([-1, -2]))
// 输出：[2,1,1,0]
// 解释：
// 5 的右侧有 2 个更小的元素 (2 和 1)
// 2 的右侧仅有 1 个更小的元素 (1)
// 6 的右侧有 1 个更小的元素 (1)
// 1 的右侧有 0 个更小的元素

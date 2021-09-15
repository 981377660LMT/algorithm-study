class BIT {
  private size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array<number>(size + 1).fill(0)
  }

  add(x: number, k: number) {
    if (x <= 0) throw Error('add操作时树状数组索引应为正整数')
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

  private lowbit(x: number) {
    return x & -x
  }
}

/**
 * @param {number[]} nums
 * @return {number}
 * @description
 * 如果 i < j 且 nums[i] > 2*nums[j] 我们就将 (i, j) 称作一个重要翻转对。
 * @summary
 * 1.离散化:离散化的时候2倍元素值也得加入离散化数组里
 * 2.动态更新 之后就是用树状数组求数集里有多少个数比2*a[i]大了
 * @link
 * https://leetcode-cn.com/problems/reverse-pairs/solution/jstsshu-zhuang-shu-zu-jie-fa-by-cao-mei-uowff/
 */
const reversePairs = function (nums: number[]): number {
  const set = new Set([...nums, ...nums.map(v => v * 2)])

  const map = new Map<number, number>()
  for (const [key, realValue] of [...set].sort((a, b) => a - b).entries()) {
    map.set(realValue, key + 1)
  }

  let res = 0
  const bit = new BIT(map.size)
  for (let i = 0; i < nums.length; i++) {
    const realValue = nums[i]
    const discretizedValue = map.get(realValue)!
    const doubleDiscretizedValue = map.get(realValue * 2)!
    // 树状数组求数集里有多少个数严格大于2*nums[j]
    res += bit.query(map.size) - bit.query(doubleDiscretizedValue)
    bit.add(discretizedValue, 1)
  }

  return res
}

console.log(reversePairs([1, 3, 2, 3, 1]))
console.log(reversePairs([2, 4, 3, 5, 1]))
// 输出: 2
export {}

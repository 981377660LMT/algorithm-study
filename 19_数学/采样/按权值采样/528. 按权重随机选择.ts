import { bisectLeft } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'

class Solution {
  private pre: number[]

  /**
   *
   * @param w 给定一个正整数数组 w,其中 w[i] 代表下标 i 的权重
   * 转化为和
   * 和为几就占几份 3 占 1 2 3
   *
   */
  constructor(w: number[]) {
    // 计算前缀和，这样可以生成一个随机数，根据数的大小对应分布的坐标
    this.pre = w.slice()
    this.pre.reduce((pre, _, index, array) => (array[index] += pre))
  }

  // 选取下标 i 的概率为 w[i] / sum(w) 。
  pickIndex(): number {
    // 第几个点
    const rand = this.randint(1, this.pre[this.pre.length - 1])
    return bisectLeft(this.pre, rand)
  }

  private randint(start: number, end: number) {
    if (start > end) throw new Error('invalid interval')
    const amplitude = end - start
    return Math.floor((amplitude + 1) * Math.random()) + start
  }
}

const solution = new Solution([1, 3])
console.log(solution.pickIndex())
console.log(solution.pickIndex())
console.log(solution.pickIndex())
console.log(solution.pickIndex())
// solution.pickIndex(); // 返回 1，返回下标 1，返回该下标概率为 3/4 。
// solution.pickIndex(); // 返回 1
// solution.pickIndex(); // 返回 1
// solution.pickIndex(); // 返回 1
// solution.pickIndex(); // 返回 0，返回下标 0，返回该下标概率为 1/4 。

export {}

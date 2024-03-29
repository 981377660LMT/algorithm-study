/* eslint-disable no-shadow */

// # 吃谷物的最短时间
// # 在一条线上有 n 个母鸡和 m 个谷粒。
// # 给出了两个整数阵列中母鸡和谷物的初始位置，母鸡和谷物的大小分别为 n 和 m。
// # 任何母鸡都可以吃谷物，如果他们在同一位置。花在这上面的时间可以忽略不计。
// # 一只母鸡也可以吃多种谷物。在1秒内，母鸡可以向左或向右移动1个单位。
// # `母鸡可以同时独立地移动`。如果母鸡表现最佳，返回吃所有谷物的最短时间。
// # n<=2e4
// # !0<=pos<=1e9

// # !注意不是最短移动距离, 而是总时间(并行)

// # 二分+排序
// # 1.每只母鸡应该按顺序匹配谷物，因此我们可以对母鸡和谷物的位置进行排序。
// # 2.二分答案,每只鸡在给定的时间内可以吃多少粒谷物
// # !3.调头一次的最短距离=min(2*leftMax+rightMax,2*rightMax+leftMax)
// # !其中leftMax为向左走的最大距离,rightMax为向右走的最大距离
// # (摘水果那道题也是这个调头公式)

function minimumTime(hens: number[], grains: number[]): number {
  hens = hens.slice().sort((a, b) => a - b)
  grains = grains.slice().sort((a, b) => a - b)

  let left = 0
  let right = 2e15
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) {
      right = mid - 1
    } else {
      left = mid + 1
    }
  }

  return left

  /**
   * mid分钟内能否并行吃完所有谷物.
   */
  function check(mid: number): boolean {
    let left = 0
    let eat = 0
    for (let hi = 0; hi < hens.length; hi++) {
      const start = hens[hi]
      let right = left
      while (right < grains.length && calDist(start, grains[left], grains[right]) <= mid) {
        right++
      }
      eat += right - left
      left = right
    }

    return eat === grains.length
  }

  /**
   * 从start位置出发,遍历[left,right]区间的最短距离(最多调头一次).
   */
  function calDist(start: number, left: number, right: number): number {
    const leftMax = Math.max(0, start - left)
    const rightMax = Math.max(0, right - start)
    return Math.min(2 * leftMax + rightMax, 2 * rightMax + leftMax)
  }
}

export { minimumTime }

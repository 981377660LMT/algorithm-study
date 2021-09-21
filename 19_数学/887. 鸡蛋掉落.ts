/**
 * @param {number} k  给你 k 枚相同的鸡蛋
 * @param {number} n  共有 n 层楼的建筑
 * @return {number}  最少要扔几次？
 * @description
 * 任何从 高于 f 的楼层落下的鸡蛋都会碎，从 f 楼层或比它低的楼层落下的鸡蛋都不会破
 * @summary
 * 转变为：有 eggs 个蛋，扔 throws 次，求可以确定 F 的个数，然后得出 N 个楼层
 */
const superEggDrop = function (k: number, n: number): number {
  /**
   * 
   * @param eggs 
   * @param throws 
   * @returns 
   * 如果只有 1 个蛋，或只有 1 次机会时，只可以确定出 T + 1 个 F
     其他情况时，递归。
     【蛋碎了减 1 个，机会减 1 次】 + 【蛋没碎，机会减 1 次】
   */
  const calFloor = (eggs: number, throws: number): number => {
    if (eggs === 1 || throws === 1) return throws + 1
    return calFloor(eggs, throws - 1) + calFloor(eggs - 1, throws - 1)
  }

  let throws = 1
  // n层建筑 floor有0到n n+1种
  while (calFloor(k, throws) < n + 1) {
    throws++
  }

  return throws
}

console.log(superEggDrop(1, 2)) // 2
console.log(superEggDrop(2, 6)) // 3

// 子问题
// 有 2 个蛋，用一座 100 层的楼，要使用最少次数测试出蛋几层会碎（F）。
// 问第一次应该从几层扔。
console.log(superEggDrop(2, 100)) // 14
// 因为最少需要 14 次，所以第 1 次扔在 14 层，如果蛋碎了，接下来 1~13 这个区间就只能一次一次尝试了。

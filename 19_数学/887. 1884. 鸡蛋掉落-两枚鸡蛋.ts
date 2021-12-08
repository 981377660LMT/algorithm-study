function twoEggDrop(n: number): number {
  let throws = 1
  // n层建筑 floor有0到n n+1种
  while (calFloor(2, throws) < n + 1) {
    throws++
  }

  return throws

  // 能确定的楼层数(f在每一层都有可能，所以要全部覆盖)
  function calFloor(eggs: number, throws: number): number {
    if (eggs === 1 || throws === 1) return throws + 1
    return calFloor(eggs, throws - 1) + calFloor(eggs - 1, throws - 1)
  }
}

console.log(twoEggDrop(100))
// 给你 2 枚相同 的鸡蛋，和一栋从第 1 层到第 n 层共有 n 层楼的建筑。
// 已知存在楼层 f ，满足 0 <= f <= n ，任何从 高于 f 的楼层落下的鸡蛋都 会碎 ，从 f 楼层或比它低 的楼层落下的鸡蛋都 不会碎 。
// 请你计算并返回要确定 f 确切的值 的 最小操作次数 是多少？

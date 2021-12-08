// 你不能购买超出购物清单指定数量的物品，即使那样会降低整体价格。任意大礼包可无限次购买。
function shoppingOffers(price: number[], special: number[][], needs: number[]): number {
  const memo = new Map<string, number>()

  const dfs = (needs: number[]): number => {
    const key = needs.join('#')
    if (memo.has(key)) return memo.get(key)!

    // 最差情况，计算所有单个买的价格
    let res = needs.reduce((pre, _, index, needs) => pre + needs[index] * price[index], 0)
    for (const sp of special) {
      // 更新needs，减去礼包能够提供的数量
      const newNeeds = Array.from(needs, (need, index) => need - sp[index])
      // 判断能否购买，会不会礼包超过需要购买数量
      if (newNeeds.every(need => need >= 0)) {
        res = Math.min(res, dfs(newNeeds) + sp[sp.length - 1])
      }
    }

    memo.set(key, res)
    return res
  }

  return dfs(needs)
}

console.log(
  shoppingOffers(
    [2, 5],
    [
      [3, 0, 5],
      [1, 2, 10],
    ],
    [3, 2]
  )
)
// 输出：14
// 解释：有 A 和 B 两种物品，价格分别为 ¥2 和 ¥5 。
// 大礼包 1 ，你可以以 ¥5 的价格购买 3A 和 0B 。
// 大礼包 2 ，你可以以 ¥10 的价格购买 1A 和 2B 。
// 需要购买 3 个 A 和 2 个 B ， 所以付 ¥10 购买 1A 和 2B（大礼包 2），以及 ¥4 购买 2A 。

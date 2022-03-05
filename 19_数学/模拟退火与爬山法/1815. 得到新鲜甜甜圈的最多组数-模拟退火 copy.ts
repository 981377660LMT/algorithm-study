// https://leetcode-cn.com/problems/maximum-number-of-groups-getting-fresh-donuts/solution/mo-ni-tui-huo-by-lianxin-tvq6/
function maxHappyGroups(batchSize: number, groups: number[]): number {
  if (batchSize == 1) return groups.length

  const tmp = []
  // 万能牌free  这波人必然会带来一波快乐   先踢出来  这样数据量小点
  let res1 = 0
  for (let i = 0; i < groups.length; i++) {
    if (groups[i] % batchSize == 0) {
      res1++
    } else {
      tmp.push(groups[i])
    }
  }

  groups = tmp

  let res2 = 0
  for (let i = 0; i < 30; i++) {
    simulateAnneal()
  }

  return res1 + res2

  // 模拟退火 退火率0.975
  function simulateAnneal(): void {
    randomShuffle()

    for (let i = 5e4; i > 1e-5; i *= 0.975) {
      const before = main()
      const index1 = Math.floor(Math.random() * groups.length)
      const index2 = Math.floor(Math.random() * groups.length)
      swap(index1, index2)
      const after = main()
      const delta = after - before

      // 变好了自动保留
      if (delta > 0) {
        continue
      }

      // 这是变坏了  变化了应该是要“恢复的”， 也就是swap回去
      // 但变化了也有一定概率保留（坏了也就坏了 因为有的case是先变坏才能变好，先恶化再变得最优）
      // Math.E ** (delta / i) 这个挺有意思的..  恶化的越多 这个值越小，越应该恢复回去， 温度高时，这个值越趋近于1，越容易维持现状
      // 如果没恶化太多，或当前温度比较高， 是有很大概率不恢复回去的，将错就错
      if (delta < 0 && Math.E ** (delta / i) < Math.random()) {
        swap(index1, index2)
      }
    }
  }

  function main(): number {
    // 第一组永远开心.
    let happyNum = 1
    let total = groups[0]
    for (let i = 1; i < groups.length; i++) {
      if (total % batchSize == 0) {
        happyNum++
      }

      total += groups[i]
    }

    res2 = Math.max(res2, happyNum)
    // 退火delta
    return happyNum
  }

  function swap(l: number, r: number): void {
    ;[groups[l], groups[r]] = [groups[r], groups[l]]
  }

  // 洗牌
  function randomShuffle(): void {
    for (let i = 0; i < groups.length; i++) {
      let random = Math.floor(Math.random() * (groups.length - i))
      swap(random, groups.length - 1 - i)
    }
  }
}

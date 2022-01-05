// 奇偶数量，返回小的那个
// chips数组里存放的是第i个筹码存放的位置
function minCostToMoveChips(position: number[]): number {
  let odd = 0,
    even = 0
  for (let i = 0; i < position.length; i++) {
    position[i] % 2 ? odd++ : even++
  }

  return Math.min(odd, even)
}

console.log(minCostToMoveChips([2, 2, 2, 3, 3]))
// 输出：2
// 解释：第四和第五个筹码移动到位置二的代价都是 1，所以最小总代价为 2。

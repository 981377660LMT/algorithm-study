// 最后会移动到1/2 比较这两个位置
function minCostToMoveChips(position: number[]): number {
  let [odd, even] = [0, 0]
  for (let i = 0; i < position.length; i++) {
    position[i] & 1 ? odd++ : even++
  }

  return Math.min(odd, even)
}

console.log(minCostToMoveChips([2, 2, 2, 3, 3]))
// 输出：2
// 解释：第四和第五个筹码移动到位置二的代价都是 1，所以最小总代价为 2。

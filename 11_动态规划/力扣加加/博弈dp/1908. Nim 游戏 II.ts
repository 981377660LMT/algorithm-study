/**
 * @param {number[]} piles
 * @return {boolean}
 */
var nimGame = function (piles: number[]): boolean {
  let xor = 0
  piles.forEach(pile => (xor ^= pile))
  return xor !== 0
}
// 共有 n 堆石头。在每个玩家的回合中，
// 玩家需要 选择 任一非空石头堆，从中移除任意 非零 数量的石头。
// 如果不能移除任意的石头，就输掉游戏，同时另一人获胜。

// 异或为0 先手必败
// 因为后手可以让异或始终为0(对应二进制位取)

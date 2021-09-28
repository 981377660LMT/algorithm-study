// 每一回合，你可以从两端之一拿起一枚石子（位置最大或最小），并将其放入两端之间的任一空闲位置
// 要使游戏结束，你可以执行的最小和最大移动次数分别是多少？
// 以长度为 2 的数组形式返回答案：answer = [minimum_moves, maximum_moves]

function numMovesStones(a: number, b: number, c: number): number[] {
  ;[a, b, c] = [a, b, c].sort((a, b) => a - b)
  if (c - a === 2) return [0, 0]
  // 若a、c任一与b差值小于等于2，则可一步到位
  if (b - a <= 2 || c - b <= 2) return [1, c - b - 1 + (b - a - 1)]
  else return [2, c - b - 1 + (b - a - 1)]
}

// 输入：a = 1, b = 2, c = 5
// 输出：[1, 2]
// 解释：将石子从 5 移动到 4 再移动到 3，或者我们可以直接将石子移动到 3。

// 输入：a = 4, b = 3, c = 2
// 输出：[0, 0]
// 解释：我们无法进行任何移动。

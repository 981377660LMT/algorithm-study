// 假定有n (n > 0) 颗石子。
// 玩家A和B轮流取石子，一次可以取1个或两个，A先取。
// 谁取到最后一颗石子即为输。
// 请问 A或者B是否有必赢策略? 如果有请返回能必赢的玩家名。没有则返回null。
function canWinStonePicking(n: number): 'A' | 'B' | null {
  return n % 3 === 1 ? 'B' : 'A'
}

function canWinStonePicking2(n: number): 'A' | 'B' | null {
  // 第n步先手赢
  const dfs = (n: number): boolean => {}
}

// 你和朋友轮流将 连续 的两个 "++" 反转成 "--"
// 计算并返回 一次有效操作 后，字符串 currentState 所有的可能状态
/**
 * @param {string} currentState
 * @return {string[]}
 */
var generatePossibleNextMoves = function (currentState) {
  const res = []
  for (let i = 0; i < currentState.length; i++) {
    if (currentState.startsWith('++', i)) {
      res.push(currentState.slice(0, i) + '--' + currentState.slice(i + 2))
    }
  }
  return res
}

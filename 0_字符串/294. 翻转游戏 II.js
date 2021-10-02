// 你和朋友轮流将 连续 的两个 "++" 反转成 "--"
// 当一方无法进行有效的翻转时便意味着游戏结束，则另一方获胜
// 请你写出一个函数来判定起始玩家 是否存在必胜的方案
/**
 * @param {string} currentState
 * @return {boolean}
 */
var canWin = function (currentState) {
  const memo = {}
  return inner(currentState)

  function inner(str) {
    if (str in memo) return memo[str]

    n = str.length
    for (let i = 0; i < n - 1; i++) {
      // 自己可以且对方不能赢
      if (str[i] === '+' && str[i + 1] === '+') {
        if (!inner(str.slice(0, i) + '--' + str.slice(i + 2))) {
          memo[str] = true
          return true
        }
      }
    }

    memo[str] = false
    return false
  }
}

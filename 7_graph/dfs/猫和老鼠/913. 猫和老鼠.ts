/**
 * @param {number[][]} graph  邻接表
 * @return {number}
 * 老鼠从结点 1 开始并率先出发，猫从结点 2 开始且随后出发，在结点 0 处有一个洞。
 */
var catMouseGame = function (graph: number[][]): 0 | 1 | 2 {
  const memo = new Map<string, 0 | 1 | 2>()
  const len = graph.length

  const dfs = (steps: number, mouse: number, cat: number): 0 | 1 | 2 => {
    if (steps === len * 2) return 0 // 猫和老鼠各走了N步 但是还没分出胜负
    if (mouse === 0) return 1 // 老鼠进0，老鼠赢
    if (mouse === cat) return 2 // 老鼠和猫在一个位置，猫赢

    const key = `${steps}#${mouse}#${cat}`
    if (memo.has(key)) return memo.get(key)!

    // 下一步是猫
    if (steps & 1) {
      for (const next of graph[cat]) {
        if (next === 0) continue
        const key = `${steps + 1}#${mouse}#${next}`
        // 猫能在任何一个下一步中赢
        if (dfs(steps + 1, mouse, next) === 2) {
          memo.set(key, 2)
          return 2
        }
      }

      for (const next of graph[cat]) {
        if (next === 0) continue
        const key = `${steps + 1}#${mouse}#${next}`
        // 平局
        if (dfs(steps + 1, mouse, next) === 0) {
          memo.set(key, 0)
          return 0
        }
      }

      memo.set(key, 1)
      return 1
    } else {
      // 下一步是老鼠
      for (const next of graph[mouse]) {
        const key = `${steps + 1}#${next}#${cat}`
        // 如果老鼠能在任何一个下一步中赢 返回老鼠赢
        if (dfs(steps + 1, next, cat) === 1) {
          memo.set(key, 1)
          return 1
        }
      }

      for (const next of graph[mouse]) {
        const key = `${steps + 1}#${next}#${cat}`
        // 如果老鼠能在任何一个下一步中赢 返回老鼠赢
        if (dfs(steps + 1, next, cat) === 0) {
          memo.set(key, 0)
          return 0
        }
      }

      memo.set(key, 2)
      return 2
    }
  }

  return dfs(0, 1, 2)
}

console.log(catMouseGame([[2, 5], [3], [0, 4, 5], [1, 4, 5], [2, 3], [0, 2, 3]]))
// 如果猫和老鼠占据相同的结点，猫获胜1。
// 如果老鼠躲入洞里，老鼠获胜2。
// 如果某一位置重复出现（即，玩家们的位置和移动顺序都与上一个回合相同），游戏平局0。

// 用t,x,y代表各个状态
// t是步数，x是老鼠的位置，y是猫的位置
////////////////////////////////////
// 一共有三种情况：
// 1.如果总步数达到2*N,意味着猫和老鼠各走了N步，
// 但是老鼠到达洞的步数最多只有N,如果过了N步老鼠还没有到达洞，
// 并且猫也没有抓住老鼠，那么就是平局，返回0
// 2.x=y,老鼠和猫在一个位置，猫赢
// 3.x=0,老鼠进0，老鼠赢

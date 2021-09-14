/**
 * @param {string} startTime
 * @param {string} finishTime
 * @return {number}
 * 以 刻钟 为周期规划若干时长为 15 分钟 的游戏对局
 * startTime finishTime 分别表示你 进入 和 退出 游戏的确切时间，
 * 请计算在整个游戏会话期间，你完成的 完整对局的对局数 。
 * 如果 finishTime 早于 startTime ，这表示你玩了个通宵（也就是从 startTime 到午夜，再从午夜到 finishTime）。
 */
var numberOfRounds = function (startTime: string, finishTime: string): number {
  const toMinute = (time: string) => {
    const [a, b] = time.split(':')
    const hour = parseInt(a)
    const minute = parseInt(b)
    return hour * 60 + minute
  }

  let t1 = toMinute(startTime)
  let t2 = toMinute(finishTime)
  if (t1 > t2) t2 += 24 * 60

  // 注意这里
  const start = Math.ceil(t1 / 15)
  const end = Math.floor(t2 / 15)

  return Math.max(0, end - start)
}

console.log(numberOfRounds('00:00', '23:59'))
console.log(numberOfRounds('20:00', '06:00'))
console.log(parseInt('20:00'))

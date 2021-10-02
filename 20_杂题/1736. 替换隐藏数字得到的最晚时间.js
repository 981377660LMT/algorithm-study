// 替换 time 中隐藏的数字，返回你可以得到的最晚有效时间。
function maximumTime(time) {
  time = time.split('')
  if (time[0] === '?') time[0] = time[1] > 3 ? '1' : '2'
  if (time[1] === '?') time[1] = time[0] > 1 ? '3' : '9'
  if (time[3] === '?') time[3] = '5'
  if (time[4] === '?') time[4] = '9'
  return time.join('')
}

console.log(maximumTime('2?:?0'))
// 输出："23:50"
// 解释：以数字 '2' 开头的最晚一小时是 23 ，以 '0' 结尾的最晚一分钟是 50 。

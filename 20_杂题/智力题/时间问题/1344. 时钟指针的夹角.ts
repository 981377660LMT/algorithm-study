// 给你两个数 hour 和 minutes 。请你返回在时钟上，由给定时间的时针和分针组成的较小角的角度（60 单位制）。
function angleClock(hour: number, minutes: number): number {
  const angleMin = (minutes / 60) * 360
  const angleHour = (hour / 12) * 360 + angleMin / 12

  return getCycleDist(angleMin, angleHour, 360)

  function getCycleDist(n1: number, n2: number, mod: number): number {
    return Math.min((n1 - n2 + mod) % mod, (n2 - n1 + mod) % mod)
  }
}

console.log(angleClock(3, 15))

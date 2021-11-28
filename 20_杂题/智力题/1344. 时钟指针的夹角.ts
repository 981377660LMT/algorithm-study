// 给你两个数 hour 和 minutes 。请你返回在时钟上，由给定时间的时针和分针组成的较小角的角度（60 单位制）。
function angleClock(hour: number, minutes: number): number {
  const angleMin = (minutes / 60) * 360
  const angleHour = (hour / 12) * 360 + angleMin / 12

  const angleDiff = Math.abs(angleMin - angleHour)
  const finalAngle = angleDiff > 180 ? 360 - angleDiff : angleDiff

  return finalAngle
}

console.log(angleClock(3, 15))

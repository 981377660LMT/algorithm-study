/**
 * @param {string} time  时间以HH:mm的格式传入。
 * @returns {number}  计算出时钟的时针和分针的角度（两个角度的较小者，四舍五入）
 */
function angle(time: string): number {
  const [hour, min] = time.split(':').map(Number)

  const h = hour >= 12 ? hour - 12 : hour
  const m = min

  const angleMin = (m / 60) * 360
  const angleHour = (h / 12) * 360 + angleMin / 12

  const angleDiff = Math.abs(angleMin - angleHour)
  const finalAngle = angleDiff > 180 ? 360 - angleDiff : angleDiff

  return Math.round(finalAngle)
}

console.log(angle('12:00'))
// 0

console.log(angle('23:30'))
// 165

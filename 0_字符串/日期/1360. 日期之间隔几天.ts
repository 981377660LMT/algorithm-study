// 给定的日期是 1971 年到 2100 年之间的有效日期。
function daysBetweenDates(date1: string, date2: string): number {
  return Math.abs(new Date(date1).getTime() - new Date(date2).getTime()) / (24 * 60 * 60 * 1000)
}

console.log(daysBetweenDates('2020-01-15', '2019-12-31'))

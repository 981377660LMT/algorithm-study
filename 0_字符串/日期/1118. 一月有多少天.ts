// 日期api
function numberOfDays(year: number, month: number): number {
  return new Date(year, month, 0).getDate()
}

console.log(numberOfDays(2000, 2))
export {}

function dayOfTheWeek(day: number, month: number, year: number): string {
  const mapping: Record<number, string> = {
    0: 'Sunday',
    1: 'Monday',
    2: 'Tuesday',
    3: 'Wednesday',
    4: 'Thursday',
    5: 'Friday',
    6: 'Saturday',
  }

  return mapping[new Date(`${year}-${month}-${day}`).getDay()]
}

console.log(dayOfTheWeek(31, 8, 2019))
// 您返回的结果必须是这几个值中的一个 {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}。

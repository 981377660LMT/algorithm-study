function reformatDate(date: string): string {
  const dict: Record<string, string> = {
    Jan: '01',
    Feb: '02',
    Mar: '03',
    Apr: '04',
    May: '05',
    Jun: '06',
    Jul: '07',
    Aug: '08',
    Sep: '09',
    Oct: '10',
    Nov: '11',
    Dec: '12',
  }

  let [day, month, year] = date.split(/\s+/)
  day = parseInt(day).toString().padStart(2, '0')
  return [year, dict[month], day].join('-')
}
console.log(reformatDate('20th Oct 2052'))
// 输出："2052-10-20"

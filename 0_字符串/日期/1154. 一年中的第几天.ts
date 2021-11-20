function dayOfYear(date: string): number {
  const firstDay = date.slice(0, 4) + '-01-01'
  return (new Date(date).getTime() - new Date(firstDay).getTime()) / (24 * 3600 * 1000) + 1
}

console.log(dayOfYear('2019-02-10'))

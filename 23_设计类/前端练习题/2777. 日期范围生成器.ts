// 2777. 日期范围生成器
// 现给定起始日期 start 、结束日期 end 和正整数 step ，返回一个生成器对象，
// 该生成器对象按照从 start 到 end（包括 start 和 end ）的范围生成日期。
// 所有日期都以字符串格式 YYYY-MM-DD 表示。
// step 的值表示连续生成的日期之间的天数间隔。

const MS_PER_DAY = 24 * 60 * 60 * 1000

function* dateRangeGenerator(start: string, end: string, step: number): Generator<string> {
  let startTimeStamp = dateToTimeStamp(start)
  const endTimeStamp = dateToTimeStamp(end)
  while (startTimeStamp <= endTimeStamp) {
    yield timeStampToDate(startTimeStamp)
    startTimeStamp += step * MS_PER_DAY
  }
}

/**
 * 时间戳转换为日期字符串'YYYY-MM-DD'.
 */
function timeStampToDate(timeStamp: number): string {
  const date = new Date(timeStamp)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`
}

/**
 * 日期字符串'YYYY-MM-DD'转换为时间戳.
 */
function dateToTimeStamp(date: string): number {
  return new Date(date).getTime()
}

/**
 * const g = dateRangeGenerator('2023-04-01', '2023-04-04', 1);
 * g.next().value; // '2023-04-01'
 * g.next().value; // '2023-04-02'
 * g.next().value; // '2023-04-03'
 * g.next().value; // '2023-04-04'
 * g.next().done; // true
 */

export {}

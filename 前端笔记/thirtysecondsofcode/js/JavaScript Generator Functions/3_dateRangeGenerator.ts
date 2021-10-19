;[...dateRangeGenerator(new Date('2021-06-01'), new Date('2021-06-04'))]
// [ 2021-06-01, 2021-06-02, 2021-06-03 ]

function* dateRangeGenerator(start: Date, end: Date, step = 1) {
  let innerStart = start
  while (innerStart < end) {
    yield new Date(innerStart)
    /** Get/Set the day-of-the-month, using local time. */
    innerStart.setDate(innerStart.getDate() + step)
  }
}

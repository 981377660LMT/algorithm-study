/** 环区间分解.*/
function divideCycle(n: number, start: number, end: number, f: (start: number, end: number, times: number) => void): void {
  if (start >= end || n <= 0) return
  const loop = Math.floor((end - start) / n)
  if (loop > 0) {
    f(0, n, loop)
  }
  if ((end - start) % n === 0) return
  start %= n
  end %= n
  if (start < end) {
    f(start, end, 1)
  } else {
    f(start, n, 1)
    if (end > 0) {
      f(0, end, 1)
    }
  }
}

export { divideCycle }

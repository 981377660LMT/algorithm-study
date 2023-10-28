/**
 * 数论分块.
 * 将 `[1,n]` 内的数分成 `O(2*sqrt(n))` 段, 每段内的 `n//i` 相同。
 */
function enumerateFloor(n: number, f: (left: number, right: number, div: number) => void): void {
  for (let l = 1, r = 0; l <= n; l = r + 1) {
    const h = Math.floor(n / l)
    r = Math.floor(n / h)
    f(l, r, h)
  }
}

/**
 * 数论分块.
 * 将 `[lower,upper]` 内的数分成 `O(2*sqrt(upper))` 段, 每段内的 `upper//i` 相同。
 */
function enumerateFloorInterval(
  lower: number,
  upper: number,
  f: (left: number, right: number, div: number) => void
): void {
  for (let l = lower, r = 0; l <= upper; l = r + 1) {
    const h = Math.floor(upper / l)
    if (h === 0) break
    r = Math.min(Math.floor(upper / h), upper)
    f(l, r, h)
  }
}

/**
 * 二维数论分块.
 * 将 `[1,n] x [1,m]` 内的数分成 `O(2*sqrt(n)*2*sqrt(m))` 段, 每段内的 `(n//i, m//i)` 相同。
 */
function enumerateFloor2D(
  n: number,
  m: number,
  f: (x1: number, x2: number, y1: number, y2: number, div1: number, div2: number) => void
): void {
  for (let x1 = 1, x2 = 0; x1 <= n; x1 = x2 + 1) {
    const hn = Math.floor(n / x1)
    x2 = Math.floor(n / hn)
    for (let y1 = 1, y2 = 0; y1 <= m; y1 = y2 + 1) {
      const hm = Math.floor(m / y1)
      y2 = Math.floor(m / hm)
      f(x1, x2, y1, y2, hn, hm)
    }
  }
}

export { enumerateFloor, enumerateFloorInterval, enumerateFloor2D }

if (require.main === module) {
  enumerateFloor(10, console.log)
  enumerateFloorInterval(19, 20, console.log)
  enumerateFloor2D(5, 5, console.log)
}

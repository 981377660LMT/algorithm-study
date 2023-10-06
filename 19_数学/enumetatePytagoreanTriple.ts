// 遍历勾股数(勾股定理)

/**
 * 遍历勾股数对 (a,b,c), 使得 a^2 + b^2 = c^2, 且 c <= cLimit.
 * @param cLimit 限制勾股数的最大值.
 * @param f 回调函数.
 * @param coprimeOnly 是否只枚举两两互质的勾股数对.默认为false.
 * cLimit = 1e8：互素对有 1.59*1e7 个, 0.4s;
 * cLimit = 1e8：全部有 2.71*1e8 个, 1s.
 */
function enumeratePytagoreanTriple(
  cLimit: number,
  f: (a: number, b: number, c: number) => void,
  coprimeOnly = false
): void {
  let stack: { a: number; b: number; c: number }[] = []
  const add = (a: number, b: number, c: number): void => {
    if (c <= cLimit) stack.push({ a, b, c })
  }
  add(3, 4, 5)
  while (stack.length) {
    const cur = stack.pop()!
    const { a, b, c } = cur
    add(a - 2 * b + 2 * c, 2 * a - b + 2 * c, 2 * a - 2 * b + 3 * c)
    add(a + 2 * b + 2 * c, 2 * a + b + 2 * c, 2 * a + 2 * b + 3 * c)
    add(-a + 2 * b + 2 * c, -2 * a + b + 2 * c, -2 * a + 2 * b + 3 * c)
    if (coprimeOnly) {
      f(a, b, c)
    } else {
      let x = a
      let y = b
      let z = c
      while (z <= cLimit) {
        f(x, y, z)
        x += a
        y += b
        z += c
      }
    }
  }
}

export { enumeratePytagoreanTriple }

if (require.main === module) {
  console.time('enumeratePytagoreanTriple')
  let count = 0
  enumeratePytagoreanTriple(
    1e8,
    (a, b, c) => {
      count++
    },
    true
  )
  console.log(count)
  console.timeEnd('enumeratePytagoreanTriple')
}

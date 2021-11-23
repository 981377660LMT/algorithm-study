/**
 *
 * @param a
 * @param b
 * @description
 * 求ax+by=gcd(a,b) 的一组解
 * @returns
 * [x,y,gcd(a,b)]
 */
function exgcd(a: number, b: number): [x: number, y: number, gcd: number] {
  if (b === 0) return [1, 0, a]
  let [x, y, gcd] = exgcd(b, a % b)
  // 根据(b，a%b)的解推出(a,b)的解
  // 辗转相除法反向推导每层a、b的因子使得gcd(a,b)=ax+by成立
  ;[x, y] = [y, x - ~~(a / b) * y]
  return [x, y, gcd]
}
if (require.main === module) {
  console.log(exgcd(3, 10))
}

export { exgcd }

import { GCD } from '../../最大公约数/gcd'
import { isPrime } from '../../素数/isPrime'
import { exgcd } from '../扩展欧几里得/扩展欧几里得'

/**
 *
 * @param num
 * @param p
 * @returns
 * 拓展欧几里得求num 模 p 的乘法逆元
 * num*x 同余 1 (mod P)
 */
function calInv1(num: number, p: number): number {
  const [x, _, gcd] = exgcd(num, p)
  if (gcd === 1) return ((x % p) + p) % p
  else return -1
}

function calInv2(num: number, p: number): number {
  if (GCD(num, p) !== 1 || !isPrime(p)) throw new Error('无法用费马小定理求逆元')
  return num ** (p - 2) % p
}

if (require.main === module) {
  console.log(calInv1(21, 25))
}

export { calInv1 }

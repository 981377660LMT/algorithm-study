import { qpow } from '../../数字/快速幂/qpow'
import { exgcd } from '../扩展欧几里得/扩展欧几里得'

/**
 *
 * @param num
 * @param p
 * @returns
 * 拓展欧几里得求num 模 p 的乘法逆元
 * num*x 同余 1 (mod P)
 * 即求num*inv(num)+k*p=1 mod p中的 inv(num)
 */
function calInv1(num: number, p: number): number {
  const [x, _, gcd] = exgcd(num, p)
  if (gcd === 1) return ((x % p) + p) % p
  else return -1
}

function calInv2(num: bigint, p: bigint): bigint {
  // if (GCD(num, p) !== 1 || !isPrime(p)) throw new Error('无法用费马小定理求逆元')
  return qpow(Number(num), Number(p - 2n), 1e9 + 7)
}

if (require.main === module) {
  console.log((2 * calInv1(2, 1e9 + 7)) % (1e9 + 7))
  console.log((2n * calInv2(2n, BigInt(1e9 + 7))) % BigInt(1e9 + 7))
}

export { calInv1, calInv2 }

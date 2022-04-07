// 题目里：hash(s, p, m) = (val(s[0]) * p0 + val(s[1]) * p1 + ... + val(s[k-1]) * pk-1) mod m.
// 注意我们的RK算法里计算哈希值的方法是左边字符权重大，题目是右边权重大
// 所以要把我们的字符串反过来，调api，哈希值相等时返回这一段的reversed

import { BigIntHasher } from '../BigIntHasher'

function subStrHash(
  s: string,
  power: number,
  modulo: number,
  k: number,
  hashValue: number
): string {
  BigIntHasher.setBASE(power)
  BigIntHasher.setMOD(modulo)
  BigIntHasher.setOFFSET(96)

  s = s.split('').reverse().join('')
  const hasher = new BigIntHasher(s)

  let res = 0
  for (let i = 0; i + k <= s.length; i++) {
    const hash = hasher.getHashOfSlice(i, i + k)
    if (Number(hash) === hashValue) res = i
  }

  return s
    .slice(res, res + k)
    .split('')
    .reverse()
    .join('')
}

console.log(subStrHash('leetcode', 7, 20, 2, 0))
console.log(subStrHash('fbxzaad', 31, 100, 3, 32))

export {}

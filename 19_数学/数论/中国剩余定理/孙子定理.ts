// 有物不知其数，三三数之剩二，五五数之剩三，七七数之剩二。问物几何？
// x 同余 2 (mod3)
// x 同余 3 (mod5)
// x 同余 2 (mod7)
// 即解三个同余方程
// 35x1 同余 2 (mod3)
// 21x2 同余 3 (mod5)
// 15x3 同余 2 (mod7)
// 即求逆元w1 w2 w3
// 35w1 同余 1 (mod3)
// 21w2 同余 1 (mod5)
// 15w3 同余 1 (mod7)

import { modularInverse } from '../扩展欧几里得/扩展欧几里得'

/**
 * x ≡ remains_i (mod mods_i)
 * @param remains 余数数组
 * @param mods 模数数组
 */
function ChineseRemainderTheorem(remains: number[], mods: number[]) {
  const n = mods.length
  let res = 0

  const p = mods.reduce((pre, cur) => pre * cur, 1)
  for (let i = 0; i < n; i++) {
    const tmp = p / mods[i]
    // modularInverse(tmp, div[i]) * tmp等于1(模p意义下) 所以 tmp*inv(tmp)*mod[i] 与 mod[i] 模p同余
    res = (res + modularInverse(tmp, mods[i]) * tmp * remains[i]) % p
  }

  return (res + p) % p
}

console.log(ChineseRemainderTheorem([2, 3, 2], [3, 5, 7]))

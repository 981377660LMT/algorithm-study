// 模运算 加/减/乘/除/幂/逆元

function modAdd(num1: number, num2: number, mod = 1e9 + 7): number {
  let cand = (num1 + num2) % mod
  if (cand < 0) cand += mod
  return cand
}

function modSub(num1: number, num2: number, mod = 1e9 + 7): number {
  return modAdd(num1, -num2, mod)
}

function modMul(num1: number, num2: number, mod = 1e9 + 7): number {
  return (((Math.floor(num1 / 65536) * num2) % mod) * 65536 + (num1 & 65535) * num2) % mod
}

function modDiv(num1: number, num2: number, mod = 1e9 + 7): number {
  return modMul(num1, modInv(num2, mod), mod)
}

/**
 * 模逆元.需要保证mod是质数.
 */
function modInv(num: number, mod = 1e9 + 7): number {
  return modPow(num, mod - 2, mod)
}

function modPow(num: number, pow: number, mod = 1e9 + 7): number {
  num = modAdd(num, 0)
  let res = 1
  while (pow) {
    if (pow & 1) res = modMul(res, num, mod)
    num = modMul(num, num, mod)
    pow = Math.floor(pow / 2)
  }
  return res
}

export { modAdd, modSub, modMul, modDiv, modInv, modPow }

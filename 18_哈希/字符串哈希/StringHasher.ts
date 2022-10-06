/* eslint-disable no-param-reassign */

import assert from 'assert'

/**
 * 哈希值计算方法：
 *
 * hash(s, p, m) = (val(s[0]) * pk-1 + val(s[1]) * pk-2 + ... + val(s[k-1]) * p0) mod m.
 * 越靠左字符权重越大
 */
function useStringHasher(s: ArrayLike<string>, mod = 10 ** 11 + 7, base = 1313131, offset = 0) {
  const bigMod = BigInt(mod)
  const bigBase = BigInt(base)
  const bigOffset = BigInt(offset)

  const n = s.length
  const prePow = new BigUint64Array(n + 1)
  prePow[0] = 1n
  const preHash = new BigUint64Array(n + 1)

  for (let i = 1; i < n + 1; i++) {
    prePow[i] = (prePow[i - 1] * bigBase) % bigMod
    preHash[i] = (preHash[i - 1] * bigBase + BigInt(s[i - 1].charCodeAt(0)) - bigOffset) % bigMod
  }

  return getSliceHash

  /**
   * !切片 s[left:right] 的哈希值
   */
  function getSliceHash(left: number, right: number): bigint {
    if (left >= right) return 0n
    left++
    return (
      (preHash[right] - ((preHash[left - 1] * prePow[right - left + 1]) % bigMod) + bigMod) % bigMod
    )
  }
}

/**
 * 哈希值计算方法：
 *
 * hash(s, p, m) = (val(s[0]) * pk-1 + val(s[1]) * pk-2 + ... + val(s[k-1]) * p0) mod m.
 * 越靠左字符权重越大
 */
function useArrayHasher(
  arr: ArrayLike<number | bigint>,
  mod = 10 ** 11 + 7,
  base = 1313131,
  offset = 0
) {
  const bigMod = BigInt(mod)
  const bigBase = BigInt(base)
  const bigOffset = BigInt(offset)

  const n = arr.length
  const prePow = new BigUint64Array(n + 1)
  prePow[0] = 1n
  const preHash = new BigUint64Array(n + 1)

  for (let i = 1; i < n + 1; i++) {
    prePow[i] = (prePow[i - 1] * bigBase) % bigMod
    preHash[i] = (preHash[i - 1] * bigBase + BigInt(arr[i - 1]) - bigOffset) % bigMod
  }

  return getSliceHash

  /**
   * !切片 arr[left:right] 的哈希值
   */
  function getSliceHash(left: number, right: number): bigint {
    if (left >= right) return 0n
    left++
    return (
      (preHash[right] - ((preHash[left - 1] * prePow[right - left + 1]) % bigMod) + bigMod) % bigMod
    )
  }
}

if (require.main === module) {
  const s1 = 'asdasd'
  const hasher = useStringHasher(s1)
  assert.strictEqual(hasher(0, 1), 97n)
  assert.strictEqual(hasher(0, 2), 97n * 1313131n + 115n)
  assert.strictEqual(hasher(0, 3), hasher(3, 6))

  const arr1 = [1, 2, 3, 4, 1, 2]
  const hasher2 = useArrayHasher(arr1)
  assert.strictEqual(hasher2(0, 1), 1n)
  assert.strictEqual(hasher(0, 2), hasher(3, 5))
}

export { useArrayHasher, useStringHasher }

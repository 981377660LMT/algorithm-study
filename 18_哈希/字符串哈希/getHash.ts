/**
 * 将 `字符串` 或 `正整数数组` 转换为哈希值.
 * 使用双哈希减小哈希冲突概率.
 */
function getHash(
  arr: string | ArrayLike<number>,
  base1 = 131,
  base2 = 13331,
  mod1 = 1e7 + 19,
  mod2 = 1e7 + 79
): [hash1: number, hash2: number] {
  if (!arr.length) return [0, 0]
  let hash1 = 0
  let hash2 = 0
  if (typeof arr === 'string') {
    for (let i = 0; i < arr.length; i++) {
      const v = arr.charCodeAt(i)
      hash1 = (hash1 * base1 + v) % mod1
      hash2 = (hash2 * base2 + v) % mod2
    }
  } else {
    for (let i = 0; i < arr.length; i++) {
      const v = arr[i]
      hash1 = (hash1 * base1 + v) % mod1
      hash2 = (hash2 * base2 + v) % mod2
    }
  }
  return [hash1, hash2]
}

export { getHash }

if (require.main === module) {
  console.log(getHash('abc'))
  console.log(getHash([1, 1, 1]))
  console.log(getHash([1, 1]))
}

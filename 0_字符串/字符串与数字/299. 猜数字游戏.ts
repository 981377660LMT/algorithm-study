/**
 * @param {string} secret
 * @param {string} guess
 * @return {string}
 * 请注意秘密数字和朋友的猜测数都可能含有重复数字，每位数字只能统计一次
 */
var getHint = function (secret: string, guess: string): string {
  let bulls = 0
  let cows = 0
  const n = secret.length
  const map = new Map<string, Set<number>>()

  for (let i = 0; i < n; i++) {
    const char = secret[i]
    map.set(char, (map.get(char) || new Set()).add(i))
  }

  for (let i = 0; i < n; i++) {
    const char = guess[i]
    if (!map.has(char)) continue
    if (map.get(char)!.has(i)) bulls++
    else cows++
  }
  console.log(map)
  return `${bulls}A${cows}B`
}

console.log(getHint('1807', '7810'))
console.log(getHint('1123', '0111'))
// 输出: "1A3B"
// 解释: 1 公牛和 3 奶牛。公牛是 8，奶牛是 0, 1 和 7。

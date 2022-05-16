// 给定一个字符串，编写一个函数判定其是否为某个回文串的排列之一。
// 字符全为小写时
function canPermutePalindrome(s: string): boolean {
  let xor = 0
  for (const char of s) {
    const ascii = 1 << (char.codePointAt(0)! - 97)
    xor ^= ascii
  }

  // 异或为0或2的幂
  return (xor & (xor - 1)) === 0
}

// 字符不全为小写
function canPermutePalindrome2(s: string): boolean {
  const set = new Set<string>()
  for (const char of s) {
    if (set.has(char)) set.delete(char)
    else set.add(char)
  }

  // 异或为0或2的幂
  return set.size <= 1
}

console.log(canPermutePalindrome('tactcoa'))

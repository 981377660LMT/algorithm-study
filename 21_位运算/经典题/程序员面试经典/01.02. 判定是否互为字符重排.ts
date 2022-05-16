function CheckPermutation(s1: string, s2: string): boolean {
  if (s1.length !== s2.length) return false

  let xor = 0
  let diff = 0
  for (let i = 0; i < s1.length; i++) {
    const key1 = s1[i].codePointAt(0)! - 97
    const key2 = s2[i].codePointAt(0)! - 97
    xor ^= key1 ^ key2
    diff += key1 - key2
  }

  return xor === 0 && diff === 0
}

console.log(CheckPermutation('abc', 'bca'))

// 移位 哈希 计算差距 考虑1
// "abc" -> "bcd"。这样，我们可以持续进行 “移位” 操作，从而生成如下移位序列：

// 能通过移位得到的必须满足每位之间差距相等
function groupStrings(strings: string[]): string[][] {
  const counter = new Map<string, string[]>()

  for (const str of strings) {
    const hash = calculateHash(str)
    if (!counter.has(hash)) counter.set(hash, [])
    counter.get(hash)!.push(str)
  }

  return [...counter.values()]

  function calculateHash(str: string): string {
    if (str.length <= 1) return '#'

    const tmp: number[] = []
    for (let index = 0; index < str.length - 1; index++) {
      const diff = str[index + 1].codePointAt(0)! - str[index].codePointAt(0)!
      tmp.push((diff + 26) % 26)
    }

    return tmp.join('#')
  }
}

console.log(groupStrings(['abc', 'bcd', 'acef', 'xyz', 'az', 'ba', 'a', 'z']))
// 输出：
// [
//   ["abc","bcd","xyz"],
//   ["az","ba"],
//   ["acef"],
//   ["a","z"]
// ]

import { BigIntHasher } from '../BigIntHasher'

// 编写一个函数来找出所有目标子串，目标子串的长度为 10，且在 DNA 字符串 s 中出现次数超过一次。
function findRepeatedDnaSequences(s: string): string[] {
  const stringHasher = new BigIntHasher(s.toLowerCase())
  const counter = new Map<bigint, number>()
  const res: string[] = []

  for (let i = 1; i + 9 <= s.length; i++) {
    const hash = stringHasher.getHashOfRange(i, i + 9)
    counter.set(hash, (counter.get(hash) || 0) + 1)
    if (counter.get(hash) === 2) res.push(s.slice(i - 1, i - 1 + 10))
  }

  return res
}

console.log(
  findRepeatedDnaSequences(
    'CGACGCAATTTAGAACGGGCCGCACTGCAACCATTGCTCAGACAACGCATGAGTTAAATTTCACAAGTGATAGTGGCTTGCGAGACGTGGGTTGGTGGTAGCGTACGCCCGCTATTCGCCCCTAACGTGACGGGATTATAAGGTCGCTTCCCGGAATGCGCAGACGAGTCTCCGGTTTAGCCTAGACGTCTCACGCGCGCAAGGCGTCAGTTCTCACTATCTCGCACAGGTGTATTCTATTAGTTATGGGTTCTCACTACAGTCGGTTACTTCCTCATCCATTTCTGCATACGGGTCAACTAACAGTGTCATGGGGTATTGGGAAGGATGCGTTTTTAAACCCTCTCAGTAGCGCGAGGATGCCCACAAATACGACGGCGGCCACGGATCTAATGCGAAGCTAGCGACGCTTTCCAGCAACGAGCGCCCCACTTATGACTGCGTGGGGCGCTCCGCTTTCCTAGAGAACATAGATGGTGTTTTCGAATCGTAACCACTTA'
  )
)

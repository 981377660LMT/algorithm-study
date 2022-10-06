import { useStringHasher } from '../StringHasher'

// 编写一个函数来找出所有目标子串，目标子串的长度为 10，且在 DNA 字符串 s 中出现次数超过一次。
function findRepeatedDnaSequences(s: string): string[] {
  const getSliceHash = useStringHasher(s)
  const counter = new Map<bigint, number>()
  const res: string[] = []

  for (let i = 0; i + 10 <= s.length; i++) {
    const hash = getSliceHash(i, i + 10)
    counter.set(hash, (counter.get(hash) || 0) + 1)
    if (counter.get(hash) === 2) res.push(s.slice(i, i + 10))
  }

  return res
}

console.log(
  findRepeatedDnaSequences(
    'CGACGCAATTTAGAACGGGCCGCACTGCAACCATTGCTCAGACAACGCATGAGTTAAATTTCACAAGTGATAGTGGCTTGCGAGACGTGGGTTGGTGGTAGCGTACGCCCGCTATTCGCCCCTAACGTGACGGGATTATAAGGTCGCTTCCCGGAATGCGCAGACGAGTCTCCGGTTTAGCCTAGACGTCTCACGCGCGCAAGGCGTCAGTTCTCACTATCTCGCACAGGTGTATTCTATTAGTTATGGGTTCTCACTACAGTCGGTTACTTCCTCATCCATTTCTGCATACGGGTCAACTAACAGTGTCATGGGGTATTGGGAAGGATGCGTTTTTAAACCCTCTCAGTAGCGCGAGGATGCCCACAAATACGACGGCGGCCACGGATCTAATGCGAAGCTAGCGACGCTTTCCAGCAACGAGCGCCCCACTTATGACTGCGTGGGGCGCTCCGCTTTCCTAGAGAACATAGATGGTGTTTTCGAATCGTAACCACTTA'
  )
)

import { useUnionFindArray } from './推荐使用并查集精简版'

// 如果交换字符串 X 中的两个不同位置的字母，使得它和字符串 Y 相等，那么称 X 和 Y 两个字符串相似
function numSimilarGroups(strs: string[]): number {
  const uf = useUnionFindArray(strs.length)

  for (let i = 0; i < strs.length; i++) {
    for (let j = i + 1; j < strs.length; j++) {
      // 相似字符串具有传递性=>并查集
      if (isSimilar(strs[i], strs[j])) uf.union(i, j)
    }
  }

  return uf.getCount()

  function isSimilar(word1: string, word2: string): boolean {
    let diff = 0

    for (let i = 0; i < word1.length; i++) {
      if (word1[i] !== word2[i]) diff++
      if (diff > 2) return false
    }

    return diff === 0 || diff === 2
  }
}

console.log(numSimilarGroups(['tars', 'rats', 'arts', 'star']))
// 它们通过相似性形成了两个关联组：{"tars", "rats", "arts"} 和 {"star"}
// 给你一个字符串列表 strs。列表中的每个字符串都是 strs 中其它所有字符串的一个字母异位词。
// 请问 strs 中有多少个相似字符串组？

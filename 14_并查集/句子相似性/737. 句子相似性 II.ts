// 相似关系是 具有 传递性的
// 相似关系是具有对称性的

import { useUnionFindArray } from '../useUnionFind'

// 即：无向图并查集
function areSentencesSimilarTwo(
  sentence1: string[],
  sentence2: string[],
  similarPairs: string[][]
): boolean {
  if (sentence1.length !== sentence2.length) return false
  const uf = useUnionFindArray(similarPairs.length * 2)
  const words = [...new Set<string>(similarPairs.flat())]
  const wordId = new Map<string, number>([...words.entries()].map(([i, v]) => [v, i]))
  for (const [u, v] of similarPairs) {
    uf.union(wordId.get(u)!, wordId.get(v)!)
  }

  for (let i = 0; i < sentence1.length; i++) {
    if (sentence1[i] === sentence2[i]) continue
    const [id1, id2] = [wordId.get(sentence1[i]), wordId.get(sentence2[i])]
    if (id1 == undefined || id2 == undefined) return false
    if (!uf.isConnected(id1, id2)) return false
  }

  return true
}

// console.log(
//   areSentencesSimilarTwo(
//     ['great', 'acting', 'skills'],
//     ['fine', 'drama', 'talent'],
//     [
//       ['great', 'fine'],
//       ['acting', 'drama'],
//       ['skills', 'talent'],
//     ]
//   )
// )
console.log(
  areSentencesSimilarTwo(
    ['I', 'love', 'leetcode'],
    ['I', 'love', 'onepiece'],
    [
      ['manga', 'onepiece'],
      ['platform', 'anime'],
      ['leetcode', 'platform'],
      ['anime', 'manga'],
    ]
  )
)
console.log(areSentencesSimilarTwo(['I', 'love', 'leetcode'], ['I', 'love', 'onepiece'], []))

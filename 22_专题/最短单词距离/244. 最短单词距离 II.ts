import { minDiff } from '../k个数组/1_两个有序数组选数差的绝对值最小'

class WordDistance {
  private indexes: Map<string, number[]>

  constructor(wordsDict: string[]) {
    this.indexes = new Map()
    for (let i = 0; i < wordsDict.length; i++) {
      const word = wordsDict[i]
      !this.indexes.has(word) && this.indexes.set(word, [])
      this.indexes.get(word)!.push(i)
    }
  }

  // 此方法将被以不同的参数调用 多次
  shortest(word1: string, word2: string): number {
    const nums1 = this.indexes.get(word1)!
    const nums2 = this.indexes.get(word2)!
    return minDiff(nums1, nums2)
  }
}

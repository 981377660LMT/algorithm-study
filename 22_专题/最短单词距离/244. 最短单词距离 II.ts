import { minDiff } from '../k个数组/1_两个有序数组选数差的绝对值最小'

class WordDistance {
  private wordLocation: Map<string, number[]>

  constructor(wordsDict: string[]) {
    this.wordLocation = new Map()
    for (let i = 0; i < wordsDict.length; i++) {
      const word = wordsDict[i]
      !this.wordLocation.has(word) && this.wordLocation.set(word, [])
      this.wordLocation.get(word)!.push(i)
    }
  }

  // 此方法将被以不同的参数调用 多次
  shortest(word1: string, word2: string): number {
    const nums1 = this.wordLocation.get(word1)!
    const nums2 = this.wordLocation.get(word2)!

    // 此处由于location列表有序，可以是双指针优化
    // for (const index1 of this.wordLocation.get(word1)!) {
    //   for (const index2 of this.wordLocation.get(word2)!) {
    //     res = Math.min(res, Math.abs(index1 - index2))
    //   }
    // }

    return minDiff(nums1, nums2)
  }
}

/** 支持新增元素、删除元素、快速判断容器内是否包含了指定的所有元素. */
class SlidingWindowContainsAll<T> {
  private _missingCount: number
  private readonly _missing: Map<T, number>

  constructor(n: number, supplier: (i: number) => T) {
    this._missing = new Map()
    for (let i = 0; i < n; i++) {
      const v = supplier(i)
      this._missing.set(v, (this._missing.get(v) || 0) + 1)
    }
    this._missingCount = this._missing.size
  }

  add(v: T): boolean {
    const c = this._missing.get(v)
    if (c === undefined) return false
    this._missing.set(v, c - 1)
    if (c === 1) this._missingCount--
    return true
  }

  discard(v: T): boolean {
    const c = this._missing.get(v)
    if (c === undefined) return false
    this._missing.set(v, c + 1)
    if (c === 0) this._missingCount++
    return true
  }

  containsAll(): boolean {
    return this._missingCount === 0
  }
}

export { SlidingWindowContainsAll }

if (require.main === module) {
  // 3298. 统计重新排列后包含另一个字符串的子字符串数目 II
  // https://leetcode.cn/problems/count-substrings-that-can-be-rearranged-to-contain-a-string-ii/description/
  // eslint-disable-next-line no-inner-declarations
  function validSubstringCount(word1: string, word2: string): number {
    const S = new SlidingWindowContainsAll(word2.length, i => word2[i])
    let res = 0
    let left = 0
    const n = word1.length
    for (let right = 0; right < n; right++) {
      S.add(word1[right])
      while (left <= right && S.containsAll()) {
        S.discard(word1[left])
        left++
      }
      res += left
    }
    return res
  }
}

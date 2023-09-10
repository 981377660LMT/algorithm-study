// https://zhuanlan.zhihu.com/p/408665473

import assert from 'assert'

/**
 * 构建某个模式串(较短串)的失配数组.
 * @param shorter 模式串或者模式串的unicode编码数组.
 * @returns
 * 失配数组.`next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度.
 * 在AC自动机中被命名为`fail`.
 */
function getNext(shorter: string | ArrayLike<number>): Uint32Array {
  const next = new Uint32Array(shorter.length)
  let j = 0
  for (let i = 1; i < shorter.length; i++) {
    while (j > 0 && shorter[i] !== shorter[j]) {
      // 新来了一个字符如果不匹配，需要上跳fail指针到最长公共后缀结尾处，看下一个(子节点)是否能匹配
      j = next[j - 1]
    }
    if (shorter[i] === shorter[j]) j++
    next[i] = j
  }
  return next
}

/**
 * `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
 */
function indexOfAll<S extends string | ArrayLike<number> = string>(
  longer: S,
  shorter: S,
  position = 0
): number[] {
  if (shorter.length === 0) return [0]
  if (longer.length < shorter.length) return []
  const res: number[] = []
  const next = getNext(shorter)
  let hitJ = 0
  for (let i = position; i < longer.length; i++) {
    while (hitJ > 0 && longer[i] !== shorter[hitJ]) hitJ = next[hitJ - 1]
    if (longer[i] === shorter[hitJ]) hitJ++
    if (hitJ === shorter.length) {
      res.push(i - shorter.length + 1)
      hitJ = next[hitJ - 1] // 不允许重叠时 hitJ = 0
    }
  }
  return res
}

class KMP<T extends string | ArrayLike<number> = string> {
  /**
   * `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
   * @param longer 搜索串.
   * @param shorter 模式串.
   * @param position 搜索的起始位置.
   * @returns 所有匹配的位置.
   */
  static indexOfAll<S extends string | ArrayLike<number> = string>(
    longer: S,
    shorter: S,
    position = 0
  ): number[] {
    if (shorter.length === 0) return [0]
    if (longer.length < shorter.length) return []
    const res: number[] = []
    const next = this.getNext(shorter)
    let hitJ = 0
    for (let i = position; i < longer.length; i++) {
      while (hitJ > 0 && longer[i] !== shorter[hitJ]) hitJ = next[hitJ - 1]
      if (longer[i] === shorter[hitJ]) hitJ++
      if (hitJ === shorter.length) {
        res.push(i - shorter.length + 1)
        hitJ = next[hitJ - 1] // 不允许重叠时 hitJ = 0
      }
    }
    return res
  }

  static getNext(shorter: string | ArrayLike<number>): Uint32Array {
    const next = new Uint32Array(shorter.length)
    let j = 0
    for (let i = 1; i < shorter.length; i++) {
      while (j > 0 && shorter[i] !== shorter[j]) j = next[j - 1]
      if (shorter[i] === shorter[j]) j++
      next[i] = j
    }
    return next
  }

  private readonly _pattern: T
  private readonly _next: Uint32Array

  /**
   * @param pattern `模式串`或者`模式串`的unicode编码数组.
   * 注意模式串是`较短`的字符串，搜索串是`较长`的字符串.
   */
  constructor(pattern: T) {
    this._pattern = pattern
    this._next = KMP.getNext(pattern)
  }

  search(searchString: T, position = 0): number {
    if (searchString.length < this._pattern.length) return -1
    let pos = 0
    for (let i = position; i < searchString.length; i++) {
      pos = this.move(pos, searchString[i])
      if (this.accept(pos)) return i - this._pattern.length + 1
    }
    return -1
  }

  /**
   * `o(n+m)`求搜索串 {@link searchString} 中所有匹配 {@link _pattern} 的位置.
   * @param searchString 搜索串.
   * @param position 搜索的起始位置.
   * @returns 所有匹配的位置.
   * @alias indexOfAll/findAll
   */
  searchAll(searchString: T, position = 0): number[] {
    if (searchString.length < this._pattern.length) return []
    const res: number[] = []
    let pos = 0
    for (let i = position; i < searchString.length; i++) {
      pos = this.move(pos, searchString[i])
      if (this.accept(pos)) {
        res.push(i - this._pattern.length + 1)
        pos = this._next[pos - 1]
      }
    }
    return res
  }

  move(pos: number, input: T[0]): number {
    if (pos < 0 || pos >= this._pattern.length) throw new RangeError(`pos: ${pos} is out of range`)
    while (pos && input !== this._pattern[pos]) {
      pos = this._next[pos - 1]
    }
    if (input === this._pattern[pos]) {
      pos++
    }
    return pos
  }

  accept(pos: number): boolean {
    return pos === this._pattern.length
  }

  /**
   * 求字符串{@link _pattern[:end]}的最短周期.如果不存在,返回0.
   * @param end 0<=end<=n.
   */
  period(end = this._pattern.length): number {
    end--
    if (end < 0 || end >= this._pattern.length) throw new RangeError(`end: ${end} is out of range`)
    const res = end + 1 - this._next[end]
    if (res && end + 1 > res && (end + 1) % res === 0) return res
    return 0
  }
}

export { getNext, indexOfAll, getNext as getLPS, KMP }

if (require.main === module) {
  assert.deepStrictEqual(getNext('ababaaa'), new Uint32Array([0, 0, 1, 2, 3, 1, 1]))
  assert.deepStrictEqual(indexOfAll('ababaaa', 'ab'), [0, 2])

  // https://leetcode.cn/problems/find-the-index-of-the-first-occurrence-in-a-string/
  // eslint-disable-next-line no-inner-declarations
  function strStr(haystack: string, needle: string): number {
    const kmp = new KMP(needle)
    return kmp.search(haystack)
  }

  const shorter = 'a'
  const kmp = new KMP(shorter)
  console.log(kmp.move(0, 'b'))

  // !无字符串拷贝(slice)的子串匹配
  const shorter2 = 'ab'
  const ords1 = new Uint32Array(shorter2.split('').map(c => c.charCodeAt(0)))
  const kmp2 = new KMP(ords1)
  const longer2 = 'ababababababaaa'
  const ords2 = new Uint32Array(longer2.split('').map(c => c.charCodeAt(0)))
  console.log(kmp2.searchAll(ords2.subarray(0, 10)))
}

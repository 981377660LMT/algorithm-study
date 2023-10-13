/* eslint-disable no-inner-declarations */

// erase, prev, next 的时间复杂度均为 O(1)
// erase: 删除某个元素
// prev: 找到前一个元素
// next: 找到后一个元素

class FinderLinkedList {
  private readonly _n: number
  private readonly _prev: Int32Array
  private readonly _next: Int32Array
  private readonly _erased: Uint8Array

  /**
   * 使用双向链表维护前驱后继.
   * 初始时, 0~n-1 个元素都是未访问过的.
   */
  constructor(n: number) {
    const prev = new Int32Array(n + 2)
    const next = new Int32Array(n + 2)
    const erased = new Uint8Array(n)
    for (let i = 1; i < n + 1; i++) {
      prev[i] = i - 1
      next[i] = i + 1
    }
    this._n = n
    this._prev = prev
    this._next = next
    this._erased = erased
  }

  /**
   * 删除元素i.
   */
  erase(i: number): boolean {
    if (this._erased[i]) return false
    this._erased[i] = 1
    i++
    this._prev[this._next[i]] = this._prev[i]
    this._next[this._prev[i]] = this._next[i]
    return true
  }

  /**
   * 判断元素i是否存在.
   * 0<=i<n.
   */
  has(i: number): boolean {
    return !this._erased[i]
  }

  /**
   * 找到i左侧第一个未被访问过的位置(包含i).
   * 如果不存在, 返回null.
   */
  prev(i: number): number | null {
    if (this.has(i)) return i
    const res = this._prev[i + 1] - 1
    return res >= 0 ? res : null
  }

  /**
   * 找到i右侧第一个未被访问过的位置(包含i).
   * 如果不存在, 返回null.
   */
  next(i: number): number | null {
    if (this.has(i)) return i
    const res = this._next[i + 1] - 1
    console.log(this._prev)
    return res < this._n ? res : null
  }
}

export { FinderLinkedList }

if (require.main === module) {
  function check(): void {
    const n = 3
    const finder1 = new NextFinder(n)
    const finder2 = new FinderLinkedList(n)
    for (let i = 0; i < n; i++) {
      const x = Math.floor(Math.random() * n)

      if (finder1.has(x) !== finder2.has(x)) throw new Error('has')

      for (let i = 0; i < 3; i++) {
        console.log(finder1.next(i), finder2.next(i), 'index', i)
      }
      if (finder1.next(x) !== finder2.next(x)) {
        console.log(x, finder1.next(x), finder2.next(x))
        throw new Error('next')
      }
      finder1.erase(x)
      finder2.erase(x)
      console.log(`erase ${x}`)
    }
  }
  check()
}

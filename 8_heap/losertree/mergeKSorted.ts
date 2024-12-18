/* eslint-disable no-inner-declarations */

/**
 * 合并K个有序数据结构.
 *
 * @description 时间复杂度`O(nlogk)`, 空间复杂度`O(logk)`.`k`为有序数据结构的个数，`n`为所有数据的总个数.
 * @throws 如果`sortedItems`为空，则抛出错误.
 */
function mergeKSorted<T>(sorted: ArrayLike<T>, merge: (a: T, b: T) => T): T {
  const n = sorted.length
  if (n === 0) throw new Error('sortedItems is empty')
  if (n === 1) return sorted[0]
  if (n === 2) return merge(sorted[0], sorted[1])

  const f = (start: number, end: number): T => {
    if (end - start === 1) return sorted[start]
    const mid = (start + end) >>> 1
    return merge(f(start, mid), f(mid, end))
  }
  return f(0, n)
}

export { mergeKSorted }

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  class ListNode {
    val: number
    next: ListNode | null
    constructor(val?: number, next?: ListNode | null) {
      this.val = val === undefined ? 0 : val
      this.next = next === undefined ? null : next
    }
  }

  function mergeKLists(lists: Array<ListNode>): ListNode | null {
    if (lists.length === 0) return null

    const merge = (a: ListNode | null, b: ListNode | null): ListNode | null => {
      const dummy = new ListNode()
      let cur = dummy
      while (a !== null && b !== null) {
        if (a.val < b.val) {
          cur.next = a
          a = a.next
        } else {
          cur.next = b
          b = b.next
        }
        cur = cur.next
      }
      cur.next = a !== null ? a : b
      return dummy.next
    }

    return mergeKSorted(lists, merge)
  }
}

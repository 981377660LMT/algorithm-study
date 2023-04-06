/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable eqeqeq */
/* eslint-disable no-inner-declarations */

/**
 * 弗洛伊德探环法(floyd探环法).
 *
 * 给定一个`首项为 s0 , s[i]=next(s[i-1]) (i>=1)` 的序列，求环的起点和周期(长度)
 * 返回值 start 为环的起点，period 为环的长度.
 * 即 `s[i] = s[i+period] (i>=start)`.
 * !O(start+period) 时间复杂度.
 *
 * @param s0 序列的首项
 * @param next 序列的下一项
 */
function floydCycleFind<T>(
  s0: T,
  next: (cur: T) => T | null
): [start: number, period: number] | [start: null, period: null] {
  let p1 = 0
  let p2 = 1
  let slow: T | null = s0
  if (slow == null) return [null, null]
  let fast: T | null = next(s0)
  if (fast == null) return [null, null]

  while (slow !== fast) {
    const n1 = next(fast)
    if (n1 == null) return [null, null]
    const n2 = next(n1)
    if (n2 == null) return [null, null]
    fast = n2
    p2 += 2
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    slow = next(slow!)
    p1++
  }

  // !has cycle now
  fast = s0
  for (let _ = 0; _ < p2 - p1; _++) {
    fast = next(fast!)
  }

  let start = 0
  slow = s0
  while (slow !== fast) {
    slow = next(slow!)
    fast = next(fast!)
    start++
  }

  let period = 1
  for (fast = next(slow!); slow !== fast; period++) {
    fast = next(fast!)
  }

  return [start, period]
}

if (require.main === module) {
  // https://leetcode.cn/problems/linked-list-cycle-ii/
  // 寻找环形链表的环的起点

  class ListNode {
    val: number
    next: ListNode | null
    constructor(val: number, next: ListNode | null) {
      this.val = val
      this.next = next
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  function detectCycle(head: ListNode | null): ListNode | null {
    if (!head) return null
    const [start] = floydCycleFind(head, cur => cur.next)
    if (start == null) return null
    for (let i = 0; i < start; i++) {
      head = head!.next
    }
    return head
  }
}

export { floydCycleFind }

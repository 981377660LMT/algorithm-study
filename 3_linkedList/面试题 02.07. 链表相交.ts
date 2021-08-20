class ListNode {
  value: number | undefined
  next: ListNode | undefined
  constructor(value?: number, next?: ListNode) {
    this.value = value
    this.next = next
  }
}

const a = new ListNode(1)
const b = new ListNode(2)
const c = new ListNode(3)
a.next = b
b.next = c
const d = new ListNode(-1)
const e = new ListNode(0)
d.next = e
e.next = b

/**
 * @param {ListNode} headA
 * @param {ListNode} headB
 * @return {ListNode}
 * 交点不是数值相等，而是指针相等
 * 你能否设计一个时间复杂度 O(n) 、仅用 O(1) 内存的解决方案
 * @summary 计算长度 让两个指针从相同index出发
 */
const getIntersectionNode = function (headA: ListNode, headB: ListNode): ListNode | null {
  const getLength = (head: ListNode) => {
    let res = 0
    let headP: ListNode | undefined = head
    while (headP) {
      res++
      headP = headP.next
    }
    return res
  }

  const headALen = getLength(headA)
  const headBLen = getLength(headB)
  console.assert(headALen === 3 && headBLen === 4)
  let longListP = headALen > headBLen ? headA : headB
  let shortListP = headALen > headBLen ? headB : headA
  const lenDiff = Math.abs(headALen - headBLen)

  for (let i = 0; i < lenDiff; i++) {
    longListP = longListP.next!
  }

  while (longListP && longListP.next) {
    if (longListP === shortListP) return longListP
    longListP = longListP.next
    shortListP = shortListP.next!
  }

  return null
}

console.dir(getIntersectionNode(a, d))

export default 1

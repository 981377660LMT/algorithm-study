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
 * @summary 
 * 当 a 到达链表的尾部时,重定位到链表 B 的头结点
   当 b 到达链表的尾部时,重定位到链表 A 的头结点。
   a, b 指针相遇的点为相交的起始节点，否则没有相交点
 */
const getIntersectionNode = function (headA: ListNode, headB: ListNode): ListNode | undefined {
  let a: ListNode | undefined = headA
  let b: ListNode | undefined = headB
  while (a !== b) {
    a = a === undefined ? headB : a.next
    b = b === undefined ? headA : b.next
  }
  return a
}

console.dir(getIntersectionNode(a, d))

export default 1

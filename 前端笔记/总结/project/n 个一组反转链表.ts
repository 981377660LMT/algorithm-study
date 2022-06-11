// n个一组反转链表
class ListNode {
  val: number
  next: ListNode | null
  constructor(val?: number, next?: ListNode | null) {
    this.val = val === undefined ? 0 : val
    this.next = next === undefined ? null : next
  }
}

function reverse(node: ListNode | null): ListNode | null {
  let pre: ListNode | null = null
  let cur: ListNode | null = node
  while (cur) {
    const next = cur.next
    cur.next = pre
    pre = cur
    cur = next
  }

  return pre
}

// 递归 三步
function reverseKGroup1(head: ListNode | null, k: number): ListNode | null {
  if (!head) return head
  let headP = head

  // 找到这一段
  for (let _ = 0; _ < k - 1; _++) {
    headP = headP.next!
    if (!headP) return head
  }

  const nextHead = headP.next
  headP.next = null
  reverse(head)

  head.next = reverseKGroup1(nextHead, k)

  return headP
}
export {}

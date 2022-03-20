class ListNode {
  val: number
  next: ListNode | null
  constructor(val?: number, next?: ListNode | null) {
    this.val = val === undefined ? 0 : val
    this.next = next === undefined ? null : next
  }
}
const a = new ListNode(1)
const b = new ListNode(2)
const c = new ListNode(3)
const d = new ListNode(4)
const e = new ListNode(5)
a.next = b
b.next = c
c.next = d
d.next = e

// 思路和两个一组翻转链表相同
// 每 k 个节点一组进行翻转
// 1.判断存在
// 2.reverseTwo
// 3.递归
const reverseKGroup = (head: ListNode | null, k: number): ListNode | null => {
  // 1.判断连续的k个节点是否存在 否则返回head
  // 边界条件:例如k=2时需要head和head.next都存在,循环完后headP指向k个数最后一个
  if (!head) return head
  let headP = head

  for (let i = 0; i < k - 1; i++) {
    headP = headP.next!
    if (!headP) return head
  }

  // 2. reverseTwo k 这一段 反转后headP在头部
  const nextHead = headP.next
  headP.next = null
  reverse(head)

  // 3.递归
  head.next = reverseKGroup(nextHead, k)

  return headP

  function reverse(head: ListNode) {
    let p1: ListNode | null = null
    let p2: ListNode | null = head
    while (p2) {
      const tmp: ListNode | null = p2.next
      p2.next = p1
      p1 = p2
      p2 = tmp
    }
    return p1
  }
}

console.dir(reverseKGroup(a, 3), { depth: null })
// 3,2,1,4,5
export {}

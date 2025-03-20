export {}

class ListNode {
  val: number
  next: ListNode | null
  constructor(val?: number, next?: ListNode | null) {
    this.val = val === undefined ? 0 : val
    this.next = next === undefined ? null : next
  }
}

// 特殊情况处理
// 如果链表为空、只有一个节点或 k 为 0，则直接返回 head。
//
// 求链表长度并构造环
// 遍历链表求得长度 len，同时找到链表尾节点 tail，然后将 tail.next 指向 head，构成环。
//
// 计算新的断开点
// k 可能非常大，所以先计算 k % len；如果结果为 0，则无需旋转。
// 新的尾节点位于位置 len - k - 1，新头节点为其 next。
//
// 断开环
// 将新的尾节点的 next 置为 null，即可得到旋转后的链表。
function rotateRight(head: ListNode | null, k: number): ListNode | null {
  if (!head || !head.next || k === 0) return head

  let len = 1
  let tail = head
  while (tail.next) {
    len++
    tail = tail.next
  }

  k %= len
  if (k === 0) return head

  // 构成环
  tail.next = head

  // 4. 找到新的尾节点，位于 len - k - 1 位置
  let newTail = head
  for (let _ = 0; _ < len - k - 1; _++) {
    newTail = newTail.next!
  }

  // 5. 新的头节点为 newTail.next，并断开环
  const newHead = newTail.next
  newTail.next = null

  return newHead
}

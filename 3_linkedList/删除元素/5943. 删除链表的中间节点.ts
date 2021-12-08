// 删除链表的中间节点
// 返回修改后的链表的头节点 head 。
function deleteMiddle(head: ListNode | null): ListNode | null {
  const dummy = new ListNode(0, head)
  let slow: ListNode = dummy
  let fast: ListNode | null = dummy
  while (fast && fast.next && fast.next.next) {
    fast = fast.next.next
    slow = slow.next!
  }

  slow.next = slow.next!.next
  return dummy.next
}

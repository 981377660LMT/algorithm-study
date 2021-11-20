// 开始时以头节点作为当前节点.
// 保留以当前节点开始的前 m 个节点.
// 删除接下来的 n 个节点.
// 重复步骤 2 和 3, 直到到达链表结尾.

function deleteNodes(head: ListNode | null, m: number, n: number): ListNode | null {
  const dummy = new ListNode(0, head)
  let dummyP: ListNode | null = dummy

  while (dummyP && dummyP.next) {
    for (let i = 0; i < m && dummyP && dummyP.next; i++) {
      dummyP = dummyP.next
    }

    for (let i = 0; i < n && dummyP && dummyP.next; i++) {
      dummyP.next = dummyP.next.next
    }
  }

  return dummy.next
}

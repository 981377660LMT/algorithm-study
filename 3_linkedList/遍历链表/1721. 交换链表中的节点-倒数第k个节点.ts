// 交换 链表正数第 k 个节点和倒数第 k 个节点的值后，返回链表的头节点（链表 从 1 开始索引）。

// 获取倒数第k个节点:快的先走k步，制造k距离；慢的快的一起走直到快的消失
function swapNodes(head: ListNode | null, k: number): ListNode | null {
  let first = head
  let last = head
  for (let i = 1; i < k; i++) {
    first = first!.next
  }

  let nullChecker = first
  while (nullChecker && nullChecker.next) {
    last = last!.next
    nullChecker = nullChecker.next
  }

  ;[first!.val, last!.val] = [last!.val, first!.val]

  return head
}

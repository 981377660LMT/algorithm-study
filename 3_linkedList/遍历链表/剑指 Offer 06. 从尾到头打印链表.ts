class Node {
  val: number
  next: Node | undefined
  constructor(val: number = 0, next?: Node) {
    this.val = val
    this.next = next
  }
}

// 输入一个链表的头节点，从尾到头反过来返回每个节点的值（用数组返回）。
function reversePrint(head: Node | null): number[] {
  const res: number[] = []
  if (head) {
    res.push(...reversePrint(head.next!))
    res.push(head.val)
  }

  return res
}

export {}

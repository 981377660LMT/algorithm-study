class Node {
  value: number | undefined
  next: Node | undefined
  constructor(value?: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

function getKthFromEnd(head: Node | undefined, k: number): Node | undefined {
  let slow: Node | undefined = head
  let fast: Node | undefined = head

  for (let i = 0; i < k; i++) {
    fast = fast!.next
  }

  while (fast) {
    slow = slow!.next
    fast = fast.next
  }

  return slow
}

export {}

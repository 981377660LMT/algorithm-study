class Node {
  value: number
  next: Node | null
  constructor(value: number = 0, next: Node | null = null) {
    this.value = value
    this.next = next
  }
}
const a1 = new Node(1)
const b1 = new Node(4)
const c1 = new Node(7)
a1.next = b1
b1.next = c1

const a2 = new Node(2)
const b2 = new Node(5)
const c2 = new Node(8)
a2.next = b2
b2.next = c2

const a3 = new Node(3)
const b3 = new Node(6)
const c3 = new Node(9)
a3.next = b3
b3.next = c3

// 这里使用递归写法
const mergeKLists = (...heads: (Node | null)[]): Node | null => {
  // 这里使用递归写法
  const mergeTwoLists = (head1: Node | null, head2: Node | null): Node | null => {
    if (!head1) return head2
    if (!head2) return head1
    if (head1.value < head2.value) {
      head1.next = mergeTwoLists(head1.next, head2)
      return head1
    } else {
      head2.next = mergeTwoLists(head1, head2.next)
      return head2
    }
  }
  return heads.reduce(mergeTwoLists)
}

console.dir(mergeKLists(a1, a3, a2), { depth: null })
export default 1

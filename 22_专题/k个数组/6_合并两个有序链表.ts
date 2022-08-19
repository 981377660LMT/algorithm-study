/* eslint-disable no-else-return */
class Node {
  value: number
  next: Node | null
  constructor(value = 0, next: Node | null = null) {
    this.value = value
    this.next = next
  }
}
const a1 = new Node(1)
const b1 = new Node(3)
const c1 = new Node(5)
a1.next = b1
b1.next = c1

const a2 = new Node(2)
const b2 = new Node(4)
const c2 = new Node(6)
a2.next = b2
b2.next = c2

// 递归写法
const mergeTwoList = (head1: Node | null, head2: Node | null): Node | null => {
  if (!head1) return head2
  if (!head2) return head1
  if (head1.value < head2.value) {
    head1.next = mergeTwoList(head1.next, head2)
    return head1
  } else {
    head2.next = mergeTwoList(head1, head2.next)
    return head2
  }
}

// 非递归写法
const merge = (head1: Node | null, head2: Node | null): Node | null => {
  if (!head1) return head2
  if (!head2) return head1
  const dummy = new Node(-1)
  let dummyP: Node | null = dummy
  let head1P: Node | null = head1
  let head2P: Node | null = head2

  while (head1P && head2P) {
    if (head1P.value <= head2P.value) {
      dummyP.next = head1P
      head1P = head1P.next
    } else {
      dummyP.next = head2P
      head2P = head2P.next
    }
    dummyP = dummyP?.next
  }

  dummyP.next = head1P || head2P
  return dummy.next
}

// console.dir(mergeTwoList(a1, a2), { depth: null })
console.dir(merge(a1, a2), { depth: null })
export default 1

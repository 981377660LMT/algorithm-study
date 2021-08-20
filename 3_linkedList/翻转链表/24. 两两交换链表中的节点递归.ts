class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
const d = new Node(4)
const e = new Node(5)
const f = new Node(6)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 三指针:head,head.next,head.next.next
// 1.判断存在
// 2.reverse
// 3.递归
const swapPairs = (head: Node): Node => {
  if (!head || !head.next) return head

  let p1 = head
  let p2 = head.next
  let p3 = head.next.next

  p2.next = p1
  p1.next = swapPairs(p3!)

  return p2
}

console.dir(swapPairs(a), { depth: null })
// 2,1,4,3,6,5
export {}

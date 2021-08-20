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

const swapPairs = (head: Node): Node => {
  let dummy = new Node(0, head)
  let dummyP: Node = dummy

  while (dummyP.next && dummyP.next.next) {
    let p1 = dummyP.next
    let p2 = dummyP.next.next
    // 三步
    p1.next = p2.next
    p2.next = p1
    dummyP.next = p2
    // 准备下一次翻转
    dummyP = p1
  }

  return dummy.next!
}

console.dir(swapPairs(a), { depth: null })
// 2,1,4,3,6,5
export {}

class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(4)
const c = new Node(3)
const d = new Node(2)
const e = new Node(5)
const f = new Node(2)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 使得所有 小于 x 的节点都出现在 大于或等于 x 的节点之前。
// 你应当 保留 两个分区中每个节点的初始相对位置。

// 思路: 使用两个指针，一对取大一队取小，最后串起来
const partition = (head: Node, x: number) => {
  let small = new Node(0)
  let big = new Node(0)
  let smallP = small
  let bigP = big
  let headP = head
  while (headP) {
    if (headP.value < x) {
      smallP.next = headP
      smallP = smallP.next
    } else {
      bigP.next = headP
      bigP = bigP.next
    }
    headP = headP.next!
  }

  smallP.next = big.next
  bigP.next = undefined

  return small.next
}

console.dir(partition(a, 3), { depth: null })
export {}

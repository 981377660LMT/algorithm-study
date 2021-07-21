class Node {
  value: number
  next?: Node
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
a.next = b
b.next = c
c.next = d
d.next = e

// 思路：deque
// 遍历节点将所有节点全部push入队然后pop shift操作
const reorderList = (head: Node): void => {
  const deque: Node[] = []

  let headP1: Node | undefined = head
  while (headP1) {
    deque.push(headP1)
    headP1 = headP1.next
  }

  let headP2: Node | undefined = head
  const len = deque.length
  for (let i = 0; i < len; i++) {
    if (i % 2 === 0) {
      headP2.next = deque.shift()
    } else {
      headP2.next = deque.pop()
    }
    headP2 = headP2.next!
  }

  headP2.next = undefined
}
reorderList(a)
console.dir(a, { depth: null })
// 1,5 2,4,3
export {}

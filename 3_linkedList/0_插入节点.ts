class Node {
  value: number | undefined
  next: Node | undefined
  constructor(value?: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
a.next = b
b.next = c

const inserted = new Node(0)
const cur = b
// 插入节点
inserted.next = b.next
cur.next = inserted

console.dir(a, { depth: null })

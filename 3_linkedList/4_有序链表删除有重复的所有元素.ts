class Node {
  value: number | undefined
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(1)
const c = new Node(1)
const d = new Node(2)
const e = new Node(3)
a.next = b
b.next = c
c.next = d
d.next = e

// 删除有重复的所有节点:思路，先遍历一次记录重复值，再遍历一次删除节点
const unique = (node: Node): void => {}

console.dir(a, { depth: null })

unique(a)

console.log(a)
// 返回2,3
export {}

// 遍历链表和删除链表节点

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
const c = new Node(2)
const d = new Node(3)
a.next = b
b.next = c
c.next = d

// 遍历链表，如果当前元素等于下个元素值，就删除下个元素值
const unique = (node: Node): void => {
  let n1: Node | undefined = node

  while (n1) {
    //@ts-ignore
    const nextNode = n1.next
    if (n1.value === nextNode?.value) {
      n1.next = nextNode?.next
    }
    n1 = n1.next
  }
}

console.dir(a, { depth: null })

unique(a)

console.log(a)

export {}

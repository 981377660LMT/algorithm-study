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

const reverseList = (head: Node) => {
  // 注意p1p2都是node
  let n1: Node | undefined = undefined
  let n2: Node | undefined = head
  while (n2) {
    //@ts-ignore
    // p2的下一个Node
    const tmp = n2.next
    // 最重要的关系，为串联节点做准备
    n2.next = n1
    n1 = n2
    n2 = tmp
  }

  return n1
}

console.log(reverseList(a))
export {}

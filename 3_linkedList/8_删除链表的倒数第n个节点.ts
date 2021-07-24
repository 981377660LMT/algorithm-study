class Node {
  value: number
  next: Node | undefined
  constructor(value: number = 0, next?: Node) {
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

// 节点添加到数组，然后重新创建
// 便利了两遍链表
const removeNthFromEnd = (head: Node | undefined, n: number) => {
  const newNode = new Node()
  let newNodeP = newNode
  const res: number[] = []
  // 遍历链表
  while (head) {
    res.push(head.value)
    head = head.next
  }

  res.splice(res.length - n, 1)

  res.forEach(num => {
    newNodeP.next = new Node(num)
    newNodeP = newNodeP.next
  })

  return newNode.next || []
}

// console.dir(removeNthFromEnd(a, 2), { depth: null })
// console.dir(removeNthFromEnd(d, 1), { depth: null })
console.dir(removeNthFromEnd(e, 1), { depth: null })

export {}

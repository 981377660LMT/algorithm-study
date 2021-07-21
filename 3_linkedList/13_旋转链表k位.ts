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

// 不是返回head,使用虚拟头节点

// 思路:
// 1.使用快慢指针找到倒数第k个节点，作为dummy的next
// 2.断开第k-1个节点与后面的连接
// 2.将第k个节点后的这一段接在原来的头上
const rotateRight = (head: Node, k: number) => {
  const dummy = new Node(0)

  const last = findLastK(head, 1)
  const lastK = findLastK(head, k)
  const lastKPre = findLastK(head, k + 1)
  dummy.next = lastK
  lastKPre!.next = undefined
  last!.next = head

  return dummy.next
}

const findLastK = (head: Node, k: number) => {
  let slow: Node | undefined = head
  let fast: Node | undefined = head

  for (let i = 0; i < k; i++) {
    fast = fast?.next
  }

  while (fast) {
    slow = slow?.next
    fast = fast?.next
  }

  return slow
}

// const k = findLastK(a, 2)
// k!.next = undefined
// console.dir(k, { depth: null })
// console.dir(a, { depth: null })
console.dir(rotateRight(a, 2), { depth: null })

export {}

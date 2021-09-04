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
  const dummy = new Node(0, head)
  let fast: Node | undefined = dummy
  let slow: Node | undefined = dummy
  let len = 0

  while (fast && fast.next) {
    fast = fast.next
    len++
  }
  console.log(len)
  k = k % len
  for (let i = 0; i <= len - k; i++) {
    slow = slow?.next
  }

  fast.next = dummy.next
  dummy.next = slow?.next
  slow!.next = undefined

  return dummy.next
}

// const k = findLastK(a, 2)
// k!.next = undefined
// console.dir(k, { depth: null })
// console.dir(a, { depth: null })
console.dir(rotateRight(a, 2), { depth: null })

export {}

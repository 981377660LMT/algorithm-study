class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

// 参考合并k个链表
// 采用队列+最小堆(顶部元素最小)的思想
const mergeTwo = (l1: Node, l2: Node) => {
  const res = new Node(0)
  let p = res
  const queue: Node[] = l1.value > l2.value ? [l2, l1] : [l1, l2]
  while (queue.length) {
    const head = queue.shift()!
    p.next = head
    p = p.next
    head.next && queue.push(head.next)
    if (queue.length == 2 && queue[0].value > queue[1].value) {
      ;[[queue[0], queue[1]]] = [[queue[1], queue[0]]]
    }
  }

  return res.next
}
const a = new Node(1)
const b = new Node(4)
const c = new Node(2)
const d = new Node(3)
const e = new Node(5)
a.next = b
c.next = d
d.next = e

console.dir(mergeTwo(a, c), { depth: null })
export {}

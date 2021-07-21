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

/**
 *
 * @param head 链表头
 * @param index index从0开始
 */
// const selectNode = (head: Node, index: number) => {
//   let p: Node | undefined = head
//   for (let i = 0; i < index; i++) {
//     p = p?.next
//   }
//   return p
// }

// console.log(selectNode(a, 3))

// 最后的样子是 首-----p1--p4(start对应节点)---p3(end对应节点)-p2-p2.next-------尾
// 需要四个指针(p1=>p3=>p4=>p2)
const reverseList = (head: Node, start: number, end: number) => {
  // 注意p1p2都是node
  let p1: Node | undefined = head
  let p2: Node | undefined = head
  let i = 1
  // 初始化起点
  while (i < start) {
    p1 = p2
    p2 = p2?.next
    i++
  }

  let p3: Node | undefined = undefined
  let p4 = p2
  while (i <= end) {
    const tmp = p2?.next
    p2!.next = p3
    p3 = p2
    p2 = tmp
    i++
  }

  p1!.next = p3
  p4!.next = p2

  // start为 1时，头被反转，新的头是p3
  return start === 1 ? p3 : head
}

console.dir(reverseList(a, 1, 4), { depth: null })
export {}

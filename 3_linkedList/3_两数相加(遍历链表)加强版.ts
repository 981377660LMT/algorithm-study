class Node {
  value: number | undefined
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(7)
const b = new Node(2)
const c = new Node(4)
const d = new Node(3)
a.next = b
b.next = c
c.next = d
const e = new Node(5)
const f = new Node(6)
const g = new Node(4)
e.next = f
f.next = g

const convertNode = (head: Node) => {
  let p1: Node | undefined = undefined
  let p2: Node | undefined = head
  while (p2) {
    // @ts-ignore
    const tmp = p2.next
    p2.next = p1
    p1 = p2
    p2 = tmp
  }

  return p1
}

// console.dir(convertNode(a), { depth: null })

// 在两数相加的基础上多了反转链表的步骤
// 如果不允许修改原来的node 则要使用stack存储 遍历链表将node的val全部push到stack中再全部pop
const addTwo = (l1: Node, l2: Node): Node => {
  const newNode = new Node(0)
  let p1 = convertNode(l1)
  let p2 = convertNode(l2)
  let p3 = newNode
  let overflow = 0

  while (p1 || p2) {
    const v1 = p1?.value || 0
    const v2 = p2?.value || 0
    const sum = v1 + v2 + overflow
    overflow = Math.floor(sum / 10)
    const res = sum % 10
    const nextNode = new Node(res)
    p3.next = nextNode

    p1 = p1?.next
    p2 = p2?.next
    p3 = p3.next
  }

  if (overflow) {
    p3.next = new Node(overflow)
  }

  return convertNode(newNode.next!)!
}

console.dir(addTwo(a, e), { depth: null })
// l1 = [7,2,4,3], l2 = [5,6,4]
// 输出:[7,8,0,7]
export {}

// O(n)

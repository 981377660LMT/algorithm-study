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

// 如果不允许修改原来的node 则要使用stack存储 遍历链表将node的val全部push到stack中再全部pop
// 先push进栈的值正好后出
const addTwo = (l1: Node, l2: Node): Node => {
  const stack1: number[] = []
  const stack2: number[] = []
  let p1: Node | undefined = l1
  let p2: Node | undefined = l2
  const newNode = new Node(0)
  let p3 = newNode
  let carry = 0

  while (p1) {
    stack1.push(p1.value!)
    p1 = p1.next!
  }
  while (p2) {
    stack2.push(p2.value!)
    p2 = p2.next!
  }

  while (stack1.length || stack2.length) {
    const v1 = stack1.pop() || 0
    const v2 = stack2.pop() || 0
    const sum = v1 + v2 + carry
    carry = Math.floor(sum / 10)
    const res = sum % 10

    // 注意这里创建的是新的头节点,不是新的next节点
    const head = new Node(res)
    head.next = p3
    p3 = head
  }

  if (carry) {
    const head = new Node(carry)
    head.next = p3
    p3 = head
  }

  return p3
}

console.dir(addTwo(a, e), { depth: null })
// l1 = [7,2,4,3], l2 = [5,6,4]
// 输出:[7,8,0,7]
export {}

// O(n)

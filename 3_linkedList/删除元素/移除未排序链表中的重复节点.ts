class Node {
  value: number
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
const e = new Node(2)
const f = new Node(3)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 递归写法
function removeDuplicateNodes(head: Node | undefined): Node | undefined {
  const helper = (head: Node | undefined, set: Set<number>): Node | undefined => {
    if (!head) return
    if (set.has(head.value)) return helper(head.next, set)
    set.add(head.value)
    head.next = helper(head.next, set)
    return head
  }

  return helper(head, new Set())
}

export default 1

console.log(removeDuplicateNodes(a))

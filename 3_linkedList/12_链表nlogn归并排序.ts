class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(3)
const c = new Node(5)
const d = new Node(4)
const e = new Node(2)
const f = new Node(6)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = f

// 自顶向下递归
// 自底向上迭代

// 类比数组
// 1.split方法 使用快慢指针寻找链表的中间节点，然后分成两份
// 2.merge方法 合并两个有序链表
const sortList = (head?: Node): Node | undefined => {
  // 长度小于等于1时返回
  if (!head?.next) return head
  const [left, right] = split(head)
  return merge(sortList(left), sortList(right))
}

const split = (head?: Node) => {
  let slow: Node | undefined = head
  let fast: Node | undefined = head
  while (fast?.next?.next) {
    slow = slow?.next
    fast = fast.next.next
  }
  // 分割为[0,slow]和[slow+1,list.size]
  const left = head
  const right = slow?.next
  // 注意要断开slow与slow之后的连接
  slow!.next = undefined

  return [left, right]
}

const merge = (left?: Node, right?: Node): Node | undefined => {
  const dummy = new Node(0)
  let dummyP = dummy
  while (left && right) {
    if (left.value > right.value) {
      dummyP.next = right
      right = right.next
    } else {
      dummyP.next = left
      left = left.next
    }
    dummyP = dummyP.next
  }

  // 连接剩余的元素
  if (left) dummyP.next = left
  if (right) dummyP.next = right

  return dummy.next
}

// console.dir(split(a), { depth: null })
// const foo = new Node(1)
// const bar = new Node(2)
// console.log(merge(foo, bar))
console.dir(sortList(a), { depth: null })

export {}

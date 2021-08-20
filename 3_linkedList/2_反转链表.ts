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
a.next = b
b.next = c

// 递归法反转和迭代法反转都要会
// 递归法反转

const reverseList1 = (head: Node | undefined) => {
  const reverse = (pre: Node | undefined, cur: Node | undefined): Node => {
    if (!cur) return pre!
    const tmp = cur.next
    cur.next = pre
    // 如下递归的写法，其实就是做了这两步
    //     // pre = cur;
    //     // cur = temp;
    return reverse(cur, tmp)
  }

  return reverse(undefined, head)
}

// 迭代法反转
const reverseList = (head: Node) => {
  // 注意p1p2都是node
  let n1: Node | undefined = undefined
  let n2: Node | undefined = head
  while (n2) {
    // n2的下一个Node;n2一开始是head,反转后head的next是undefined
    const tmp: Node | undefined = n2.next
    // 最重要的关系，为串联节点做准备
    n2.next = n1
    n1 = n2
    n2 = tmp
  }

  return n1
}

console.log(reverseList(a))
export {}

class Node {
  value: number
  next?: Node
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(3)
const c = new Node(5)
const d = new Node(7)
const e = new Node(9)
a.next = b
b.next = c
c.next = d
d.next = e
e.next = a

/**
 *
 * @param head
 * @param insertVal
 * 给定循环升序列表中的一个点，写一个函数向这个列表中插入一个新元素 insertVal ，使这个列表仍然是循环非降序的。
 * 如果node.val== insertVal，则将node插入node后面
   如果node.val<insertVal, node.next.val>=insertVal，将node插入node后面
   链表完全遍历一次后，如果未找到满足上面条件，说明insertVal超过最大值或小于最小值，这种情况下将insertVal插入到maxNode后面。

 */
function insert(head: Node | null, insertVal: number): Node | null {
  const insertNode = new Node(insertVal)
  if (!head) {
    insertNode.next = insertNode
    return insertNode
  }

  let cur = head
  let maxNode = head
  while (true) {
    if (maxNode.value <= cur.value) maxNode = cur
    if (insertVal === cur.value || (cur.value < insertVal && insertVal <= cur.next!.value)) {
      insertAfter(cur, insertNode)
      return head
    }

    cur = cur.next!
    if (cur === head) break
  }

  insertAfter(maxNode, insertNode)
  return head

  function insertAfter(pre: Node, insertNode: Node) {
    insertNode.next = pre.next
    pre.next = insertNode
  }
}

console.dir(insert(a, 6), { depth: null })

export {}

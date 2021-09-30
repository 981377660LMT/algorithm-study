class RandomNode {
  value: number
  next: RandomNode | undefined
  random: RandomNode | undefined
  constructor(value: number, next?: RandomNode, random?: RandomNode) {
    this.value = value
    this.next = next
    this.random = random
  }
}

const a = new RandomNode(1)
const b = new RandomNode(2)
const c = new RandomNode(3)
a.next = b
a.random = b
b.next = c
c.random = a

// 构造这个链表的 深拷贝。
// 启示：深拷贝可以 建立单项映射关系 遍历两次
const copyRandomList = (head: RandomNode) => {
  if (!head) return
  const record = new WeakMap<RandomNode, RandomNode | undefined>()

  let headP: RandomNode | undefined = head
  while (headP) {
    record.set(headP, new RandomNode(headP.value))
    headP = headP.next
  }

  headP = head
  while (headP) {
    headP.next && (record.get(headP)!.next = record.get(headP.next))
    headP.random && (record.get(headP)!.random = record.get(headP.random))
    headP = headP.next
  }

  return record.get(head)
}

console.dir(copyRandomList(a), { depth: null })

export {}

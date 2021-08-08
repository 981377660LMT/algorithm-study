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
// 遍历两次
const copyRandomList = (head: RandomNode) => {
  if (!head) return
  const map = new WeakMap<RandomNode, RandomNode | undefined>()

  let headP: RandomNode | undefined = head
  while (headP) {
    map.set(headP, new RandomNode(headP.value))
    headP = headP.next
  }

  headP = head
  while (headP) {
    map.get(headP)!.next = map.get(headP.next!) || undefined
    map.get(headP)!.random = map.get(headP.random!) || undefined
    headP = headP.next
  }

  return map.get(head)
}

console.dir(copyRandomList(a), { depth: null })

export {}

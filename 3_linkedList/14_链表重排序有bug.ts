class Node {
  value: number
  next?: Node
  constructor(value: number, next?: Node) {
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

// 思路：双指针
// 类比数组的双指针
const reorderList = (head: Node) => {
  const dummy = new Node(0)
  let dummyP = dummy
  let l = dummy
  let lIndex = 0
  let [r, rIndex] = findLastK(head, 1)

  if (typeof rIndex === 'number') {
    while (lIndex <= rIndex) {
      dummyP.next = l
      dummyP.next.next = r

      lIndex++
      rIndex--
      l = l.next!
      const [rPre, index] = findLastK(head, rIndex + 1)
      rIndex = index
      r = rPre

      dummyP = dummyP.next.next!
    }
  }

  return dummy.next
}

const findLastK = (head: Node, k: number) => {
  let index = 0
  let slow: Node | undefined = head
  let fast: Node | undefined = head

  for (let i = 0; i < k; i++) {
    fast = fast?.next
  }

  while (fast) {
    slow = slow?.next
    fast = fast?.next
    index++
  }

  return [slow, index] as const
}

console.dir(reorderList(a), { depth: null })
// 1,5 2,4,3
export {}
export {}

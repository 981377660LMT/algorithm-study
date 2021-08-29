class Node {
  value: number
  next?: Node
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(4)
const b = new Node(2)
const c = new Node(1)
const d = new Node(3)
a.next = b
b.next = c
c.next = d

const bubbleSort = (head: Node | undefined) => {
  if (!head || !head.next) return head

  while (true) {
    let isSorted = true
    let headP: Node | undefined = head
    while (headP && headP.next) {
      if (headP.value > headP.next.value) {
        isSorted = false
        ;[headP.value, headP.next.value] = [headP.next.value, headP.value]
      }
      headP = headP.next
    }
    if (isSorted) break
  }

  return head
}

console.dir(bubbleSort(a), { depth: null })
export {}

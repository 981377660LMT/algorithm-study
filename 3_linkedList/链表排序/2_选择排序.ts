class Node {
  value: number
  next?: Node
  constructor(value = 0, next?: Node) {
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

const selectSort = (head: Node | undefined) => {
  if (!head || !head.next) return head

  // 相当于i
  let headP: Node | undefined = head

  while (headP) {
    // 相当于minIndex
    let minNode = headP
    // 相当于j
    let candidate = minNode.next

    while (candidate) {
      if (minNode.value > candidate.value) {
        minNode = candidate
      }
      candidate = candidate.next
    }

    ;[headP.value, minNode.value] = [minNode.value, headP.value]
    headP = headP.next
  }

  return head
}

console.dir(selectSort(a), { depth: null })
export {}

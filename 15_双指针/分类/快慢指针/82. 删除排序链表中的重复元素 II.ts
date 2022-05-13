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

const deleteDuplicates = function (head: Node) {
  if (!head) return null
  const K = 2
  // preKValue队列记录前K个数
  let preKValue = Array<number>(K).fill(Infinity)
  let slow: Node | undefined = head
  let fast: Node | undefined = head

  while (fast) {
    if (fast.value !== preKValue[0]) {
      slow!.value = fast.value
      preKValue[0] = fast.value
      fast.next && (slow = slow!.next)
      preKValue.shift()
      preKValue.push(fast.value)
    }
    fast = fast.next
  }

  console.log(slow?.value)
  slow!.next = undefined
  return head
}

console.dir(deleteDuplicates(a), { depth: null })
export {}

class Node {
  value: number
  next: Node | undefined
  constructor(value: number = 0, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(2)
const b = new Node(1)
const c = new Node(5)
a.next = b
b.next = c

/**
 * @param {Node} head
 * @return {number[]}
 */
function nextLargerNodes(head: Node): number[] {
  const res: number[] = []
  const stack: [number, number][] = []
  let headP: Node | undefined = head
  let index = 0

  while (headP) {
    while (stack.length && stack[stack.length - 1][1] < headP.value) {
      const [i, _] = stack.pop()!
      res[i] = headP.value
    }
    stack.push([index, headP.value])
    headP = headP.next
    index++
  }

  for (const [i, _] of stack) {
    res[i] = 0
  }

  return res
}

console.log(nextLargerNodes(a))
export {}

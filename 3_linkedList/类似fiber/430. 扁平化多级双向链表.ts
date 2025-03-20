export {}

class _Node {
  val: number
  prev: _Node | null
  next: _Node | null
  child: _Node | null

  constructor(val?: number, prev?: _Node, next?: _Node, child?: _Node) {
    this.val = val === undefined ? 0 : val
    this.prev = prev === undefined ? null : prev
    this.next = next === undefined ? null : next
    this.child = child === undefined ? null : child
  }
}

function flatten(head: _Node | null): _Node | null {
  if (!head) return null
  const dummy = new _Node()
  dummy.next = head

  const stack: _Node[] = [head]
  let pre = dummy
  while (stack.length) {
    const cur = stack.pop()!
    pre.next = cur
    cur.prev = pre
    if (cur.next) stack.push(cur.next)
    if (cur.child) {
      stack.push(cur.child)
      cur.child = null
    }
    pre = cur
  }

  dummy.next!.prev = null
  return dummy.next
}

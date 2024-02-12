/* eslint-disable prefer-destructuring */
/* eslint-disable @typescript-eslint/no-this-alias */
/* eslint-disable max-len */

const LOGM = 30
const M = 1 << LOGM
const MSTEP = (M / LOGM) | 0

class UpperNode {
  label: number
  prev: UpperNode | undefined
  next: UpperNode | undefined

  constructor(label: number, prev: UpperNode | undefined, next: UpperNode | undefined) {
    this.label = label
    this.prev = prev
    this.next = next
  }

  insertAfter(): UpperNode {
    let adjustCount = 1
    let cur = this.next
    const startLabel = this.label
    let gap = 1
    for (; cur && cur.label - startLabel <= adjustCount * adjustCount; ++adjustCount, cur = cur.next) {
      gap = Math.ceil((cur.label - startLabel) / adjustCount) | 0
    }
    cur = this.next
    for (let k = 1; k < adjustCount; ++k, cur = cur!.next) {
      cur!.label = (startLabel + gap * k) | 0
    }
    const newLabel = (startLabel + this.next!.label) >>> 1
    const res = new UpperNode(newLabel, this.next, this)
    this.next!.prev = res
    this.next = res
    return res
  }

  remove(): void {
    this.next!.prev = this.prev
    this.prev!.next = this.next
    this.next = undefined
    this.prev = undefined
  }

  split(): UpperNode {
    if (!this.next) return createUpperList()
    const other = this.next
    this.next = new UpperNode(M - 1, undefined, undefined)
    other.prev = new UpperNode(0, undefined, undefined)
    return other
  }

  compare(other: UpperNode): number {
    return this.label - other.label
  }

  getOrderedKey(): number {
    return this.label
  }
}

class LowerNode<V> {
  upper: UpperNode
  label: number
  next: LowerNode<V> | undefined
  prev: LowerNode<V> | undefined
  value: V

  constructor(upper: UpperNode, label: number, next: LowerNode<V> | undefined, prev: LowerNode<V> | undefined, value: V) {
    this.upper = upper
    this.label = label
    this.next = next
    this.prev = prev
    this.value = value
  }

  insertAfter(value: V): LowerNode<V> {
    let n = M
    // Create node and link it.
    let res = new LowerNode(this.upper, -1, this.next, this, value)
    if (this.next) {
      n = this.next.label
      this.next.prev = res
    }
    this.next = res
    // Update labels.
    if (n === this.label + 1) {
      // Scan to extents of subtree.
      let begin: LowerNode<V> = this
      while (begin.prev && begin.prev.upper === this.upper) {
        begin = begin.prev
      }
      let end: LowerNode<V> = this
      while (end.next && end.next.upper === this.upper) {
        end = end.next
      }
      end = end.next!
      // Redistribute nodes.
      let upper = this.upper
      let cur = begin
      while (true) {
        // Relabel nodes.
        let label = 0
        for (let j = 0; j < LOGM; ++j, cur = cur.next!, label += MSTEP) {
          if (cur === end) {
            return res
          }
          cur.label = label
          cur.upper = upper
        }
        upper = upper.insertAfter()
      }
    } else {
      res.label = Math.min((this.label + n) >>> 1, this.label + LOGM) | 0
    }
    return res
  }

  remove(): void {
    let uniqueUpper = true
    if (this.next) {
      this.next.prev = this.prev
      uniqueUpper = this.next.upper !== this.upper
    }
    if (this.prev) {
      this.prev.next = this.next
      uniqueUpper = uniqueUpper && this.prev.upper !== this.upper
    }
    if (uniqueUpper) {
      this.upper.remove()
    }
  }

  split(): LowerNode<V> | undefined {
    if (!this.next) return undefined
    const other = this.next
    const newUpper = this.upper.split()
    if (newUpper.prev) {
      newUpper.prev
    }
    this.next = undefined
    other.prev = undefined
    for (let cur = other; cur && cur.upper === this.upper; cur = cur.next!) {
      cur.upper = newUpper
    }
    return other
  }

  /**
   * !比较的 key 为 `(upper.label, lower.label)`.
   */
  compare(other: LowerNode<V>): number {
    return this.upper.compare(other.upper) || this.label - other.label
  }

  getOrderedKey(): [upperKey: number, lowerKey: number] {
    return [this.upper.label, this.label]
  }
}

function createUpperList(): UpperNode {
  const begin = new UpperNode(0, undefined, undefined)
  const end = new UpperNode(M - 1, undefined, begin)
  begin.next = end
  return begin
}

/**
 * A data structure for ordered list maintenance.
 * This generalizes a linked list, except it adds the ability to query the order of any two elements in the list in constant time.
 * This implementation is based on Bender's O(1) amortized algorithm using two-level indirection.
 * The upper list is a linked list with a length of O(n/logn),
 * and each child of the upper list, the lower list, is a linked list with a length of O(logn).
 *
 * {@link https://github.com/mikolalysenko/order-maintenance}
 * {@link https://courses.csail.mit.edu/6.851/spring12/lectures/L08.html}
 */
function createOrderMaintainer<V>(...args: V[]): LowerNode<V> {
  const root = new LowerNode(createUpperList(), 0, undefined, undefined, args[0])
  if (args.length < 1) return root
  for (let i = args.length - 1; i >= 0; --i) {
    root.insertAfter(args[i])
  }
  const res = root.next
  root.remove()
  return res!
}

export { createOrderMaintainer }

if (require.main === module) {
  const head = createOrderMaintainer(1, 0, -5, 10)

  const p = head.next!.next!

  console.log(p.value)

  console.log(head.compare(p))
  console.log(p.compare(head))
  console.log(p.compare(p))
}

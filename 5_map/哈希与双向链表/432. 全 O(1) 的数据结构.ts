// 双向链表 + map
// 链表找最大最小
// key到count count到Node 的映射

// node存当前值为count的所有key
class ListNode {
  pre!: ListNode
  next!: ListNode
  count: number
  // 记录该cnt(计数)下key包括哪些
  keySet: Set<string>

  constructor(count: number = 1) {
    this.count = count
    this.keySet = new Set()
  }
}

class AllOne {
  // 记录头尾 便于求最小值最大值
  private head: ListNode
  private tail: ListNode
  private keyToCount: Map<string, number>
  private countToNode: Map<number, ListNode>

  constructor() {
    this.head = new ListNode(-Infinity)
    this.tail = new ListNode(Infinity)
    this.head.next = this.tail
    this.tail.pre = this.head
    // key对应的值
    this.keyToCount = new Map()
    // 值对应的节点
    this.countToNode = new Map()
  }

  inc(key: string): void {
    if (this.keyToCount.has(key)) {
      this.changeKey(key, 1)
    } else {
      this.keyToCount.set(key, 1)
      // 说明没有计数为1的节点,在self.head后面加入
      if (this.head.next.count !== 1) this.appendNode(new ListNode(1), this.head)

      this.head.next.keySet.add(key)
      this.countToNode.set(1, this.head.next)
    }
  }

  dec(key: string): void {
    if (this.keyToCount.has(key)) {
      const count = this.keyToCount.get(key)!
      if (count === 1) {
        this.keyToCount.delete(key)
        this.removeKeyInNode(this.countToNode.get(count)!, key)
      } else {
        this.changeKey(key, -1)
      }
    }
  }

  getMaxKey(): string {
    console.dir(this.head, { depth: null })
    return this.tail.pre === this.head ? '' : this.tail.pre.keySet.keys().next().value
  }

  getMinKey(): string {
    return this.head.next === this.tail ? '' : this.head.next.keySet.keys().next().value
  }

  // key加1或者减1
  private changeKey(key: string, offset: 1 | -1) {
    const count = this.keyToCount.get(key)!
    const node = this.countToNode.get(count)!
    const newCount = count + offset
    this.keyToCount.set(key, newCount)

    let newNode: ListNode
    if (this.countToNode.has(newCount)) {
      newNode = this.countToNode.get(newCount)!
    } else {
      newNode = new ListNode(newCount)
      this.countToNode.set(newCount, newNode)
      this.appendNode(newNode, offset === 1 ? node : node.pre)
    }
    newNode.keySet.add(key)
    this.removeKeyInNode(node, key)
  }

  // 在链表删除node
  private removeNode(node: ListNode) {
    const preNode = node.pre
    const nextNode = node.next
    preNode.next = nextNode
    nextNode.pre = preNode
  }

  // 在node删除key
  private removeKeyInNode(node: ListNode, key: string) {
    node.keySet.delete(key)
    if (node.keySet.size === 0) {
      this.removeNode(node)
      this.countToNode.delete(node.count)
    }
  }

  // 在preNode后面加入node
  private appendNode(node: ListNode, preNode: ListNode) {
    const next = preNode.next
    preNode.next = node
    node.pre = preNode
    node.next = next
    next.pre = node
  }
}

const allOne = new AllOne()

allOne.inc('hello')
allOne.inc('hello')
console.log(allOne.getMaxKey())
console.log(allOne.getMinKey())
allOne.inc('leet')
console.log(allOne.getMaxKey())
console.log(allOne.getMinKey())

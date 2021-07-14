class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}
// 1.插入 insert
// 2.删除堆顶 pop
// 3.获取堆顶 peek
// 4.获取堆大小 size

class MinHeap {
  constructor(private heap: Node[], private volumn?: number) {
    this.heap = heap
    this.volumn = volumn
  }

  /**
   *
   * @param val 插入的值
   * @description 将值插入数组(堆)的尾部，然后上移直至父节点不超过它
   * @description 时间复杂度为`O(log(k))`
   */
  insert(val: Node) {
    if (this.volumn !== undefined && this.heap.length >= this.volumn) {
      this.pop()
    }

    this.heap.push(val)
    this.shiftUp(this.heap.length - 1)

    return this
  }

  /**
   * @description 用数组尾部元素替换堆顶(直接删除会破坏堆结构),然后下移动直至子节点都大于新堆顶
   * @description 时间复杂度为`O(log(k))`
   */
  pop() {
    const head = this.heap.pop()
    this.heap[0] = head!
    this.shiftDown(0)
    return head
  }

  peek() {
    return this.heap[0]
  }

  size() {
    return this.heap.length
  }

  // 上移
  private shiftUp(index: number) {
    if (index <= 0) return
    const parentIndex = this.getParentIndex(index)
    while (this.heap[parentIndex].value > this.heap[index].value) {
      this.swap(parentIndex, index)
      this.shiftUp(parentIndex)
    }
  }

  // 下移
  private shiftDown(index: number) {
    const leftChildIndex = this.getLeftChildIndex(index)
    const rightChildIndex = this.getRightChildIndex(index)
    if (this.heap[leftChildIndex].value < this.heap[index].value) {
      this.swap(leftChildIndex, index)
      this.shiftDown(leftChildIndex)
    }
    if (this.heap[rightChildIndex].value < this.heap[index].value) {
      this.swap(rightChildIndex, index)
      this.shiftDown(rightChildIndex)
    }
  }

  private getParentIndex(index: number) {
    // 二进制数向右移动一位，相当于Math.floor((index-1)/2)
    return (index - 1) >> 1
  }

  private getLeftChildIndex(index: number) {
    return index * 2 + 1
  }

  private getRightChildIndex(index: number) {
    return index * 2 + 2
  }

  private swap(parentIndex: number, index: number) {
    ;[this.heap[parentIndex], this.heap[index]] = [this.heap[index], this.heap[parentIndex]]
  }
}

// 新的链表的下一个节点是k个链表头中的最小节点
const mergeKLists = (lists: Node[]) => {
  if (lists.length === 0) return
  const newNode = new Node(0)
  let p = newNode
  const k = lists.length
  const minHeap = new MinHeap([], k)

  lists.forEach(node => {
    minHeap.insert(node)
  })

  while (minHeap.size()) {
    const head = minHeap.pop()
    p.next = head
    p = p.next!
    head?.next && minHeap.insert(head.next)
  }

  return newNode.next
}

const a = new Node(1)
const b = new Node(3)
const c = new Node(5)
a.next = b
b.next = c
const d = new Node(2)
const e = new Node(4)
const f = new Node(5)
d.next = e
e.next = f

console.log(mergeKLists([a, d]))
export {}

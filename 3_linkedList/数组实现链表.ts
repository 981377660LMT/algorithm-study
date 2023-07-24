// 数组实现链表

function linkedListDemo(cap: number): void {
  let head = -1 // 链表头
  let ptr = 0 // 当前用到了哪个结点
  const value = Array(cap) // 链表值
  const next = new Int32Array(cap) // 链表指针

  /**
   * 在链表头插入一个数x.
   */
  function insert(x: number): void {
    value[ptr] = x
    next[ptr] = head
    head = ptr
    ptr++
  }

  /**
   * 删除链表头.需要保证链表非空.
   */
  function remove(): void {
    head = next[head]
  }
}

function doubleLinkedList(cap: number): void {
  const value = Array(cap) // 链表值
  const prev = new Int32Array(cap) // 链表指针
  const next = new Int32Array(cap) // 链表指针
  // 0是左端点, 1是右端点
  next[0] = 1
  prev[1] = 0
  let ptr = 2 // 当前用到了哪个结点

  /**
   * 将结点x插入到node的右边.
   */
  function insert(node: number, x: number): void {
    value[ptr] = x
    prev[ptr] = node
    next[ptr] = next[node]
    prev[next[node]] = ptr
    next[node] = ptr
    ptr++
  }

  /**
   * 删除结点x.
   */
  function remove(x: number): void {
    next[prev[x]] = next[x]
    prev[next[x]] = prev[x]
  }
}

export {}

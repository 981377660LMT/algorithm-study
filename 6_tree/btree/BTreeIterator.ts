// https://github.com/tidwall/btree/blob/51838063d453a243f6c4df75c5278efba5c6ffc8/btreeg.go#L1124

interface IBTreeNode<T> {
  isLeaf(): boolean
  items(): T[]
  children(): IBTreeNode<T>[]
}

interface IBTree<N extends IBTreeNode<T>, T> {
  root(): N | undefined
  find(node: N, item: T, depth: number): { index: number; found: boolean }
}

class BTreeIterator<N extends IBTreeNode<T>, T> {
  private _tree: IBTree<N, T> | undefined

  private readonly _stack: { node: N; index: number }[] = []
  private _item: T | undefined = undefined
  private _seeked = false
  private _atStart = false
  private _atEnd = false

  constructor(tree: IBTree<N, T> | undefined) {
    this._tree = tree
  }

  first(): boolean {
    if (!this._tree) return false
    this._atEnd = false
    this._atStart = false
    this._seeked = true
    this._stack.length = 0
    const root = this._tree.root()
    if (!root) return false

    let node = root
    for (;;) {
      this._stack.push({ node, index: 0 })
      if (node.isLeaf()) {
        break
      }
      node = node.children()[0] as N
    }
    const s = this._stack[this._stack.length - 1]
    this._item = s.node.items()[s.index]
    return true
  }

  last(): boolean {
    if (!this._tree) return false
    this._seeked = true
    this._stack.length = 0
    const root = this._tree.root()
    if (!root) return false

    let node = root
    for (;;) {
      this._stack.push({ node, index: node.items().length })
      if (node.isLeaf()) {
        this._stack[this._stack.length - 1].index--
        break
      }
      node = node.children()[node.items().length] as N
    }
    const s = this._stack[this._stack.length - 1]
    this._item = s.node.items()[s.index]
    return true
  }

  seek(item: T): boolean {
    if (!this._tree) return false
    this._seeked = true
    this._stack.length = 0
    const root = this._tree.root()
    if (!root) return false

    let node = root
    let depth = 0
    for (;;) {
      const { index, found } = this._tree.find(node, item, depth)
      this._stack.push({ node, index })
      if (found) {
        this._item = node.items()[index]
        return true
      }
      if (node.isLeaf()) {
        this._stack[this._stack.length - 1].index--
        return this.next()
      }
      node = node.children()[index] as N
      depth++
    }
  }

  prev(): boolean {
    if (!this._tree) return false
    if (!this._seeked) return false
    if (this._stack.length === 0) {
      if (this._atEnd) {
        return this.last() && this.prev()
      }
      return false
    }

    let s = this._stack[this._stack.length - 1]
    if (s.node.isLeaf()) {
      s.index--
      if (s.index === -1) {
        for (;;) {
          this._stack.length--
          if (this._stack.length === 0) {
            this._atStart = true
            return false
          }
          s = this._stack[this._stack.length - 1]
          s.index--
          if (s.index > -1) {
            break
          }
        }
      }
    } else {
      let node = s.node.children()[s.index] as N
      for (;;) {
        this._stack.push({ node, index: node.items().length })
        if (node.isLeaf()) {
          this._stack[this._stack.length - 1].index--
          break
        }
        node = node.children()[node.items().length] as N
      }
    }

    s = this._stack[this._stack.length - 1]
    this._item = s.node.items()[s.index]
    return true
  }

  next(): boolean {
    if (!this._tree) return false
    if (!this._seeked) return this.first()
    if (this._stack.length === 0) {
      if (this._atStart) {
        return this.first() && this.next()
      }
      return false
    }

    let s = this._stack[this._stack.length - 1]
    s.index++
    if (s.node.isLeaf()) {
      if (s.index === s.node.items().length) {
        for (;;) {
          this._stack.length--
          if (this._stack.length === 0) {
            this._atEnd = true
            return false
          }
          s = this._stack[this._stack.length - 1]
          if (s.index < s.node.items().length) {
            break
          }
        }
      }
    } else {
      let node = s.node.children()[s.index] as N
      for (;;) {
        this._stack.push({ node, index: 0 })
        if (node.isLeaf()) {
          break
        }
        node = node.children()[0] as N
      }
    }

    s = this._stack[this._stack.length - 1]
    this._item = s.node.items()[s.index]
    return true
  }

  item(): T | undefined {
    return this._item
  }

  release(): void {
    if (!this._tree) return
    this._stack.length = 0
    this._tree = undefined
    this._item = undefined
  }
}

export type { IBTree, IBTreeNode }
export { BTreeIterator }

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  class Node implements IBTreeNode<number> {
    private readonly _items: number[]
    private readonly _children: Node[]

    constructor(items: number[], children: Node[]) {
      this._items = items
      this._children = children
    }

    isLeaf(): boolean {
      return this._children.length === 0
    }

    items(): number[] {
      return this._items
    }

    children(): Node[] {
      return this._children
    }
  }

  class Tree implements IBTree<Node, number> {
    private readonly _root: Node | undefined

    constructor(root: Node | undefined) {
      this._root = root
    }

    root(): Node | undefined {
      return this._root
    }

    find(node: Node, item: number): { index: number; found: boolean } {
      const items = node.items()
      let low = 0
      let high = items.length - 1
      while (low <= high) {
        const mid = low + ((high - low) >>> 1)
        if (items[mid] === item) return { index: mid, found: true }
        if (items[mid] < item) {
          low = mid + 1
        } else {
          high = mid - 1
        }
      }
      return { index: low, found: false }
    }
  }

  const root = new Node(
    [1, 3, 5, 7, 9],
    [
      new Node([0, 1], []),
      new Node([2, 3], []),
      new Node([4, 5], []),
      new Node([6, 7], []),
      new Node([8, 9], []),
      new Node([10, 11], [])
    ]
  )

  const tree = new Tree(root)
  const iter = new BTreeIterator(tree)

  // output items >= 4
  iter.seek(4)
  do {
    console.log(iter.item())
  } while (iter.next())
}

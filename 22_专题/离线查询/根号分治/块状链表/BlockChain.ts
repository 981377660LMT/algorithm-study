// https://www.acwing.com/blog/content/28060/
// TODO: 查询聚合值有问题，谨慎使用

/**
 * S: 聚合类型.
 * U: 更新类型.
 * E: 元素类型.
 * B: 子类块类型.
 */
abstract class AbstractBlock<S, U, V, B extends AbstractBlock<S, U, V, B>> {
  /** 分裂整块，block1包含前k个元素，block2包含后面的元素.*/
  abstract split(k: number): { block1: B; block2: B }
  abstract merge(other: B): B
  /** 在index位置之前插入元素e，0表示第一个.*/
  abstract insertBefore(index: number, e: V): void
  abstract delete(index: number): void
  abstract get(index: number): V
  abstract reverse(): void

  abstract fullyQuery(sum: S): void
  abstract partialQuery(index: number, sum: S): void
  abstract fullyUpdate(lazy: U): void
  abstract partialUpdate(index: number, lazy: U): void
  beforePartialQuery(): void {}
  afterPartialUpdate(): void {}
}

class BlockChain<S, U, V, B extends AbstractBlock<S, U, V, B>> {
  /**
   * @param blockSize 块大小.默认为`2 * (1 + Math.sqrt(n))`.
   */
  static create<S, U, V, B extends AbstractBlock<S, U, V, B>>(
    n: number,
    blockSupplier: (start: number, end: number) => B,
    blockSize = 2 * (1 + (Math.sqrt(n) | 0))
  ): BlockChain<S, U, V, B> {
    const res = new BlockChain<S, U, V, B>()
    res._b = blockSize
    res._size = n
    LinkedNode.link(res.head, res.tail)
    for (let start = 0; start < n; start += blockSize) {
      let end = start + blockSize
      if (end > n) end = n
      const block = blockSupplier(start, end)
      const node = new LinkedNode<B>()
      node.data = block
      node.size = end - start
      LinkedNode.link(res.tail.prev, node)
      LinkedNode.link(node, res.tail)
    }
    return res
  }

  private static new2<S, U, V, B extends AbstractBlock<S, U, V, B>>(
    b: number,
    supplier: () => B
  ): BlockChain<S, U, V, B> {
    const res = new BlockChain<S, U, V, B>()
    const block = supplier()
    const node = new LinkedNode<B>()
    node.data = block
    node.size = 0
    res._b = b
    LinkedNode.link(res.tail.prev, node)
    LinkedNode.link(node, res.tail)
    return res
  }

  private static new3<S, U, V, B extends AbstractBlock<S, U, V, B>>(
    b: number,
    begin: LinkedNode<B>,
    end: LinkedNode<B>
  ): BlockChain<S, U, V, B> {
    // add an empty node
    const res = new BlockChain<S, U, V, B>()
    res._b = b
    LinkedNode.link(res.head, begin)
    LinkedNode.link(end, res.tail)
    res._maintain()
    return res
  }

  head = new LinkedNode<B>()
  tail = new LinkedNode<B>()
  private _b = 0
  private _size = 0

  private constructor() {}

  get(index: number): V {
    if (index < 0) index += this._size
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      if (node.size <= index) {
        index -= node.size
        continue
      }
      node.data!.beforePartialQuery()
      return node.data!.get(index)
    }
    throw new Error('Index out of bounds')
  }

  prefixSize(block: B, include: boolean): number {
    let res = 0
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      if (node.data === block) {
        if (include) {
          res += node.size
        }
        break
      }
      res += node.size
    }
    return res
  }

  split(
    k: number,
    blockSupplier: () => B
  ): { first: BlockChain<S, U, V, B>; second: BlockChain<S, U, V, B> } {
    k++
    if (k === 1) {
      return { first: BlockChain.new2(this._b, blockSupplier), second: this }
    }
    if (k > this._size) {
      return { first: this, second: BlockChain.new2(this._b, blockSupplier) }
    }
    const head = this._splitKth(k)
    const end = this.tail.prev
    LinkedNode.link(head.prev, this.tail)
    const b = BlockChain.new3<S, U, V, B>(this._b, head, end)
    this._maintain()
    return { first: this, second: b }
  }

  mergeDestructively(other: BlockChain<S, U, V, B>): BlockChain<S, U, V, B> {
    LinkedNode.link(this.tail.prev, other.head.next)
    this.tail = other.tail
    this._size += other._size
    return this
  }

  insertBefore(index: number, e: V): void {
    if (index < 0) index += this._size
    if (index < 0) index = 0
    if (index > this._size) index = this._size
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      if (node.size < index) {
        index -= node.size
        continue
      }
      node.data!.insertBefore(index, e)
      node.size++
      break
    }
    this._maintain()
  }

  delete(index: number): void {
    if (index < 0) index += this._size
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      if (node.size <= index) {
        index -= node.size
        continue
      }
      node.data!.delete(index)
      node.size--
      break
    }
    this._maintain()
  }

  update(start: number, end: number, update: U): void {
    if (start < 0) start += this._size
    if (end > this._size) end = this._size
    if (start >= end) return
    end--
    let offset = 0
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      const left = offset
      const right = offset + node.size - 1
      offset += node.size
      if (this._enter(start, end, left, right)) {
        node.data!.fullyUpdate(update)
      } else if (this._leave(start, end, left, right)) {
        continue
      } else {
        for (let i = Math.max(left, start), to = Math.min(right, end); i <= to; i++) {
          node.data!.partialUpdate(i - left, update)
        }
        node.data!.afterPartialUpdate()
      }
    }
  }

  query(start: number, end: number, sum: S): void {
    if (start < 0) start += this._size
    if (end > this._size) end = this._size
    if (start >= end) return
    end--
    let offset = 0
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      const left = offset
      const right = offset + node.size - 1
      offset += node.size
      if (this._enter(start, end, left, right)) {
        node.data!.fullyQuery(sum)
      } else if (this._leave(start, end, left, right)) {
        continue
      } else {
        node.data!.beforePartialQuery()
        for (let i = Math.max(left, start), to = Math.min(right, end); i <= to; i++) {
          node.data!.partialQuery(i - left, sum)
        }
      }
    }
  }

  /** 向左旋转k个位置.*/
  rotateLeft(k: number): void {
    if (k < 0) k += this._size
    if (k >= this._size) k %= this._size
    if (k === 0) return
    k++
    const node = this._splitKth(k)
    const h1 = this.head.next
    const e1 = node.prev
    const h2 = node
    const e2 = this.tail.prev
    LinkedNode.link(this.head, h2)
    LinkedNode.link(e2, h1)
    LinkedNode.link(e1, this.tail)
    this._maintain()
  }

  reverse(start = 0, end = this._size): void {
    if (start >= end) return
    end--
    const left = this._splitKth(start + 1)
    const right = this._splitKth(end + 2).prev
    const begin = left.prev
    const endNode = right.next
    right.next = LinkedNode._NIL
    this._reverse(left, LinkedNode._NIL)
    LinkedNode.link(begin, right)
    LinkedNode.link(left, endNode)
    this._maintain()
  }

  getAll(): V[] {
    const res: V[] = []
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      for (let i = 0; i < node.size; i++) {
        res.push(node.data!.get(i))
      }
    }
    return res
  }

  enumerateNode(f: (node: LinkedNode<B>) => void): void {
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      f(node)
    }
  }

  get length(): number {
    return this._size
  }

  private _reverse(root: LinkedNode<B>, p: LinkedNode<B>): void {
    if (root === LinkedNode._NIL) {
      return
    }
    this._reverse(root.next, root)
    root.data!.reverse()
    root.prev = root.next
    root.next = p
  }

  private _split(node: LinkedNode<B>, k: number): void {
    const post = new LinkedNode<B>()
    const { block1, block2 } = node.data!.split(k)
    LinkedNode.link(post, node.next)
    LinkedNode.link(node, post)
    post.data = block2
    post.size = node.size - k
    node.data = block1
    node.size = k
  }

  private _mergeNode(a: LinkedNode<B>, b: LinkedNode<B>): void {
    LinkedNode.link(a, b.next)
    a.data = a.data!.merge(b.data!)
    a.size += b.size
  }

  private _maintain(): void {
    this._size = 0
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      this._size += node.size
      if (node.size >= 2 * this._b) {
        this._split(node, this._b)
      } else if (node.prev !== this.head && node.size + node.prev.size <= this._b) {
        this._mergeNode(node.prev, node)
      }
    }
  }

  /**
   * 拆分成前k个元素和后面的元素，返回后面的元素.
   */
  private _splitKth(k: number): LinkedNode<B> {
    for (let node = this.head.next; node !== this.tail; node = node.next) {
      if (node.size < k) {
        k -= node.size
        continue
      }
      if (k !== 1) {
        this._split(node, k - 1)
        node = node.next
      }
      return node
    }
    return this.tail
  }

  private _enter(L: number, R: number, l: number, r: number): boolean {
    return L <= l && r <= R
  }

  private _leave(L: number, R: number, l: number, r: number): boolean {
    return l > R || r < L
  }
}

class LinkedNode<E> {
  static readonly _NIL = new LinkedNode<any>()

  static link<E>(a: LinkedNode<E>, b: LinkedNode<E>): void {
    b.prev = a
    a.next = b
  }

  prev: LinkedNode<E> = LinkedNode._NIL
  next: LinkedNode<E> = LinkedNode._NIL
  data: E | undefined = undefined
  size = 0
}

export { BlockChain, AbstractBlock }

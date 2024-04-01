import { ArrayDeque } from '../../../2_queue/Deque/ArrayDeque'

// 左偏树
// https://taodaling.github.io/blog/2019/06/28/%E4%BC%98%E5%85%88%E9%98%9F%E5%88%97/
class PersistenteLeftistTree<K> {
  private static _NIL = new PersistenteLeftistTree<any>(undefined)

  static {
    this._NIL._left = this._NIL
    this._NIL._right = this._NIL
    this._NIL._dist = -1
  }

  private _left: PersistenteLeftistTree<K> = PersistenteLeftistTree._NIL
  private _right: PersistenteLeftistTree<K> = PersistenteLeftistTree._NIL
  private _dist = 0
  private readonly _key: K

  constructor(key: K) {
    this._key = key
  }

  static createFromIterable<K>(
    iterable: Iterable<PersistenteLeftistTree<K>>,
    comparator: (a: K, b: K) => number
  ): PersistenteLeftistTree<K> {
    return this.createFromDeque(new ArrayDeque(iterable), comparator)
  }

  static createFromDeque<K>(
    deque: ArrayDeque<PersistenteLeftistTree<K>>,
    comparator: (a: K, b: K) => number
  ): PersistenteLeftistTree<K> {
    while (deque.length > 1) {
      deque.push(this.merge(deque.shift()!, deque.shift()!, comparator))
    }
    return deque.shift()!
  }

  static merge<K>(
    a: PersistenteLeftistTree<K>,
    b: PersistenteLeftistTree<K>,
    comparator: (a: K, b: K) => number
  ): PersistenteLeftistTree<K> {
    if (a === this._NIL) return b
    if (b === this._NIL) return a
    if (comparator(a._key, b._key) > 0) {
      const tmp = a
      a = b
      b = tmp
    }
    a = a.clone()
    a._right = this.merge(a._right, b, comparator)
    if (a._left._dist < a._right._dist) {
      const tmp = a._left
      a._left = a._right
      a._right = tmp
    }
    a._dist = a._right._dist + 1
    return a
  }

  static pop<K>(
    root: PersistenteLeftistTree<K>,
    comparator: (a: K, b: K) => number
  ): PersistenteLeftistTree<K> {
    return this.merge(root._left, root._right, comparator)
  }

  static asIterator<V>(
    heap: PersistenteLeftistTree<V>,
    comparator: (a: V, b: V) => number
  ): Iterator<V> {
    return new PersistentLeftistTreeIteratorAdapter(heap, comparator)
  }

  clone(): PersistenteLeftistTree<K> {
    const clone = new PersistenteLeftistTree(this._key)
    clone._left = this._left
    clone._right = this._right
    clone._dist = this._dist
    return clone
  }

  isEmpty(): boolean {
    return this === PersistenteLeftistTree._NIL
  }

  peek(): K {
    return this._key
  }

  toString(): string {
    const builder: string[] = []
    this._toStringDfs(builder)
    return JSON.stringify(builder)
  }

  private _toStringDfs(builder: string[]): void {
    if (this === PersistenteLeftistTree._NIL) return
    builder.push(JSON.stringify(this._key))
    builder.push(' ')
    this._left._toStringDfs(builder)
    this._right._toStringDfs(builder)
  }
}

class PersistentLeftistTreeIteratorAdapter<V> implements Iterator<V> {
  constructor(
    private tree: PersistenteLeftistTree<V>,
    private comparator: (a: V, b: V) => number
  ) {}

  next(): IteratorResult<V> {
    if (this.tree.isEmpty()) {
      return { done: true, value: undefined }
    }
    const res = this.tree.peek()
    this.tree = PersistenteLeftistTree.pop(this.tree, this.comparator)
    return { done: false, value: res }
  }
}

export {}

if (require.main === module) {
  const iter = PersistenteLeftistTree.asIterator(new PersistenteLeftistTree(1), (a, b) => a - b)
  console.log(iter.next())
  console.log(iter.next())
}

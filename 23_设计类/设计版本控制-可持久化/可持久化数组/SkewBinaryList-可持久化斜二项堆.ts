/* eslint-disable @typescript-eslint/no-non-null-assertion */
// SkewBinaryList(斜二项堆)
// 纯函数式堆（纯函数式优先级队列）
// https://scrapbox.io/data-structures/Skew_Binary_List
// https://noshi91.github.io/Library/data_structure/persistent_skew_binary_random_access_list.cpp

// API:
//  get(index): 访问下标.
//  front(): 访问第一个元素.
//  set(index, x): 更新下标的元素, 返回新的List.
//  appendLeft(x): 在左边添加元素, 返回新的List.
//  popLeft(): 弹出第一个元素, 返回新的List.
//  empty(): 判断是否为空.

// !比golang慢很多(1000ms vs 100ms), 因为js处理指针很慢???
// 感觉js遍历和操作数组非常快,但是使用对象就非常慢.

/**
 * 可持久化斜二项堆.
 */
class PersistentSkewBinaryRandomAccessList<T> {
  private readonly _root: _PDigit<T> | undefined = undefined
  constructor(root?: _PDigit<T>) {
    this._root = root
  }

  empty(): boolean {
    return !this._root
  }

  get(index: number): T | undefined {
    return this._root ? this._root.lookup(index) : undefined
  }

  front(): T | undefined {
    return this._root ? this._root.tree.value : undefined
  }

  set(index: number, x: T): PersistentSkewBinaryRandomAccessList<T> {
    if (!this._root) throw new Error('root is undefined')
    return new PersistentSkewBinaryRandomAccessList(this._root.update(index, x))
  }

  appendLeft(x: T): PersistentSkewBinaryRandomAccessList<T> {
    if (this._root && this._root.next && this._root.size === this._root.next.size) {
      return new PersistentSkewBinaryRandomAccessList(
        new _PDigit(
          1 + this._root.size + this._root.next.size,
          new _PTree(x, this._root.tree, this._root.next.tree),
          this._root.next.next
        )
      )
    }
    return new PersistentSkewBinaryRandomAccessList(new _PDigit(1, new _PTree(x), this._root))
  }

  popLeft(): PersistentSkewBinaryRandomAccessList<T> {
    if (!this._root) throw new Error('empty list')
    if (this._root.size === 1) {
      return new PersistentSkewBinaryRandomAccessList(this._root.next)
    }
    const chSize = this._root.size >> 1
    return new PersistentSkewBinaryRandomAccessList(
      new _PDigit(chSize, this._root.tree, new _PDigit(chSize, this._root.tree, this._root.next))
    )
  }

  toString(): string {
    const sb: string[] = []
    let i = 0
    try {
      while (true) {
        sb.push(this.get(i)!.toString())
        i++
      }
      // eslint-disable-next-line no-empty
    } catch (_) {}

    return `List{${sb.join(', ')}}`
  }
}

class _PTree<T> {
  value: T
  left: _PTree<T> | undefined
  right: _PTree<T> | undefined
  constructor(value: T, left?: _PTree<T>, right?: _PTree<T>) {
    this.value = value
    this.left = left
    this.right = right
  }

  lookup(size: number, index: number): T {
    if (!index) {
      return this.value
    }
    const remIndex = index - 1
    const chSize = size >> 1
    if (remIndex < chSize) {
      return this.left!.lookup(chSize, remIndex)
    }
    return this.right!.lookup(chSize, remIndex - chSize)
  }

  update(size: number, index: number, x: T): _PTree<T> {
    if (!index) {
      return new _PTree(x, this.left, this.right)
    }
    const remIndex = index - 1
    const chSize = size >> 1
    if (remIndex < chSize) {
      return new _PTree(this.value, this.left!.update(chSize, remIndex, x), this.right)
    }
    return new _PTree(this.value, this.left, this.right!.update(chSize, remIndex - chSize, x))
  }
}

class _PDigit<T> {
  size: number
  tree: _PTree<T>
  next: _PDigit<T> | undefined
  constructor(size: number, tree: _PTree<T>, next: _PDigit<T> | undefined) {
    this.size = size
    this.tree = tree
    this.next = next
  }

  lookup(index: number): T | undefined {
    if (index < this.size) {
      return this.tree.lookup(this.size, index)
    }
    return this.next!.lookup(index - this.size)
  }

  update(index: number, x: T): _PDigit<T> {
    if (index < this.size) {
      return new _PDigit(this.size, this.tree.update(this.size, index, x), this.next)
    }
    return new _PDigit(this.size, this.tree, this.next!.update(index - this.size, x))
  }
}

export { PersistentSkewBinaryRandomAccessList }

if (require.main === module) {
  let list = new PersistentSkewBinaryRandomAccessList<number>()
  console.log(list.empty())
  console.log(list.get(0))
  console.log(list.front())
  list = list.appendLeft(1)
  console.log(list.empty())
  console.log(list.get(0))
  console.log(list.front())
  list = list.appendLeft(2)
  console.log(list.empty())
  console.log(list.get(1))
  list = list.set(1, 10)
  console.log(list.toString())
  list = list.popLeft()
  console.log(list.toString())

  console.time('append')
  for (let i = 0; i < 1e5; i++) {
    list = list.appendLeft(i)
    list = list.set(i, i)
  }
  for (let i = 0; i < 1e5; i++) {
    list = list.popLeft()
  }
  console.timeEnd('append')
}

interface NestedInteger<T> {
  isInteger: () => boolean
  getInteger: () => T
  setInteger: (value: T) => void
  add: (elem: NestedInteger<T>) => void
  getList: () => NestedInteger<T>[]
}

/**
 * 直接遍历迭代器所有元素，然后保存到数组中.
 * 不太好的方法.
 */
class NestedIterator1<T> {
  private readonly _data: NestedInteger<T>[] = []
  private _ptr = 0

  constructor(nestedList: NestedInteger<T>[]) {
    this._flatten(nestedList)
  }

  next(): T {
    return this._data[this._ptr++].getInteger()!
  }

  hasNext(): boolean {
    return this._ptr < this._data.length
  }

  // 核心
  private _flatten(nestedList: NestedInteger<T>[]) {
    nestedList.forEach(v => {
      if (v.isInteger()) {
        this._data.push(v)
      } else {
        this._flatten(v.getList())
      }
    })
  }
}

/**
 * 生成器函数延迟flatten.
 */
class NestedIterator<T> {
  private readonly _gen: Generator<T | undefined>
  private _nextItem: T | undefined // 预先取出下一个元素,便于判断hasNext.

  constructor(nestedList: NestedInteger<T>[]) {
    this._gen = gen(nestedList)
    this._nextItem = this._gen.next().value
  }

  hasNext(): boolean {
    // eslint-disable-next-line eqeqeq
    return this._nextItem != undefined
  }

  next(): T | undefined {
    const res = this._nextItem
    this._nextItem = this._gen.next().value
    return res
  }
}

/**
 * lazyFlatten/flattenLazy.
 */
function* gen<T>(nestedList: NestedInteger<T>[]): Generator<T> {
  for (let i = 0; i < nestedList.length; i++) {
    const item = nestedList[i]
    if (item.isInteger()) {
      yield item.getInteger()
    } else {
      yield* gen(item.getList())
    }
  }
}

/**
 * 先将所有的 NestedInteger 逆序放到栈中，当需要展开的时候才进行展开。
 */
class NestedIterator3<T> {
  private readonly _stack: NestedInteger<T>[] = []

  constructor(nestedList: NestedInteger<T>[]) {
    for (let i = nestedList.length - 1; ~i; i--) {
      this._stack.push(nestedList[i])
    }
  }

  next(): T | undefined {
    return this.hasNext() ? this._stack.pop()!.getInteger() : undefined
  }

  hasNext(): boolean {
    if (this._stack.length === 0) return false
    const item = this._stack[this._stack.length - 1]
    if (item.isInteger()) return true
    const last = this._stack.pop()!
    const list = last.getList()
    for (let i = list.length - 1; ~i; i--) {
      this._stack.push(list[i])
    }
    return this.hasNext()
  }
}

export {}

// 输入：nestedList = [[1,1],2,[1,1]]
// 输出：[1,1,2,1,1]
// 解释：通过重复调用 next 直到 hasNext 返回 false，next 返回的元素的顺序应该是: [1,1,2,1,1]。

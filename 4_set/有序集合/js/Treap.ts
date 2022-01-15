import assert from 'assert'

type CompareFunction<T, R extends 'number' | 'boolean'> = (
  a: T,
  b: T
) => R extends 'number' ? number : boolean

// 无法支持islice 因为Treap结点随机旋转
interface ITreapMultiSet<T> extends Iterable<T> {
  add: (value: T) => this
  has: (value: T) => boolean
  delete: (value: T) => void

  bisectLeft: (value: T) => number
  bisectRight: (value: T) => number

  indexOf: (value: T) => number
  lastIndexOf: (value: T) => number

  at: (index: number) => T | undefined
  first: () => T | undefined
  last: () => T | undefined

  lower: (value: T) => T | undefined
  higher: (value: T) => T | undefined
  floor: (value: T) => T | undefined
  ceil: (value: T) => T | undefined

  shift: () => T | undefined
  pop: (index?: number) => T | undefined

  count(value: T): number

  keys: () => Generator<T, any, any>
  values: () => Generator<T, any, any>
  rvalues: () => Generator<T, any, any>
  entries(): IterableIterator<[number, T]>

  readonly size: number
}

class TreapNode<T = number> {
  value: T
  count: number
  size: number
  priority: number
  left: TreapNode<T> | null
  right: TreapNode<T> | null

  constructor(value: T) {
    this.value = value
    this.count = 1
    this.size = 1
    this.priority = Math.random()
    this.left = null
    this.right = null
  }

  static getSize(node: TreapNode<any> | null): number {
    return node?.size ?? 0
  }

  static getFac(node: TreapNode<any> | null): number {
    return node?.priority ?? 0
  }

  pushUp(): void {
    let tmp = this.count
    tmp += TreapNode.getSize(this.left)
    tmp += TreapNode.getSize(this.right)
    this.size = tmp
  }

  rotateRight(): TreapNode<T> {
    let node: TreapNode<T> = this
    const left = node.left
    node.left = left?.right ?? null
    left && (left.right = node)
    left && (node = left)
    node.right?.pushUp()
    node.pushUp()
    return node
  }

  rotateLeft(): TreapNode<T> {
    let node: TreapNode<T> = this
    const right = node.right
    node.right = right?.left ?? null
    right && (right.left = node)
    right && (node = right)
    node.left?.pushUp()
    node.pushUp()
    return node
  }
}

class TreapMultiSet<T = number> implements ITreapMultiSet<T> {
  private root: TreapNode<T>
  private compareFn: CompareFunction<T, 'number'>
  private leftBound: T
  private rightBound: T

  /**
   *
   * @param compareFn A compare function which returns boolean or number
   * @param leftBound defalut value is `-Infinity`
   * @param rightBound defalut value is `Infinity`
   * @description
   * create a `MultiSet`, compare elements using `compareFn`, which is increasing order by default.
   * @example
   * ```ts
   * interface Person {
      name: string
      age: number
    }

    const leftBound = {
      name: 'Alice',
      age: -Infinity,
    }

    const rightBound = {
      name: 'Bob',
      age: Infinity,
    }

    const sortedList = new TreapMultiSet<Person>(
      (a: Person, b: Person) => a.age - b.age,
      leftBound,
      rightBound
    )
   * ```
   */
  constructor(compareFn?: CompareFunction<T, 'number'>)
  constructor(compareFn: CompareFunction<T, 'number'>, leftBound: T, rightBound: T)
  constructor(
    compareFn: CompareFunction<T, any> = (a: any, b: any) => a - b,
    leftBound: any = -Infinity,
    rightBound: any = Infinity
  ) {
    this.root = new TreapNode<T>(rightBound)
    this.root.priority = Infinity
    this.root.left = new TreapNode<T>(leftBound)
    this.root.left!.priority = -Infinity
    this.root.pushUp()

    this.leftBound = leftBound
    this.rightBound = rightBound
    this.compareFn = compareFn
  }

  get size(): number {
    return this.root.size - 2
  }

  get height(): number {
    const getHeight = (node: TreapNode<T> | null): number => {
      if (node == null) return 0
      return 1 + Math.max(getHeight(node.left), getHeight(node.right))
    }

    return getHeight(this.root)
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Returns true if value is a member.
   */
  has(value: T): boolean {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): boolean => {
      if (node == null) return false
      if (compare(node.value, value) === 0) return true
      if (compare(node.value, value) < 0) return dfs(node.right, value)
      return dfs(node.left, value)
    }

    return dfs(this.root, value)
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Add value to sorted set.
   */
  add(value: T): this {
    const compare = this.compareFn
    const dfs = (
      node: TreapNode<T> | null,
      value: T,
      parent: TreapNode<T>,
      direction: 'left' | 'right'
    ): void => {
      if (node == null) return
      if (compare(node.value, value) === 0) {
        node.count++
        node.pushUp()
      } else if (compare(node.value, value) > 0) {
        if (node.left) {
          dfs(node.left, value, node, 'left')
        } else {
          node.left = new TreapNode(value)
          node.pushUp()
        }

        if (TreapNode.getFac(node.left) > node.priority) {
          parent[direction] = node.rotateRight()
        }
      } else if (compare(node.value, value) < 0) {
        if (node.right) {
          dfs(node.right, value, node, 'right')
        } else {
          node.right = new TreapNode(value)
          node.pushUp()
        }

        if (TreapNode.getFac(node.right) > node.priority) {
          parent[direction] = node.rotateLeft()
        }
      }
      parent.pushUp()
    }

    dfs(this.root.left, value, this.root, 'left')
    return this
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Remove value from sorted set if it is a member.
   * If value is not a member, do nothing.
   */
  delete(value: T): void {
    const compare = this.compareFn
    const dfs = (
      node: TreapNode<T> | null,
      value: T,
      parent: TreapNode<T>,
      direction: 'left' | 'right'
    ): void => {
      if (node == null) return

      if (compare(node.value, value) === 0) {
        if (node.count > 1) {
          node.count--
          node?.pushUp()
        } else if (node.left == null && node.right == null) {
          parent[direction] = null
        } else {
          // 旋到根节点
          if (node.right == null || TreapNode.getFac(node.left) > TreapNode.getFac(node.right)) {
            parent[direction] = node.rotateRight()
            dfs(parent[direction]?.right ?? null, value, parent[direction]!, 'right')
          } else {
            parent[direction] = node.rotateLeft()
            dfs(parent[direction]?.left ?? null, value, parent[direction]!, 'left')
          }
        }
      } else if (compare(node.value, value) > 0) {
        dfs(node.left, value, node, 'left')
      } else if (compare(node.value, value) < 0) {
        dfs(node.right, value, node, 'right')
      }

      parent?.pushUp()
    }

    dfs(this.root.left, value, this.root, 'left')
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Returns an index to insert value in the sorted set.
   * If the value is already present, the insertion point will be before (to the left of) any existing values.
   */
  bisectLeft(value: T): number {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0

      if (compare(node.value, value) === 0) {
        return TreapNode.getSize(node.left)
      } else if (compare(node.value, value) > 0) {
        return dfs(node.left, value)
      } else if (compare(node.value, value) < 0) {
        return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个lowerBound 所以-1
    return dfs(this.root, value) - 1
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Returns an index to insert value in the sorted set.
   * If the value is already present, the insertion point will be before (to the right of) any existing values.
   */
  bisectRight(value: T): number {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0

      if (compare(node.value, value) === 0) {
        return TreapNode.getSize(node.left) + node.count
      } else if (compare(node.value, value) > 0) {
        return dfs(node.left, value)
      } else if (compare(node.value, value) < 0) {
        return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个lowerBound 所以-1
    return dfs(this.root, value) - 1
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Returns the index of the first occurrence of a value in the set, or -1 if it is not present.
   */
  indexOf(value: T): number {
    const compare = this.compareFn
    let isExist = false

    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0

      if (compare(node.value, value) === 0) {
        isExist = true
        return TreapNode.getSize(node.left)
      } else if (compare(node.value, value) > 0) {
        return dfs(node.left, value)
      } else if (compare(node.value, value) < 0) {
        return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个lowerBound 所以-1
    const res = dfs(this.root, value) - 1
    return isExist ? res : -1
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Returns the index of the last occurrence of a value in the set, or -1 if it is not present.
   */
  lastIndexOf(value: T): number {
    const compare = this.compareFn
    let isExist = false

    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0

      if (compare(node.value, value) === 0) {
        isExist = true
        return TreapNode.getSize(node.left) + node.count - 1
      } else if (compare(node.value, value) > 0) {
        return dfs(node.left, value)
      } else if (compare(node.value, value) < 0) {
        return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个lowerBound 所以-1
    const res = dfs(this.root, value) - 1
    return isExist ? res : -1
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Returns the item located at the specified index.
   * @param index The zero-based index of the desired code unit. A negative index will count back from the last item.
   */
  at(index: number): T | undefined {
    if (index < 0) index += this.size
    if (index < 0 || index >= this.size) return undefined

    const dfs = (node: TreapNode<T> | null, rank: number): T | undefined => {
      if (node == null) return undefined

      if (TreapNode.getSize(node.left) >= rank) {
        return dfs(node.left, rank)
      } else if (TreapNode.getSize(node.left) + node.count >= rank) {
        return node.value
      } else {
        return dfs(node.right, rank - TreapNode.getSize(node.left) - node.count)
      }
    }

    // 因为有个-Infinity 所以 + 2
    const res = dfs(this.root, index + 2)
    return ([this.leftBound, this.rightBound] as any[]).includes(res) ? undefined : res
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Find and return the element less than `val`, return `undefined` if no such element found.
   */
  lower(value: T): T | undefined {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.value, value) >= 0) return dfs(node.left, value)

      const tmp = dfs(node.right, value)
      if (tmp == null || compare(node.value, tmp) > 0) {
        return node.value
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, value) as any
    return res === this.leftBound ? undefined : res
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Find and return the element greater than `val`, return `undefined` if no such element found.
   */
  higher(value: T): T | undefined {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.value, value) <= 0) return dfs(node.right, value)

      const tmp = dfs(node.left, value)

      if (tmp == null || compare(node.value, tmp) < 0) {
        return node.value
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, value) as any
    return res === this.rightBound ? undefined : res
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Find and return the element less than or equal to `val`, return `undefined` if no such element found.
   */
  floor(value: T): T | undefined {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.value, value) === 0) return node.value
      if (compare(node.value, value) >= 0) return dfs(node.left, value)

      const tmp = dfs(node.right, value)
      if (tmp == null || compare(node.value, tmp) > 0) {
        return node.value
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, value) as any
    return res === this.leftBound ? undefined : res
  }

  /**
   *
   * @complexity `O(logn)`
   * @description Find and return the element greater than or equal to `val`, return `undefined` if no such element found.
   */
  ceil(value: T): T | undefined {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.value, value) === 0) return node.value
      if (compare(node.value, value) <= 0) return dfs(node.right, value)

      const tmp = dfs(node.left, value)

      if (tmp == null || compare(node.value, tmp) < 0) {
        return node.value
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, value) as any
    return res === this.rightBound ? undefined : res
  }

  /**
   * @complexity `O(logn)`
   * @description
   * Returns the last element from set.
   * If the set is empty, undefined is returned.
   */
  first(): T | undefined {
    const iter = this.inOrder()
    iter.next()
    const res = iter.next().value
    return res === this.rightBound ? undefined : res
  }

  /**
   * @complexity `O(logn)`
   * @description
   * Returns the last element from set.
   * If the set is empty, undefined is returned .
   */
  last(): T | undefined {
    const iter = this.reverseInOrder()
    iter.next()
    const res = iter.next().value
    return res === this.leftBound ? undefined : res
  }

  /**
   * @complexity `O(logn)`
   * @description
   * Removes the first element from an set and returns it.
   * If the set is empty, undefined is returned and the set is not modified.
   */
  shift(): T | undefined {
    const first = this.first()
    if (first == undefined) return undefined
    this.delete(first)
    return first
  }

  /**
   * @complexity `O(logn)`
   * @description
   * Removes the last element from an set and returns it.
   * If the set is empty, undefined is returned and the set is not modified.
   */
  pop(index?: number): T | undefined {
    if (index == null) {
      const last = this.last()
      if (last == undefined) return undefined
      this.delete(last)
      return last
    }

    const toDelete = this.at(index)
    if (toDelete == null) return
    this.delete(toDelete)
  }

  /**
   *
   * @complexity `O(logn)`
   * @description
   * Returns number of occurrences of value in the sorted set.
   */
  count(value: T): number {
    const compare = this.compareFn
    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0
      if (compare(node.value, value) === 0) return node.count
      if (compare(node.value, value) < 0) return dfs(node.right, value)
      return dfs(node.left, value)
    }

    return dfs(this.root, value)
  }

  *[Symbol.iterator](): Generator<T, any, any> {
    yield* this.values()
  }

  /**
   * @description
   * Returns an iterable of keys in the set.
   */
  *keys(): Generator<T, any, any> {
    yield* this.values()
  }

  /**
   * @description
   * Returns an iterable of values in the set.
   */
  *values(): Generator<T, any, any> {
    const iter = this.inOrder()
    iter.next()
    let steps = this.size
    for (let _ = 0; _ < steps; _++) {
      yield iter.next().value
    }
  }

  /**
   * @description
   * Returns a generator for reversed order traversing the set.
   */
  *rvalues(): Generator<T, any, any> {
    const iter = this.reverseInOrder()
    iter.next()
    let steps = this.size
    for (let _ = 0; _ < steps; _++) {
      yield iter.next().value
    }
  }

  /**
   * @description
   * Returns an iterable of key, value pairs for every entry in the set.
   */
  *entries(): IterableIterator<[number, T]> {
    const iter = this.inOrder()
    iter.next()
    let steps = this.size
    for (let i = 0; i < steps; i++) {
      yield [i, iter.next().value]
    }
  }

  private *inOrder(root: TreapNode<T> | null = this.root): Generator<T, any, any> {
    if (root == null) return
    yield* this.inOrder(root.left)
    const count = root.count
    for (let _ = 0; _ < count; _++) {
      yield root.value
    }
    yield* this.inOrder(root.right)
  }

  private *reverseInOrder(root: TreapNode<T> | null = this.root): Generator<T, any, any> {
    if (root == null) return
    yield* this.reverseInOrder(root.right)
    const count = root.count
    for (let _ = 0; _ < count; _++) {
      yield root.value
    }
    yield* this.reverseInOrder(root.left)
  }
}

if (require.main === module) {
  const treap = new TreapMultiSet()

  // has add delete
  treap.add(1)
  assert.strictEqual(treap.has(1), true)
  assert.strictEqual(treap.has(2), false)
  treap.delete(1)
  assert.strictEqual(treap.size, 0)
  treap.add(1)
  treap.add(2)
  treap.add(3)
  treap.add(3)
  treap.add(3)
  assert.strictEqual(treap.size, 5)
  treap.delete(3)
  assert.strictEqual(treap.size, 4)
  // console.dir(treap, { depth: null })

  // upper lower ceil floor
  treap.add(3)
  assert.strictEqual(treap.higher(2), 3)
  assert.strictEqual(treap.higher(1.9), 2)
  assert.strictEqual(treap.higher(3), undefined)
  assert.strictEqual(treap.ceil(2), 2)
  assert.strictEqual(treap.ceil(1.9), 2)
  assert.strictEqual(treap.ceil(3.1), undefined)
  assert.strictEqual(treap.lower(2), 1)
  assert.strictEqual(treap.lower(0), undefined)
  assert.strictEqual(treap.lower(3.1), 3)
  assert.strictEqual(treap.floor(0.9), undefined)
  assert.strictEqual(treap.floor(1), 1)
  assert.strictEqual(treap.floor(3), 3)

  // indexOf lastIndexOf
  // [1,2,3,3,3,4]
  treap.add(4)
  assert.strictEqual(treap.indexOf(1), 0)
  assert.strictEqual(treap.indexOf(1.2), -1)
  assert.strictEqual(treap.indexOf(2), 1)
  assert.strictEqual(treap.indexOf(2.2), -1)
  assert.strictEqual(treap.indexOf(3), 2)
  assert.strictEqual(treap.indexOf(4), 5)
  assert.strictEqual(treap.indexOf(5), -1)
  assert.strictEqual(treap.lastIndexOf(1), 0)
  assert.strictEqual(treap.lastIndexOf(1.2), -1)
  assert.strictEqual(treap.lastIndexOf(2), 1)
  assert.strictEqual(treap.lastIndexOf(3), 4)
  assert.strictEqual(treap.lastIndexOf(4), 5)
  assert.strictEqual(treap.lastIndexOf(5), -1)
  treap.delete(4)

  // keys values rvalues entries
  assert.deepStrictEqual([...treap.keys()], [1, 2, 3, 3, 3])
  assert.deepStrictEqual([...treap.values()], [1, 2, 3, 3, 3])
  assert.deepStrictEqual([...treap.rvalues()], [3, 3, 3, 2, 1])
  assert.deepStrictEqual(
    [...treap.entries()],
    [
      [0, 1],
      [1, 2],
      [2, 3],
      [3, 3],
      [4, 3],
    ]
  )

  // at
  assert.strictEqual(treap.at(0), 1)
  assert.strictEqual(treap.at(1), 2)
  assert.strictEqual(treap.at(2), 3)
  assert.strictEqual(treap.at(3), 3)
  assert.strictEqual(treap.at(4), 3)
  assert.strictEqual(treap.at(-1), 3)
  assert.strictEqual(treap.at(-100), undefined)

  // first last shift pop
  assert.strictEqual(treap.shift(), 1)
  assert.strictEqual(treap.first(), 2)
  assert.strictEqual(treap.pop(), 3)
  assert.strictEqual(treap.last(), 3)
  assert.strictEqual(treap.size, 3)

  // bisectLeft bisectRight
  // [2,3,3]
  assert.strictEqual(treap.bisectLeft(1.9), 0)
  assert.strictEqual(treap.bisectLeft(2), 0)
  assert.strictEqual(treap.bisectLeft(2.5), 1)
  assert.strictEqual(treap.bisectLeft(3), 1)
  assert.strictEqual(treap.bisectLeft(4), 3)
  assert.strictEqual(treap.bisectRight(1.9), 0)
  assert.strictEqual(treap.bisectRight(2), 1)
  assert.strictEqual(treap.bisectRight(2.5), 1)
  assert.strictEqual(treap.bisectRight(3), 3)
  assert.strictEqual(treap.bisectRight(4), 3)

  // count
  assert.strictEqual(treap.count(1), 0)
  assert.strictEqual(treap.count(2), 1)
  assert.strictEqual(treap.count(3), 2)
}

export { TreapMultiSet }

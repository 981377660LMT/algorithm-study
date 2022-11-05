/* eslint-disable generator-star-spacing */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable prefer-destructuring */
// !有选treap

type Comparator<T, R extends 'number' | 'boolean'> = (
  a: T,
  b: T
) => R extends 'number' ? number : boolean

interface ITreapMultiSet<T> extends Iterable<T> {
  add: (...value: T[]) => this
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

  count: (value: T) => number

  keys: () => IterableIterator<T>
  values: () => IterableIterator<T>
  rvalues: () => IterableIterator<T>
  entries: () => IterableIterator<[number, T]>

  readonly size: number
}

class TreapNode<T = number> {
  static getSize(node: TreapNode<unknown> | null): number {
    return node ? node.size : 0
  }

  static getFac(node: TreapNode<unknown> | null): number {
    return node ? node.priority : 0
  }

  value: T
  priority: number
  count = 1
  size = 1
  left: TreapNode<T> | null = null
  right: TreapNode<T> | null = null

  constructor(value: T, priority: number) {
    this.value = value
    this.priority = priority
  }

  pushUp(): void {
    this.size = this.count + TreapNode.getSize(this.left) + TreapNode.getSize(this.right)
  }

  rotateRight(): TreapNode<T> {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    let node: TreapNode<T> = this
    const left = node.left
    node.left = left ? left.right : null
    left && (left.right = node)
    left && (node = left)
    node.right && node.right.pushUp()
    node.pushUp()
    return node
  }

  rotateLeft(): TreapNode<T> {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    let node: TreapNode<T> = this
    const right = node.right
    node.right = right ? right.left : null
    right && (right.left = node)
    right && (node = right)
    node.left && node.left.pushUp()
    node.pushUp()
    return node
  }
}

class TreapMultiSet<T = number> implements ITreapMultiSet<T> {
  private readonly _root: TreapNode<T>
  private readonly _compareFn: Comparator<T, 'number'>
  private readonly _leftBound: T
  private readonly _rightBound: T
  private _seed = (Math.floor(Date.now() / 2) + 1) >>> 0

  /**
   * create a `MultiSet`, compare elements using `compareFn`, which is increasing order by default.
   * @param compareFn A compare function which returns boolean or number
   * @param leftBound defalut value is `-Infinity`
   * @param rightBound defalut value is `Infinity`
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
  constructor(compareFn?: Comparator<T, 'number'>)
  constructor(compareFn: Comparator<T, 'number'>, leftBound: T, rightBound: T)
  constructor(
    compareFn: Comparator<T, any> = (a: any, b: any) => a - b,
    leftBound: any = -Infinity,
    rightBound: any = Infinity
  ) {
    this._root = new TreapNode<T>(rightBound, this._fastRandom())
    this._root.priority = Infinity
    this._root.left = new TreapNode<T>(leftBound, this._fastRandom())
    this._root.left.priority = -Infinity
    this._root.pushUp()

    this._leftBound = leftBound
    this._rightBound = rightBound
    this._compareFn = compareFn
  }

  /**
   * Returns true if value is a member.
   * @complexity `O(logn)`
   */
  has(value: T): boolean {
    return this._has(this._root, value)
  }

  private _has(node: TreapNode<T> | null, value: T): boolean {
    if (node == null) return false
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) return true
    if (cmp < 0) return this._has(node.right, value)
    return this._has(node.left, value)
  }

  /**
   * Add value to sorted set.
   * @complexity `O(logn)`
   */
  add(...values: T[]): this {
    values.forEach(value => this._add(this._root.left, value, this._root, 'left'))
    return this
  }

  private _add(
    node: TreapNode<T> | null,
    value: T,
    parent: TreapNode<T>,
    direction: 'left' | 'right'
  ): void {
    if (node == null) return

    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) {
      node.count++
      node.pushUp()
    } else if (cmp > 0) {
      if (node.left) {
        this._add(node.left, value, node, 'left')
      } else {
        node.left = new TreapNode(value, this._fastRandom())
        node.pushUp()
      }

      if (TreapNode.getFac(node.left) > node.priority) {
        parent[direction] = node.rotateRight()
      }
    } else if (cmp < 0) {
      if (node.right) {
        this._add(node.right, value, node, 'right')
      } else {
        node.right = new TreapNode(value, this._fastRandom())
        node.pushUp()
      }

      if (TreapNode.getFac(node.right) > node.priority) {
        parent[direction] = node.rotateLeft()
      }
    }

    parent.pushUp()
  }

  /**
   * Remove value from sorted set if it is a member.
   * If value is not a member, do nothing.
   * @complexity `O(logn)`
   */
  delete(value: T): void {
    this._delete(this._root.left, value, this._root, 'left')
  }

  private _delete(
    node: TreapNode<T> | null,
    value: T,
    parent: TreapNode<T>,
    direction: 'left' | 'right'
  ): void {
    if (node == null) return
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) {
      if (node.count > 1) {
        node.count--
        node.pushUp()
      } else if (node.left == null && node.right == null) {
        parent[direction] = null
      } else {
        // 旋到根节点
        // eslint-disable-next-line no-lonely-if
        if (node.right == null || TreapNode.getFac(node.left) > TreapNode.getFac(node.right)) {
          parent[direction] = node.rotateRight()
          this._delete(parent[direction]?.right ?? null, value, parent[direction]!, 'right')
        } else {
          parent[direction] = node.rotateLeft()
          this._delete(parent[direction]?.left ?? null, value, parent[direction]!, 'left')
        }
      }
    } else if (cmp > 0) {
      this._delete(node.left, value, node, 'left')
    } else {
      this._delete(node.right, value, node, 'right')
    }

    parent?.pushUp()
  }

  /**
   * Returns an index to insert value in the sorted set.
   * If the value is already present, the insertion point will
   * be before (to the left of) any existing values.
   * @complexity `O(logn)`
   */
  bisectLeft(value: T): number {
    // 因为有个lowerBound 所以-1
    return this._bisectLeft(this._root, value) - 1
  }

  private _bisectLeft(node: TreapNode<T> | null, value: T): number {
    if (node == null) return 0
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) {
      return TreapNode.getSize(node.left)
    }
    if (cmp > 0) {
      return this._bisectLeft(node.left, value)
    }
    return this._bisectLeft(node.right, value) + TreapNode.getSize(node.left) + node.count
  }

  /**
   * Returns an index to insert value in the sorted set.
   * If the value is already present, the insertion point will
   * be before (to the right of) any existing values.
   * @complexity `O(logn)`
   */
  bisectRight(value: T): number {
    // 因为有个lowerBound 所以-1
    return this._bisectRight(this._root, value) - 1
  }

  private _bisectRight(node: TreapNode<T> | null, value: T): number {
    if (node == null) return 0
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) {
      return TreapNode.getSize(node.left) + node.count
    }
    if (cmp > 0) {
      return this._bisectRight(node.left, value)
    }
    return this._bisectRight(node.right, value) + TreapNode.getSize(node.left) + node.count
  }

  /**
   * Returns the index of the first occurrence of a value in the set, or -1 if it is not present.
   * @complexity `O(logn)`
   */
  indexOf(value: T): number {
    let isExist = false

    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0
      const cmp = this._compareFn(node.value, value)
      if (cmp === 0) {
        isExist = true
        return TreapNode.getSize(node.left)
      }
      if (cmp > 0) {
        return dfs(node.left, value)
      }
      return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
    }

    // 因为有个lowerBound 所以-1
    const res = dfs(this._root, value) - 1
    return isExist ? res : -1
  }

  /**
   * Returns the index of the last occurrence of a value in the set, or -1 if it is not present.
   * @complexity `O(logn)`
   */
  lastIndexOf(value: T): number {
    let isExist = false

    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0
      const cmp = this._compareFn(node.value, value)
      if (cmp === 0) {
        isExist = true
        return TreapNode.getSize(node.left) + node.count - 1
      }
      if (cmp > 0) {
        return dfs(node.left, value)
      }
      return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
    }

    const res = dfs(this._root, value) - 1
    return isExist ? res : -1
  }

  /**
   * Returns the item located at the specified index.
   * @param index The zero-based index of the desired code unit.
   * A negative index will count back from the last item.
   * @complexity `O(logn)`
   */
  at(index: number): T | undefined {
    if (index < 0) index += this.size
    if (index < 0 || index >= this.size) return undefined

    const res = this._at(this._root, index + 2)
    return ([this._leftBound, this._rightBound] as any[]).includes(res) ? undefined : res
  }

  private _at(node: TreapNode<T> | null, rank: number): T | undefined {
    if (node == null) return undefined
    const leftSize = TreapNode.getSize(node.left)
    if (leftSize >= rank) {
      return this._at(node.left, rank)
    }
    if (leftSize + node.count >= rank) {
      return node.value
    }
    return this._at(node.right, rank - leftSize - node.count)
  }

  /**
   * Find and return the element less than `val`, return `undefined` if no such element found.
   * @complexity `O(logn)`
   */
  lower(value: T): T | undefined {
    const res = this._lower(this._root, value)
    return res === this._leftBound ? undefined : res
  }

  private _lower(node: TreapNode<T> | null, value: T): T | undefined {
    if (node == null) return undefined
    if (this._compareFn(node.value, value) >= 0) return this._lower(node.left, value)

    const tmp = this._lower(node.right, value)
    if (tmp == null || this._compareFn(node.value, tmp) > 0) {
      return node.value
    }
    return tmp
  }

  /**
   * Find and return the element greater than `val`.
   * return `undefined` if no such element found.
   * @complexity `O(logn)`
   */
  higher(value: T): T | undefined {
    const res = this._higher(this._root, value)
    return res === this._rightBound ? undefined : res
  }

  private _higher(node: TreapNode<T> | null, value: T): T | undefined {
    if (node == null) return undefined
    if (this._compareFn(node.value, value) <= 0) return this._higher(node.right, value)

    const tmp = this._higher(node.left, value)
    if (tmp == null || this._compareFn(node.value, tmp) < 0) {
      return node.value
    }
    return tmp
  }

  /**
   * Find and return the element less than or equal to `val`.
   * return `undefined` if no such element found.
   * @complexity `O(logn)`
   */
  floor(value: T): T | undefined {
    const res = this._floor(this._root, value)
    return res === this._leftBound ? undefined : res
  }

  private _floor(node: TreapNode<T> | null, value: T): T | undefined {
    if (node == null) return undefined
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) return node.value
    if (cmp >= 0) return this._floor(node.left, value)

    const tmp = this._floor(node.right, value)
    if (tmp == null || this._compareFn(node.value, tmp) > 0) {
      return node.value
    }
    return tmp
  }

  /**
   * Find and return the element greater than or equal to `val`.
   * return `undefined` if no such element found.
   * @complexity `O(logn)`
   */
  ceil(value: T): T | undefined {
    const res = this._ceil(this._root, value)
    return res === this._rightBound ? undefined : res
  }

  private _ceil(node: TreapNode<T> | null, value: T): T | undefined {
    if (node == null) return undefined
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) return node.value
    if (cmp <= 0) return this._ceil(node.right, value)

    const tmp = this._ceil(node.left, value)
    if (tmp == null || this._compareFn(node.value, tmp) < 0) {
      return node.value
    }
    return tmp
  }

  /**
   * Returns the last element from set.
   * If the set is empty, undefined is returned.
   * @complexity `O(logn)`
   */
  first(): T | undefined {
    const iter = this.inOrder()
    iter.next()
    const res = iter.next().value
    return res === this._rightBound ? undefined : res
  }

  /**
   * Returns the last element from set.
   * If the set is empty, undefined is returned .
   * @complexity `O(logn)`
   */
  last(): T | undefined {
    const iter = this.reverseInOrder()
    iter.next()
    const res = iter.next().value
    return res === this._leftBound ? undefined : res
  }

  /**
   * Removes the first element from an set and returns it.
   * If the set is empty, undefined is returned and the set is not modified.
   * @complexity `O(logn)`
   */
  shift(): T | undefined {
    const first = this.first()
    if (first === undefined) return undefined
    this.delete(first)
    return first
  }

  /**
   * Removes the last element from an set and returns it.
   * If the set is empty, undefined is returned and the set is not modified.
   * @complexity `O(logn)`
   */
  pop(index?: number): T | undefined {
    if (index == null) {
      const last = this.last()
      if (last === undefined) return undefined
      this.delete(last)
      return last
    }

    const toDelete = this.at(index)
    if (toDelete == null) return undefined
    this.delete(toDelete)
    return toDelete
  }

  /**
   * Returns number of occurrences of value in the sorted set.
   * @complexity `O(logn)`
   */
  count(value: T): number {
    return this._count(this._root, value)
  }

  private _count(node: TreapNode<T> | null, value: T): number {
    if (node == null) return 0
    const cmp = this._compareFn(node.value, value)
    if (cmp === 0) return node.count
    if (cmp < 0) return this._count(node.right, value)
    return this._count(node.left, value)
  }

  /**
   * Returns an iterable of keys in the set.
   */
  *keys(): Generator<T, any, any> {
    yield* this.values()
  }

  /**
   * Returns an iterable of values in the set.
   */
  *values(): Generator<T, any, any> {
    const iter = this.inOrder()
    iter.next()
    const steps = this.size
    for (let _ = 0; _ < steps; _++) {
      yield iter.next().value
    }
  }

  /**
   * Returns a generator for reversed order traversing the set.
   */
  *rvalues(): Generator<T, any, any> {
    const iter = this.reverseInOrder()
    iter.next()
    const steps = this.size
    for (let _ = 0; _ < steps; _++) {
      yield iter.next().value
    }
  }

  /**
   * Returns an iterable of key, value pairs for every entry in the set.
   */
  *entries(): IterableIterator<[number, T]> {
    const iter = this.inOrder()
    iter.next()
    const steps = this.size
    for (let i = 0; i < steps; i++) {
      yield [i, iter.next().value]
    }
  }

  *[Symbol.iterator](): Generator<T, any, any> {
    yield* this.values()
  }

  get size(): number {
    return this._root.size - 2
  }

  get height(): number {
    const getHeight = (node: TreapNode<T> | null): number => {
      if (node == null) return 0
      return 1 + Math.max(getHeight(node.left), getHeight(node.right))
    }

    return getHeight(this._root)
  }

  private _fastRandom(): number {
    this._seed ^= this._seed << 13
    this._seed ^= this._seed >>> 17
    this._seed ^= this._seed << 5
    return this._seed >>> 0
  }

  private *inOrder(root: TreapNode<T> | null = this._root): Generator<T, any, any> {
    if (root == null) return
    yield* this.inOrder(root.left)
    const count = root.count
    for (let _ = 0; _ < count; _++) {
      yield root.value
    }
    yield* this.inOrder(root.right)
  }

  private *reverseInOrder(root: TreapNode<T> | null = this._root): Generator<T, any, any> {
    if (root == null) return
    yield* this.reverseInOrder(root.right)
    const count = root.count
    for (let _ = 0; _ < count; _++) {
      yield root.value
    }
    yield* this.reverseInOrder(root.left)
  }
}

export { TreapMultiSet }

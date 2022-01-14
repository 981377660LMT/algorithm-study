import assert from 'assert'

type CompareFunction<T, R extends 'number' | 'boolean'> = (
  a: T,
  b: T
) => R extends 'number' ? number : boolean

interface ITreapMultiSet<T> extends Iterable<T> {
  add: (value: T) => this
  has: (value: T) => boolean
  delete: (value: T) => void
  bisectLeft: (value: T) => number
  bisectRight: (value: T) => number
  getRankByValue: (value: T) => number
  at: (index: number) => T | undefined
  lower: (value: T) => T | undefined
  higher: (value: T) => T | undefined
  floor: (value: T) => T | undefined
  ceil: (value: T) => T | undefined
  first: () => T | undefined
  last: () => T | undefined
  shift: () => T | undefined
  pop: () => T | undefined
  keys: () => Generator<T, void, void>
  values: () => Generator<T, void, void>
  rvalues: () => Generator<T, void, void>
  readonly size: number
}

class TreapNode<T = number> {
  value: T
  count: number
  size: number
  fac: number
  left: TreapNode<T> | null
  right: TreapNode<T> | null

  constructor(value: T) {
    this.value = value
    this.count = 1
    this.size = 1
    this.fac = Math.random()
    this.left = null
    this.right = null
  }

  static getSize(node: TreapNode<any> | null): number {
    return node?.size ?? 0
  }

  static getFac(node: TreapNode<any> | null): number {
    return node?.fac ?? 0
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
  private compare: CompareFunction<T, 'number'>
  private lowerBound: T
  private upperBound: T

  constructor(compare?: CompareFunction<T, 'number'>)
  // constructor(compare?: CompareFunction<T, 'boolean'>)
  constructor(compare: CompareFunction<T, 'number'>, left: T, right: T)
  // constructor(compare: CompareFunction<T, 'boolean'>, left: T, right: T)
  constructor(
    compare: CompareFunction<T, any> = (a: any, b: any) => a - b,
    left: any = -Infinity,
    right: any = Infinity
  ) {
    this.root = new TreapNode<T>(right)
    this.root.fac = Infinity
    this.root.left = new TreapNode<T>(left)
    this.root.left!.fac = -Infinity
    this.root.pushUp()

    this.compare = compare
    this.lowerBound = left
    this.upperBound = right
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

  has(value: T): boolean {
    const compare = this.compare
    const dfs = (node: TreapNode<T> | null, value: T): boolean => {
      if (node == null) return false
      if (compare(node.value, value) === 0) return true
      if (compare(node.value, value) < 0) return dfs(node.right, value)
      return dfs(node.left, value)
    }

    return dfs(this.root, value)
  }

  add(value: T): this {
    const compare = this.compare
    // js 里没 & 这种引用  所以要带着parent和上次的方向  在c++里直接 Tree &rt 就可以了
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

        if (TreapNode.getFac(node.left) > node.fac) {
          parent[direction] = node.rotateRight()
        }
      } else if (compare(node.value, value) < 0) {
        if (node.right) {
          dfs(node.right, value, node, 'right')
        } else {
          node.right = new TreapNode(value)
          node.pushUp()
        }

        if (TreapNode.getFac(node.right) > node.fac) {
          parent[direction] = node.rotateLeft()
        }
      }
      parent.pushUp()
    }

    dfs(this.root.left, value, this.root, 'left')
    return this
  }

  delete(value: T): void {
    const compare = this.compare

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
   * @param value
   * @returns 当前元素位于第几位，rank从0开始
   */
  getRankByValue(value: T): number {
    const compare = this.compare

    const dfs = (node: TreapNode<T> | null, value: T): number => {
      if (node == null) return 0

      if (compare(node.value, value) === 0) {
        return TreapNode.getSize(node.left) + 1
      } else if (compare(node.value, value) > 0) {
        return dfs(node.left, value)
      } else if (compare(node.value, value) < 0) {
        return dfs(node.right, value) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个-Infinity 所以-1
    return dfs(this.root, value) - 1
  }

  bisectLeft(value: T): number {
    const compare = this.compare

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

    // 因为有个-Infinity 所以-1
    return dfs(this.root, value) - 1
  }

  bisectRight(value: T): number {
    const compare = this.compare

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

    // 因为有个-Infinity 所以-1
    return dfs(this.root, value) - 1
  }

  /**
   *
   * @param index 支持负索引
   * @description 时间复杂度O(logN)
   */
  at(index: number): T | undefined {
    if (index < 0) index += this.size

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
    return ([this.lowerBound, this.upperBound] as any[]).includes(res) ? undefined : res
  }

  /**
   *
   * @param value
   * @returns 严格小于val的第一个数
   */
  lower(value: T): T | undefined {
    const compare = this.compare

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
    return res === this.lowerBound ? undefined : res
  }

  /**
   *
   * @param value
   * @returns 严格大于val的第一个数
   */
  higher(value: T): T | undefined {
    const compare = this.compare

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
    return res === this.upperBound ? undefined : res
  }

  /**
   *
   * @param value
   * @returns 小于等于val的第一个数
   */
  floor(value: T): T | undefined {
    const compare = this.compare

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
    return res === this.lowerBound ? undefined : res
  }

  /**
   *
   * @param value
   * @returns 大于等于val的第一个数
   */
  ceil(value: T): T | undefined {
    const compare = this.compare

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
    return res === this.upperBound ? undefined : res
  }

  first(): T | undefined {
    const iter = this.inOrder()
    iter.next()
    const res = iter.next().value
    return res === this.upperBound ? undefined : res
  }

  last(): T | undefined {
    const iter = this.reverseInOrder()
    iter.next()
    const res = iter.next().value
    return res === this.lowerBound ? undefined : res
  }

  shift(): T | undefined {
    const first = this.first()
    if (first == undefined) return undefined
    this.delete(first)
    return first
  }

  pop(): T | undefined {
    const last = this.last()
    if (last == undefined) return undefined
    this.delete(last)
    return last
  }

  *[Symbol.iterator](): Generator<T, void, void> {
    yield* this.values()
  }

  *keys(): Generator<T, void, void> {
    yield* this.values()
  }

  *values(): Generator<T, void, void> {
    const iter = this.inOrder()
    iter.next()
    let remain = this.size
    for (let _ = 0; _ < remain; _++) {
      yield iter.next().value
    }
  }

  /**
   * Return a generator for reverse order traversing the multi-set
   */
  *rvalues(): Generator<T, void, void> {
    const iter = this.reverseInOrder()
    iter.next()
    let remain = this.size
    for (let _ = 0; _ < remain; _++) {
      yield iter.next().value
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

  // getRankByValue
  assert.strictEqual(treap.getRankByValue(1.2), 1)
  assert.strictEqual(treap.getRankByValue(2.2), 2)
  assert.strictEqual(treap.getRankByValue(3), 3)
  assert.strictEqual(treap.getRankByValue(4), 5)

  // at
  assert.strictEqual(treap.at(0), 1)
  assert.strictEqual(treap.at(1), 2)
  assert.strictEqual(treap.at(2), 3)
  assert.strictEqual(treap.at(3), 3)
  assert.strictEqual(treap.at(4), 3)
  assert.strictEqual(treap.at(-1), 3)
  assert.strictEqual(treap.at(-100), undefined)

  // values rvalues
  assert.deepStrictEqual([...treap.values()], [1, 2, 3, 3, 3])
  assert.deepStrictEqual([...treap.rvalues()], [3, 3, 3, 2, 1])

  // first last shift pop
  assert.deepStrictEqual(treap.shift(), 1)
  assert.deepStrictEqual(treap.first(), 2)
  assert.deepStrictEqual(treap.pop(), 3)
  assert.deepStrictEqual(treap.last(), 3)
  assert.deepStrictEqual(treap.size, 3)

  //getRankByValue bisectLeft bisectRight
  // [2,3,3]
  assert.deepStrictEqual(treap.getRankByValue(1.9), 0)
  assert.deepStrictEqual(treap.getRankByValue(2), 1)
  assert.deepStrictEqual(treap.getRankByValue(2.5), 1)
  assert.deepStrictEqual(treap.getRankByValue(3), 2)
  assert.deepStrictEqual(treap.getRankByValue(4), 3)
  assert.deepStrictEqual(treap.bisectLeft(1.9), 0)
  assert.deepStrictEqual(treap.bisectLeft(2), 0)
  assert.deepStrictEqual(treap.bisectLeft(2.5), 1)
  assert.deepStrictEqual(treap.bisectLeft(3), 1)
  assert.deepStrictEqual(treap.bisectLeft(4), 3)
  assert.deepStrictEqual(treap.bisectRight(1.9), 0)
  assert.deepStrictEqual(treap.bisectRight(2), 1)
  assert.deepStrictEqual(treap.bisectRight(2.5), 1)
  assert.deepStrictEqual(treap.bisectRight(3), 3)
  assert.deepStrictEqual(treap.bisectRight(4), 3)
}

export { TreapMultiSet }

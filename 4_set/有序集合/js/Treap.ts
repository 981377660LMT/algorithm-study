import assert from 'assert'

class TreapNode<T = number> {
  val: T
  count: number
  size: number
  fac: number
  left: TreapNode<T> | null
  right: TreapNode<T> | null

  constructor(val: T) {
    this.val = val
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

  rotateLeft() {
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

type CompareFunction<T, R extends 'number' | 'boolean'> = (
  a: T,
  b: T
) => R extends 'number' ? number : boolean

class TreapMultiSet<T = number> {
  private root: TreapNode<T>
  private compare: CompareFunction<T, 'number'>

  constructor(
    compare: CompareFunction<T, 'number'> = (a: any, b: any) => a - b,
    left = -Infinity,
    right = Infinity
  ) {
    this.root = new TreapNode<any>(right)
    this.root.fac = Infinity
    this.root.left = new TreapNode<any>(left)
    this.root.left!.fac = -Infinity
    this.root.pushUp()

    this.compare = compare
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

  has(val: T): boolean {
    const compare = this.compare
    const dfs = (node: TreapNode<T> | null, val: T): boolean => {
      if (node == null) return false
      if (compare(node.val, val) === 0) return true
      if (compare(node.val, val) < 0) return dfs(node.right, val)
      return dfs(node.left, val)
    }

    return dfs(this.root, val)
  }

  add(val: T): void {
    const compare = this.compare

    // js 里没 & 这种引用  所以要带着parent和上次的方向  在c++里直接 Tree &rt 就可以了
    const dfs = (
      node: TreapNode<T> | null,
      val: T,
      parent: TreapNode<T>,
      direction: 'left' | 'right'
    ): void => {
      if (node == null) return
      if (compare(node.val, val) === 0) {
        node.count++
        node.pushUp()
      } else if (compare(node.val, val) > 0) {
        if (node.left) {
          dfs(node.left, val, node, 'left')
        } else {
          node.left = new TreapNode(val)
          node.pushUp()
        }

        if (TreapNode.getFac(node.left) > node.fac) {
          parent[direction] = node.rotateRight()
        }
      } else if (compare(node.val, val) < 0) {
        if (node.right) {
          dfs(node.right, val, node, 'right')
        } else {
          node.right = new TreapNode(val)
          node.pushUp()
        }

        if (TreapNode.getFac(node.right) > node.fac) {
          parent[direction] = node.rotateLeft()
        }
      }
      parent.pushUp()
    }

    dfs(this.root.left, val, this.root, 'left')
  }

  delete(val: T): void {
    const compare = this.compare

    const dfs = (
      node: TreapNode<T> | null,
      val: T,
      parent: TreapNode<T>,
      direction: 'left' | 'right'
    ): void => {
      if (node == null) return

      if (compare(node.val, val) === 0) {
        if (node.count > 1) {
          node.count--
          node?.pushUp()
        } else if (node.left == null && node.right == null) {
          parent[direction] = null
        } else {
          // 旋到根节点
          if (node.right == null || TreapNode.getFac(node.left) > TreapNode.getFac(node.right)) {
            parent[direction] = node.rotateRight()
            dfs(parent[direction]?.right ?? null, val, parent[direction]!, 'right')
          } else {
            parent[direction] = node.rotateLeft()
            dfs(parent[direction]?.left ?? null, val, parent[direction]!, 'left')
          }
        }
      } else if (compare(node.val, val) > 0) {
        dfs(node.left, val, node, 'left')
      } else if (compare(node.val, val) < 0) {
        dfs(node.right, val, node, 'right')
      }

      parent?.pushUp()
    }

    dfs(this.root.left, val, this.root, 'left')
  }

  /**
   *
   * @param val
   * @returns 当前元素位于第几位，rank从0开始
   */
  getRankByValue(val: T): number {
    const compare = this.compare

    const dfs = (node: TreapNode<T> | null, val: T): number => {
      if (node == null) return 0

      if (compare(node.val, val) === 0) {
        return TreapNode.getSize(node.left) + 1
      } else if (compare(node.val, val) > 0) {
        return dfs(node.left, val)
      } else if (compare(node.val, val) < 0) {
        return dfs(node.right, val) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个-Infinity 所以-1
    return dfs(this.root, val) - 1
  }

  bisectLeft(val: T): number {
    const compare = this.compare

    const dfs = (node: TreapNode<T> | null, val: T): number => {
      if (node == null) return 0

      if (compare(node.val, val) === 0) {
        return TreapNode.getSize(node.left)
      } else if (compare(node.val, val) > 0) {
        return dfs(node.left, val)
      } else if (compare(node.val, val) < 0) {
        return dfs(node.right, val) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个-Infinity 所以-1
    return dfs(this.root, val) - 1
  }

  bisectRight(val: T): number {
    const compare = this.compare

    const dfs = (node: TreapNode<T> | null, val: T): number => {
      if (node == null) return 0

      if (compare(node.val, val) === 0) {
        return TreapNode.getSize(node.left) + node.count
      } else if (compare(node.val, val) > 0) {
        return dfs(node.left, val)
      } else if (compare(node.val, val) < 0) {
        return dfs(node.right, val) + TreapNode.getSize(node.left) + node.count
      }

      return 0
    }

    // 因为有个-Infinity 所以-1
    return dfs(this.root, val) - 1
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
        return node.val
      } else {
        return dfs(node.right, rank - TreapNode.getSize(node.left) - node.count)
      }
    }

    // 因为有个-Infinity 所以 + 2
    const res = dfs(this.root, index + 2)
    return ([Infinity, -Infinity] as any[]).includes(res) ? undefined : res
  }

  /**
   *
   * @param val
   * @returns 严格小于val的第一个数
   */
  lower(val: T): T | undefined {
    if (val == null) return undefined
    const compare = this.compare

    const dfs = (node: TreapNode<T> | null, val: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.val, val) >= 0) return dfs(node.left, val)

      const tmp = dfs(node.right, val)
      if (tmp == null || compare(node.val, tmp) > 0) {
        return node.val
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, val) as any
    return res === -Infinity ? undefined : res
  }

  /**
   *
   * @param val
   * @returns 严格大于val的第一个数
   */
  upper(val: T): T | undefined {
    if (val == null) return undefined
    const compare = this.compare
    const dfs = (node: TreapNode<T> | null, val: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.val, val) <= 0) return dfs(node.right, val)

      const tmp = dfs(node.left, val)

      if (tmp == null || compare(node.val, tmp) < 0) {
        return node.val
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, val) as any
    return res === Infinity ? undefined : res
  }

  /**
   *
   * @param val
   * @returns 小于等于val的第一个数
   */
  floor(val: T): T | undefined {
    if (val == null) return undefined
    const compare = this.compare

    const dfs = (node: TreapNode<T> | null, val: T): T | undefined => {
      if (node == null) return undefined
      if (compare(node.val, val) === 0) return node.val
      if (compare(node.val, val) >= 0) return dfs(node.left, val)

      const tmp = dfs(node.right, val)
      if (tmp == null || compare(node.val, tmp) > 0) {
        return node.val
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, val) as any
    return res === -Infinity ? undefined : res
  }

  /**
   *
   * @param val
   * @returns 大于等于val的第一个数
   */
  ceil(val: T): T | undefined {
    if (val == null) return undefined
    const compare = this.compare
    function dfs(node: TreapNode<T> | null, val: T): T | undefined {
      if (node == null) return undefined
      if (compare(node.val, val) === 0) return node.val
      if (compare(node.val, val) <= 0) return dfs(node.right, val)

      const tmp = dfs(node.left, val)

      if (tmp == null || compare(node.val, tmp) < 0) {
        return node.val
      } else {
        return tmp
      }
    }

    const res = dfs(this.root, val) as any
    return res === Infinity ? undefined : res
  }

  first(): T | undefined {
    const iter = this.inOrder()
    iter.next()
    const res = iter.next().value
    return res === Infinity ? undefined : res
  }

  last(): T | undefined {
    const iter = this.reverseInOrder()
    iter.next()
    const res = iter.next().value
    return res === -Infinity ? undefined : res
  }

  shift(): T | undefined {
    const first = this.first()
    if (first === undefined) return undefined
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
      yield root.val
    }
    yield* this.inOrder(root.right)
  }

  private *reverseInOrder(root: TreapNode<T> | null = this.root): Generator<T, any, any> {
    if (root == null) return
    yield* this.reverseInOrder(root.right)
    const count = root.count
    for (let _ = 0; _ < count; _++) {
      yield root.val
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
  console.dir(treap, { depth: null })

  // upper lower ceil floor
  treap.add(3)
  assert.strictEqual(treap.upper(2), 3)
  assert.strictEqual(treap.upper(1.9), 2)
  assert.strictEqual(treap.upper(3), undefined)
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

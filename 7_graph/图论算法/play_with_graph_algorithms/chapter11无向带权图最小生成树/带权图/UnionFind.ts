class UnionFind<U = unknown> {
  static rootSymbol = Symbol.for('UnionFind_Root')
  // 记录每个节点的父节点
  // 如果节点互相连通（从一个节点可以达到另一个节点），那么他们的祖先是相同的。
  private parent: Map<U, U | symbol>
  // 记录无向图连通域数量
  private count: number

  constructor() {
    this.parent = new Map()
    this.count = 0
  }

  get size() {
    return this.count
  }

  isConnected(key1: U, key2: U) {
    const root1 = this.find(key1)
    const root2 = this.find(key2)
    return root1 !== undefined && root2 !== undefined && this.find(key1) === this.find(key2)
  }

  /**
   *
   * @param key 把一个新节点添加到并查集中，它的父节点应该为UnionFind.rootSymbol。
   */
  add(key: U): this {
    if (!this.parent.has(key)) {
      this.parent.set(key, UnionFind.rootSymbol)
      this.count++
    }
    return this
  }

  /**
   *
   * @description 如果两个节点是连通的，那么就要把他们合并，也就是他们的祖先是相同的。
   * @example
   * ```js
   * const union = new UnionFind<number>()
   * union.add(1).add(2).add(3).add(4).union(2, 3).union(4, 3)
   * console.dir(union, { depth: null })
   *
   * // output:
   * UnionFind {
   *   parent: Map(4) { 1 => undefined, 2 => 3, 3 => undefined, 4 => 3 }
   * }
   * ```
   */
  union(key1: U, key2: U): this {
    const root1 = this.find(key1)
    const root2 = this.find(key2)
    if (root1 !== undefined && root2 !== undefined && root1 !== root2) {
      this.parent.set(root1, root2)
      this.count--
    }
    return this
  }

  // /**
  //  * @description 判断两个节点是否处于同一个连通分量的时候，就要判断他们的祖先是否相同。
  //  */
  // isConnected(key1: U, key2: U): boolean {
  //   return this.find(key1) === this.find(key2)
  // }

  /**
   *
   * @param key 查找祖先；如果节点的父节点不为空或者symbol，那就不断迭代。
   * @returns 返回undefined代表key不在并查集中
   */
  find(key: U): U | undefined {
    let root = key as any
    if (!this.parent.has(root)) return undefined
    while (this.parent.get(root) !== UnionFind.rootSymbol) {
      root = this.parent.get(root)
    }
    return root
  }
}

export { UnionFind }

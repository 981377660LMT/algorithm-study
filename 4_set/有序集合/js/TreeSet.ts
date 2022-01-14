import assert from 'assert'

class RBTreeNode<T = number> {
  data: T
  color: number
  left: RBTreeNode<T> | null
  right: RBTreeNode<T> | null
  parent: RBTreeNode<T> | null

  constructor(data: T) {
    this.data = data
    this.left = this.right = this.parent = null
    this.color = 0
  }

  sibling(): RBTreeNode<T> | null {
    if (!this.parent) return null // sibling null if no parent
    return this.isOnLeft() ? this.parent.right : this.parent.left
  }

  hasRedChild(): boolean {
    return Boolean(this.left?.color === 0) || Boolean(this.right?.color === 0)
  }

  isOnLeft(): boolean {
    return this === this.parent?.left
  }
}

type CompareFunction<T, R extends 'number' | 'boolean'> = (
  a: T,
  b: T
) => R extends 'number' ? number : boolean

class RBTree<T = number> {
  root: RBTreeNode<T> | null
  private compare: CompareFunction<T, 'boolean'>
  static defaultCompare = (a: any, b: any) => a - b

  constructor(compare: CompareFunction<T, 'number'> = RBTree.defaultCompare) {
    this.root = null
    this.compare = (a: any, b: any) => {
      const diff = compare(a, b)
      return diff < 0
    }
  }

  insert(data: T): boolean {
    const node = new RBTreeNode(data)
    const parent = this.search(data)
    if (!parent) this.root = node
    else if (this.compare(node.data, parent.data)) parent.left = node
    else if (this.compare(parent.data, node.data)) parent.right = node
    else return false
    node.parent = parent
    this.fixAfterInsert(node)
    return true
  }

  find(data: T): RBTreeNode<T> | null {
    const node = this.search(data)
    return node && node.data === data ? node : null
  }

  deleteByValue(val: T): boolean {
    const node = this.search(val)
    if (node?.data !== val) return false
    this.deleteNode(node)
    return true
  }

  *inOrder(root = this.root): Generator<T, any, any> {
    if (root == null) return
    yield* this.inOrder(root.left)
    yield root.data
    yield* this.inOrder(root.right)
  }

  *reverseInOrder(root = this.root): Generator<T, any, any> {
    if (root == null) return
    yield* this.reverseInOrder(root.right)
    yield root.data
    yield* this.reverseInOrder(root.left)
  }

  private rotateLeft(pt: RBTreeNode<T>): void {
    const right = pt.right
    pt.right = right?.left ?? null
    if (pt.right) pt.right.parent = pt
    right!.parent = pt.parent
    if (!pt.parent) this.root = right
    else if (pt === pt.parent.left) pt.parent.left = right
    else pt.parent.right = right
    right!.left = pt
    pt.parent = right
  }

  private rotateRight(pt: RBTreeNode<T>): void {
    const left = pt.left
    pt.left = left?.right ?? null
    if (pt.left) pt.left.parent = pt
    left!.parent = pt.parent
    if (!pt.parent) this.root = left
    else if (pt === pt.parent.left) pt.parent.left = left
    else pt.parent.right = left
    left!.right = pt
    pt.parent = left
  }

  private swapColor(p1: RBTreeNode<T>, p2: RBTreeNode<T>): void {
    const tmp = p1.color
    p1.color = p2.color
    p2.color = tmp
  }

  private swapData(p1: RBTreeNode<T>, p2: RBTreeNode<T>): void {
    const tmp = p1.data
    p1.data = p2.data
    p2.data = tmp
  }

  private fixAfterInsert(pt: RBTreeNode<T>): void {
    let parent: RBTreeNode<T> | null = null
    let grandParent: RBTreeNode<T> | null = null
    while (pt !== this.root && pt.color !== 1 && pt.parent?.color === 0) {
      parent = pt.parent
      grandParent = pt.parent.parent
      /*  Case : A
                    Parent of pt is left child of Grand-parent of pt */
      if (parent === grandParent?.left) {
        const uncle = grandParent.right
        /* Case : 1
                         The uncle of pt is also red
                         Only Recoloring required */
        if (uncle && uncle.color === 0) {
          grandParent.color = 0
          parent.color = 1
          uncle.color = 1
          pt = grandParent
        } else {
          /* Case : 2
                               pt is right child of its parent
                               Left-rotation required */
          if (pt === parent.right) {
            this.rotateLeft(parent)
            pt = parent
            parent = pt.parent
          }
          /* Case : 3
                               pt is left child of its parent
                               Right-rotation required */
          this.rotateRight(grandParent)
          this.swapColor(parent!, grandParent)
          pt = parent!
        }
      } else {
        /* Case : B
                   Parent of pt is right child of Grand-parent of pt */
        const uncle = grandParent!.left
        /*  Case : 1
                          The uncle of pt is also red
                          Only Recoloring required */
        if (uncle != null && uncle.color === 0) {
          grandParent!.color = 0
          parent.color = 1
          uncle.color = 1
          pt = grandParent!
        } else {
          /* Case : 2
                               pt is left child of its parent
                               Right-rotation required */
          if (pt === parent.left) {
            this.rotateRight(parent)
            pt = parent
            parent = pt.parent
          }
          /* Case : 3
                               pt is right child of its parent
                               Left-rotation required */
          this.rotateLeft(grandParent!)
          this.swapColor(parent!, grandParent!)
          pt = parent!
        }
      }
    }
    this.root!.color = 1
  }

  // searches for given value
  // if found returns the node (used for delete)
  // else returns the last node while traversing (used in insert)
  private search(val: T): RBTreeNode<T> | null {
    let p = this.root
    while (p) {
      if (this.compare(val, p.data)) {
        if (!p.left) break
        else p = p.left
      } else if (this.compare(p.data, val)) {
        if (!p.right) break
        else p = p.right
      } else break
    }
    return p
  }

  private deleteNode(v: RBTreeNode<T>): void {
    const u = BSTreplace(v)
    // True when u and v are both black
    const uvBlack = (u == null || u.color === 1) && v.color === 1
    const parent = v.parent
    if (!u) {
      // u is null therefore v is leaf
      if (v === this.root) this.root = null
      // v is root, making root null
      else {
        if (uvBlack) {
          // u and v both black
          // v is leaf, fix double black at v
          this.fixDoubleBlack(v)
        } else {
          // u or v is red
          if (v.sibling()) {
            // sibling is not null, make it red"
            v.sibling()!.color = 0
          }
        }
        // delete v from the tree
        if (v.isOnLeft()) parent!.left = null
        else parent!.right = null
      }
      return
    }
    if (!v.left || !v.right) {
      // v has 1 child
      if (v === this.root) {
        // v is root, assign the value of u to v, and delete u
        v.data = u.data
        v.left = v.right = null
      } else {
        // Detach v from tree and move u up
        if (v.isOnLeft()) parent!.left = u
        else parent!.right = u
        u.parent = parent
        if (uvBlack) this.fixDoubleBlack(u)
        // u and v both black, fix double black at u
        else u.color = 1 // u or v red, color u black
      }
      return
    }
    // v has 2 children, swap data with successor and recurse
    this.swapData(u, v)
    this.deleteNode(u)
    // find node that replaces a deleted node in BST
    function BSTreplace(x: RBTreeNode<T>) {
      // when node have 2 children
      if (x.left && x.right) return successor(x.right)
      // when leaf
      if (!x.left && !x.right) return null
      // when single child
      return x.left ?? x.right
    }
    // find node that do not have a left child
    // in the subtree of the given node
    function successor(x: RBTreeNode<T>) {
      let temp = x
      while (temp.left) temp = temp.left
      return temp
    }
  }

  private fixDoubleBlack(x: RBTreeNode<T>): void {
    if (x === this.root) return // Reached root
    const sibling = x.sibling()
    const parent = x.parent as RBTreeNode<T>
    if (!sibling) {
      // No sibiling, double black pushed up
      this.fixDoubleBlack(parent)
    } else {
      if (sibling.color === 0) {
        // Sibling red
        parent!.color = 0
        sibling.color = 1
        if (sibling.isOnLeft()) this.rotateRight(parent)
        // left case
        else this.rotateLeft(parent) // right case
        this.fixDoubleBlack(x)
      } else {
        // Sibling black
        if (sibling.hasRedChild()) {
          // at least 1 red children
          if (sibling.left && sibling.left.color === 0) {
            if (sibling.isOnLeft()) {
              // left left
              sibling.left.color = sibling.color
              sibling.color = parent.color
              this.rotateRight(parent)
            } else {
              // right left
              sibling.left.color = parent.color
              this.rotateRight(sibling)
              this.rotateLeft(parent)
            }
          } else {
            if (sibling.isOnLeft()) {
              // left right
              sibling.right!.color = parent.color
              this.rotateLeft(sibling)
              this.rotateRight(parent)
            } else {
              // right right
              sibling.right!.color = sibling.color
              sibling.color = parent.color
              this.rotateLeft(parent)
            }
          }
          parent.color = 1
        } else {
          // 2 black children
          sibling.color = 0
          if (parent.color === 1) this.fixDoubleBlack(parent)
          else parent.color = 1
        }
      }
    }
  }
}

/**
 * @description C++ 里的set
 */
class TreeSet<T = number> {
  private _size: number
  private tree: RBTree<T>
  private compare: CompareFunction<T, 'boolean'>

  constructor(collection: Iterable<T> = [], compare = RBTree.defaultCompare) {
    this._size = 0
    this.tree = new RBTree(compare)
    this.compare = (a: any, b: any) => {
      const isSmaller = compare(a, b)
      return isSmaller < 0
    }
    for (const val of collection) this.add(val)
  }

  get size() {
    return this._size
  }

  has(val: T): boolean {
    return !!this.tree.find(val)
  }

  add(val: T): boolean {
    const added = this.tree.insert(val)
    this._size += added ? 1 : 0
    return added
  }

  delete(val: T): boolean {
    const deleted = this.tree.deleteByValue(val)
    this._size -= deleted ? 1 : 0
    return deleted
  }

  /**
   *
   * @param val
   * @returns 大于等于val的第一个数
   */
  ceiling(val: T): T | undefined {
    let p = this.tree.root
    let higher = null
    while (p) {
      if (!this.compare(p.data, val)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  /**
   *
   * @param val
   * @returns 小于等于val的第一个数
   */
  floor(val: T): T | undefined {
    let p = this.tree.root
    let lower = null
    while (p) {
      if (!this.compare(val, p.data)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  /**
   *
   * @param val
   * @returns 严格大于val的第一个数
   */
  higher(val: T): T | undefined {
    let p = this.tree.root
    let higher = null
    while (p) {
      if (this.compare(val, p.data)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  /**
   *
   * @param val
   * @returns 严格小于val的第一个数
   */
  lower(val: T): T | undefined {
    let p = this.tree.root
    let lower = null
    while (p) {
      if (this.compare(p.data, val)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  first(): T | undefined {
    return this.tree.inOrder().next().value
  }

  last(): T | undefined {
    return this.tree.reverseInOrder().next().value
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

  *values(): Generator<T, undefined, void> {
    yield* this.tree.inOrder()
    return undefined
  }

  /**
   * Return a generator for reverse order traversing the set
   */
  *rvalues(): Generator<T, undefined, void> {
    yield* this.tree.inOrder()
    return undefined
  }
}

/**
 * @description C++ 里的multiset
 */
class TreeMultiSet<T = number> {
  private _size: number
  private tree: RBTree<T>
  private counts: Map<T, number>
  private compare: CompareFunction<T, 'boolean'>

  constructor(
    collection: Iterable<T> = [],
    compare: CompareFunction<T, 'number'> = RBTree.defaultCompare
  ) {
    this._size = 0
    this.tree = new RBTree(compare)
    this.counts = new Map()
    this.compare = (a: any, b: any) => {
      const isSmaller = compare(a, b)
      return isSmaller < 0
    }
    for (const val of collection) this.add(val)
  }

  get size() {
    return this._size
  }

  has(val: T): boolean {
    return !!this.tree.find(val)
  }

  add(val: T): boolean {
    const added = this.tree.insert(val)
    this.increase(val)
    this._size++
    return added
  }

  delete(val: T): boolean {
    if (!this.has(val)) return false
    this.decrease(val)
    if (this.count(val) === 0) {
      this.tree.deleteByValue(val)
    }
    this._size--
    return true
  }

  count(val: T): number {
    return this.counts.get(val) ?? 0
  }

  ceiling(val: T): T | undefined {
    let p = this.tree.root
    let higher = null
    while (p) {
      if (!this.compare(p.data, val)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  floor(val: T): T | undefined {
    let p = this.tree.root
    let lower = null
    while (p) {
      if (!this.compare(val, p.data)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  higher(val: T): T | undefined {
    let p = this.tree.root
    let higher = null
    while (p) {
      if (this.compare(val, p.data)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  lower(val: T): T | undefined {
    let p = this.tree.root
    let lower = null
    while (p) {
      if (this.compare(p.data, val)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  first(): T | undefined {
    return this.tree.inOrder().next().value
  }

  last(): T | undefined {
    return this.tree.reverseInOrder().next().value
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

  *values(): Generator<T, undefined, void> {
    for (const val of this.tree.inOrder()) {
      let count = this.count(val)
      while (count--) yield val
    }
    return undefined
  }

  /**
   * Return a generator for reverse order traversing the multi-set
   */
  *rvalues(): Generator<T, undefined, void> {
    for (const val of this.tree.reverseInOrder()) {
      let count = this.count(val)
      while (count--) yield val
    }
    return undefined
  }

  private decrease(val: T): void {
    this.counts.set(val, this.count(val) - 1)
  }

  private increase(val: T): void {
    this.counts.set(val, this.count(val) + 1)
  }
}

if (require.main === module) {
  const treeSet = new TreeSet()

  // add
  treeSet.add(1)
  treeSet.add(2)
  assert.strictEqual(treeSet.size, 2)
  assert.strictEqual(treeSet.has(1), true)
  assert.strictEqual(treeSet.has(2), true)
  assert.strictEqual(treeSet.has(3), false)
  treeSet.add(2)
  assert.strictEqual(treeSet.size, 2)

  // delete
  treeSet.delete(2)
  assert.strictEqual(treeSet.size, 1)

  // keys
  assert.deepStrictEqual([...treeSet.keys()], [1])

  // first last
  treeSet.add(2)
  assert.strictEqual(treeSet.first(), 1)
  assert.strictEqual(treeSet.last(), 2)

  // lower higher floor ceiling
  treeSet.add(1)
  treeSet.add(2)
  treeSet.add(3)
  treeSet.add(4)
  treeSet.add(5)
  assert.strictEqual(treeSet.higher(2), 3)
  assert.strictEqual(treeSet.higher(5), undefined)
  assert.strictEqual(treeSet.lower(2), 1)
  assert.strictEqual(treeSet.lower(1), undefined)
  assert.strictEqual(treeSet.floor(1), 1)
  assert.strictEqual(treeSet.floor(0.9), undefined)
  assert.strictEqual(treeSet.ceiling(5), 5)
  assert.strictEqual(treeSet.ceiling(5.1), undefined)

  // shift pop
  assert.strictEqual(treeSet.shift(), 1)
  assert.strictEqual(treeSet.pop(), 5)
  assert.strictEqual(treeSet.size, 3)

  const treeMultiSet = new TreeMultiSet<number>([], (a: number, b: number) => a - b)

  // add
  treeMultiSet.add(1)
  treeMultiSet.add(1)
  treeMultiSet.add(1)
  treeMultiSet.add(2)
  treeMultiSet.add(2)
  treeMultiSet.add(3)
  treeMultiSet.add(3)
  treeMultiSet.add(4)
  treeMultiSet.add(4)
  assert.strictEqual(treeMultiSet.size, 9)

  // delete
  treeMultiSet.delete(1)
  treeMultiSet.delete(0)
  assert.strictEqual(treeMultiSet.size, 8)

  // keys
  assert.deepStrictEqual([...treeMultiSet.keys()], [1, 1, 2, 2, 3, 3, 4, 4])

  // first last
  assert.strictEqual(treeMultiSet.first(), 1)
  assert.strictEqual(treeMultiSet.last(), 4)

  // lower higher floor ceiling
  assert.strictEqual(treeMultiSet.higher(2), 3)
  assert.strictEqual(treeMultiSet.higher(5), undefined)
  assert.strictEqual(treeMultiSet.lower(2), 1)
  assert.strictEqual(treeMultiSet.lower(1), undefined)
  assert.strictEqual(treeMultiSet.floor(1), 1)
  assert.strictEqual(treeMultiSet.floor(0.9), undefined)
  assert.strictEqual(treeMultiSet.ceiling(4), 4)
  assert.strictEqual(treeMultiSet.ceiling(4.1), undefined)

  // shift pop`
  assert.strictEqual(treeMultiSet.shift(), 1)
  assert.strictEqual(treeMultiSet.pop(), 4)
  assert.strictEqual(treeMultiSet.pop(), 4)
  assert.strictEqual(treeMultiSet.pop(), 3)
  assert.strictEqual(treeMultiSet.size, 4)
}

export { TreeSet, TreeMultiSet }

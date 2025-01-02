/* eslint-disable no-lone-blocks */
/**
 * 返回要遍历的下一个子节点.
 * 如果没有更多子节点, 返回 undefined.
 */
type NextChildFunc<N> = () => N | undefined

// interface INode {
//   isLeaf(): boolean
//   createTraverseFunc(reverse: boolean): TraverseFunc
// }

// interface ITree {
//   root(): INode
//   version(): number
// }

interface ITreeIterator<N> {
  hasNext(): boolean
  next(): N | undefined
}

interface IOperations<T, N> {
  getRoot(tree: T): N
  getVersion(tree: T): number

  isLeaf(node: N): boolean
  getNextChildFunc(node: N, reverse: boolean): NextChildFunc<N>
}

interface ITraverseOptions {
  traverseLeaf: boolean
  traverseNode: boolean
  reverse: boolean
}

class Context<N> {
  private readonly _nextChild: NextChildFunc<N>

  constructor(nextChild: NextChildFunc<N>) {
    this._nextChild = nextChild
  }

  next(): N | undefined {
    return this._nextChild()
  }
}

class State<N> {
  private readonly _stacks: Context<N>[] = []

  push(context: Context<N>): void {
    this._stacks.push(context)
  }

  pop(): void {
    this._stacks.pop()
  }

  current(): Context<N> | undefined {
    return this._stacks.length ? this._stacks[this._stacks.length - 1] : undefined
  }
}

class TreeIterator<T, N> implements ITreeIterator<N> {
  private readonly _tree: T
  private readonly _operations: IOperations<T, N>
  private readonly _options: ITraverseOptions
  private readonly _version: number

  private readonly _state: State<N>
  private _nextNode: N | undefined

  /**
   * @param tree 要遍历的树
   *
   * @param operations 树、节点的操作
   *
   * @param options 遍历选项
   * @param options.traverseLeaf 是否遍历叶子节点. 默认为 true.
   * @param options.traverseNode 是否遍历非叶子节点. 默认为 false.
   * @param options.reverse 是否反向遍历. 默认为 false.
   */
  constructor(tree: T, operations: IOperations<T, N>, options?: Partial<ITraverseOptions>) {
    const defaultOptions: ITraverseOptions = {
      traverseLeaf: true,
      traverseNode: false,
      reverse: false
    }
    const mergedOptions = { ...defaultOptions, ...options }
    const root = operations.getRoot(tree)

    this._tree = tree
    this._operations = operations
    this._options = mergedOptions
    this._version = operations.getVersion(tree)
    this._state = new State()
    this._nextNode = root

    this._state.push(this._createContext(root))
  }

  hasNext(): boolean {
    return !!this._nextNode
  }

  next(): N | undefined {}

  private _peek(): void {
    for (;;) {
      const next = this._next()
      if (!next) return
      if (this._matchesFilter()) return
    }
  }

  private _next(): N | undefined {
    const move = () => {
      for (;;) {
        const context = this._state.current()
        if (!context) {
          this._nextNode = undefined
          return
        }

        const nextNode = context.next()
        if (nextNode) {
          this._nextNode = nextNode
          this._state.push(this._createContext(nextNode))
          return
        }

        this._state.pop()
      }
    }

    if (!this._nextNode) return undefined
    if (this._hasConcurrentModification()) return undefined
    const res = this._nextNode
    move()
    return res
  }

  private _matchesFilter(): boolean {
    if (!this._nextNode) return false
    const isLeaf = this._operations.isLeaf(this._nextNode)
    if (isLeaf && this._options.traverseLeaf) return true
    if (!isLeaf && this._options.traverseNode) return true
    return false
  }

  private _createContext(node: N): Context<N> {
    const nextChildFunc = this._operations.getNextChildFunc(node, !!this._options.reverse)
    return new Context(nextChildFunc)
  }

  private _hasConcurrentModification(): boolean {
    const treeVersion = this._operations.getVersion(this._tree)
    return this._version !== treeVersion
  }
}

export type { ITreeIterator }
export { TreeIterator }

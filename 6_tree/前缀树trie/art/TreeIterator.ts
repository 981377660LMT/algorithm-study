/* eslint-disable no-lone-blocks */

/**
 * 返回要遍历的下一个子节点.
 * 如果没有更多子节点, 返回 undefined.
 */
type NextChildFunc<N> = () => N | undefined

interface ITreeIterator<N> {
  hasNext(): boolean
  next(): N | undefined
}

interface IOperations<N> {
  getNextChildFunc(node: N, reverse: boolean): NextChildFunc<N>
}

interface ITraverseOptions<N> {
  reverse: boolean
  filter: (node: N) => boolean
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

class TreeIterator<N> implements ITreeIterator<N> {
  private readonly _operations: IOperations<N>
  private readonly _options: ITraverseOptions<N>

  private readonly _state: State<N>
  private _nextNode: N | undefined

  /**
   * @param root 树的根节点
   *
   * @param operations 树、节点的操作
   *
   * @param options 遍历选项
   * @param options.reverse 是否反向遍历. 默认为 false.
   * @param options.filter 过滤器. 默认为始终返回 true 的函数.
   */
  constructor(root: N, operations: IOperations<N>, options?: Partial<ITraverseOptions<N>>) {
    const defaultOptions: ITraverseOptions<N> = {
      reverse: false,
      filter: () => true
    }
    const mergedOptions = { ...defaultOptions, ...options }

    this._operations = operations
    this._options = mergedOptions
    this._state = new State()
    this._nextNode = root

    this._state.push(this._createContext(root))
    if (!this._matchesFilter()) this._moveUntilMatch()
  }

  hasNext(): boolean {
    return !!this._nextNode
  }

  next(): N | undefined {
    const res = this._nextNode
    this._moveUntilMatch()
    return res
  }

  private _moveUntilMatch(): void {
    for (;;) {
      if (!this._nextNode) return
      this._move()
      if (this._matchesFilter()) return
    }
  }

  private _move(): void {
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

  private _matchesFilter(): boolean {
    if (!this._nextNode) return false
    return this._options.filter(this._nextNode)
  }

  private _createContext(node: N): Context<N> {
    const nextChildFunc = this._operations.getNextChildFunc(node, this._options.reverse)
    return new Context(nextChildFunc)
  }
}

export type { ITreeIterator }
export { TreeIterator }

if (require.main === module) {
  class Node {
    // eslint-disable-next-line no-useless-constructor
    constructor(public value: number, public children: Node[] = []) {}
  }

  class Operations implements IOperations<Node> {
    getNextChildFunc(node: Node, reverse: boolean): NextChildFunc<Node> {
      let index = reverse ? node.children.length : -1
      return () => {
        index += reverse ? -1 : 1
        if (index < 0 || index >= node.children.length) return undefined
        return node.children[index]
      }
    }
  }

  //      1
  //     / \
  //    2   5
  //   / \
  //  3   4
  const root = new Node(1, [new Node(2, [new Node(3), new Node(4)]), new Node(5)])
  const isLeaf = (node: Node) => !node.children.length
  const iter = new TreeIterator(root, new Operations(), { reverse: true, filter: isLeaf })
  while (iter.hasNext()) {
    const node = iter.next()!
    console.log(node.value)
  }
}

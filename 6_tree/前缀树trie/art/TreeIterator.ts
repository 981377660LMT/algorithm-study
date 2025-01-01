interface INode {
  allChildren(): INode[]
}

interface ITree {
  root(): INode
  version(): number
}

interface ITraverseOptions {
  traverseLeaf?: boolean
  traverseNode?: boolean
  reverse?: boolean
}

/** 返回要遍历的下一个子节点的索引. */
type TraverseFunc = () => number | undefined

class State {
  private readonly _stacks: IteratorContext[] = []

  push(context: IteratorContext): void {
    this._stacks.push(context)
  }

  pop(): void {
    this._stacks.pop()
  }

  current(): IteratorContext | undefined {
    return this._stacks.length ? this._stacks[this._stacks.length - 1] : undefined
  }
}

class IteratorContext {
  private readonly _nextChildFn: TraverseFunc
  private readonly _children: INode[]

  constructor(node: INode, reverse: boolean) {
    this._nextChildFn = createTraverseFunc(node, reverse)
    this._children = node.allChildren()
  }

  next(): INode | undefined {
    for (;;) {
      const index = this._nextChildFn()
      if (index === undefined) return undefined
      const child = this._children[index]
      if (child) return child
    }
  }
}

class TreeIterator {
  private readonly _tree: ITree
  private readonly _state: State
  private readonly _reverse: boolean
  private _nextNode: INode | undefined
  private _version = 0

  constructor(tree: ITree, reverse = false) {
    const state = new State()
    state.push(new IteratorContext(tree.root(), reverse))
    this._tree = tree
    this._state = state
    this._reverse = reverse
    this._nextNode = tree.root()
    this._version = tree.version()
  }

  hasNext(): boolean {
    return !!this._nextNode
  }

  next(): INode | undefined {
    if (!this._nextNode) return undefined
    if (this._hasConcurrentModification()) return undefined
    const res = this._nextNode
    this._next()
    return res
  }

  private _next(): void {
    for (;;) {
      const context = this._state.current()
      if (!context) {
        this._nextNode = undefined
        return
      }

      const nextNode = context.next()
      if (nextNode) {
        this._nextNode = nextNode
        this._state.push(new IteratorContext(nextNode, this._reverse))
        return
      }

      this._state.pop()
    }
  }

  private _hasConcurrentModification(): boolean {
    return this._tree.version() !== this._version
  }
}

function createTraverseFunc(node: INode, reverse: boolean): TraverseFunc {
  throw new Error('TODO')
}

// TreeManager4: 扁平树结构

// 参考 TreeManager2.
export class FlattedTree<TNode> {
  private _rootId: string
  private readonly _getId: (node: TNode) => string
  private readonly _getChildrenIds: (node: TNode) => string[]

  private readonly _idToNode = new Map<string, TNode>()
  private readonly _idToParentId = new Map<string, string>()

  constructor(
    rootId: string,
    nodes: TNode[],
    operations: {
      getId: (node: TNode) => string
      getChildrenIds: (node: TNode) => string[]
    }
  ) {
    this._rootId = rootId
    this._getId = operations.getId
    this._getChildrenIds = operations.getChildrenIds

    nodes.forEach(node => {
      const id = this._getId(node)
      this._idToNode.set(id, node)
    })
  }

  dispose(): void {
    this._idToNode.clear()
  }

  /**
   * 将节点追加到父节点的子节点列表末尾。
   */
  append(node: TNode, parentId: string): void {
    this._insert(node, parentId, childrenIds => {
      childrenIds.push(this._getId(node))
    })
  }

  /**
   * 将节点添加到父节点的子节点列表开头。
   */
  prepend(node: TNode, parentId: string): void {
    this._insert(node, parentId, childrenIds => {
      childrenIds.unshift(this._getId(node))
    })
  }

  /**
   * 在目标节点之前插入新节点。
   */
  insertBefore(node: TNode, targetId: string): void {
    const parentId = this._idToParentId.get(targetId)
    if (!parentId) {
      throw new Error(`Target node ${targetId} has no parent or does not exist`)
    }

    this._insert(node, parentId, childrenIds => {
      const targetIndex = childrenIds.indexOf(targetId)
      if (targetIndex === -1) {
        throw new Error(`Target node ${targetId} not found in parent's children`)
      }
      childrenIds.splice(targetIndex, 0, this._getId(node))
    })
  }

  /**
   * 在目标节点之后插入新节点。
   */
  insertAfter(node: TNode, targetId: string): void {
    const parentId = this._idToParentId.get(targetId)
    if (!parentId) {
      throw new Error(`Target node ${targetId} has no parent or does not exist`)
    }

    this._insert(node, parentId, childrenIds => {
      const targetIndex = childrenIds.indexOf(targetId)
      if (targetIndex === -1) {
        throw new Error(`Target node ${targetId} not found in parent's children`)
      }
      childrenIds.splice(targetIndex + 1, 0, this._getId(node))
    })
  }

  /**
   * 移除节点及其所有子节点.
   */
  remove(id: string): void {
    if (this._isRoot(id)) {
      throw new Error('Cannot remove root node')
    }
    if (!this.has(id)) {
      throw new Error(`Node ${id} does not exist`)
    }

    const parentId = this._idToParentId.get(id)
    if (parentId) {
      const siblings = this._getChildrenIds(this._idToNode.get(parentId)!)
      const index = siblings.indexOf(id)
      if (index !== -1) {
        siblings.splice(index, 1)
      }
    }

    const nodesToRemove: string[] = []
    const dfs = (cur: string) => {
      nodesToRemove.push(cur)
      const children = this._getChildrenIds(this._idToNode.get(cur)!)
      for (const childId of children) {
        dfs(childId)
      }
    }
    dfs(id)
  }

  /**
   * 将节点移动到目标节点之前。
   */
  moveBefore(nodeId: string, targetId: string): void {
    this._move(nodeId, targetId, true)
  }

  /**
   * 将节点移动到目标节点之后。
   */
  moveAfter(nodeId: string, targetId: string): void {
    this._move(nodeId, targetId, false)
  }

  getNodes(): TNode[] {
    return [...this._idToNode.values()]
  }

  get(id: string): TNode | undefined {
    return this._idToNode.get(id)
  }

  getRoot(): TNode | undefined {
    return this._idToNode.get(this._rootId)
  }

  getParent(id: string): TNode | undefined {
    const parentId = this._idToParentId.get(id)
    return parentId ? this._idToNode.get(parentId) : undefined
  }

  getChildren(id: string): TNode[] {
    const node = this._idToNode.get(id)
    if (!node) return []
    const childrenIds = this._getChildrenIds(node)
    return childrenIds.map(childId => this._idToNode.get(childId)).filter(Boolean) as TNode[]
  }

  getSize(): number {
    return this._idToNode.size
  }

  has(id: string): boolean {
    return this._idToNode.has(id)
  }

  /**
   * 前序遍历树。
   * @param f 遍历回调函数，返回true时停止遍历。
   */
  enumerate(f: (node: TNode, path: TNode[]) => boolean | void): void {
    const root = this.getRoot()
    if (!root) return

    const dfs = (cur: TNode, path: TNode[]): boolean => {
      if (f(cur, path)) {
        return true
      }
      path.push(cur)
      for (const child of this.getChildren(this._getId(cur))) {
        if (dfs(child, path)) {
          path.pop()
          return true
        }
      }
      path.pop()
      return false
    }
    dfs(root, [])
  }

  print(options?: { format?: (node: TNode) => string; output?: (message: string) => void }): void {
    const format = options?.format || ((node: TNode) => this._getId(node))
    const output = options?.output || console.log
    this.enumerate((cur, path) => {
      const indent = '  '.repeat(path.length)
      const message = `${indent}- ${format(cur)}`
      output(message)
    })
  }

  private _insert(node: TNode, parentId: string, f: (childrenIds: string[]) => void): void {
    const nodeId = this._getId(node)
    if (this.has(nodeId)) {
      throw new Error(`Node with ID ${nodeId} already exists`)
    }
    const parent = this._idToNode.get(parentId)
    if (!parent) {
      throw new Error(`Parent node ${parentId} does not exist`)
    }

    this._idToNode.set(nodeId, node)

    const childrenIds = this._getChildrenIds(parent)
    f(childrenIds)
  }

  private _move(nodeId: string, targetId: string, before: boolean): void {
    // 基本验证
    if (nodeId === targetId) {
      throw new Error('Cannot move node to itself')
    }
    if (!this.has(nodeId) || !this.has(targetId)) {
      throw new Error('Source or target node does not exist')
    }
    // 防止将节点移动到自己的子孙节点
    if (this._isAncestor(nodeId, targetId)) {
      throw new Error('Cannot move node to its descendant')
    }

    const targetParentId = this._idToParentId.get(targetId)
    if (!targetParentId) {
      throw new Error(`Target node ${targetId} has no parent`)
    }

    // 从原位置移除
    this._removeFromParent(nodeId)

    // 插入到新位置
    const targetParent = this._idToNode.get(targetParentId)!
    const childrenIds = this._getChildrenIds(targetParent).slice()
    const targetIndex = childrenIds.indexOf(targetId)

    if (targetIndex === -1) {
      throw new Error(`Target node ${targetId} not found in parent's children`)
    }

    const insertIndex = before ? targetIndex : targetIndex + 1
    childrenIds.splice(insertIndex, 0, nodeId)

    // 更新父子关系
    this._idToParentId.set(nodeId, targetParentId)
    this._updateNodeChildren(targetParent, childrenIds)
  }

  private _isAncestor(ancestorId: string, nodeId: string): boolean {
    let curId = this._idToParentId.get(nodeId)
    while (curId) {
      if (curId === ancestorId) {
        return true
      }
      curId = this._idToParentId.get(curId)
    }
    return false
  }

  private _isRoot(id: string): boolean {
    return this._rootId === id
  }
}

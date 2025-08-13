// insert
// remove
// moveBefore
// moveAfter
// getParent
// enumerate

// TODO：

/**
 * 基于id的树形结构管理工具.
 *
 * 设计改进建议：
 * 1. 移除 getChildren 回调，内部维护完整的树结构
 * 2. 添加必要的数据验证和错误处理
 * 3. 提供更丰富的查询方法
 */
class TreeManager2<Id = string> {
  readonly rootId: Id
  private readonly parentMap = new Map<Id, Id>()
  private readonly childrenMap = new Map<Id, Id[]>()

  constructor(rootId: Id) {
    this.rootId = rootId
    this.childrenMap.set(rootId, [])
  }

  // 获取子节点（替代外部回调）
  getChildren(nodeId: Id): Id[] {
    return this.childrenMap.get(nodeId) || []
  }

  insert(nodeId: Id, parentId: Id): void {
    if (this.parentMap.has(nodeId) || nodeId === this.rootId) {
      throw new Error(`Node ${nodeId} already exists`)
    }

    if (!this.childrenMap.has(parentId)) {
      throw new Error(`Parent node ${parentId} does not exist`)
    }

    this.parentMap.set(nodeId, parentId)
    this.childrenMap.set(nodeId, [])

    const siblings = this.childrenMap.get(parentId)!
    siblings.push(nodeId)
  }

  remove(nodeId: Id): void {
    if (nodeId === this.rootId) {
      throw new Error('Cannot remove root node')
    }

    const parentId = this.parentMap.get(nodeId)
    if (!parentId) {
      throw new Error(`Node ${nodeId} does not exist`)
    }

    // 递归删除子节点
    const children = [...(this.childrenMap.get(nodeId) || [])]
    children.forEach(childId => this.remove(childId))

    // 从父节点移除
    const siblings = this.childrenMap.get(parentId)!
    const index = siblings.indexOf(nodeId)
    if (index >= 0) siblings.splice(index, 1)

    // 清理映射
    this.parentMap.delete(nodeId)
    this.childrenMap.delete(nodeId)
  }

  moveBefore(nodeId: Id, referenceId: Id): void {
    this._moveNode(nodeId, referenceId, 'before')
  }

  moveAfter(nodeId: Id, referenceId: Id): void {
    this._moveNode(nodeId, referenceId, 'after')
  }

  private _moveNode(nodeId: Id, referenceId: Id, position: 'before' | 'after'): void {
    if (nodeId === this.rootId || nodeId === referenceId) {
      throw new Error('Invalid move operation')
    }

    const nodeParent = this.parentMap.get(nodeId)
    const refParent = this.parentMap.get(referenceId)

    if (!nodeParent || !refParent) {
      throw new Error('Node not found')
    }

    if (nodeParent !== refParent) {
      throw new Error('Nodes must have same parent')
    }

    const siblings = this.childrenMap.get(nodeParent)!
    const nodeIndex = siblings.indexOf(nodeId)
    const refIndex = siblings.indexOf(referenceId)

    // 移除节点
    if (nodeIndex >= 0) siblings.splice(nodeIndex, 1)

    // 重新插入
    const newRefIndex = siblings.indexOf(referenceId)
    const insertIndex = position === 'before' ? newRefIndex : newRefIndex + 1
    siblings.splice(insertIndex, 0, nodeId)
  }

  getParent(nodeId: Id): Id | undefined {
    return this.parentMap.get(nodeId)
  }

  enumerate(f: (nodeId: Id, path: number[]) => boolean | void): void {
    const dfs = (nodeId: Id, path: number[]): boolean => {
      if (f(nodeId, path) === false) return false

      const children = this.getChildren(nodeId)
      for (let i = 0; i < children.length; i++) {
        if (dfs(children[i], [...path, i]) === false) return false
      }
      return true
    }

    dfs(this.rootId, [])
  }

  // 额外的实用方法
  exists(nodeId: Id): boolean {
    return nodeId === this.rootId || this.parentMap.has(nodeId)
  }

  getDepth(nodeId: Id): number {
    let depth = 0
    let current = nodeId
    while (current !== this.rootId) {
      const parent = this.parentMap.get(current)
      if (!parent) return -1
      current = parent
      depth++
    }
    return depth
  }

  isAncestor(ancestorId: Id, nodeId: Id): boolean {
    let current = nodeId
    while (current !== this.rootId && current !== undefined) {
      const parent = this.parentMap.get(current)
      if (parent === ancestorId) return true
      current = parent!
    }
    return ancestorId === this.rootId
  }
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  // 测试代码
  const tree = new TreeManager2('root')

  tree.insert('a', 'root')
  tree.insert('b', 'root')
  tree.insert('c', 'a')

  console.log('Root children:', tree.getChildren('root')) // ['a', 'b']
  console.log('A children:', tree.getChildren('a')) // ['c']
  console.log('C parent:', tree.getParent('c')) // 'a'

  tree.enumerate((id, path) => {
    console.log(`Node: ${id}, Path: [${path.join(', ')}]`)
  })
}
export { TreeManager2 }

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
}

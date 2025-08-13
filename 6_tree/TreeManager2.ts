/* eslint-disable no-lone-blocks */
/* eslint-disable no-inner-declarations */
/**
 * 基于id的树形结构管理工具.
 */
class TreeManager2<TNode extends { id: string }> {
  private _root: TNode | undefined
  private readonly _idToNode = new Map<string, TNode>()

  private readonly _idToParentId = new Map<string, string>()
  private readonly _idToChildrenId = new Map<string, string[]>()

  constructor(root: TNode) {
    this._root = root
    this._idToNode.set(root.id, root)
    this._idToChildrenId.set(root.id, [])
  }

  dispose(): void {
    this._root = undefined
    this._idToNode.clear()
    this._idToParentId.clear()
    this._idToChildrenId.clear()
  }

  append(node: TNode, parentId: string): void {
    this._insert(node, parentId, children => {
      children.push(node.id)
    })
  }

  prepend(node: TNode, parentId: string): void {
    this._insert(node, parentId, children => {
      children.unshift(node.id)
    })
  }

  insertBefore(node: TNode, targetId: string): void {
    const parentId = this._idToParentId.get(targetId)
    if (!parentId) {
      throw new Error('Cannot insert before root node')
    }

    this._insert(node, parentId, children => {
      const index = children.indexOf(targetId)
      if (index === -1) {
        throw new Error(`Reference node ${targetId} not found in parent's children`)
      }
      children.splice(index, 0, node.id)
    })
  }

  insertAfter(node: TNode, targetId: string): void {
    const parentId = this._idToParentId.get(targetId)
    if (!parentId) {
      throw new Error('Cannot insert after root node')
    }

    this._insert(node, parentId, children => {
      const index = children.indexOf(targetId)
      if (index === -1) {
        throw new Error(`Reference node ${targetId} not found in parent's children`)
      }
      children.splice(index + 1, 0, node.id)
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
      const siblings = this._idToChildrenId.get(parentId) || []
      const index = siblings.indexOf(id)
      if (index !== -1) {
        siblings.splice(index, 1)
      }
    }

    const nodesToRemove: string[] = []
    const dfs = (cur: string) => {
      nodesToRemove.push(cur)
      const children = this._idToChildrenId.get(cur) || []
      for (const childId of children) {
        dfs(childId)
      }
    }
    dfs(id)
    for (const v of nodesToRemove) {
      this._idToNode.delete(v)
      this._idToParentId.delete(v)
      this._idToChildrenId.delete(v)
    }
  }

  replace(id: string, newNode: TNode): void {
    if (!this.has(id)) {
      throw new Error(`Node ${id} does not exist`)
    }

    if (id !== newNode.id && this.has(newNode.id)) {
      throw new Error(`New node ${newNode.id} already exists`)
    }

    const isReplaceRoot = id === this._root?.id

    if (id === newNode.id) {
      this._idToNode.set(id, newNode)
      if (isReplaceRoot) {
        this._root = newNode
      }
      return
    }

    const parentId = this._idToParentId.get(id)
    const children = this._idToChildrenId.get(id) || []

    for (const childId of children) {
      this._idToParentId.set(childId, newNode.id)
    }

    if (parentId) {
      const siblings = this._idToChildrenId.get(parentId) || []
      const index = siblings.indexOf(id)
      if (index !== -1) {
        siblings[index] = newNode.id
      }
    }

    this._idToNode.delete(id)
    this._idToParentId.delete(id)
    this._idToChildrenId.delete(id)

    this._idToNode.set(newNode.id, newNode)
    if (parentId) {
      this._idToParentId.set(newNode.id, parentId)
    }
    this._idToChildrenId.set(newNode.id, children)

    if (isReplaceRoot) {
      this._root = newNode
    }
  }

  moveBefore(nodeId: string, targetId: string): void {
    this._move(nodeId, targetId, true)
  }

  moveAfter(nodeId: string, targetId: string): void {
    this._move(nodeId, targetId, false)
  }

  get(id: string): TNode | undefined {
    return this._idToNode.get(id)
  }

  getParent(id: string): TNode | undefined {
    const parentId = this._idToParentId.get(id)
    return parentId ? this._idToNode.get(parentId) : undefined
  }

  getChildren(id: string): TNode[] {
    const res: TNode[] = []
    for (const childId of this._idToChildrenId.get(id) || []) {
      const child = this._idToNode.get(childId)
      if (child) {
        res.push(child)
      }
    }
    return res
  }

  getRoot(): TNode | undefined {
    return this._root
  }

  previousSibling(id: string): TNode | undefined {
    const parentId = this._idToParentId.get(id)
    if (!parentId) return undefined
    const siblings = this._idToChildrenId.get(parentId) || []
    const index = siblings.indexOf(id)
    if (index <= 0) return undefined
    const prevSiblingId = siblings[index - 1]
    return this._idToNode.get(prevSiblingId)
  }

  nextSibling(id: string): TNode | undefined {
    const parentId = this._idToParentId.get(id)
    if (!parentId) return undefined
    const siblings = this._idToChildrenId.get(parentId) || []
    const index = siblings.indexOf(id)
    if (index === -1 || index === siblings.length - 1) return undefined
    const nextSiblingId = siblings[index + 1]
    return this._idToNode.get(nextSiblingId)
  }

  has(id: string): boolean {
    return this._idToNode.has(id)
  }

  enumerate(f: (id: string, path: string[]) => boolean | void): void {
    const dfs = (id: string, path: string[]): boolean | void => {
      if (f(id, path)) {
        return true
      }
      const children = this._idToChildrenId.get(id) || []
      for (const childId of children) {
        if (dfs(childId, [...path, id])) {
          return true
        }
      }
      return false
    }

    if (this._root) {
      dfs(this._root.id, [])
    }
  }

  print(options?: {
    toString?: (node: TNode) => string
    output?: (message: string) => void
  }): void {
    const { output = console.log, toString = (node: TNode) => node.id } = options || {}
    this.enumerate((id, path) => {
      const node = this._idToNode.get(id)
      if (node) {
        const message = `${' '.repeat(path.length * 2)}- ${toString(node)}`
        output(message)
      }
    })
  }

  get size(): number {
    return this._idToNode.size
  }

  private _insert(node: TNode, parentId: string, f: (children: string[]) => void): void {
    if (this.has(node.id)) {
      throw new Error(`Node ${node.id} already exists`)
    }

    if (!this.has(parentId)) {
      throw new Error(`Parent node ${parentId} does not exist`)
    }

    this._idToNode.set(node.id, node)
    this._idToParentId.set(node.id, parentId)
    this._idToChildrenId.set(node.id, [])

    const children = this._idToChildrenId.get(parentId) || []
    f(children)
    this._idToChildrenId.set(parentId, children)
  }

  private _move(nodeId: string, targetId: string, before: boolean): void {
    if (this._isRoot(nodeId) || this._isRoot(targetId)) {
      throw new Error('Cannot move root node')
    }
    if (!this.has(nodeId) || !this.has(targetId)) {
      throw new Error('Node does not exist')
    }
    if (nodeId === targetId) {
      throw new Error('Cannot move node to itself')
    }
    if (this._isAncestor(nodeId, targetId)) {
      throw new Error('Cannot move node to its descendant')
    }

    const newParentId = this._idToParentId.get(targetId)
    if (!newParentId) {
      throw new Error(`Target node ${targetId} does not have a parent`)
    }

    const oldParentId = this._idToParentId.get(nodeId)
    if (oldParentId) {
      const oldSiblings = this._idToChildrenId.get(oldParentId) || []
      const oldIndex = oldSiblings.indexOf(nodeId)
      if (oldIndex !== -1) {
        oldSiblings.splice(oldIndex, 1)
      }
    }

    const newSiblings = this._idToChildrenId.get(newParentId) || []
    const targetIndex = newSiblings.indexOf(targetId)
    const insertIndex = before ? targetIndex : targetIndex + 1
    newSiblings.splice(insertIndex, 0, nodeId)
    this._idToChildrenId.set(newParentId, newSiblings)

    this._idToParentId.set(nodeId, newParentId)
  }

  private _isAncestor(ancestorId: string, nodeId: string): boolean {
    let curId: string | undefined = nodeId
    while (curId) {
      if (curId === ancestorId) {
        return true
      }
      curId = this._idToParentId.get(curId)
    }
    return false
  }

  private _isRoot(id: string): boolean {
    return this._root?.id === id
  }
}

if (typeof require !== 'undefined') {
  function assert(condition: boolean, message: string): void {
    if (!condition) {
      throw new Error(`断言失败: ${message}`)
    }
  }

  function arrayEquals<T>(arr1: T[], arr2: T[]): boolean {
    return arr1.length === arr2.length && arr1.every((val, i) => val === arr2[i])
  }

  function runTests(): void {
    console.log('开始运行 TreeManager2 测试...')

    // 1. 基础树操作测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })

      // 添加节点测试
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'root')
      tree.prepend({ id: 'c', name: 'Node C' }, 'root')
      tree.append({ id: 'd', name: 'Node D' }, 'a')
      tree.insertBefore({ id: 'e', name: 'Node E' }, 'b')
      tree.insertAfter({ id: 'f', name: 'Node F' }, 'a')

      // 验证树结构
      const rootChildren = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(rootChildren, ['c', 'a', 'f', 'e', 'b']),
        `根节点子节点顺序错误，期望: ['c', 'a', 'f', 'e', 'b']，实际: ${JSON.stringify(
          rootChildren
        )}`
      )

      const aChildren = tree.getChildren('a').map(n => n.id)
      assert(
        arrayEquals(aChildren, ['d']),
        `节点a的子节点错误，期望: ['d']，实际: ${JSON.stringify(aChildren)}`
      )

      tree.dispose()
    }

    // 2. 节点查询测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'a')

      assert(tree.get('a')?.name === 'Node A', '获取节点a失败')
      assert(tree.getParent('a')?.id === 'root', 'a的父节点应该是root')
      assert(tree.getParent('b')?.id === 'a', 'b的父节点应该是a')
      assert(tree.getParent('root') === undefined, 'root不应该有父节点')
      assert(tree.getRoot()?.id === 'root', '根节点应该是root')

      tree.dispose()
    }

    // 3. 节点关系测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'root')
      tree.append({ id: 'c', name: 'Node C' }, 'root')

      assert(tree.previousSibling('b')?.id === 'a', 'b的前一个兄弟节点应该是a')
      assert(tree.nextSibling('b')?.id === 'c', 'b的后一个兄弟节点应该是c')
      assert(tree.previousSibling('a') === undefined, 'a不应该有前一个兄弟节点')
      assert(tree.nextSibling('c') === undefined, 'c不应该有后一个兄弟节点')

      tree.dispose()
    }

    // 4. 移动节点测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'root')
      tree.append({ id: 'c', name: 'Node C' }, 'root')

      tree.moveBefore('c', 'a')
      let children = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(children, ['c', 'a', 'b']),
        `moveBefore后顺序错误，期望: ['c', 'a', 'b']，实际: ${JSON.stringify(children)}`
      )

      tree.moveAfter('c', 'b')
      children = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(children, ['a', 'b', 'c']),
        `moveAfter后顺序错误，期望: ['a', 'b', 'c']，实际: ${JSON.stringify(children)}`
      )

      tree.dispose()
    }

    // 5. 删除节点测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'a')
      tree.append({ id: 'c', name: 'Node C' }, 'a')

      tree.remove('a') // 删除a及其子节点
      assert(!tree.has('a'), '节点a应该被删除')
      assert(!tree.has('b'), '节点b应该被删除')
      assert(!tree.has('c'), '节点c应该被删除')
      assert(tree.getChildren('root').length === 0, 'root应该没有子节点')

      tree.dispose()
    }

    // 6. 替换节点测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'a')

      // 同ID替换
      tree.replace('a', { id: 'a', name: 'New Node A' })
      assert(tree.get('a')?.name === 'New Node A', '同ID替换失败')
      assert(tree.getChildren('a').length === 1, '替换后子节点应该保持')

      // 不同ID替换
      tree.replace('a', { id: 'x', name: 'Node X' })
      assert(!tree.has('a'), '原节点a应该不存在')
      assert(tree.has('x'), '新节点x应该存在')
      assert(tree.getParent('b')?.id === 'x', 'b的父节点应该变为x')

      tree.dispose()
    }

    // 7. 错误情况测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')

      // 重复添加节点
      try {
        tree.append({ id: 'a', name: 'Duplicate A' }, 'root')
        throw new Error('应该抛出重复节点错误')
      } catch (e: any) {
        assert(e.message.includes('already exists'), '应该抛出节点已存在错误')
      }

      // 删除根节点
      try {
        tree.remove('root')
        throw new Error('应该抛出不能删除根节点错误')
      } catch (e: any) {
        assert(e.message.includes('Cannot remove root node'), '应该抛出不能删除根节点错误')
      }

      // 移动根节点
      try {
        tree.moveBefore('root', 'a')
        throw new Error('应该抛出不能移动根节点错误')
      } catch (e: any) {
        assert(e.message.includes('Cannot move root node'), '应该抛出不能移动根节点错误')
      }

      // 在不存在的父节点下添加节点
      try {
        tree.append({ id: 'b', name: 'Node B' }, 'nonexistent')
        throw new Error('应该抛出父节点不存在错误')
      } catch (e: any) {
        assert(e.message.includes('does not exist'), '应该抛出父节点不存在错误')
      }

      tree.dispose()
    }

    // 8. 枚举测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'root')
      tree.append({ id: 'c', name: 'Node C' }, 'a')

      const visitedNodes: string[] = []
      tree.enumerate((id, path) => {
        visitedNodes.push(id)
      })

      assert(
        arrayEquals(visitedNodes, ['root', 'a', 'c', 'b']),
        `枚举顺序错误，期望: ['root', 'a', 'c', 'b']，实际: ${JSON.stringify(visitedNodes)}`
      )

      tree.dispose()
    }

    // 9. has方法测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')

      assert(tree.has('root'), '应该包含root节点')
      assert(tree.has('a'), '应该包含a节点')
      assert(!tree.has('nonexistent'), '不应该包含不存在的节点')

      tree.dispose()
    }

    // 10. 根节点替换测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })
      tree.append({ id: 'a', name: 'Node A' }, 'root')

      // 替换根节点（同ID）
      tree.replace('root', { id: 'root', name: 'New Root' })
      assert(tree.getRoot()?.name === 'New Root', '根节点同ID替换失败')

      // 替换根节点（不同ID）
      tree.replace('root', { id: 'newroot', name: 'New Root Node' })
      assert(tree.getRoot()?.id === 'newroot', '根节点不同ID替换失败')
      assert(tree.getParent('a')?.id === 'newroot', 'a的父节点应该更新为新根节点')

      tree.dispose()
    }

    // 4. 复杂移动节点测试
    {
      const tree = new TreeManager2({ id: 'root', name: 'Root Node' })

      // 构建复杂树结构
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'root')
      tree.append({ id: 'c', name: 'Node C' }, 'root')
      tree.append({ id: 'a1', name: 'Node A1' }, 'a')
      tree.append({ id: 'a2', name: 'Node A2' }, 'a')
      tree.append({ id: 'b1', name: 'Node B1' }, 'b')
      tree.append({ id: 'b2', name: 'Node B2' }, 'b')
      tree.append({ id: 'c1', name: 'Node C1' }, 'c')
      tree.append({ id: 'a1a', name: 'Node A1A' }, 'a1')
      tree.append({ id: 'a1b', name: 'Node A1B' }, 'a1')

      // 4.1 同级节点移动测试
      tree.moveBefore('c', 'a')
      let rootChildren = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(rootChildren, ['c', 'a', 'b']),
        `同级moveBefore失败，期望: ['c', 'a', 'b']，实际: ${JSON.stringify(rootChildren)}`
      )

      tree.moveAfter('c', 'b')
      rootChildren = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(rootChildren, ['a', 'b', 'c']),
        `同级moveAfter失败，期望: ['a', 'b', 'c']，实际: ${JSON.stringify(rootChildren)}`
      )

      // 4.2 跨层级移动测试 - 从深层移动到浅层
      tree.moveBefore('a1a', 'b')
      rootChildren = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(rootChildren, ['a', 'a1a', 'b', 'c']),
        `跨层级moveBefore失败，期望: ['a', 'a1a', 'b', 'c']，实际: ${JSON.stringify(rootChildren)}`
      )

      let a1Children = tree.getChildren('a1').map(n => n.id)
      assert(
        arrayEquals(a1Children, ['a1b']),
        `a1a移动后a1子节点错误，期望: ['a1b']，实际: ${JSON.stringify(a1Children)}`
      )

      // 4.3 跨层级移动测试 - 从浅层移动到深层
      tree.moveAfter('a1a', 'a2')
      let aChildren = tree.getChildren('a').map(n => n.id)
      assert(
        arrayEquals(aChildren, ['a1', 'a2', 'a1a']),
        `a1a移动到a下失败，期望: ['a1', 'a2', 'a1a']，实际: ${JSON.stringify(aChildren)}`
      )

      rootChildren = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(rootChildren, ['a', 'b', 'c']),
        `a1a移动后root子节点，期望: ['a', 'b', 'c']，实际: ${JSON.stringify(rootChildren)}`
      )

      // 4.4 带子树的节点移动
      tree.moveBefore('a1', 'c1')
      let cChildren = tree.getChildren('c').map(n => n.id)
      assert(
        arrayEquals(cChildren, ['a1', 'c1']),
        `带子树移动失败，期望: ['a1', 'c1']，实际: ${JSON.stringify(cChildren)}`
      )

      // 验证a1的子节点依然存在
      a1Children = tree.getChildren('a1').map(n => n.id)
      assert(
        arrayEquals(a1Children, ['a1b']),
        `移动后a1子节点丢失，期望: ['a1b']，实际: ${JSON.stringify(a1Children)}`
      )

      // 4.5 复杂位置移动 - 移动到兄弟节点的最后
      tree.append({ id: 'b3', name: 'Node B3' }, 'b')
      tree.moveAfter('b1', 'b3')
      let bChildren = tree.getChildren('b').map(n => n.id)
      assert(
        arrayEquals(bChildren, ['b2', 'b3', 'b1']),
        `移动到最后位置失败，期望: ['b2', 'b3', 'b1']，实际: ${JSON.stringify(bChildren)}`
      )

      // 4.6 移动错误情况测试

      // 尝试移动节点到其后代节点
      try {
        tree.moveBefore('a', 'a1a')
        throw new Error('应该抛出不能移动到后代节点错误')
      } catch (e: any) {
        assert(
          e.message.includes('descendant') || e.message.includes('后代'),
          '应该抛出不能移动到后代节点错误'
        )
      }

      // 尝试移动节点到自己
      try {
        tree.moveBefore('a', 'a')
        throw new Error('应该抛出不能移动到自己错误')
      } catch (e: any) {
        assert(
          e.message.includes('itself') || e.message.includes('自己'),
          '应该抛出不能移动到自己错误'
        )
      }

      // 尝试移动根节点
      try {
        tree.moveBefore('root', 'a')
        throw new Error('应该抛出不能移动根节点错误')
      } catch (e: any) {
        assert(e.message.includes('root') || e.message.includes('根'), '应该抛出不能移动根节点错误')
      }

      // 尝试移动不存在的节点
      try {
        tree.moveBefore('nonexistent', 'a')
        throw new Error('应该抛出节点不存在错误')
      } catch (e: any) {
        assert(
          e.message.includes('does not exist') || e.message.includes('不存在'),
          '应该抛出节点不存在错误'
        )
      }

      // 4.7 边界情况测试 - 移动到第一个位置
      tree.moveBefore('b2', 'b3')
      bChildren = tree.getChildren('b').map(n => n.id)
      assert(
        arrayEquals(bChildren, ['b2', 'b3', 'b1']),
        `移动到第一位置失败，期望: ['b2', 'b3', 'b1']，实际: ${JSON.stringify(bChildren)}`
      )

      // 4.8 验证移动后的父子关系
      assert(tree.getParent('a1')?.id === 'c', 'a1的父节点应该是c')
      assert(tree.getParent('a1b')?.id === 'a1', 'a1b的父节点应该还是a1')
      assert(tree.getParent('b1')?.id === 'b', 'b1的父节点应该是b')

      // 4.9 多次连续移动测试
      tree.moveBefore('b1', 'b2')
      tree.moveAfter('b3', 'b1')
      tree.moveBefore('b2', 'b1')
      bChildren = tree.getChildren('b').map(n => n.id)
      assert(
        arrayEquals(bChildren, ['b2', 'b1', 'b3']),
        `多次移动后顺序错误，期望: ['b2', 'b1', 'b3']，实际: ${JSON.stringify(bChildren)}`
      )

      tree.dispose()
    }

    // 6. 复杂替换节点测试
    {
      const tree = new TreeManager2<any>({ id: 'root', name: 'Root Node' })

      // 构建复杂树结构
      tree.append({ id: 'a', name: 'Node A' }, 'root')
      tree.append({ id: 'b', name: 'Node B' }, 'root')
      tree.append({ id: 'c', name: 'Node C' }, 'root')
      tree.append({ id: 'a1', name: 'Node A1' }, 'a')
      tree.append({ id: 'a2', name: 'Node A2' }, 'a')
      tree.append({ id: 'a3', name: 'Node A3' }, 'a')
      tree.append({ id: 'b1', name: 'Node B1' }, 'b')
      tree.append({ id: 'c1', name: 'Node C1' }, 'c')
      tree.append({ id: 'a1a', name: 'Node A1A' }, 'a1')
      tree.append({ id: 'a1b', name: 'Node A1B' }, 'a1')
      tree.append({ id: 'a2a', name: 'Node A2A' }, 'a2')

      // 6.1 同ID替换 - 叶子节点
      const originalA1aParent = tree.getParent('a1a')?.id
      tree.replace('a1a', { id: 'a1a', name: 'New A1A', data: 'extra' })
      assert(tree.get('a1a')?.name === 'New A1A', '叶子节点同ID替换失败')
      assert(tree.getParent('a1a')?.id === originalA1aParent, '叶子节点替换后父节点应保持不变')

      // 6.2 同ID替换 - 有子节点的节点
      const originalA1Children = tree.getChildren('a1').map(n => n.id)
      tree.replace('a1', { id: 'a1', name: 'New A1', type: 'updated' })
      assert(tree.get('a1')?.name === 'New A1', '有子节点的同ID替换失败')
      const newA1Children = tree.getChildren('a1').map(n => n.id)
      assert(
        arrayEquals(originalA1Children, newA1Children),
        `替换后子节点丢失，期望: ${JSON.stringify(originalA1Children)}，实际: ${JSON.stringify(
          newA1Children
        )}`
      )

      // 6.3 同ID替换 - 根节点
      const originalRootChildren = tree.getChildren('root').map(n => n.id)
      tree.replace('root', { id: 'root', name: 'New Root', version: 2 })
      assert(tree.getRoot()?.name === 'New Root', '根节点同ID替换失败')
      const newRootChildren = tree.getChildren('root').map(n => n.id)
      assert(
        arrayEquals(originalRootChildren, newRootChildren),
        `根节点替换后子节点丢失，期望: ${JSON.stringify(
          originalRootChildren
        )}，实际: ${JSON.stringify(newRootChildren)}`
      )

      // 6.4 不同ID替换 - 叶子节点
      tree.replace('a1b', { id: 'a1x', name: 'A1X', status: 'renamed' })
      assert(!tree.has('a1b'), '原节点a1b应该不存在')
      assert(tree.has('a1x'), '新节点a1x应该存在')
      assert(tree.getParent('a1x')?.id === 'a1', 'a1x的父节点应该是a1')

      // 验证兄弟节点顺序
      const a1Children = tree.getChildren('a1').map(n => n.id)
      assert(
        arrayEquals(a1Children, ['a1a', 'a1x']),
        `a1子节点顺序错误，期望: ['a1a', 'a1x']，实际: ${JSON.stringify(a1Children)}`
      )

      // 6.5 不同ID替换 - 有子节点的中间节点
      const a2Children = tree.getChildren('a2').map(n => n.id)
      tree.replace('a2', { id: 'a2new', name: 'A2 New', category: 'renamed' })
      assert(!tree.has('a2'), '原节点a2应该不存在')
      assert(tree.has('a2new'), '新节点a2new应该存在')
      assert(tree.getParent('a2new')?.id === 'a', 'a2new的父节点应该是a')

      // 验证子节点的父节点更新
      for (const childId of a2Children) {
        assert(tree.getParent(childId)?.id === 'a2new', `子节点${childId}的父节点应该更新为a2new`)
      }

      // 验证在父节点中的位置
      const aChildren = tree.getChildren('a').map(n => n.id)
      assert(
        arrayEquals(aChildren, ['a1', 'a2new', 'a3']),
        `a的子节点顺序错误，期望: ['a1', 'a2new', 'a3']，实际: ${JSON.stringify(aChildren)}`
      )

      // 6.6 不同ID替换 - 根节点
      const allRootChildren = tree.getChildren('root').map(n => n.id)
      tree.replace('root', { id: 'newroot', name: 'Brand New Root', level: 0 })
      assert(!tree.has('root'), '原根节点root应该不存在')
      assert(tree.has('newroot'), '新根节点newroot应该存在')
      assert(tree.getRoot()?.id === 'newroot', 'getRoot应该返回新根节点')

      // 验证所有原子节点的父节点都更新了
      for (const childId of allRootChildren) {
        assert(
          tree.getParent(childId)?.id === 'newroot',
          `根节点子节点${childId}的父节点应该更新为newroot`
        )
      }

      // 6.7 复杂嵌套替换测试
      // 先添加更多层级
      tree.append({ id: 'a1a1', name: 'A1A1' }, 'a1a')
      tree.append({ id: 'a1a2', name: 'A1A2' }, 'a1a')

      // 替换中间层级节点
      tree.replace('a1a', { id: 'a1z', name: 'A1Z', depth: 3 })
      assert(!tree.has('a1a'), '原节点a1a应该不存在')
      assert(tree.has('a1z'), '新节点a1z应该存在')
      assert(tree.getParent('a1z')?.id === 'a1', 'a1z的父节点应该是a1')
      assert(tree.getParent('a1a1')?.id === 'a1z', 'a1a1的父节点应该更新为a1z')
      assert(tree.getParent('a1a2')?.id === 'a1z', 'a1a2的父节点应该更新为a1z')

      // 6.8 替换错误情况测试

      // 尝试替换不存在的节点
      try {
        tree.replace('nonexistent', { id: 'new', name: 'New' })
        throw new Error('应该抛出节点不存在错误')
      } catch (e: any) {
        assert(
          e.message.includes('does not exist') || e.message.includes('不存在'),
          '应该抛出节点不存在错误'
        )
      }

      // 尝试用已存在的ID替换（不同ID情况）
      try {
        tree.replace('a1', { id: 'a3', name: 'Conflict' })
        throw new Error('应该抛出ID冲突错误')
      } catch (e: any) {
        assert(
          e.message.includes('already exists') || e.message.includes('已存在'),
          '应该抛出ID已存在错误'
        )
      }

      // 6.9 链式替换测试
      tree.append({ id: 'temp1', name: 'Temp1' }, 'newroot')
      tree.append({ id: 'temp2', name: 'Temp2' }, 'temp1')
      tree.append({ id: 'temp3', name: 'Temp3' }, 'temp2')

      // 从下往上替换
      tree.replace('temp3', { id: 'final3', name: 'Final3' })
      tree.replace('temp2', { id: 'final2', name: 'Final2' })
      tree.replace('temp1', { id: 'final1', name: 'Final1' })

      assert(tree.getParent('final3')?.id === 'final2', '链式替换父子关系错误1')
      assert(tree.getParent('final2')?.id === 'final1', '链式替换父子关系错误2')
      assert(tree.getParent('final1')?.id === 'newroot', '链式替换父子关系错误3')

      // 6.10 大批量替换测试
      for (let i = 1; i <= 5; i++) {
        tree.append({ id: `batch${i}`, name: `Batch ${i}` }, 'newroot')
        for (let j = 1; j <= 3; j++) {
          tree.append({ id: `batch${i}_${j}`, name: `Batch ${i}_${j}` }, `batch${i}`)
        }
      }

      // 批量替换所有batch节点
      for (let i = 1; i <= 5; i++) {
        tree.replace(`batch${i}`, { id: `newbatch${i}`, name: `New Batch ${i}`, updated: true })

        // 验证子节点父节点更新
        for (let j = 1; j <= 3; j++) {
          assert(
            tree.getParent(`batch${i}_${j}`)?.id === `newbatch${i}`,
            `批量替换后子节点父节点错误: batch${i}_${j}`
          )
        }
      }

      // 6.11 验证树的完整性
      let nodeCount = 0
      tree.enumerate((id, path) => {
        nodeCount++
        const node = tree.get(id)
        assert(node !== undefined, `枚举到的节点${id}不存在`)

        // 验证父子关系一致性（除根节点外）
        if (id !== tree.getRoot()?.id) {
          const parent = tree.getParent(id)
          assert(parent !== undefined, `节点${id}没有父节点`)
          const siblings = tree.getChildren(parent!.id).map(n => n.id)
          assert(siblings.includes(id), `节点${id}不在其父节点的子节点列表中`)
        }
      })

      assert(nodeCount === tree.size, `枚举节点数量${nodeCount}与size${tree.size}不匹配`)

      tree.dispose()
    }

    console.log('所有测试通过！')
  }

  runTests()
}

export { TreeManager2 }

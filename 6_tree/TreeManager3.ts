/* eslint-disable no-lone-blocks */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-lone-blocks */
/* eslint-disable no-inner-declarations */

// 本身是无状态的（除了根节点引用），不依赖于额外的内存来存储索引.
// 缺点是几乎所有操作的时间复杂度都会从 O(1) 或 O(log n) 变为 O(n)，其中 n 是树中的节点总数.
// 因为每次操作都可能需要遍历整棵树。
// !如果需要快速检索，需要维护 idToNode、idToParent.

/**
 * 内部的修改操作会直接影响节点.
 */
class TreeManager3<TNode> {
  private _root: TNode
  private readonly _getId: (node: TNode) => string
  private readonly _getChildren: (node: TNode) => TNode[]

  constructor(
    root: TNode,
    operators: {
      getId: (node: TNode) => string
      getChildren: (node: TNode) => TNode[]
    }
  ) {
    this._root = root
    this._getId = operators.getId
    this._getChildren = operators.getChildren
  }

  dispose(): void {
    // @ts-ignore
    this._root = undefined
  }

  /**
   * 将节点追加到父节点的子节点列表末尾。
   */
  append(node: TNode, parentId: string): void {
    this._insert(node, parentId, children => children.push(node))
  }

  /**
   * 将节点添加到父节点的子节点列表开头。
   */
  prepend(node: TNode, parentId: string): void {
    this._insert(node, parentId, children => children.unshift(node))
  }

  insertBefore(node: TNode, targetId: string): void {
    const { node: targetNode, parent: parentNode } = this._findNodeAndParent(targetId)
    if (!targetNode || !parentNode) {
      throw new Error(`Target node ${targetId} or its parent not found.`)
    }
    const parentId = this._getId(parentNode)
    this._insert(node, parentId, children => {
      const index = children.indexOf(targetNode)
      children.splice(index, 0, node)
    })
  }

  insertAfter(node: TNode, targetId: string): void {
    const { node: targetNode, parent: parentNode } = this._findNodeAndParent(targetId)
    if (!targetNode || !parentNode) {
      throw new Error(`Target node ${targetId} or its parent not found.`)
    }
    const parentId = this._getId(parentNode)
    this._insert(node, parentId, children => {
      const index = children.indexOf(targetNode)
      children.splice(index + 1, 0, node)
    })
  }

  /**
   * 移除节点及其所有子节点.
   */
  remove(id: string): void {
    if (this._isRoot(id)) {
      throw new Error('Cannot remove the root node.')
    }
    const { node, parent } = this._findNodeAndParent(id)
    if (!node || !parent) {
      throw new Error(`Node with ID ${id} not found or has no parent.`)
    }
    const siblings = this._getChildren(parent)
    const index = siblings.indexOf(node)
    if (index !== -1) {
      siblings.splice(index, 1)
    }
  }

  /**
   * 将一个节点替换为新节点。
   * 新节点将占据旧节点的位置，但不会继承其子树。旧节点及其子树将被移除。
   */
  replace(oldId: string, newNode: TNode): void {
    const newId = this._getId(newNode)
    if (oldId !== newId && this.has(newId)) {
      throw new Error(`Node with new ID ${newId} already exists.`)
    }

    const { node: oldNode, parent: parentNode } = this._findNodeAndParent(oldId)
    if (!oldNode) {
      throw new Error(`Node to be replaced with ID ${oldId} not found.`)
    }

    if (!parentNode) {
      // 替换根节点
      this._root = newNode
    } else {
      const siblings = this._getChildren(parentNode)
      const index = siblings.indexOf(oldNode)
      if (index !== -1) {
        siblings.splice(index, 1, newNode)
      }
    }
  }

  update(id: string, f: (node: TNode) => TNode): void {
    const node = this.get(id)
    if (!node) {
      throw new Error(`Node with ID ${id} not found.`)
    }
    const newNode = f(node)
    if (this._getId(newNode) !== id) {
      throw new Error('Updater function cannot change the node ID.')
    }
    // 如果返回了一个新对象，则进行替换
    if (newNode !== node) {
      this.replace(id, newNode)
    }
  }

  moveBefore(id: string, targetId: string): void {
    this._move(id, targetId, false)
  }

  moveAfter(id: string, targetId: string): void {
    this._move(id, targetId, true)
  }

  get(id: string): TNode | undefined {
    return this._findNodeAndParent(id).node
  }

  getParent(id: string): TNode | undefined {
    return this._findNodeAndParent(id).parent
  }

  getChildren(id: string): TNode[] {
    const node = this.get(id)
    return node ? this._getChildren(node) : []
  }

  getRoot(): TNode {
    return this._root
  }

  getSize(): number {
    let res = 0
    this.enumerate(() => {
      res++
    })
    return res
  }

  previousSibling(id: string): TNode | undefined {
    const { node, parent } = this._findNodeAndParent(id)
    if (!node || !parent) return undefined
    const siblings = this._getChildren(parent)
    const index = siblings.indexOf(node)
    return index > 0 ? siblings[index - 1] : undefined
  }

  nextSibling(id: string): TNode | undefined {
    const { node, parent } = this._findNodeAndParent(id)
    if (!node || !parent) return undefined
    const siblings = this._getChildren(parent)
    const index = siblings.indexOf(node)
    return index !== -1 && index < siblings.length - 1 ? siblings[index + 1] : undefined
  }

  has(id: string): boolean {
    return !!this.get(id)
  }

  /**
   * 前序遍历树。
   * @param f 遍历回调函数，返回true时停止遍历。
   */
  enumerate(f: (cur: TNode, path: TNode[]) => boolean | void): void {
    const dfs = (cur: TNode, path: TNode[]): boolean => {
      if (f(cur, path)) {
        return true
      }
      path.push(cur)
      for (const child of this._getChildren(cur)) {
        if (dfs(child, path)) {
          path.pop()
          return true
        }
      }
      path.pop()
      return false
    }
    dfs(this._root, [])
  }

  print(options?: { format?: (node: TNode) => string; output?: (message: string) => void }): void {
    const format = options?.format || ((node: TNode) => this._getId(node))
    const output = options?.output || console.log
    this.enumerate((cur, path) => {
      const message = `${'  '.repeat(path.length)}- ${format(cur)}`
      output(message)
    })
  }

  private _insert(node: TNode, parentId: string, f: (children: TNode[]) => void): void {
    if (this.has(this._getId(node))) {
      throw new Error(`Node with ID ${this._getId(node)} already exists.`)
    }
    const parentNode = this.get(parentId)
    if (!parentNode) {
      throw new Error(`Parent node with ID ${parentId} does not exist.`)
    }
    const children = this._getChildren(parentNode)
    f(children)
  }

  /**
   * 查找节点及其父节点。
   *
   * @returns
   * 返回一个包含找到的节点和其父节点的对象。
   * 如果未找到节点，node为undefined。
   * 如果找到的节点是根节点，parent为undefined。
   */
  private _findNodeAndParent(id: string): { node: TNode | undefined; parent: TNode | undefined } {
    if (this._isRoot(id)) {
      return { node: this._root, parent: undefined }
    }

    const dfs = (
      node: TNode,
      parent: TNode | undefined
    ): { node: TNode | undefined; parent: TNode | undefined } => {
      if (this._getId(node) === id) {
        return { node, parent }
      }
      for (const child of this._getChildren(node)) {
        const res = dfs(child, node)
        if (res.node) {
          return res
        }
      }
      return { node: undefined, parent: undefined }
    }
    return dfs(this._root, undefined)
  }

  private _move(id: string, targetId: string, after: boolean): void {
    if (id === targetId) throw new Error('Cannot move a node relative to itself.')
    if (this._isAncestor(id, targetId)) throw new Error('Cannot move a node into its descendant.')

    const nodeToMove = this.get(id)
    if (!nodeToMove) throw new Error(`Node to move with ID ${id} not found.`)

    this.remove(id)
    after ? this.insertAfter(nodeToMove, targetId) : this.insertBefore(nodeToMove, targetId)
  }

  private _isAncestor(ancestorId: string, descendantId: string): boolean {
    if (ancestorId === descendantId) return false

    const ancestorNode = this.get(ancestorId)
    if (!ancestorNode) return false

    const stack: TNode[] = [ancestorNode]
    while (stack.length) {
      const cur = stack.pop()!
      for (const child of this._getChildren(cur)) {
        if (this._getId(child) === descendantId) return true
        stack.push(child)
      }
    }
    return false
  }

  private _isRoot(id: string): boolean {
    return this._getId(this._root) === id
  }
}

export { TreeManager3 }

if (typeof require !== 'undefined') {
  // ==================== 测试框架 ====================
  function assert(condition: boolean, message: string): void {
    if (!condition) {
      throw new Error(`断言失败: ${message}`)
    }
  }

  function assertThrows(fn: () => void, message: string): void {
    try {
      fn()
      assert(false, `预期抛出错误，但没有抛出: ${message}`)
    } catch (e: any) {
      // 可选：检查错误消息 e.message
    }
  }

  // ==================== 测试数据结构 ====================
  interface TestNode {
    id: string
    children: TestNode[]
    data?: string
  }

  const getTestNodeId = (node: TestNode) => node.id
  const getTestNodeChildren = (node: TestNode) => node.children

  // 创建一个标准的、可复用的树结构
  const treeFactory = (): TestNode => ({
    id: 'root',
    children: [
      {
        id: 'child1',
        children: [
          { id: 'grandchild1A', children: [] },
          { id: 'grandchild1B', children: [] }
        ]
      },
      {
        id: 'child2',
        children: [{ id: 'grandchild2A', children: [] }]
      },
      { id: 'child3', children: [] }
    ]
  })

  // ==================== 测试主体 ====================
  function runTests(): void {
    console.log('开始运行 TreeManager3 测试...')

    // --- 1. 构造与基本读取测试 ---
    console.log('  - 测试: 构造与基本读取...')
    let root = treeFactory()
    let tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })

    assert(tm.getRoot().id === 'root', 'getRoot() 应返回根节点')
    assert(tm.get('child2')?.id === 'child2', 'get() 应能获取节点')
    assert(tm.get('non-existent') === undefined, 'get() 对不存在的ID应返回undefined')
    assert(tm.has('grandchild1A'), 'has() 应能判断节点存在')
    assert(!tm.has('non-existent'), 'has() 应能判断节点不存在')
    assert(tm.getParent('child1')?.id === 'root', 'getParent() 应能获取父节点')
    assert(tm.getParent('root') === undefined, 'getParent() 对根节点应返回undefined')
    assert(tm.getChildren('child1').length === 2, 'getChildren() 应能获取子节点列表')
    assert(tm.getChildren('child1')[0].id === 'grandchild1A', 'getChildren() 子节点顺序应正确')
    assert(tm.getSize() === 7, 'getSize() 应返回正确的节点总数')

    // --- 2. 兄弟节点测试 ---
    console.log('  - 测试: 兄弟节点...')
    assert(tm.previousSibling('child2')?.id === 'child1', 'previousSibling() 应能获取前一个兄弟')
    assert(
      tm.previousSibling('child1') === undefined,
      'previousSibling() 对第一个子节点应返回undefined'
    )
    assert(tm.nextSibling('child2')?.id === 'child3', 'nextSibling() 应能获取后一个兄弟')
    assert(tm.nextSibling('child3') === undefined, 'nextSibling() 对最后一个子节点应返回undefined')
    assert(tm.nextSibling('root') === undefined, 'nextSibling() 对根节点应返回undefined')

    // --- 3. 插入操作测试 ---
    console.log('  - 测试: 插入操作...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    const newNode: TestNode = { id: 'newChild', children: [] }
    tm.append(newNode, 'root')
    assert(root.children[3].id === 'newChild', 'append() 应能追加到末尾')
    assert(tm.getSize() === 8, 'append() 后 getSize() 应正确')
    assert(tm.getParent('newChild')?.id === 'root', 'append() 后 getParent() 应正确')

    const newNode2: TestNode = { id: 'newChild2', children: [] }
    tm.prepend(newNode2, 'child1')
    assert(root.children[0].children[0].id === 'newChild2', 'prepend() 应能添加到开头')
    assert(tm.getSize() === 9, 'prepend() 后 getSize() 应正确')

    const newNode3: TestNode = { id: 'newChild3', children: [] }
    tm.insertBefore(newNode3, 'grandchild1B')
    tm.print()
    assert(root.children[0].children[2].id === 'newChild3', 'insertBefore() 应能插入到目标之前')
    assert(tm.getSize() === 10, 'insertBefore() 后 getSize() 应正确')

    const newNode4: TestNode = { id: 'newChild4', children: [] }
    tm.insertAfter(newNode4, 'child2')
    assert(root.children[2].id === 'newChild4', 'insertAfter() 应能插入到目标之后')
    assert(tm.getSize() === 11, 'insertAfter() 后 getSize() 应正确')

    // --- 4. 删除操作测试 ---
    console.log('  - 测试: 删除操作...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    tm.remove('child2')
    assert(root.children.length === 2, 'remove() 后父节点的子节点数应减少')
    assert(root.children[1].id === 'child3', 'remove() 后兄弟节点应连接')
    assert(!tm.has('child2'), 'remove() 后节点应不存在')
    assert(!tm.has('grandchild2A'), 'remove() 后子孙节点也应被移除 (逻辑上)')
    assert(tm.getSize() === 5, 'remove() 后 getSize() 应正确')

    // --- 5. 替换操作测试 ---
    console.log('  - 测试: 替换操作...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    const replacementNode: TestNode = { id: 'replacedChild', children: [] }
    tm.replace('child1', replacementNode)
    assert(root.children[0].id === 'replacedChild', 'replace() 应能替换节点')
    assert(root.children[0].children.length === 0, 'replace() 不应继承子节点')
    assert(!tm.has('child1'), 'replace() 后旧节点应不存在')
    assert(!tm.has('grandchild1A'), 'replace() 后旧节点的子孙应不存在 (逻辑上)')
    assert(tm.getSize() === 5, 'replace() 后 getSize() 应正确')

    // 替换根节点
    const newRoot: TestNode = { id: 'newRoot', children: [] }
    tm.replace('root', newRoot)
    assert(tm.getRoot().id === 'newRoot', 'replace() 应能替换根节点')
    assert(tm.getSize() === 1, 'replace() 根节点后 getSize() 应为1')

    // --- 6. 更新操作测试 ---
    console.log('  - 测试: 更新操作...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    tm.update('child1', node => {
      node.data = 'updated'
      return node
    })
    assert((tm.get('child1') as any).data === 'updated', 'update() 应能修改节点数据')

    // 更新返回新对象
    tm.update('child2', node => ({ ...node, data: 'updated again' }))
    assert((tm.get('child2') as any).data === 'updated again', 'update() 应能处理返回新对象的情况')
    assert(root.children[1].data === 'updated again', 'update() 返回新对象时应在树中生效')

    // --- 7. 移动操作测试 ---
    console.log('  - 测试: 移动操作...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    tm.moveAfter('child3', 'child1')
    assert(root.children[1].id === 'child3', 'moveAfter() 应能移动节点')
    assert(root.children.length === 3, 'moveAfter() 后原父节点的子节点数应正确')
    assert(tm.getParent('child3')?.id === 'root', 'moveAfter() 后节点的父节点应更新')

    tm.moveBefore('grandchild1B', 'grandchild2A')
    assert(tm.getParent('grandchild1B')?.id === 'child2', 'moveBefore() 应能跨父节点移动')
    assert(root.children[0].children.length === 1, 'moveBefore() 后原父节点的子节点数应减少')
    assert(root.children[2].children.length === 2, 'moveBefore() 后新父节点的子节点数应增加')
    assert(root.children[2].children[0].id === 'grandchild1B', 'moveBefore() 应移动到正确位置')

    // --- 8. 遍历与打印测试 ---
    console.log('  - 测试: 遍历与打印...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    const enumeratedIds: string[] = []
    tm.enumerate((node, path) => {
      enumeratedIds.push(node.id)
      if (node.id === 'child1') {
        assert(path.length === 1 && path[0].id === 'root', 'enumerate() path参数应正确')
      }
    })
    const expectedOrder = [
      'root',
      'child1',
      'grandchild1A',
      'grandchild1B',
      'child2',
      'grandchild2A',
      'child3'
    ]
    assert(
      JSON.stringify(enumeratedIds) === JSON.stringify(expectedOrder),
      'enumerate() 应按深度优先顺序遍历'
    )

    // 测试提前终止
    enumeratedIds.length = 0
    tm.enumerate(node => {
      enumeratedIds.push(node.id)
      return node.id === 'grandchild1A' // 在此中断
    })
    assert(
      JSON.stringify(enumeratedIds) === JSON.stringify(['root', 'child1', 'grandchild1A']),
      'enumerate() 应能被提前终止'
    )

    // 打印测试 (仅确保不抛出错误)
    tm.print()

    // --- 9. 错误与边界条件测试 ---
    console.log('  - 测试: 错误与边界条件...')
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    assertThrows(
      () => tm.append({ id: 'child1', children: [] }, 'root'),
      '插入已存在的ID应抛出错误'
    )
    assertThrows(
      () => tm.append({ id: 'newNode', children: [] }, 'non-existent'),
      '插入到不存在的父节点应抛出错误'
    )
    assertThrows(() => tm.remove('root'), '移除根节点应抛出错误')
    assertThrows(() => tm.moveBefore('child1', 'grandchild1A'), '不能将节点移动到其子孙节点内部')
    assertThrows(
      () => tm.update('child1', node => ({ ...node, id: 'newId' })),
      'update() 不能修改ID'
    )

    // --- 10. Dispose 测试 ---
    console.log('  - 测试: Dispose...')
    root = treeFactory() // 重新创建以进行 dispose 测试
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    tm.dispose()
    assertThrows(() => tm.getSize(), 'dispose() 后调用方法应失败 (因为_root为undefined)')
    assertThrows(() => tm.get('root'), 'dispose() 后调用方法应失败')

    // --- 11. 更多移动和结构修改测试 ---
    console.log('  - 测试: 更多移动和结构修改...')
    // 测试移动一个带有子树的节点
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    tm.moveAfter('child1', 'child3') // 将 child1 及其子树移动到 child3 之后
    assert(root.children.length === 3, '移动后根节点的子节点数应不变')
    assert(root.children[2].id === 'child1', '带有子树的节点应被移动到正确位置')
    assert(tm.getParent('child1')?.id === 'root', '节点的父节点引用应正确')
    assert(tm.has('grandchild1A'), '节点的子孙应随之移动')
    assert(tm.getParent('grandchild1A')?.id === 'child1', '子孙节点的父节点引用应保持正确')
    assert(tm.getSize() === 7, '移动带有子树的节点后getSize()应不变')

    // 测试替换一个节点，新节点使用相同的ID
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    const sameIdReplacement: TestNode = { id: 'child2', children: [], data: 'is a replacement' }
    tm.replace('child2', sameIdReplacement)
    assert(tm.get('child2')?.data === 'is a replacement', 'replace() 使用相同ID时应能成功')
    assert(!tm.has('grandchild2A'), 'replace() 使用相同ID时，旧节点的子孙也应被移除')
    assert(tm.getSize() === 6, 'replace() 使用相同ID后 getSize() 应正确')

    // --- 12. 更多边界条件测试 ---
    console.log('  - 测试: 更多边界条件...')
    // 测试只含根节点的树
    let singleNodeRoot: TestNode = { id: 'singleRoot', children: [] }
    let singleNodeTm = new TreeManager3(singleNodeRoot, {
      getId: getTestNodeId,
      getChildren: getTestNodeChildren
    })
    assert(singleNodeTm.getSize() === 1, '单节点树的getSize()应为1')
    singleNodeTm.append({ id: 'newLeaf', children: [] }, 'singleRoot')
    assert(singleNodeTm.getSize() === 2, '单节点树append后getSize()应为2')
    assert(singleNodeRoot.children[0].id === 'newLeaf', '单节点树append后子节点应正确')
    singleNodeTm.remove('newLeaf')
    assert(singleNodeTm.getSize() === 1, '单节点树remove后getSize()应为1')

    // 更多移动操作的错误条件
    root = treeFactory()
    tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })
    assertThrows(() => tm.moveBefore('child1', 'child1'), '不能将节点移动到其自身相对位置')
    assertThrows(() => tm.moveAfter('child1', 'child1'), '不能将节点移动到其自身相对位置')
    assertThrows(() => tm.moveBefore('non-existent', 'child1'), '移动不存在的节点应抛出错误')
    assertThrows(() => tm.moveBefore('child1', 'non-existent'), '移动到不存在的目标应抛出错误')

    // 测试移除叶子节点
    tm = new TreeManager3(treeFactory(), { getId: getTestNodeId, getChildren: getTestNodeChildren })
    tm.remove('grandchild1A')
    assert(!tm.has('grandchild1A'), '移除叶子节点后应不存在')
    assert(tm.get('child1')?.children.length === 1, '移除叶子节点后其父节点的子节点数应减少')
    assert(tm.get('child1')?.children[0].id === 'grandchild1B', '移除叶子节点后其兄弟应不受影响')
    assert(tm.getSize() === 6, '移除叶子节点后getSize()应正确')

    {
      console.log('  - 测试: 移动操作...')
      root = treeFactory()
      tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })

      /* ---------- 基础场景 ---------- */
      tm.moveAfter('child3', 'child1')
      assert(
        root.children.map(n => n.id).join() === 'child1,child3,child2',
        'moveAfter() 应改变顺序为 child1,child3,child2'
      )

      tm.moveBefore('child3', 'child2')
      assert(
        root.children.map(n => n.id).join() === 'child1,child3,child2',
        'moveBefore() 在相邻节点间再次移动应保持顺序不变'
      )

      /* ---------- 移动到列表开头 ---------- */
      tm.moveBefore('child2', 'child1')
      assert(root.children[0].id === 'child2', 'moveBefore() 应能把节点移到列表开头')
      assert(tm.previousSibling('child2') === undefined, '列表开头的前置兄弟应为 undefined')

      /* ---------- 移动到列表末尾 ---------- */
      tm.moveAfter('child2', 'child3')
      assert(root.children.at(-1)!.id === 'child2', 'moveAfter() 应能把节点移到列表末尾')
      assert(tm.nextSibling('child2') === undefined, '列表末尾的后置兄弟应为 undefined')

      /* ---------- 跨父节点移动，且带子树 ---------- */
      tm.moveAfter('grandchild1A', 'grandchild2A') // child1 → child2
      assert(tm.getParent('grandchild1A')?.id === 'child2', '跨父节点移动父引用应更新')
      assert(
        root.children[2].children.map(n => n.id).join() === 'grandchild2A,grandchild1A',
        '目标父节点的子列表顺序应正确'
      )
      assert(root.children[0].children.length === 1, '原父节点子数组应减少一个元素')

      /* ---------- 还原 & 复杂链式移动 ---------- */
      root = treeFactory()
      tm = new TreeManager3(root, { getId: getTestNodeId, getChildren: getTestNodeChildren })

      // 把 child1 移到最后，再把 grandchild1B 移到根最前
      tm.moveAfter('child1', 'child3')
      tm.moveBefore('grandchild1B', 'child1')
      assert(
        root.children.map(n => n.id).join() === 'child2,child3,grandchild1B,child1',
        '链式移动后根节点子顺序应符合预期'
      )
      assert(tm.getParent('grandchild1B') === root, 'grandchild1B 被提升为根的直接子节点')

      // 迭代深度确认
      const orderAfterMoves: string[] = []
      tm.enumerate(node => {
        orderAfterMoves.push(node.id)
      })
      assert(
        orderAfterMoves.includes('grandchild1A') && tm.getParent('grandchild1A')?.id === 'child1',
        '移动后仍保持 grandchild1A 的父子关系'
      )
    }

    console.log('所有测试通过！')
  }

  runTests()
}

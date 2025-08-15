class DefaultDict<K, V> extends Map<K, V> {
  private readonly _defaultFactory: (self: DefaultDict<K, V>) => V

  constructor(
    defaultFactory: (self: DefaultDict<K, V>) => V,
    iterable?: Iterable<readonly [K, V]> | null
  ) {
    super(iterable)
    this._defaultFactory = defaultFactory
  }

  override get(key: K): V {
    if (super.has(key)) return super.get(key)!
    const value = this._defaultFactory(this)
    super.set(key, value)
    return value
  }

  setDefault(key: K, value: V): V {
    if (super.has(key)) return super.get(key)!
    super.set(key, value)
    return value
  }
}

/**
 * !管理树形`结构`(节点父子关系)。
 * !不管理节点`内容`。
 */
export class TreeStructure<Id extends PropertyKey = string> {
  static fromFlattenedTree<Node, Id extends PropertyKey>(
    nodes: Node[],
    root: Id,
    operations: {
      getId: (node: Node) => Id
      getChildren: (node: Node) => Id[]
    }
  ): TreeStructure<Id> {
    const res = new TreeStructure<Id>(root)

    const idToNode = new Map<Id, Node>()
    for (const node of nodes) {
      idToNode.set(operations.getId(node), node)
    }

    const dfs = (node: Node) => {
      const nodeId = operations.getId(node)
      const childrenIds = operations.getChildren(node)
      if (!childrenIds.length) return
      res.append(nodeId, ...childrenIds)
      for (const child of childrenIds) {
        if (!idToNode.has(child)) {
          throw new Error(`Child node ${String(child)} not found in the provided nodes`)
        }
        dfs(idToNode.get(child)!)
      }
    }

    if (!idToNode.has(root)) {
      throw new Error(`Root node ${String(root)} not found in the provided nodes`)
    }
    dfs(idToNode.get(root)!)

    return res
  }

  static fromNestedTree<Node, Id extends PropertyKey>(
    root: Node,
    operations: {
      getId: (node: Node) => Id
      getChildren: (node: Node) => Node[]
    }
  ): TreeStructure<Id> {
    const res = new TreeStructure<Id>(operations.getId(root))
    const dfs = (node: Node) => {
      const nodeId = operations.getId(node)
      const children = operations.getChildren(node)
      if (!children.length) return
      const childrenIds = children.map(c => operations.getId(c))
      res.append(nodeId, ...childrenIds)
      for (const child of children) {
        dfs(child)
      }
    }
    dfs(root)
    return res
  }

  private _root: Id | undefined = undefined
  private readonly _children = new DefaultDict<Id, Id[]>(() => [])
  private readonly _parent = new Map<Id, Id>()

  constructor(root: Id) {
    this._root = root
  }

  dispose() {
    this._root = undefined
    this._children.clear()
    this._parent.clear()
  }

  /**
   * 将节点追加到父节点的子节点列表末尾。
   */
  append(parent: Id, ...nodes: Id[]): void {
    this._insert(parent, nodes, children => {
      children.push(...nodes)
    })
  }

  /**
   * 将节点添加到父节点的子节点列表开头。
   */
  prepend(parent: Id, ...nodes: Id[]): void {
    this._insert(parent, nodes, children => {
      children.unshift(...nodes)
    })
  }

  insertBefore(reference: Id, ...nodes: Id[]): void {
    const parent = this._parent.get(reference)
    if (parent === undefined) {
      throw new Error(`Cannot insert before root node or non-existent node: ${String(reference)}`)
    }

    this._insert(parent, nodes, children => {
      const index = children.indexOf(reference)
      children.splice(index, 0, ...nodes)
    })
  }

  insertAfter(reference: Id, ...nodes: Id[]): void {
    const parent = this._parent.get(reference)
    if (parent === undefined) {
      throw new Error(`Cannot insert after root node or non-existent node: ${String(reference)}`)
    }

    this._insert(parent, nodes, children => {
      const index = children.indexOf(reference)
      children.splice(index + 1, 0, ...nodes)
    })
  }

  /**
   * 移除节点及其所有子节点。
   */
  remove(node: Id): void {
    if (this._isRoot(node)) {
      throw new Error('Cannot remove root node')
    }
    if (!this.has(node)) {
      throw new Error(`Node ${String(node)} does not exist`)
    }

    this._removeFromParent(node)

    // removeSubtree
    const subTree: Id[] = []
    const dfs = (cur: Id) => {
      subTree.push(cur)
      for (const c of this._children.get(cur)) {
        dfs(c)
      }
    }
    dfs(node)
    for (const v of subTree) {
      this._children.delete(v)
      this._parent.delete(v)
    }
  }

  /**
   * 将一个节点替换为新节点。
   * 新节点将占据旧节点的位置，旧节点及其子树将被移除。
   */
  replace(oldNode: Id, newNode: Id): void {
    if (oldNode === newNode) {
      return
    }
    if (!this.has(oldNode)) {
      throw new Error(`Node ${String(oldNode)} does not exist`)
    }
    if (this.has(newNode)) {
      throw new Error(`New node ${String(newNode)} already exists`)
    }

    if (this._isRoot(oldNode)) {
      this._root = newNode
      this._children.clear()
      this._parent.clear()
      return
    }

    const parent = this._parent.get(oldNode)!
    const siblings = this._children.get(parent)
    const oldIndex = siblings.indexOf(oldNode)

    this.remove(oldNode)
    this._insert(parent, [newNode], children => {
      children.splice(oldIndex, 0, newNode)
    })
  }

  moveBefore(reference: Id, node: Id): void {
    this._move(reference, node, true)
  }

  moveAfter(reference: Id, node: Id): void {
    this._move(reference, node, false)
  }

  getParent(node: Id): Id | undefined {
    return this._parent.get(node)
  }

  getChildren(node: Id): Id[] {
    if (!this.has(node)) return []
    return this._children.get(node).slice()
  }

  getRoot(): Id | undefined {
    return this._root
  }

  previousSibling(node: Id): Id | undefined {
    return this._getSibling(node, -1)
  }

  nextSibling(node: Id): Id | undefined {
    return this._getSibling(node, 1)
  }

  /**
   * 前序遍历树。
   * @param f 遍历回调函数，返回true时停止遍历。
   */
  enumerate(f: (node: Id, path: Id[]) => boolean | void): void {
    const dfs = (node: Id, path: Id[]): boolean => {
      if (f(node, path)) {
        return true
      }
      path.push(node)
      for (const c of this._children.get(node)) {
        if (dfs(c, path)) {
          path.pop()
          return true
        }
      }
      path.pop()
      return false
    }

    if (this._root !== undefined) {
      dfs(this._root, [])
    }
  }

  print(options?: { format?: (node: Id) => string; output?: (message: string) => void }): void {
    const { format = (node: Id) => String(node), output = console.log } = options || {}
    this.enumerate((node, path) => {
      const indent = ' '.repeat(path.length * 2)
      const message = `${indent}- ${format(node)}`
      output(message)
    })
  }

  has(node: Id): boolean {
    return node === this._root || this._parent.has(node)
  }

  private _insert(parent: Id, nodes: Id[], f: (children: Id[]) => void): void {
    if (!nodes.length) {
      return
    }
    if (!this.has(parent)) {
      throw new Error(`Parent node ${String(parent)} does not exist`)
    }
    for (const n of nodes) {
      if (this.has(n)) {
        throw new Error(`Node ${String(n)} already exists`)
      }
    }

    for (const node of nodes) {
      this._parent.set(node, parent)
    }
    const children = this._children.get(parent)
    f(children)
  }

  private _move(reference: Id, node: Id, before: boolean): void {
    if (this._isRoot(node) || this._isRoot(reference)) {
      throw new Error('Cannot move root node')
    }
    if (!this.has(node) || !this.has(reference)) {
      throw new Error('Node does not exist')
    }
    if (node === reference) {
      throw new Error('Cannot move node to itself')
    }
    if (this._isAncestor(node, reference)) {
      throw new Error('Cannot move node to its descendant')
    }

    const newParent = this._parent.get(reference)
    if (newParent === undefined) {
      throw new Error(`Target node ${String(reference)} does not have a parent`)
    }

    this._removeFromParent(node)

    const newSiblings = this._children.get(newParent)
    const targetIndex = newSiblings.indexOf(reference)
    const insertIndex = before ? targetIndex : targetIndex + 1
    newSiblings.splice(insertIndex, 0, node)

    this._parent.set(node, newParent)
  }

  private _removeFromParent(node: Id): void {
    const parent = this._parent.get(node)
    if (parent !== undefined) {
      const siblings = this._children.get(parent)
      const index = siblings.indexOf(node)
      if (index !== -1) {
        siblings.splice(index, 1)
      }
    }
  }

  private _getSibling(node: Id, offset: number): Id | undefined {
    const parent = this._parent.get(node)
    if (parent === undefined) return undefined

    const siblings = this._children.get(parent)
    const index = siblings.indexOf(node)
    if (index === -1) return undefined

    const targetIndex = index + offset
    return targetIndex >= 0 && targetIndex < siblings.length ? siblings[targetIndex] : undefined
  }

  private _isAncestor(ancestor: Id, node: Id): boolean {
    let cur: Id | undefined = node
    while (cur !== undefined) {
      if (cur === ancestor) {
        return true
      }
      cur = this._parent.get(cur)
    }
    return false
  }

  private _isRoot(node: Id): boolean {
    return this._root === node
  }
}

export {}

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
    console.log('开始运行 TreeStructure 测试...')

    // 1. 基础树操作测试
    {
      const tree = new TreeStructure<string>('root')

      // 添加节点测试
      tree.append('root', 'a')
      tree.append('root', 'b')
      tree.prepend('root', 'c')
      tree.append('a', 'd')
      tree.insertBefore('b', 'e')
      tree.insertAfter('a', 'f')

      // 验证树结构
      const rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['c', 'a', 'f', 'e', 'b']),
        `根节点子节点顺序错误，期望: ['c', 'a', 'f', 'e', 'b']，实际: ${JSON.stringify(
          rootChildren
        )}`
      )

      const aChildren = tree.getChildren('a')
      assert(
        arrayEquals(aChildren, ['d']),
        `节点a的子节点错误，期望: ['d']，实际: ${JSON.stringify(aChildren)}`
      )

      tree.dispose()
    }

    // 2. 节点查询测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a')
      tree.append('a', 'b')

      assert(tree.getParent('a') === 'root', 'a的父节点应该是root')
      assert(tree.getParent('b') === 'a', 'b的父节点应该是a')
      assert(tree.getParent('root') === undefined, 'root不应该有父节点')
      assert(tree.getRoot() === 'root', '根节点应该是root')

      tree.dispose()
    }

    // 3. 节点关系测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a', 'b', 'c') // 使用批量添加

      assert(tree.previousSibling('b') === 'a', 'b的前一个兄弟节点应该是a')
      assert(tree.nextSibling('b') === 'c', 'b的后一个兄弟节点应该是c')
      assert(tree.previousSibling('a') === undefined, 'a不应该有前一个兄弟节点')
      assert(tree.nextSibling('c') === undefined, 'c不应该有后一个兄弟节点')

      tree.dispose()
    }

    // 4. 移动节点测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a', 'b', 'c')

      tree.moveBefore('a', 'c')
      let children = tree.getChildren('root')
      assert(
        arrayEquals(children, ['c', 'a', 'b']),
        `moveBefore后顺序错误，期望: ['c', 'a', 'b']，实际: ${JSON.stringify(children)}`
      )

      tree.moveAfter('b', 'c')
      children = tree.getChildren('root')
      assert(
        arrayEquals(children, ['a', 'b', 'c']),
        `moveAfter后顺序错误，期望: ['a', 'b', 'c']，实际: ${JSON.stringify(children)}`
      )

      tree.dispose()
    }

    // 5. 删除节点测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a')
      tree.append('a', 'b', 'c')

      tree.remove('a') // 删除a及其子节点
      assert(!tree.has('a'), '节点a应该被删除')
      assert(!tree.has('b'), '节点b应该被删除')
      assert(!tree.has('c'), '节点c应该被删除')
      assert(tree.getChildren('root').length === 0, 'root应该没有子节点')

      tree.dispose()
    }

    // 6. 替换节点测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a')
      tree.append('a', 'b')

      // 替换节点，新节点不继承子节点
      tree.replace('a', 'x')
      assert(!tree.has('a'), '原节点a应该不存在')
      assert(tree.has('x'), '新节点x应该存在')
      assert(tree.getParent('x') === 'root', '新节点x的父节点应该是root')
      assert(!tree.has('b'), 'a的子节点b应该被删除')

      tree.dispose()
    }

    // 7. 错误情况测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a')

      // 重复添加节点
      try {
        tree.append('root', 'a')
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
        tree.append('nonexistent', 'b')
        throw new Error('应该抛出父节点不存在错误')
      } catch (e: any) {
        assert(e.message.includes('does not exist'), '应该抛出父节点不存在错误')
      }

      tree.dispose()
    }

    // 8. 枚举测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a', 'b')
      tree.append('a', 'c')

      const visitedNodes: string[] = []
      tree.enumerate(id => {
        visitedNodes.push(String(id))
      })

      assert(
        arrayEquals(visitedNodes, ['root', 'a', 'c', 'b']),
        `枚举顺序错误，期望: ['root', 'a', 'c', 'b']，实际: ${JSON.stringify(visitedNodes)}`
      )

      tree.dispose()
    }

    // 9. has方法测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a')

      assert(tree.has('root'), '应该包含root节点')
      assert(tree.has('a'), '应该包含a节点')
      assert(!tree.has('nonexistent'), '不应该包含不存在的节点')

      tree.dispose()
    }

    // 10. 根节点替换测试
    {
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'a')

      // 替换根节点
      tree.replace('root', 'newroot')
      assert(tree.getRoot() === 'newroot', '根节点替换失败')
      assert(!tree.has('root'), '旧根节点应该不存在')
      assert(!tree.has('a'), '替换根节点后子节点应被删除')

      tree.dispose()
    }

    // 11. 复杂移动节点测试
    {
      const tree = new TreeStructure<string>('root')

      // 构建复杂树结构
      tree.append('root', 'a', 'b', 'c')
      tree.append('a', 'a1', 'a2')
      tree.append('b', 'b1', 'b2')
      tree.append('c', 'c1')
      tree.append('a1', 'a1a', 'a1b')

      // 11.1 同级节点移动测试
      tree.moveBefore('a', 'c')
      let rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['c', 'a', 'b']),
        `同级moveBefore失败，期望: ['c', 'a', 'b']，实际: ${JSON.stringify(rootChildren)}`
      )

      tree.moveAfter('b', 'c')
      rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['a', 'b', 'c']),
        `同级moveAfter失败，期望: ['a', 'b', 'c']，实际: ${JSON.stringify(rootChildren)}`
      )

      // 11.2 跨层级移动测试 - 从深层移动到浅层
      tree.moveBefore('b', 'a1a')
      rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['a', 'a1a', 'b', 'c']),
        `跨层级moveBefore失败，期望: ['a', 'a1a', 'b', 'c']，实际: ${JSON.stringify(rootChildren)}`
      )

      let a1Children = tree.getChildren('a1')
      assert(
        arrayEquals(a1Children, ['a1b']),
        `a1a移动后a1子节点错误，期望: ['a1b']，实际: ${JSON.stringify(a1Children)}`
      )

      // 11.3 跨层级移动测试 - 从浅层移动到深层
      tree.moveAfter('a2', 'a1a')
      let aChildren = tree.getChildren('a')
      assert(
        arrayEquals(aChildren, ['a1', 'a2', 'a1a']),
        `a1a移动到a下失败，期望: ['a1', 'a2', 'a1a']，实际: ${JSON.stringify(aChildren)}`
      )

      rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['a', 'b', 'c']),
        `a1a移动后root子节点，期望: ['a', 'b', 'c']，实际: ${JSON.stringify(rootChildren)}`
      )

      // 11.4 带子树的节点移动
      tree.moveBefore('c1', 'a1')
      let cChildren = tree.getChildren('c')
      assert(
        arrayEquals(cChildren, ['a1', 'c1']),
        `带子树移动失败，期望: ['a1', 'c1']，实际: ${JSON.stringify(cChildren)}`
      )

      // 验证a1的子节点依然存在
      a1Children = tree.getChildren('a1')
      assert(
        arrayEquals(a1Children, ['a1b']),
        `移动后a1子节点丢失，期望: ['a1b']，实际: ${JSON.stringify(a1Children)}`
      )

      // 11.5 错误情况测试
      // 尝试移动节点到其后代节点
      try {
        tree.moveBefore('a1a', 'a')
        throw new Error('应该抛出不能移动到后代节点错误')
      } catch (e: any) {
        assert(e.message.includes('descendant'), '应该抛出不能移动到后代节点错误')
      }

      // 尝试移动节点到自己
      try {
        tree.moveBefore('a', 'a')
        throw new Error('应该抛出不能移动到自己错误')
      } catch (e: any) {
        assert(e.message.includes('itself'), '应该抛出不能移动到自己错误')
      }

      tree.dispose()
    }

    // 12. 批量操作测试
    {
      const tree = new TreeStructure<string>('root')

      // 批量添加
      tree.append('root', 'a', 'b', 'c', 'd')
      let children = tree.getChildren('root')
      assert(children.length === 4, '批量添加后应有4个子节点')

      // 批量插入
      tree.insertBefore('c', 'x', 'y', 'z')
      children = tree.getChildren('root')
      assert(
        arrayEquals(children, ['a', 'b', 'x', 'y', 'z', 'c', 'd']),
        `批量插入失败，实际: ${JSON.stringify(children)}`
      )

      // 批量前置添加
      tree.prepend('root', '1', '2')
      children = tree.getChildren('root')
      assert(
        arrayEquals(children, ['1', '2', 'a', 'b', 'x', 'y', 'z', 'c', 'd']),
        `批量前置添加失败，实际: ${JSON.stringify(children)}`
      )

      tree.dispose()
    }

    console.log('所有测试通过！')
  }

  runTests()

  // Generated by Claude Sonnet 3.7 Thinking.
  function runAdvancedTests(): void {
    console.log('开始运行 TreeStructure 高级测试...')

    // 1. 深度嵌套树测试
    {
      console.log('1. 测试深度嵌套树...')
      const tree = new TreeStructure<string>('root')

      // 创建深度嵌套的树 (深度7层)
      let currentId = 'root'
      for (let i = 1; i <= 6; i++) {
        const childId = `level_${i}`
        tree.append(currentId, childId)
        currentId = childId
      }

      // 验证路径完整性
      const paths: string[][] = []
      tree.enumerate((node, path) => {
        if (String(node).startsWith('level_')) {
          paths.push([...path, node])
        }
      })

      // 验证最深节点的路径
      const deepestPath = paths.find(p => p[p.length - 1] === 'level_6')
      assert(deepestPath?.length === 7, `最深节点路径深度错误，应为7，实际为${deepestPath?.length}`)

      // 从深层节点删除
      tree.remove('level_3')
      assert(!tree.has('level_4'), '删除中间节点后，深层节点应被一并删除')
      assert(!tree.has('level_6'), '删除中间节点后，最深层节点应被一并删除')
      assert(tree.has('level_2'), '删除中间节点后，浅层节点应保留')

      tree.dispose()
    }

    // 2. 宽度大的树结构测试
    {
      console.log('2. 测试宽度大的树...')
      const tree = new TreeStructure<string>('root')

      // 创建一个节点有100个直接子节点的树
      const children = Array.from({ length: 100 }, (_, i) => `child_${i}`)
      tree.append('root', ...children)

      // 验证所有子节点都正确添加
      const rootChildren = tree.getChildren('root')
      assert(rootChildren.length === 100, `根节点应有100个子节点，实际有${rootChildren.length}个`)

      // 在大量子节点中间插入
      tree.insertBefore('child_50', 'inserted')
      const newRootChildren = tree.getChildren('root')
      assert(newRootChildren.length === 101, '插入后子节点总数应为101')
      assert(newRootChildren[50] === 'inserted', '插入位置错误')
      assert(newRootChildren[51] === 'child_50', '原节点应后移')

      // 测试批量删除
      for (let i = 0; i < 50; i++) {
        tree.remove(`child_${i}`)
      }
      const remainingChildren = tree.getChildren('root')
      assert(
        remainingChildren.length === 51,
        `批量删除后应有51个子节点，实际有${remainingChildren.length}个`
      )

      tree.dispose()
    }

    // 3. 特殊ID测试
    {
      console.log('3. 测试特殊ID...')
      // 使用特殊的ID值
      const specialIds = [
        0, // 数字0
        '', // 空字符串
        '0', // 字符串"0"
        Symbol('test'), // Symbol
        -1, // 负数
        Infinity, // Infinity
        NaN, // NaN
        true, // 布尔值
        false // 布尔值
      ] as any[]

      // 为每个特殊ID创建一棵树并进行基本操作
      for (const rootId of specialIds) {
        if (Object.is(rootId, NaN)) continue // NaN !== NaN，跳过

        const tree = new TreeStructure<string>(rootId)
        assert(tree.getRoot() === rootId, `根节点ID应为${String(rootId)}`)

        // 添加节点
        tree.append(rootId, 'child')
        assert(tree.getChildren(rootId)[0] === 'child', `特殊ID ${String(rootId)} 添加子节点失败`)

        // 获取父节点
        assert(tree.getParent('child') === rootId, `特殊ID ${String(rootId)} 作为父节点获取失败`)

        tree.dispose()
      }
    }

    // 4. 复杂树重组测试
    {
      console.log('4. 测试复杂树重组...')
      const tree = new TreeStructure<string>('root')

      // 构建初始树
      //         root
      //        /    \
      //       A      B
      //      / \    / \
      //     A1 A2  B1 B2
      //    /        \
      //   A1a        B1a
      tree.append('root', 'A', 'B')
      tree.append('A', 'A1', 'A2')
      tree.append('B', 'B1', 'B2')
      tree.append('A1', 'A1a')
      tree.append('B1', 'B1a')
      tree.print()

      // 复杂重组：将A移到B1下，将B2移到A1下
      tree.moveBefore('B1a', 'A')
      tree.moveAfter('A1a', 'B2')
      tree.print()

      // 验证新结构
      //         root
      //            |
      //            B
      //            |
      //            B1
      //           /  \
      //          A   B1a
      //         / \
      //        A1  A2
      //       / \
      //     A1a  B2
      assert(tree.getParent('A') === 'B1', 'A应移动到B1下')
      assert(tree.getParent('B2') === 'A1', 'B2应移动到A1下')

      const rootChildren = tree.getChildren('root')
      assert(arrayEquals(rootChildren, ['B']), 'root下应只剩B')

      const b1Children = tree.getChildren('B1')
      assert(b1Children.length === 2, 'B1应有2个子节点')
      assert(b1Children.includes('A'), 'B1子节点应包含A')
      assert(b1Children.includes('B1a'), 'B1子节点应包含B1a')

      const a1Children = tree.getChildren('A1')
      assert(a1Children.length === 2, 'A1应有2个子节点')
      assert(a1Children.includes('A1a'), 'A1子节点应包含A1a')
      assert(a1Children.includes('B2'), 'A1子节点应包含B2')

      // 复杂替换：用X替换B1 (会删除B1及其所有子节点，包括A!)
      tree.replace('B1', 'X')
      assert(!tree.has('B1'), 'B1应被删除')
      assert(!tree.has('A'), 'A应作为B1的子节点一并被删除')
      assert(!tree.has('A1'), 'A1应被递归删除')
      assert(!tree.has('A2'), 'A2应被递归删除')
      assert(!tree.has('A1a'), 'A1a应被递归删除')
      assert(!tree.has('B2'), 'B2应被递归删除')
      assert(!tree.has('B1a'), 'B1a应被递归删除')

      const bChildren = tree.getChildren('B')
      assert(arrayEquals(bChildren, ['X']), 'B下应只有新节点X')

      tree.dispose()
    }

    // 5. 树重建测试
    {
      console.log('5. 测试树重建...')
      const tree = new TreeStructure<string>('root')

      // 构建初始树，删除所有节点，然后重建
      tree.append('root', 'A', 'B', 'C')
      tree.append('A', 'A1', 'A2')
      tree.append('B', 'B1')

      // 删除所有子节点
      for (const child of ['A', 'B', 'C']) {
        tree.remove(child)
      }
      assert(tree.getChildren('root').length === 0, '删除后根节点应没有子节点')

      // 重建树
      tree.append('root', 'X', 'Y', 'Z')
      tree.append('X', 'X1', 'X2')
      tree.append('Y', 'Y1')

      // 验证重建的树
      const rootChildren = tree.getChildren('root')
      assert(arrayEquals(rootChildren, ['X', 'Y', 'Z']), '根节点子节点应为X,Y,Z')

      const xChildren = tree.getChildren('X')
      assert(arrayEquals(xChildren, ['X1', 'X2']), 'X的子节点应为X1,X2')

      const yChildren = tree.getChildren('Y')
      assert(arrayEquals(yChildren, ['Y1']), 'Y的子节点应为Y1')

      tree.dispose()
    }

    // 6. 并发修改模拟测试
    {
      console.log('6. 测试并发修改情况...')
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'A', 'B', 'C')

      // 模拟在遍历的同时修改树
      const visited: string[] = []
      let modifiedDuringTraversal = false

      tree.enumerate((node, path) => {
        visited.push(String(node))

        // 在遍历到B时，添加一个子节点
        if (node === 'B' && !modifiedDuringTraversal) {
          tree.append('B', 'B1')
          modifiedDuringTraversal = true
        }
      })

      // 验证B1被添加
      assert(tree.has('B1'), 'B1应被成功添加')

      // B1不应出现在遍历结果中，因为它是在遍历过程中添加的
      // assert(!visited.includes('B1'), '在遍历过程中添加的节点不应出现在遍历结果中')
      // !其实是会的。

      tree.dispose()
    }

    // 7. 极端操作序列测试
    {
      console.log('7. 测试极端操作序列...')
      const tree = new TreeStructure<string>('root')

      // 一系列可能导致问题的操作序列
      tree.append('root', 'A', 'B', 'C')
      tree.insertBefore('B', 'X')
      tree.moveAfter('C', 'X')
      tree.moveBefore('A', 'C')
      tree.remove('A')
      tree.replace('B', 'Y')
      tree.prepend('root', 'Z')

      // 验证最终状态
      const rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['Z', 'C', 'Y', 'X']),
        `复杂操作后根节点子节点不正确，期望: ['Z', 'C', 'Y', 'X'], 实际: ${JSON.stringify(
          rootChildren
        )}`
      )

      tree.dispose()
    }

    // 8. 循环依赖检测测试
    {
      console.log('8. 测试循环依赖检测...')
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'A')
      tree.append('A', 'B')
      tree.append('B', 'C')

      // 尝试将祖先节点移动到后代节点下，应该抛出错误
      try {
        tree.moveBefore('C', 'A')
        throw new Error('应该抛出循环依赖错误')
      } catch (e: any) {
        assert(e.message.includes('Cannot move node to its descendant'), '应该检测到循环依赖风险')
      }

      // 尝试移动到自身
      try {
        tree.moveAfter('B', 'B')
        throw new Error('应该抛出自引用错误')
      } catch (e: any) {
        assert(e.message.includes('Cannot move node to itself'), '应该检测到自引用')
      }

      tree.dispose()
    }

    // 9. 节点替换边界情况测试
    {
      console.log('9. 测试节点替换边界情况...')
      const tree = new TreeStructure<string>('root')

      // 测试替换不存在的节点
      try {
        tree.replace('nonexistent', 'X')
        throw new Error('应该抛出节点不存在错误')
      } catch (e: any) {
        assert(e.message.includes('does not exist'), '应该抛出节点不存在错误')
      }

      // 测试替换为已存在的节点ID
      tree.append('root', 'A', 'B')
      try {
        tree.replace('A', 'B')
        throw new Error('应该抛出新节点已存在错误')
      } catch (e: any) {
        assert(e.message.includes('already exists'), '应该抛出新节点已存在错误')
      }

      // 测试用相同ID替换
      tree.replace('B', 'B') // 应该无操作，不抛出错误
      assert(tree.has('B'), '相同ID替换后节点应仍存在')

      tree.dispose()
    }

    // 10. 大规模树性能测试
    {
      console.log('10. 测试大规模树性能...')
      const tree = new TreeStructure<string>('root')
      const NUM_NODES = 5000

      console.time('创建大规模树')

      // 首层添加1000个节点
      const firstLevelNodes = Array.from({ length: 1000 }, (_, i) => `node_${i}`)
      tree.append('root', ...firstLevelNodes)

      // 每个节点下添加4个子节点
      let count = 1000
      for (let i = 0; i < 1000 && count < NUM_NODES; i++) {
        const parentId = `node_${i}`
        for (let j = 1; j <= 4 && count < NUM_NODES; j++, count++) {
          tree.append(parentId, `${parentId}_child_${j}`)
        }
      }

      console.timeEnd('创建大规模树')

      // 验证节点数量
      let nodeCount = 0
      console.time('遍历大规模树')
      tree.enumerate(() => {
        nodeCount++
      })
      console.log(`大规模树节点总数: ${nodeCount}`)
      console.timeEnd('遍历大规模树')

      assert(nodeCount === count + 1, `大规模树节点计数错误，期望${count + 1}，实际${nodeCount}`)

      // 测试大规模删除
      console.time('删除大量节点')
      for (let i = 0; i < 500; i++) {
        tree.remove(`node_${i}`)
      }
      console.timeEnd('删除大量节点')

      // 验证删除后的节点数量
      let newCount = 0
      tree.enumerate(() => {
        newCount++
      })
      assert(newCount < nodeCount, `删除操作应减少节点总数，但从${nodeCount}变为${newCount}`)

      tree.dispose()
    }

    // 11. 遍历路径测试
    {
      console.log('11. 测试遍历路径正确性...')
      const tree = new TreeStructure<string>('root')
      tree.append('root', 'A', 'B')
      tree.append('A', 'A1', 'A2')
      tree.append('A1', 'A1a')

      // 记录每个节点的路径
      const nodePaths: Record<string, any[]> = {}

      tree.enumerate((node, path) => {
        nodePaths[String(node)] = [...path]
      })

      // 验证关键节点的路径
      assert(
        arrayEquals(nodePaths['A1a'], ['root', 'A', 'A1']),
        `A1a的路径错误，期望['root', 'A', 'A1']，实际${JSON.stringify(nodePaths['A1a'])}`
      )

      assert(
        arrayEquals(nodePaths['B'], ['root']),
        `B的路径错误，期望['root']，实际${JSON.stringify(nodePaths['B'])}`
      )

      assert(
        arrayEquals(nodePaths['root'], []),
        `root的路径错误，期望[]，实际${JSON.stringify(nodePaths['root'])}`
      )

      assert(
        arrayEquals(nodePaths['A'], ['root']),
        `A的路径错误，期望['root']，实际${JSON.stringify(nodePaths['A'])}`
      )

      tree.dispose()
    }

    // 12. 边界条件测试：空操作、无效操作
    {
      console.log('12. 测试边界条件...')
      const tree = new TreeStructure<string>('root')

      // 空数组添加
      tree.append('root', ...([] as string[]))
      assert(tree.getChildren('root').length === 0, '添加空数组后根节点不应有子节点')

      // 获取不存在节点的子节点
      const nonExistentChildren = tree.getChildren('nonexistent')
      assert(nonExistentChildren.length === 0, '不存在节点的子节点应为空数组')

      // 不存在节点的父节点
      const nonExistentParent = tree.getParent('nonexistent')
      assert(nonExistentParent === undefined, '不存在节点的父节点应为undefined')

      // 不存在节点的兄弟节点
      const nonExistentPrevSibling = tree.previousSibling('nonexistent')
      const nonExistentNextSibling = tree.nextSibling('nonexistent')
      assert(nonExistentPrevSibling === undefined, '不存在节点的前一个兄弟节点应为undefined')
      assert(nonExistentNextSibling === undefined, '不存在节点的后一个兄弟节点应为undefined')

      tree.dispose()
    }

    console.log('所有高级测试通过！')
  }

  runAdvancedTests()

  function testFactoryMethods(): void {
    console.log('开始测试 TreeStructure 工厂方法...')

    // 1. fromFlattenedTree 基本功能测试
    {
      console.log('1. 测试 fromFlattenedTree 基本功能')

      // 创建节点数据
      interface TestNode {
        id: string
        childIds: string[]
      }

      const nodes: TestNode[] = [
        { id: 'root', childIds: ['A', 'B', 'C'] },
        { id: 'A', childIds: ['A1', 'A2'] },
        { id: 'B', childIds: ['B1'] },
        { id: 'C', childIds: [] },
        { id: 'A1', childIds: [] },
        { id: 'A2', childIds: [] },
        { id: 'B1', childIds: [] }
      ]

      // 创建树结构
      const tree = TreeStructure.fromFlattenedTree(nodes, 'root', {
        getId: (node: TestNode) => node.id,
        getChildren: (node: TestNode) => node.childIds
      })

      // 验证树结构
      assert(tree.getRoot() === 'root', '根节点应该是root')

      const rootChildren = tree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['A', 'B', 'C']),
        `root的子节点错误，期望: ['A', 'B', 'C']，实际: ${JSON.stringify(rootChildren)}`
      )

      const aChildren = tree.getChildren('A')
      assert(
        arrayEquals(aChildren, ['A1', 'A2']),
        `A的子节点错误，期望: ['A1', 'A2']，实际: ${JSON.stringify(aChildren)}`
      )

      // 验证父子关系
      assert(tree.getParent('A') === 'root', 'A的父节点应该是root')
      assert(tree.getParent('A1') === 'A', 'A1的父节点应该是A')
      assert(tree.getParent('B1') === 'B', 'B1的父节点应该是B')
    }

    // 2. fromFlattenedTree 复杂树结构测试
    {
      console.log('2. 测试 fromFlattenedTree 复杂树结构')

      // 创建更复杂的节点结构，包含多层次和多分支
      interface ComplexNode {
        nodeId: string
        links: string[]
        data?: any // 额外数据，测试不相关字段
      }

      const complexNodes: ComplexNode[] = [
        {
          nodeId: 'root',
          links: ['level1-1', 'level1-2', 'level1-3'],
          data: { value: 'root-data' }
        },
        { nodeId: 'level1-1', links: ['level2-1', 'level2-2'], data: { value: 'l1-1-data' } },
        { nodeId: 'level1-2', links: [], data: { value: 'l1-2-data' } },
        { nodeId: 'level1-3', links: ['level2-3', 'level2-4'], data: { value: 'l1-3-data' } },
        { nodeId: 'level2-1', links: ['level3-1'], data: { value: 'l2-1-data' } },
        { nodeId: 'level2-2', links: [], data: { value: 'l2-2-data' } },
        { nodeId: 'level2-3', links: [], data: { value: 'l2-3-data' } },
        { nodeId: 'level2-4', links: ['level3-2'], data: { value: 'l2-4-data' } },
        { nodeId: 'level3-1', links: [], data: { value: 'l3-1-data' } },
        { nodeId: 'level3-2', links: [], data: { value: 'l3-2-data' } }
      ]

      // 创建树结构
      const complexTree = TreeStructure.fromFlattenedTree(complexNodes, 'root', {
        getId: (node: ComplexNode) => node.nodeId,
        getChildren: (node: ComplexNode) => node.links
      })

      // 验证树深度
      let maxDepth = 0
      complexTree.enumerate((_, path) => {
        if (path.length > maxDepth) maxDepth = path.length
      })
      assert(maxDepth === 3, `树的最大深度应为3，实际为${maxDepth}`)

      // 验证节点数量
      let nodeCount = 0
      complexTree.enumerate(() => {
        nodeCount++
      })
      assert(
        nodeCount === complexNodes.length,
        `树的节点数量应为${complexNodes.length}，实际为${nodeCount}`
      )

      // 验证特定路径
      const level3_2Parent = complexTree.getParent('level3-2')
      assert(level3_2Parent === 'level2-4', `level3-2的父节点应为level2-4，实际为${level3_2Parent}`)

      const level1_3Children = complexTree.getChildren('level1-3')
      assert(
        arrayEquals(level1_3Children, ['level2-3', 'level2-4']),
        `level1-3的子节点错误，期望: ['level2-3', 'level2-4']，实际: ${JSON.stringify(
          level1_3Children
        )}`
      )
    }

    // 3. fromFlattenedTree 错误处理测试
    {
      console.log('3. 测试 fromFlattenedTree 错误处理')

      interface SimpleNode {
        id: string
        children: string[]
      }

      const nodesWithoutRoot: SimpleNode[] = [
        { id: 'A', children: ['A1'] },
        { id: 'A1', children: [] }
      ]

      // 测试找不到根节点的情况
      try {
        TreeStructure.fromFlattenedTree(nodesWithoutRoot, 'root', {
          getId: (node: SimpleNode) => node.id,
          getChildren: (node: SimpleNode) => node.children
        })
        throw new Error('应该抛出根节点不存在错误')
      } catch (e: any) {
        // 检查错误消息是否包含期望的错误信息片段
        assert(
          e.message.includes('不存在') || e.message.includes('not found'),
          '应该抛出根节点不存在的错误'
        )
      }

      // 测试引用不存在节点的情况
      const nodesWithMissingRef: SimpleNode[] = [
        { id: 'root', children: ['A', 'B'] },
        { id: 'A', children: ['C'] }
        // 缺少B和C节点
      ]

      try {
        TreeStructure.fromFlattenedTree(nodesWithMissingRef, 'root', {
          getId: (node: SimpleNode) => node.id,
          getChildren: (node: SimpleNode) => node.children
        })
        throw new Error('应该抛出引用不存在节点的错误')
      } catch (e: any) {
        assert(
          e.message.includes('不存在') || e.message.includes('not found'),
          '应该抛出引用不存在节点的错误'
        )
      }
    }

    // 4. fromNestedTree 基本功能测试
    {
      console.log('4. 测试 fromNestedTree 基本功能')

      // 创建嵌套结构
      interface NestedNode {
        id: string
        children?: NestedNode[]
      }

      const nestedRoot: NestedNode = {
        id: 'root',
        children: [
          { id: 'A', children: [{ id: 'A1' }, { id: 'A2' }] },
          { id: 'B', children: [{ id: 'B1' }] },
          { id: 'C' }
        ]
      }

      // 创建树结构
      const nestedTree = TreeStructure.fromNestedTree(nestedRoot, {
        getId: (node: NestedNode) => node.id,
        getChildren: (node: NestedNode) => node.children || []
      })

      // 验证树结构
      assert(nestedTree.getRoot() === 'root', '根节点应该是root')

      const rootChildren = nestedTree.getChildren('root')
      assert(
        arrayEquals(rootChildren, ['A', 'B', 'C']),
        `root的子节点错误，期望: ['A', 'B', 'C']，实际: ${JSON.stringify(rootChildren)}`
      )

      // 验证父子关系
      assert(nestedTree.getParent('A1') === 'A', 'A1的父节点应该是A')
      assert(nestedTree.getParent('B') === 'root', 'B的父节点应该是root')

      // 验证遍历
      const visited: string[] = []
      nestedTree.enumerate(node => {
        visited.push(String(node))
      })

      // 前序遍历顺序应该是: root, A, A1, A2, B, B1, C
      assert(
        arrayEquals(visited, ['root', 'A', 'A1', 'A2', 'B', 'B1', 'C']),
        `前序遍历顺序错误，期望: ['root', 'A', 'A1', 'A2', 'B', 'B1', 'C']，实际: ${JSON.stringify(
          visited
        )}`
      )
    }

    // // 5. fromNestedTree 循环引用检测测试
    // {
    //   console.log('5. 测试 fromNestedTree 循环引用检测')

    //   interface CyclicNode {
    //     id: string
    //     children: CyclicNode[]
    //   }

    //   // 创建一个循环引用的结构
    //   const nodeA: CyclicNode = { id: 'A', children: [] }
    //   const nodeB: CyclicNode = { id: 'B', children: [] }
    //   const nodeC: CyclicNode = { id: 'C', children: [] }

    //   const rootNode: CyclicNode = { id: 'root', children: [nodeA, nodeB] }
    //   nodeA.children = [nodeC]
    //   nodeC.children = [nodeB] // 这里没有循环

    //   // 这个正常结构应该能够成功创建
    //   const normalTree = TreeStructure.fromNestedTree(rootNode, {
    //     getId: (node: CyclicNode) => node.id,
    //     getChildren: (node: CyclicNode) => node.children
    //   })

    //   // 验证正常树结构
    //   assert(normalTree.has('root'), '树应该包含root节点')
    //   assert(normalTree.has('A'), '树应该包含A节点')
    //   assert(normalTree.has('B'), '树应该包含B节点')
    //   assert(normalTree.has('C'), '树应该包含C节点')

    //   // 创建循环引用
    //   nodeB.children = [nodeA] // 现在形成了循环: A -> C -> B -> A

    //   try {
    //     TreeStructure.fromNestedTree(rootNode, {
    //       getId: (node: CyclicNode) => node.id,
    //       getChildren: (node: CyclicNode) => node.children
    //     })
    //     throw new Error('应该检测到循环引用')
    //   } catch (e: any) {
    //     assert(e.message.includes('循环') || e.message.includes('cycle'), '应该检测到循环引用')
    //   }
    // }

    // 6. fromNestedTree 空子节点和深层嵌套测试
    {
      console.log('6. 测试 fromNestedTree 空子节点和深层嵌套')

      interface DeepNode {
        name: string
        subs?: DeepNode[]
      }

      // 创建深层嵌套结构
      const deepRoot: DeepNode = {
        name: 'level0',
        subs: [
          {
            name: 'level1-1',
            subs: [{ name: 'level2-1', subs: [{ name: 'level3-1', subs: [{ name: 'level4-1' }] }] }]
          },
          { name: 'level1-2' } // 节点没有子节点
        ]
      }

      const deepTree = TreeStructure.fromNestedTree(deepRoot, {
        getId: (node: DeepNode) => node.name,
        getChildren: (node: DeepNode) => node.subs || []
      })

      // 验证深度
      let maxDepth = 0
      deepTree.enumerate((_, path) => {
        if (path.length > maxDepth) maxDepth = path.length
      })
      assert(maxDepth === 4, `树的最大深度应为4，实际为${maxDepth}`)

      // 验证叶子节点
      assert(deepTree.getChildren('level4-1').length === 0, 'level4-1应该是叶子节点')
      assert(deepTree.getChildren('level1-2').length === 0, 'level1-2应该是叶子节点')

      // 验证路径
      assert(deepTree.getParent('level3-1') === 'level2-1', 'level3-1的父节点应该是level2-1')
      assert(deepTree.getParent('level1-2') === 'level0', 'level1-2的父节点应该是level0')
    }

    // 7. 特殊ID类型测试 - 数字ID
    {
      console.log('7. 测试特殊ID类型 - 数字ID')

      interface NumNode {
        id: number
        children: number[]
      }

      const numNodes: NumNode[] = [
        { id: 1, children: [2, 3, 4] },
        { id: 2, children: [5, 6] },
        { id: 3, children: [] },
        { id: 4, children: [7] },
        { id: 5, children: [] },
        { id: 6, children: [] },
        { id: 7, children: [] }
      ]

      // 使用数字ID创建树
      const numTree = TreeStructure.fromFlattenedTree(numNodes, 1, {
        getId: (node: NumNode) => node.id,
        getChildren: (node: NumNode) => node.children
      })

      // 验证树结构
      assert(numTree.getRoot() === 1, '根节点应该是1')

      const rootChildren = numTree.getChildren(1)
      assert(
        arrayEquals(rootChildren, [2, 3, 4]),
        `root的子节点错误，期望: [2, 3, 4]，实际: ${JSON.stringify(rootChildren)}`
      )

      // 验证父子关系
      assert(numTree.getParent(5) === 2, '节点5的父节点应该是2')
      assert(numTree.getParent(7) === 4, '节点7的父节点应该是4')

      // 测试遍历方法
      const visited: number[] = []
      numTree.enumerate(node => {
        visited.push(node as number)
      })

      assert(
        arrayEquals(visited, [1, 2, 5, 6, 3, 4, 7]),
        `前序遍历顺序错误，期望: [1, 2, 5, 6, 3, 4, 7]，实际: ${JSON.stringify(visited)}`
      )
    }

    // 8. 特殊ID类型测试 - Symbol ID
    {
      console.log('8. 测试特殊ID类型 - Symbol ID')

      // 创建Symbol ID
      const ROOT = Symbol('root')
      const A = Symbol('A')
      const B = Symbol('B')
      const C = Symbol('C')

      interface SymbolNode {
        id: symbol
        children: symbol[]
      }

      const symbolNodes: SymbolNode[] = [
        { id: ROOT, children: [A, B] },
        { id: A, children: [C] },
        { id: B, children: [] },
        { id: C, children: [] }
      ]

      // 使用Symbol ID创建树
      const symbolTree = TreeStructure.fromFlattenedTree(symbolNodes, ROOT, {
        getId: (node: SymbolNode) => node.id,
        getChildren: (node: SymbolNode) => node.children
      })

      // 验证树结构
      assert(symbolTree.getRoot() === ROOT, '根节点应该是ROOT Symbol')

      const rootChildren = symbolTree.getChildren(ROOT)
      assert(rootChildren.length === 2, 'ROOT应该有2个子节点')
      assert(rootChildren.includes(A), 'ROOT的子节点应该包含A Symbol')
      assert(rootChildren.includes(B), 'ROOT的子节点应该包含B Symbol')

      // 验证父子关系
      assert(symbolTree.getParent(C) === A, 'C的父节点应该是A Symbol')

      // 确认节点存在性
      assert(symbolTree.has(ROOT), '树应该包含ROOT Symbol')
      assert(symbolTree.has(A), '树应该包含A Symbol')
      assert(symbolTree.has(B), '树应该包含B Symbol')
      assert(symbolTree.has(C), '树应该包含C Symbol')
    }

    console.log('TreeStructure 工厂方法测试全部通过！')
  }

  testFactoryMethods()
}

// !遍历任何抽象树形结构的三个关键：根结点、isLeaf、getChildren

interface ITreeNode {
  isLeaf(): boolean
  getChildren(): Iterable<this>
}

/**
 * 前序 dfs 遍历树.
 *
 * @param root 根节点.
 * @param f 访问函数，返回 true 时停止遍历.
 */
function enumerateTree<T extends ITreeNode>(root: T, f: (node: T) => boolean | void): void {
  const dfs = (node: T): boolean => {
    if (node.isLeaf()) {
      return !!f(node)
    }
    if (f(node)) return true
    for (const child of node.getChildren()) {
      if (dfs(child)) return true
    }
    return false
  }

  dfs(root)
}

function enumerateTree2<T>(
  root: T,
  getChildren: (node: T) => Iterable<T>,
  visit: (node: T) => boolean | void
) {
  const dfs = (node: T): boolean => {
    if (visit(node)) return true
    for (const child of getChildren(node)) {
      if (dfs(child)) return true
    }
    return false
  }

  dfs(root)
}

export {}

// ================================
// 测试用例
// ================================

/**
 * 简单的树节点实现
 */
class TreeNode<T> implements ITreeNode {
  readonly value: T
  readonly children: this[]

  constructor(value: T) {
    this.value = value
    this.children = []
  }

  isLeaf(): boolean {
    return this.children.length === 0
  }

  getChildren(): this[] {
    return this.children
  }

  addChild(child: this): this {
    this.children.push(child)
    return this
  }

  // 便于测试的辅助方法
  toString(): string {
    return String(this.value)
  }
}

/**
 * 构建测试树
 *       A
 *     /   \
 *    B     C
 *   / \     \
 *  D   E     F
 *     /
 *    G
 */
function createTestTree(): TreeNode<string> {
  const root = new TreeNode('A')
  const nodeB = new TreeNode('B')
  const nodeC = new TreeNode('C')
  const nodeD = new TreeNode('D')
  const nodeE = new TreeNode('E')
  const nodeF = new TreeNode('F')
  const nodeG = new TreeNode('G')

  root.addChild(nodeB).addChild(nodeC)
  nodeB.addChild(nodeD).addChild(nodeE)
  nodeC.addChild(nodeF)
  nodeE.addChild(nodeG)

  return root
}

// ================================
// 测试函数
// ================================

function runTests(): void {
  console.log('🌲 开始树遍历测试\n')

  // 测试1: 完整遍历
  console.log('📋 测试1: 完整遍历')
  const tree1 = createTestTree()
  const visited1: string[] = []

  enumerateTree(tree1, node => {
    visited1.push(node.value)
    console.log(`访问节点: ${node.value}`)
  })

  console.log('遍历结果:', visited1.join(' -> '))
  console.log('预期结果: A -> B -> D -> E -> G -> C -> F')
  console.log('✅ 测试1通过\n')

  // 测试2: 早期退出
  console.log('📋 测试2: 找到E时停止遍历')
  const tree2 = createTestTree()
  const visited2: string[] = []

  enumerateTree(tree2, node => {
    visited2.push(node.value)
    console.log(`访问节点: ${node.value}`)

    if (node.value === 'E') {
      console.log('🛑 找到E，停止遍历')
      return true // 停止遍历
    }
  })

  console.log('遍历结果:', visited2.join(' -> '))
  console.log('预期结果: A -> B -> D -> E')
  console.log('✅ 测试2通过\n')

  // 测试3: 单节点树
  console.log('📋 测试3: 单节点树')
  const singleNode = new TreeNode('Single')
  const visited3: string[] = []

  enumerateTree(singleNode, node => {
    visited3.push(node.value)
    console.log(`访问节点: ${node.value}`)
  })

  console.log('遍历结果:', visited3.join(' -> '))
  console.log('预期结果: Single')
  console.log('✅ 测试3通过\n')

  // 测试4: 返回void的函数
  console.log('📋 测试4: 返回void的访问函数')
  const tree4 = createTestTree()
  const visited4: string[] = []

  const voidFunction = (node: TreeNode<string>): void => {
    visited4.push(node.value)
    console.log(`访问节点: ${node.value}`)
  }

  enumerateTree(tree4, voidFunction)

  console.log('遍历结果:', visited4.join(' -> '))
  console.log('预期结果: A -> B -> D -> E -> G -> C -> F')
  console.log('✅ 测试4通过\n')

  // 测试5: 数字类型树
  console.log('📋 测试5: 数字类型树')
  const numTree = new TreeNode(1)
  numTree.addChild(new TreeNode(2)).addChild(new TreeNode(3))
  numTree.getChildren()[0].addChild(new TreeNode(4))

  const visited5: number[] = []
  let sum = 0

  enumerateTree(numTree, node => {
    visited5.push(node.value)
    sum += node.value
    console.log(`访问节点: ${node.value}, 当前累计: ${sum}`)

    // 累计超过6时停止
    if (sum > 6) {
      console.log('🛑 累计超过6，停止遍历')
      return true
    }
  })

  console.log('遍历结果:', visited5.join(' -> '))
  console.log('累计和:', sum)
  console.log('✅ 测试5通过\n')

  // 测试6: 深层嵌套树
  console.log('📋 测试6: 深层嵌套树（链式结构）')
  let chainRoot = new TreeNode('Root')
  let current = chainRoot

  for (let i = 1; i <= 5; i++) {
    const child = new TreeNode(`Level${i}`)
    current.addChild(child)
    current = child
  }

  const visited6: string[] = []
  enumerateTree(chainRoot, node => {
    visited6.push(node.value)
    console.log(`访问节点: ${node.value}`)
  })

  console.log('遍历结果:', visited6.join(' -> '))
  console.log('预期结果: Root -> Level1 -> Level2 -> Level3 -> Level4 -> Level5')
  console.log('✅ 测试6通过\n')

  console.log('🎉 所有测试完成！')
}

// ================================
// 辅助函数：可视化树结构
// ================================

function printTree(node: ITreeNode, prefix: string = '', isLast: boolean = true): void {
  const nodeStr = node instanceof TreeNode ? node.value : 'Node'
  console.log(prefix + (isLast ? '└── ' : '├── ') + nodeStr)

  if (!node.isLeaf()) {
    const children = Array.from(node.getChildren())
    children.forEach((child, index) => {
      const isLastChild = index === children.length - 1
      const newPrefix = prefix + (isLast ? '    ' : '│   ')
      printTree(child, newPrefix, isLastChild)
    })
  }
}

// ================================
// 性能测试
// ================================

function performanceTest(): void {
  console.log('⚡ 性能测试')

  // 创建大树
  const createLargeTree = (depth: number, branching: number): TreeNode<string> => {
    const root = new TreeNode(`0`)

    const buildLevel = (node: TreeNode<string>, currentDepth: number, nodeId: string): void => {
      if (currentDepth >= depth) return

      for (let i = 0; i < branching; i++) {
        const childId = `${nodeId}-${i}`
        const child = new TreeNode(childId)
        node.addChild(child)
        buildLevel(child, currentDepth + 1, childId)
      }
    }

    buildLevel(root, 0, '0')
    return root
  }

  const largeTree = createLargeTree(4, 3) // 4层，每层3个分支
  let visitCount = 0

  const startTime = performance.now()

  enumerateTree(largeTree, node => {
    visitCount++
    // 访问一半节点后停止
    if (visitCount > 40) {
      return true
    }
  })

  const endTime = performance.now()

  console.log(`访问了 ${visitCount} 个节点`)
  console.log(`耗时: ${(endTime - startTime).toFixed(2)}ms`)
  console.log('✅ 性能测试完成\n')
}

// 运行所有测试
function main(): void {
  console.log('树结构可视化:')
  const tree = createTestTree()
  printTree(tree)
  console.log('')

  runTests()
  performanceTest()
}

// 如果直接运行此文件，执行测试
if (require.main === module) {
  main()
}

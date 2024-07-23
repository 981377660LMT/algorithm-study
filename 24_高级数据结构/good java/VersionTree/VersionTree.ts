/* eslint-disable @typescript-eslint/no-empty-function */

interface IOperation {
  apply(): void
  undo(): void
}

const EMPTY_OPERATION: IOperation = { apply: () => {}, undo: () => {} }
class TreeNode {
  operation: IOperation
  children: TreeNode[] = []

  constructor(operation: IOperation) {
    this.operation = operation
  }
}

/**
 * @alias OfflinePersistentTree 版本树
 */
class VersionTree {
  private readonly _nodes: TreeNode[]
  private _now = 0
  private _version = 0

  /**
   * 初始时版本号为0.
   * @param maxOperation 最大操作(版本)数.
   */
  constructor(maxOperation: number) {
    this._nodes = Array(maxOperation + 1)
    this._nodes[0] = new TreeNode(EMPTY_OPERATION)
  }

  /**
   * 在当前版本上添加一个操作，返回新版本号.
   */
  apply(operation: IOperation): number {
    this._nodes[++this._now] = new TreeNode(operation)
    this._nodes[this._version].children.push(this._nodes[this._now])
    this._version = this._now
    return this._version
  }

  /**
   * 切换到指定版本.
   */
  switchVersion(version: number): void {
    this._version = version
  }

  /**
   * 应用所有操作.
   */
  run(): void {
    this._dfs(this._nodes[0])
  }

  getVersion(): number {
    return this._version
  }

  private _dfs(root: TreeNode): void {
    root.operation.apply()
    const { children } = root
    for (let i = 0; i < children.length; i++) {
      this._dfs(children[i])
    }
    root.operation.undo()
  }
}

export { VersionTree }

if (require.main === module) {
  const tree = new VersionTree(10)
  tree.apply({
    apply: () => console.log('apply 1'),
    undo: () => console.log('undo 1')
  })
  tree.apply({
    apply: () => console.log('apply 2'),
    undo: () => console.log('undo 2')
  })

  tree.run()
}

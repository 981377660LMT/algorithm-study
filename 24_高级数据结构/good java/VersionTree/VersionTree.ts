/* eslint-disable @typescript-eslint/no-empty-function */

interface IStep {
  apply(): void
  undo(): void
}

type QueryFunc = (context: { kth: number; version: number }) => void

const DUMMY_STEP: IStep = { apply: () => {}, undo: () => {} }

class Node {
  step: IStep
  queries: (() => void)[] = []
  children: Node[] = []
  constructor(step: IStep) {
    this.step = step
  }
}

/**
 * 版本树是一种数据结构，用于`离线`管理和跟踪数据的版本历史.
 * 版本指的是数据结构在特定时间点的状态.
 * 每次变更都会生成一个新的版本，版本是immutable的.
 * 每个版本内可以有多个查询，但是只有一个变更(原子操作).
 * 这种方式允许我们快速回溯到历史状态，实现撤销（Undo）和重做（Redo）操作.
 *
 * @alias OfflinePersistentTree 版本树/操作树
 */
class VersionTree {
  private readonly _nodes: Node[]
  private _nodePtr = 0
  private _version = 0
  private _queryCount = 0

  /**
   * 初始时版本号为0.
   * @param maxOperation 最大操作(版本)数.
   */
  constructor(maxOperation: number) {
    this._nodes = Array(maxOperation + 1)
    this._nodes[0] = new Node(DUMMY_STEP)
  }

  /**
   * 在当前版本上添加一个修改，返回新版本号.
   */
  addStep(step: IStep): number {
    this._nodes[++this._nodePtr] = new Node(step)
    this._nodes[this._version].children.push(this._nodes[this._nodePtr])
    this._version = this._nodePtr
    return this._version
  }

  addQuery(query: QueryFunc): void {
    const context = { kth: this._queryCount, version: this._version }
    this._nodes[this._version].queries.push(() => query(context))
    this._queryCount++
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
    this._run(this._nodes[0])
  }

  getVersion(): number {
    return this._version
  }

  private _run(root: Node): void {
    root.step.apply()
    root.queries.forEach(query => query())
    root.children.forEach(child => this._run(child))
    root.step.undo()
  }
}

export { VersionTree }

if (require.main === module) {
  const tree = new VersionTree(10)
  const arr = Array.from({ length: 1e5 }, (_, i) => i)

  // !undo 可以是 update，也可以是 set
  // 如果操作可逆，则用 update，否则用 set
  tree.addStep({
    apply: () => {
      arr.push(1)
    },
    undo: () => {
      arr.pop()
    }
  })

  tree.addQuery(({ kth, version }) => {
    console.log('query', { kth, version })
    console.log(arr[arr.length - 1])
  })

  tree.switchVersion(0)

  tree.addQuery(({ kth, version }) => {
    console.log('query', { kth, version })
    console.log(arr[arr.length - 1])
  })

  tree.run()
}

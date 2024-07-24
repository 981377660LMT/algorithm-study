/* eslint-disable @typescript-eslint/no-empty-function */

interface IStep {
  /**
   * @returns 操作是否成功.
   */
  apply(): boolean
  invert(): void
}

type QueryFunc = (context: { kth: number; version: number }) => void

const DUMMY_STEP: IStep = { apply: () => false, invert: () => {} }

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
 * !初始时版本号为0(没有任何修改)，第一次操作后版本号为1，以此类推.
 *
 * @alias OfflinePersistentTree 版本树/操作树
 */
class VersionTree {
  private readonly _nodes: Node[] = [new Node(DUMMY_STEP)]
  private _version = 0
  private _queryCount = 0

  /**
   * 在当前版本上添加一个修改，返回新版本号.
   */
  addStep(step: IStep): number {
    const newNode = new Node(step)
    this._nodes.push(newNode)
    this._nodes[this._version].children.push(newNode)
    this._version = this._nodes.length - 1
    return this._version
  }

  /**
   * !在当前版本上添加一个切换版本的操作，视为一次修改操作.
   */
  addSwitchVersionStep(version: number): number {
    this._checkVersion(version)
    const newNode = new Node(DUMMY_STEP)
    this._nodes.push(newNode)
    this._nodes[version].children.push(newNode)
    this._version = version
    return this._nodes.length - 1
  }

  /**
   * !切换到指定版本，不视为一次修改操作.
   */
  switchVersion(version: number): void {
    this._checkVersion(version)
    this._version = version
  }

  addQuery(query: QueryFunc): void {
    const context = { kth: this._queryCount++, version: this._version }
    this._nodes[this._version].queries.push(() => query(context))
  }

  /**
   * 应用所有操作.
   */
  commit(): void {
    this._commit(this._nodes[0])
  }

  private _commit(root: Node): void {
    const ok = root.step.apply()
    root.queries.forEach(query => query())
    root.children.forEach(child => this._commit(child))
    ok && root.step.invert()
  }

  private _checkVersion(version: number): void {
    if (version < 0 || version >= this._nodes.length) {
      throw new RangeError(`Invalid version: ${version}`)
    }
  }
}

export { VersionTree }

if (require.main === module) {
  const tree = new VersionTree()
  const arr = Array.from({ length: 1e5 }, (_, i) => i)

  // !invert 可以是 update，也可以是 set
  // 如果操作可逆，则用 update，否则用 set
  tree.addStep({
    apply: () => {
      arr.push(1)
      return true
    },
    invert: () => {
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

  tree.commit()
}

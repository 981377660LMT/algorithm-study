interface IDelegator {
  mutate(): void
  query(): void
  pushDown(): void
}

/**
 * 使用版本机制延迟更新子树的算法模型.
 * !控制逻辑与计算逻辑分离.
 */
class Node {
  private readonly _parent: Node | undefined
  private readonly _delegator: IDelegator
  private _version = 0
  private _parentVersion = 0

  constructor(parent: Node | undefined, delegator: IDelegator) {
    this._parent = parent
    if (parent) this._parentVersion = parent._version
    this._delegator = delegator
  }

  mutate() {
    this._version++
    this._delegator.mutate()
  }

  query(): void {
    this._sync()
    this._delegator.query()
  }

  private _sync(): void {
    if (!this._parent) return
    this._parent._sync()
    if (this._parentVersion !== this._parent._version) {
      this._parentVersion = this._parent._version
      this._parent._delegator.pushDown()
    }
  }
}

export {}

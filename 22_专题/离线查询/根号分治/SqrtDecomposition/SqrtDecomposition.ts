import { createBlock } from './useBlock'

interface Block<E, Id, Q = unknown> {
  /**
   * 在创建块时调用，用于初始化块的值.
   */
  created(): void

  /**
   * 在部分更新块之后时调用，用于将`整个块`的值更新到最新.
   */
  updated(): void

  /**
   * 查询块内左闭右开区间 [start, end) 的值. 0 <= start <= end <= blockSize.
   *
   * !注意查询时需要把块内的懒标记计入影响.
   */
  queryPart(start: number, end: number, queryArg?: Q): E

  /**
   * 更新块内左闭右开区间 [start, end) 的值. 0 <= start <= end <= blockSize.
   */
  updatePart(start: number, end: number, lazy: Id): void

  /**
   * 查询块内所有值.
   */
  queryAll(queryArg?: Q): E

  /**
   * 更新块内所有值,打上懒标记.
   */
  updateAll(lazy: Id): void
}

/**
 * 当区间需要维护的数据难以用半群来描述时，可以考虑根号分块。
 * @deprecated 使用 {@link createBlock}.
 */
class SqrtDecomposition<E, Id, Q = unknown> {
  private readonly _blockSize: number
  private readonly _blocks: Block<E, Id, Q>[]
  private readonly _belong: Uint16Array

  /**
   * @param n 区间长度.
   * @param createBlock 块的工厂函数.
   * @param blockSize 分块大小，一般取 `Math.sqrt(n) + 1`.
   */
  constructor(
    n: number,
    createBlock: (id: number, start: number, end: number) => Block<E, Id, Q>,
    blockSize = (Math.sqrt(n) + 1) | 0
  ) {
    this._blockSize = blockSize
    this._blocks = Array(1 + ((n / blockSize) | 0))
    this._belong = new Uint16Array(n)
    for (let i = 0; i < n; i++) this._belong[i] = (i / blockSize) | 0
    for (let i = 0; i < this._blocks.length; i++) {
      this._blocks[i] = createBlock(i, i * blockSize, Math.min((i + 1) * blockSize, n))
      this._blocks[i].created()
    }
  }

  /**
   * 更新左闭右开区间 `[start, end)` 的值.
   *
   * 0 <= start <= end <= n.
   */
  update(start: number, end: number, lazy: Id): void {
    if (start >= end) {
      return
    }
    const id1 = this._belong[start]
    const id2 = this._belong[end - 1]
    const pos1 = start - id1 * this._blockSize
    const pos2 = end - id2 * this._blockSize
    if (id1 === id2) {
      this._blocks[id1].updatePart(pos1, pos2, lazy)
      this._blocks[id1].updated()
    } else {
      this._blocks[id1].updatePart(pos1, this._blockSize, lazy)
      this._blocks[id1].updated()
      for (let i = id1 + 1; i < id2; i++) this._blocks[i].updateAll(lazy)
      this._blocks[id2].updatePart(0, pos2, lazy)
      this._blocks[id2].updated()
    }
  }

  /**
   * 查询左闭右开区间 `[start, end)` 的值.
   *
   * 0 <= start <= end <= n.
   * @param forEach 遍历每个块的结果.
   */
  query(start: number, end: number, forEach: (blockRes: E) => void, queryArg?: Q) {
    if (start >= end) {
      return
    }
    const id1 = this._belong[start]
    const id2 = this._belong[end - 1]
    const pos1 = start - id1 * this._blockSize
    const pos2 = end - id2 * this._blockSize
    if (id1 === id2) {
      forEach(this._blocks[id1].queryPart(pos1, pos2, queryArg))
      return
    }
    forEach(this._blocks[id1].queryPart(pos1, this._blockSize, queryArg))
    for (let i = id1 + 1; i < id2; i++) forEach(this._blocks[i].queryAll(queryArg))
    forEach(this._blocks[id2].queryPart(0, pos2, queryArg))
  }
}

export { SqrtDecomposition, Block }

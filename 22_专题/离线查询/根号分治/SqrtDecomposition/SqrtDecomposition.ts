interface Block<E, Id> {
  /**
   * 在创建块时调用，用于初始化块的值.
   */
  created(): void

  /**
   * 在部分更新块之后时调用，用于将`整个块`的值更新到最新.
   */
  updated(): void

  /**
   * 查询块内区间 [left, right] 的值. 0 <= left < right < blockSize.
   *
   * !注意查询时需要把块内的懒标记计入影响.
   */
  queryPart(left: number, right: number): E

  /**
   * 更新块内区间 [left, right] 的值. 0 <= left < right < blockSize.
   */
  updatePart(left: number, right: number, lazy: Id): void

  /**
   * 查询块内所有值.
   */
  queryAll(): E

  /**
   * 更新块内所有值,打上懒标记.
   */
  updateAll(lazy: Id): void
}

/**
 * 当区间需要维护的数据难以用半群来描述时，可以考虑根号分块。
 */
class SqrtDecomposition<E, Id> {
  private readonly _blocks: Block<E, Id>[]
  private readonly _left: Uint32Array
  private readonly _right: Uint32Array

  /**
   * @param n 区间长度.
   * @param createBlock 块的工厂函数. 0 <= leftBound < rightBound < n.
   * @param blockSize 分块大小，一般取 `Math.sqrt(n)`.
   */
  constructor(
    n: number,
    /**
     * @param id 当前块的编号.
     * @param leftBound 当前块在原数组中的左边界.0 <= leftBound < n.
     * @param rightBound 当前块在原数组中的右边界.0 <= rightBound < n.
     */
    createBlock: (id: number, leftBound: number, rightBound: number) => Block<E, Id>,
    blockSize = Math.floor(Math.sqrt(n))
  ) {
    const blockCount = Math.floor((n - 1) / blockSize) + 1
    this._blocks = Array(blockCount).fill(null)
    this._left = new Uint32Array(blockCount)
    this._right = new Uint32Array(blockCount)
    for (let i = 0; i < this._blocks.length; i++) {
      const left = i * blockSize
      const right = Math.min((i + 1) * blockSize - 1, n - 1)
      this._left[i] = left
      this._right[i] = right
      this._blocks[i] = createBlock(i, left, right)
      this._blocks[i].created()
    }
  }

  /**
   * 更新区间 `[left, right]` 的值.
   *
   * 0 <= left <= right < n.
   */
  update(left: number, right: number, lazy: Id): void {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._right[i] < left) {
        continue
      }
      if (this._left[i] > right) {
        break
      }
      if (left <= this._left[i] && this._right[i] <= right) {
        block.updateAll(lazy)
      } else {
        const bl = Math.max(this._left[i], left)
        const br = Math.min(this._right[i], right)
        block.updatePart(bl - this._left[i], br - this._left[i], lazy)
        block.updated()
      }
    }
  }

  /**
   * 查询闭区间 `[left, right]` 的值.
   *
   * 0 <= left <= right < n.
   * @param forEach 遍历每个块的结果.
   */
  query(left: number, right: number, forEach: (blockRes: E) => void) {
    for (let i = 0; i < this._blocks.length; i++) {
      const block = this._blocks[i]
      if (this._right[i] < left) {
        continue
      }
      if (this._left[i] > right) {
        break
      }
      if (left <= this._left[i] && this._right[i] <= right) {
        forEach(block.queryAll())
      } else {
        const bl = Math.max(this._left[i], left)
        const br = Math.min(this._right[i], right)
        forEach(block.queryPart(bl - this._left[i], br - this._left[i]))
      }
    }
  }
}

export { SqrtDecomposition, Block }

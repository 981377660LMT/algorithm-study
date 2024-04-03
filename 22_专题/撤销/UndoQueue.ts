/**
 * 满足交换律的操作，即`(op1, op2)` 和 `(op2, op1)` 作用效果相同.
 */
interface IFlaggedCommutativeOperation {
  flag: boolean
  apply(): void
  undo(): void
}

class SimpleQueue<T> {
  private readonly _data: T[] = []
  private _head = 0

  push(value: T): void {
    this._data.push(value)
  }

  pop(): T | undefined {
    return this._data.pop()
  }

  shift(): T | undefined {
    return this._head < this._data.length ? this._data[this._head++] : undefined
  }

  clear(): void {
    this._data.length = 0
    this._head = 0
  }

  toString(): string {
    return this._data.slice(this._head).toString()
  }

  get length(): number {
    return this._data.length - this._head
  }
}

/**
 * 支持撤销操作的队列，每个操作会被应用和撤销 `O(logn)` 次.
 */
class UndoQueue {
  static createFlaggedCommutativeOperation(
    apply: () => void,
    undo: () => void
  ): IFlaggedCommutativeOperation {
    return {
      flag: false,
      apply,
      undo
    }
  }

  private readonly _dq: IFlaggedCommutativeOperation[] = []
  private readonly _bufA: IFlaggedCommutativeOperation[] = []
  private readonly _bufB: SimpleQueue<IFlaggedCommutativeOperation> = new SimpleQueue()

  append(op: IFlaggedCommutativeOperation): void {
    op.flag = false
    this._pushAndDo(op)
  }

  popLeft(): IFlaggedCommutativeOperation | undefined {
    if (!this._dq[this._dq.length - 1].flag) {
      this._popAndUndo()
      while (this._dq.length && this._bufB.length !== this._bufA.length) {
        this._popAndUndo()
      }
      if (this._dq.length === 0) {
        while (this._bufB.length) {
          const res = this._bufB.shift()!
          res.flag = true
          this._pushAndDo(res)
        }
      } else {
        while (this._bufB.length) {
          const res = this._bufB.pop()!
          this._pushAndDo(res)
        }
      }
      while (this._bufA.length) {
        this._pushAndDo(this._bufA.pop()!)
      }
    }

    const res = this._dq.pop()!
    res.undo()
    return res
  }

  empty(): boolean {
    return !this._dq.length
  }

  clear(): void {
    const n = this._dq.length
    for (let _ = 0; _ < n; _++) this.popLeft()
    this._bufA.length = 0
    this._bufB.clear()
  }

  get length(): number {
    return this._dq.length
  }

  private _pushAndDo(op: IFlaggedCommutativeOperation): void {
    this._dq.push(op)
    op.apply()
  }

  private _popAndUndo(): void {
    const res = this._dq.pop()!
    res.undo()
    res.flag ? this._bufA.push(res) : this._bufB.push(res)
  }
}

export { UndoQueue }

if (require.main === module) {
  const queue = new SimpleQueue<number>()
  queue.push(1)
  queue.push(2)
  queue.push(3)
  console.log(queue.shift())
  console.log(queue.shift())
  console.log(queue.shift())
  console.log(queue.shift())
  console.log(queue.length)
}

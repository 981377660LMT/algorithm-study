interface IOperation {
  apply(): void
  undo(): void
}

class UndoStack {
  private readonly _stack: IOperation[] = []

  push(op: IOperation): void {
    this._stack.push(op)
    op.apply()
  }

  pop(): IOperation | undefined {
    const op = this._stack.pop()
    op && op.undo()
    return op
  }

  empty(): boolean {
    return !this._stack.length
  }

  clear(): void {
    const n = this._stack.length
    for (let _ = 0; _ < n; _++) this.pop()
  }

  get length(): number {
    return this._stack.length
  }
}

export {}

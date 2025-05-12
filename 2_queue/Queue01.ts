class Queue01<T> {
  private _q0: T[] = []
  private _q1: T[] = []

  appendLeft(item: T): void {
    this._q0.push(item)
  }

  append(item: T): void {
    this._q1.push(item)
  }

  popLeft(): T | undefined {
    if (this._q0.length) return this._q0.pop()
    const tmp = this._q0
    this._q0 = this._q1
    this._q1 = tmp
    return this._q0.pop()
  }

  get length(): number {
    return this._q0.length + this._q1.length
  }
}

export {}

if (require.main === module) {
  // 用于实现01bfs
  const q = new Queue01<number>()
  q.appendLeft(1)
  q.appendLeft(2)
  q.append(3)
  console.log(q.popLeft()) // 2
  console.log(q.popLeft()) // 1
  console.log(q.popLeft()) // 3
  console.log(q.length) // 0
}

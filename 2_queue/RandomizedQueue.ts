// RandomizedQueue
// 随机队列

/**
 * 随机队列.
 */
class RandomizedQueue<T> {
  private readonly _data: T[] = []

  front(): T | undefined {
    this._select()
    return this._data[this._data.length - 1]
  }

  push(x: T): void {
    this._data.push(x)
  }

  shift(): T | undefined {
    this._select()
    return this._data.pop()
  }

  get size(): number {
    return this._data.length
  }

  private _select(): void {
    const randint = ~~(Math.random() * this._data.length)
    const last = this._data.length - 1
    const tmp = this._data[last]
    this._data[last] = this._data[randint]
    this._data[randint] = tmp
  }
}

export {}

if (require.main === module) {
  const q = new RandomizedQueue<number>()
  q.push(1)
  q.push(2)
  q.push(3)
  q.push(4)
  q.push(5)
  console.log(q.shift())
  console.log(q.shift())
  console.log(q.shift())
  console.log(q.shift())
  console.log(q.shift())
  console.log(q.shift())
}

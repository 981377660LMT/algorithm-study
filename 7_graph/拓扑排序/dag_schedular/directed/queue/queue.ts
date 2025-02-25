export class Queue<T> {
  private _data: T[]
  private _head: number
  private _tail: number
  private _size: number

  get size() {
    return this._size
  }

  constructor(initialCapacity = 0) {
    this._data = new Array(initialCapacity)
    this._head = 0
    this._tail = 0
    this._size = 0
  }

  enqueue(value: T): void {
    if (this._size === this._data.length) {
      // if buffer is full, resize to twice the size
      const newData = new Array(this._data.length * 2)

      for (let i = 0; i < this._size; i++) {
        newData[i] = this._data[(this._head + i) % this._data.length]
      }

      this._data = newData
      this._head = 0
      this._tail = this._size
    }

    this._data[this._tail] = value
    this._tail = (this._tail + 1) % this._data.length
    this._size++
  }

  dequeue(): T | undefined {
    if (this._size === 0) {
      return undefined
    }

    const value = this._data[this._head]
    this._head = (this._head + 1) % this._data.length
    this._size--

    // if the buffer is less than a quarter of the size, resize to half
    if (this._size > 0 && this._size <= this._data.length / 4) {
      const newData = new Array(Math.floor(this._data.length / 2))

      for (let i = 0; i < this._size; i++) {
        newData[i] = this._data[(this._head + i) % this._data.length]
      }

      this._data = newData
      this._head = 0
      this._tail = this._size
    }

    return value
  }

  peek(): T | undefined {
    if (this._size === 0) {
      return undefined
    }

    return this._data[this._head]
  }

  isEmpty(): boolean {
    return this._size === 0
  }
}

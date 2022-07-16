import { SimpleQueue } from './Deque/Queue'

class MovingAverage {
  private readonly _queue = new SimpleQueue()
  private _sum = 0
  private _size: number

  constructor(size: number) {
    this._size = size
  }

  // 计算滑动窗口平均值
  next(val: number): number {
    if (this._queue.length >= this._size) this._sum -= this._queue.shift()!
    this._sum += val
    this._queue.push(val)
    return this._sum / this._queue.length
  }
}

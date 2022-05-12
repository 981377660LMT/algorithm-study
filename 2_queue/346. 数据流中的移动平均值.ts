import { ArrayDeque } from './Deque/ArrayDeque'

class MovingAverage {
  private queue: ArrayDeque
  private size: number
  private sum: number

  constructor(size: number) {
    this.queue = new ArrayDeque(size)
    this.size = size
    this.sum = 0
  }

  next(val: number): number {
    if (this.queue.length >= this.size) {
      this.sum -= this.queue.shift()!
    }
    this.sum += val
    this.queue.push(val)
    return this.sum / this.queue.length
  }
}

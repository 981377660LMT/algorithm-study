/**
 * @description 特点
 * @description 应用场景
 * @description 时间复杂度 O(n)
 * @description 空间复杂度 O(n)
 */
class Counter {
  constructor(private queue: number[] = []) {}

  ping(time: number) {
    this.queue.push(time)
    // 不满足的queue[0]全部出队列
    while (this.queue[0] + 3000 < time) {
      this.queue.shift()
    }

    return this.queue.length
  }
}

const counter = new Counter()
console.log(counter.ping(1))
console.log(counter.ping(2))
console.log(counter.ping(3001))
console.log(counter.ping(3002))

export {}

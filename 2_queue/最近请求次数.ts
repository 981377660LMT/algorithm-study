import { LinkedList } from '../3_linkedList/LinkedList'
import { ArrayDeque } from './Deque'

/**
 * @description 特点
 * @description 应用场景
 * @description 时间复杂度 O(n)
 * @description 空间复杂度 O(n)
 */
class RecentCounter {
  constructor(private queue = new ArrayDeque(10000)) {}

  ping(time: number) {
    this.queue.push(time)
    // 不满足的queue[0]全部出队列
    while (this.queue.length && this.queue.front()! + 3000 < time) {
      this.queue.shift()
    }

    return this.queue.length
  }
}
// class RecentCounter {
//   constructor(private queue = new LinkedList()) {}

//   ping(time: number) {
//     this.queue.push(time)
//     // 不满足的queue[0]全部出队列
//     while (this.queue.length && this.queue.first! + 3000 < time) {
//       this.queue.shift()
//     }

//     return this.queue.length
//   }
// }

const counter = new RecentCounter()
console.log(counter.ping(1))
console.log(counter.ping(2))
console.log(counter.ping(3001))
console.log(counter.ping(3002))

export {}

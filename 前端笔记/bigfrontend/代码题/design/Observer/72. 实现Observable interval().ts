import type { Func } from '../../typings'
import { Observable } from './57. 实现Observable'

/**
 * @param {number} period
 * @return {Observable}
 * 创建一个Observable，这个Observable会在设定的interval中产生value。
 */
function interval(period: number): Observable {
  let count = 0
  return new Observable(subscriber => {
    setInterval(() => {
      subscriber.next(count++)
    }, period)
  })
}

interval(1000).subscribe(console.log)
// 上述代码会间隔一秒地输出 0, 1, 2 ....

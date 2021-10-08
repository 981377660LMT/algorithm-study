// You can use Observer which is bundled to your code
import { from } from './70. 实现Observable.from()'

import type { Id, Observable, Subscriber } from './typings'

class Subject implements Observable, Subscriber {
  private id: number
  private subRecord: Map<Id, Subscriber>

  constructor() {
    this.id = 0
    this.subRecord = new Map()
  }

  subscribe(subscriber: Subscriber) {
    const id = this.id++
    this.subRecord.set(id, subscriber)
    return {
      unsubscribe: () => {
        this.subRecord.delete(id)
      },
    }
  }

  next(value: any) {
    ;[...this.subRecord.values()].forEach(subscriber => {
      subscriber.next(value)
    })
  }

  error(error: Error) {
    ;[...this.subRecord.values()].forEach(subscriber => {
      subscriber.error(error)
    })
  }

  complete() {
    ;[...this.subRecord.values()].forEach(subscriber => {
      subscriber.complete()
    })
  }
}

if (require.main === module) {
  // 你可以从log上看出来两个订阅是相互独立的，所以log也是分开的。
  const observable = from([1, 2, 3])
  observable.subscribe(console.log)
  observable.subscribe(console.log)
  // 1
  // 2
  // 3
  // 1
  // 2
  // 3

  // 如果用了Subject，则变得更像DOM的事件机制了。
  // 因为Subject首先按照observer工作，获取值，然后像一个Observable工作，把value 传递给多个不同的observer。
  const subject = new Subject()
  subject.subscribe({ next: console.log, error: console.error, complete: console.log })
  subject.subscribe({ next: console.log, error: console.error, complete: console.log })
  subject.next(1)
  subject.next(2)
  subject.next(3)
  // 1
  // 1
  // 2
  // 2
  // 3
  // 3
}

// observable vs subject
// Observable 只有subscribe接口  只能接收观察者
// Subject就像一个群 兼有subscribe接口和observer的三个回调接口 既能被观察者订阅 又能随时触发订阅自己的观察者(订阅者)

// 发展顺序即observer=>observable=>subscriber=>subscription=>subject
